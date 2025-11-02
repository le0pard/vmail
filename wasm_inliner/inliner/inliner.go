package inliner

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	parse "github.com/tdewolff/parse/v2"
	css "github.com/tdewolff/parse/v2/css"

	"github.com/andybalholm/cascadia"

	"golang.org/x/net/html"
	a "golang.org/x/net/html/atom"
)

const (
	WHITESPACE = " \t\r\n\f"
)

var (
	mediaSplitInlineRe = regexp.MustCompile(`(?i)[\s]+|,`)
	mediaInlineRe      = regexp.MustCompile(`(?i)(screen|handheld|all)`)
	resetSelectors     = regexp.MustCompile(`(?i)^(\#outlook|body.*|\.ReadMsgBody|\.ExternalClass|img|table|td|p|\#backgroundTable|\#bodyTable)`) // email reset styles
)

type StylesheetsTags struct {
	Parent  *html.Node
	Node    *html.Node
	Content string
}

type CSSGroupSelectors struct {
	Key      string
	NotApply bool
}

type CSSSelectors struct {
	Selectors       []CSSGroupSelectors
	Attributes      map[string]string
	AttributesOrder []string
}

type InlineEngine struct {
	// group for parallel processing
	wg sync.WaitGroup
	// lock for report
	mx sync.RWMutex
}

func InitInliner() *InlineEngine {
	return &InlineEngine{}
}

func indexOf(data []string, element string) int {
	for i, v := range data {
		if element == v {
			return i
		}
	}
	return -1 //not found
}

func extractHeadBodyAndStylesheets(doc *html.Node) (*html.Node, *html.Node, string, error) {
	var (
		externalStylesheets []StylesheetsTags
		contents            []string
		head                *html.Node = nil
		body                *html.Node = nil
		crawler             func(*html.Node)
		netClient           = &http.Client{
			Timeout: 5 * time.Second,
		}
	)

	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode {
			if node.Data == "head" && head == nil {
				head = node
			}
			if node.Data == "body" && body == nil {
				body = node
			}

			if node.Data == "style" {
				var (
					isMediaGood bool = true
				)

				for _, v := range node.Attr {
					key := strings.ToLower(strings.Trim(v.Key, WHITESPACE))
					val := strings.Trim(v.Val, WHITESPACE)

					if key == "media" {
						if len(val) > 0 {
							isMediaGood = false
							for _, mediaPart := range mediaSplitInlineRe.Split(strings.ToLower(val), -1) {
								if mediaInlineRe.MatchString(strings.Trim(mediaPart, WHITESPACE)) {
									isMediaGood = true
								}
							}
						}
					}
				}

				if isMediaGood && node.FirstChild.Type == html.TextNode && len(node.FirstChild.Data) > 0 {
					externalStylesheets = append(externalStylesheets, StylesheetsTags{
						Parent:  node.Parent,
						Node:    node,
						Content: node.FirstChild.Data,
					})
				}
			}

			if node.Data == "link" {
				var (
					isStylesheet bool   = false
					isMediaGood  bool   = true
					foundedHref  string = ""
				)

				for _, v := range node.Attr {
					key := strings.ToLower(strings.Trim(v.Key, WHITESPACE))
					val := strings.Trim(v.Val, WHITESPACE)

					if key == "rel" && strings.ToLower(val) == "stylesheet" {
						isStylesheet = true
					}

					if key == "media" {
						if len(val) > 0 {
							isMediaGood = false
							for _, mediaPart := range mediaSplitInlineRe.Split(strings.ToLower(val), -1) {
								if mediaInlineRe.MatchString(strings.Trim(mediaPart, WHITESPACE)) {
									isMediaGood = true
								}
							}
						}
					}

					if key == "href" {
						foundedHref = val
					}
				}

				if isStylesheet && isMediaGood && len(foundedHref) > 0 {
					sheetUrl, err := url.ParseRequestURI(foundedHref)
					if err != nil {
						return
					}

					resp, err := netClient.Get(sheetUrl.String())
					if err != nil {
						return
					}
					defer resp.Body.Close()

					body, err := io.ReadAll(resp.Body)
					if err != nil {
						return
					}

					if len(body) > 0 {
						externalStylesheets = append(externalStylesheets, StylesheetsTags{
							Parent:  node.Parent,
							Node:    node,
							Content: string(body),
						})
					}
				}
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)

	for _, item := range externalStylesheets {
		contents = append(contents, item.Content)
		item.Parent.RemoveChild(item.Node)
	}

	return head, body, strings.Join(contents, "\n"), nil
}

func (inlr *InlineEngine) collectStyles(p *css.Parser) (string, error) {
	var (
		collectedCSS      string = ""
		countClosedStyles int    = 0
	)

	for {
		gt, _, data := p.Next()

		// log.Printf("[collectStyles]: %v - %v - %v\n", gt, string(data), p.Values())

		if gt == css.ErrorGrammar {
			return collectedCSS, nil
		}

		switch gt {
		case css.BeginRulesetGrammar, css.BeginAtRuleGrammar:
			countClosedStyles += 1
			collectedCSS += string(data)
			for _, val := range p.Values() {
				collectedCSS += string(val.Data)
			}
			collectedCSS += "{"
		case css.QualifiedRuleGrammar:
			collectedCSS += string(data)
			for _, val := range p.Values() {
				collectedCSS += string(val.Data)
			}
			collectedCSS += ","
		case css.EndRulesetGrammar, css.EndAtRuleGrammar:
			if countClosedStyles <= 0 {
				return collectedCSS, nil
			} else {
				collectedCSS += string(data)
				countClosedStyles -= 1
			}
		case css.DeclarationGrammar, css.CustomPropertyGrammar:
			collectedCSS += fmt.Sprintf("%s:", string(data))
			for _, val := range p.Values() {
				collectedCSS += string(val.Data)
			}
			collectedCSS += ";"
		}
	}
}

func converCssAttributesToString(orderKeys []string, attrs map[string]string) string {
	output := ""
	for _, k := range orderKeys {
		output += k
		output += ":"
		output += strings.ReplaceAll(attrs[k], "!important", "")
		output += ";"
	}
	return output
}

func converCssSelectorToString(selector string, orderKeys []string, attrs map[string]string) string {
	output := selector
	output += "{"
	output += converCssAttributesToString(orderKeys, attrs)
	output += "}"
	return output
}

func (inlr *InlineEngine) inlineRulesetToTags(doc *html.Node, cssStore CSSSelectors) (string, error) {
	var (
		additionalCSS string = ""
	)

	for _, selectorGroup := range cssStore.Selectors {
		if selectorGroup.NotApply {
			additionalCSS += converCssSelectorToString(selectorGroup.Key, cssStore.AttributesOrder, cssStore.Attributes)
			continue
		}

		if resetSelectors.MatchString(selectorGroup.Key) {
			additionalCSS += converCssSelectorToString(selectorGroup.Key, cssStore.AttributesOrder, cssStore.Attributes)
			continue
		}

		// log.Printf("[selectorGroup]: %v\n", selectorGroup)

		selector, err := cascadia.ParseGroup(selectorGroup.Key)
		if err != nil {
			continue
		}

		// log.Printf("[selector]: %v\n", selector)

		for _, node := range cascadia.Selector(selector.Match).MatchAll(doc) {
			// switch node.DataAtom {
			// case a.Style:
			// }

			stylesStr := ""
			newAttr := []html.Attribute{}
			for _, attr := range node.Attr {
				if strings.ToLower(attr.Key) == "style" {
					stylesStr = attr.Val
				} else {
					newAttr = append(newAttr, attr)
				}
			}

			if len(stylesStr) == 0 {
				newAttr = append(newAttr, html.Attribute{
					Key: "style",
					Val: converCssAttributesToString(cssStore.AttributesOrder, cssStore.Attributes),
				})
				node.Attr = newAttr
				continue
			}

			p := css.NewParser(parse.NewInput(bytes.NewBufferString(stylesStr)), true)
			applyOrders := []string{}
			applyAttrs := make(map[string]string, len(cssStore.AttributesOrder))
			for _, k := range cssStore.AttributesOrder {
				applyOrders = append(applyOrders, k)
				applyAttrs[k] = cssStore.Attributes[k]
			}
			newAttrStr := ""
			for {
				gt, _, data := p.Next()

				if gt == css.ErrorGrammar {
					newAttrStr += converCssAttributesToString(applyOrders, applyAttrs)
					break
				} else if gt == css.AtRuleGrammar || gt == css.BeginAtRuleGrammar || gt == css.BeginRulesetGrammar {
					newAttrStr += strings.ToLower(string(data))
					for _, val := range p.Values() {
						newAttrStr += strings.ToLower(string(val.Data))
					}
					if gt == css.BeginAtRuleGrammar || gt == css.BeginRulesetGrammar {
						newAttrStr += "{"
					} else if gt == css.AtRuleGrammar {
						newAttrStr += ";"
					}
				} else if gt == css.DeclarationGrammar {
					attrKey := strings.ToLower(string(data))
					if newVal, ok := applyAttrs[attrKey]; ok {
						applyOrderIndex := indexOf(applyOrders, attrKey)

						if applyOrderIndex >= 0 {
							newAttrStr += converCssAttributesToString(
								[]string{
									attrKey,
								},
								map[string]string{
									attrKey: newVal,
								},
							)

							applyOrders = append(
								applyOrders[:applyOrderIndex],
								applyOrders[applyOrderIndex+1:]...,
							)
							delete(applyAttrs, attrKey)
						}
						continue
					}

					newAttrStr += attrKey
					newAttrStr += ":"
					for _, val := range p.Values() {
						newAttrStr += strings.ToLower(string(val.Data))
					}
					newAttrStr += ";"
				} else {
					newAttrStr += strings.ToLower(string(data))
				}
			}

			newAttr = append(newAttr, html.Attribute{
				Key: "style",
				Val: newAttrStr,
			})
			node.Attr = newAttr
		}
	}

	return additionalCSS, nil
}

func (inlr *InlineEngine) inlineStyleSheetContent(doc *html.Node, sheetContent string) (string, error) {
	var (
		cssStore CSSSelectors = CSSSelectors{
			Selectors:       []CSSGroupSelectors{},
			Attributes:      make(map[string]string),
			AttributesOrder: []string{},
		}
		notAppliedCss string = ""
	)

	// log.Printf("[sheetContent]: %v\n", sheetContent)

	p := css.NewParser(parse.NewInput(bytes.NewBufferString(sheetContent)), false)
	for {
		gt, _, data := p.Next()

		// log.Printf("[inlineStyleSheetContent]: %v - %v - %v\n", gt, string(data), p.Values())

		if gt == css.ErrorGrammar {
			return notAppliedCss, nil
		}

		switch gt {
		case css.AtRuleGrammar:
			notAppliedCss += string(data)
			for _, val := range p.Values() {
				notAppliedCss += string(val.Data)
			}
			notAppliedCss += ";"
		case css.BeginAtRuleGrammar:
			notAppliedCss += string(data)
			for _, val := range p.Values() {
				notAppliedCss += string(val.Data)
			}
			notAppliedCss += "{"
			additionalCss, err := inlr.collectStyles(p)
			if err == nil {
				notAppliedCss += additionalCss
			}
			notAppliedCss += "}"
		case css.QualifiedRuleGrammar, css.BeginRulesetGrammar:
			qselector := string(data)
			notApply := false
			for _, val := range p.Values() {
				if val.TokenType == css.CommaToken {
					cssStore.Selectors = append(cssStore.Selectors, CSSGroupSelectors{
						Key:      qselector,
						NotApply: notApply,
					})
					qselector = string(data)
					notApply = false
				} else {
					qselector += string(val.Data)
				}
				if val.TokenType == css.ColonToken {
					notApply = true
				}
			}
			cssStore.Selectors = append(cssStore.Selectors, CSSGroupSelectors{
				Key:      qselector,
				NotApply: notApply,
			})
		case css.DeclarationGrammar, css.CustomPropertyGrammar:
			cssKey := strings.ToLower(string(data))
			cssVal := ""
			for _, val := range p.Values() {
				cssVal += string(val.Data)
			}
			cssStore.AttributesOrder = append(cssStore.AttributesOrder, cssKey)
			cssStore.Attributes[cssKey] = strings.ToLower(cssVal)
		case css.EndRulesetGrammar:
			if len(cssStore.Selectors) > 0 {
				additionalCss, err := inlr.inlineRulesetToTags(doc, cssStore)
				if err == nil {
					notAppliedCss += additionalCss
				}
			}
			cssStore = CSSSelectors{
				Selectors:       []CSSGroupSelectors{},
				Attributes:      make(map[string]string),
				AttributesOrder: []string{},
			}
		}
	}
	return notAppliedCss, nil
}

func (inlr *InlineEngine) addNonAppliedCssToDom(doc *html.Node, sheetContent string) {
	styleContent := &html.Node{
		Type: html.TextNode,
		Data: sheetContent,
	}

	styleTag := &html.Node{
		Type:     html.ElementNode,
		Data:     "style",
		DataAtom: a.Style,
		Attr: []html.Attribute{
			html.Attribute{
				Key: "type",
				Val: "text/css",
			},
		},
	}

	styleTag.AppendChild(styleContent)
	doc.AppendChild(styleTag)
}

func (inlr *InlineEngine) InlineCss(htmlDoc []byte) ([]byte, error) {
	var (
		doc *html.Node
		err error
	)

	if len(htmlDoc) == 0 {
		return htmlDoc, nil // empty doc
	}

	if doc, err = html.Parse(bytes.NewReader(htmlDoc)); err != nil {
		return []byte{}, err
	}

	head, body, stylesheetContents, err := extractHeadBodyAndStylesheets(doc)
	if err != nil {
		return []byte{}, err
	}

	if head == nil {
		head = doc // no head, use html as root
	}

	if body == nil {
		body = doc // no body, use html as root
	}

	notAppliedCss, err := inlr.inlineStyleSheetContent(body, stylesheetContents)
	if err != nil {
		return []byte{}, err
	}

	if len(notAppliedCss) > 0 {
		inlr.addNonAppliedCssToDom(head, notAppliedCss)
	}

	inlr.wg.Wait() // wait all jobs

	buf := new(bytes.Buffer)
	if err = html.Render(buf, doc); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func InlineCssInHTML(htmlDoc []byte) ([]byte, error) {
	inliner := InitInliner()
	newHtmlDoc, err := inliner.InlineCss(htmlDoc)
	if err != nil {
		return nil, err
	}

	return newHtmlDoc, nil
}
