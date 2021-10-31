package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"sync"

	css "github.com/gorilla/css/scanner"

	"golang.org/x/net/html"
	a "golang.org/x/net/html/atom"

	_ "embed"
)

//go:embed caniuse.json
var caniuseJSON []byte

const (
	WHITESPACE            = " \t\r\n\f"
	CSS_SCOPE_STYLE_ATTR  = "style_attribute"
	CSS_SCOPE_STYLE_TAG   = "style_tag"
	LIMIT_REPORT_LINES    = 30
	TWO_KEYS_MERGE_FORMAT = "%s||%s"
)

// css selectors
type CssSelectorType int

const (
	ADJACENT_SIBLING_COMBINATOR_TYPE CssSelectorType = iota
	CHAINING_SELECTORS_TYPE
	CHILD_COMBINATOR_TYPE
	CLASS_SELECTOR_TYPE
	DESCENDANT_COMBINATOR_TYPE
	GENERAL_SIBLING_COMBINATOR_TYPE
	GROUPING_SELECTORS_TYPE
	ID_SELECTOR_TYPE
	TYPE_SELECTOR_TYPE
	UNIVERSAL_SELECTOR_STAR_TYPE
)

var (
	// css selectors types
	adjacentSiblingCombinatorRe = regexp.MustCompile(`[\w\s-]\+\s?[#\.\w-]`)
	chainingSelectorsRe         = regexp.MustCompile(`([\w-])(\.|#)([\w-])`)
	childCombinatorRe           = regexp.MustCompile(`[\w\s-]>\s?[#\.\w-]`)
	descendantCombinatorRe      = regexp.MustCompile(`[\w\s-] \s?[#\.\w-]`)
	generalSiblingCombinatorRe  = regexp.MustCompile(`[\w\s-]~\s?[#\.\w-]`)
	idSelectorRe                = regexp.MustCompile(`#\w\+`)
	typeSelectorRe              = regexp.MustCompile(`(^|[\s\+~>])[\w-]`)
	universalSelectorStarRe     = regexp.MustCompile(`\*[^=]`)
)

// json config structs begin

type CaniuseDBHTMLTag struct {
	Notes map[string]string                         `json:"notes"`
	Stats map[string]map[string]map[string][]string `json:"stats"`
	Url   string                                    `json:"url"`
}

type CaniuseDB struct {
	HtmlTags map[string]map[string]CaniuseDBHTMLTag `json:"html_tags"`
}

var rulesDB CaniuseDB

// json config structs end

// result structure begin

type HTMLTagReport struct {
	Rules     CaniuseDBHTMLTag `json:"rules"`
	Lines     map[int]int      `json:"lines"`
	MoreLines bool             `json:"more_lines"`
}

type ParseReport struct {
	HtmlTags map[string]map[string]HTMLTagReport `json:"html_tags"`
}

// result structure end

type ParserEngine struct {
	// array for bytes -> lines
	bytesToLine []int
	// group for parallel processing
	wg sync.WaitGroup
	// lock for report
	mx sync.RWMutex
	// report itself
	pr ParseReport
	// parse time states
	isStyleTagOpen  bool
	styleTagContent string
	styleTagLine    int
}

func InitParser() (*ParserEngine, error) {
	return &ParserEngine{
		bytesToLine:     []int{},
		isStyleTagOpen:  false,
		styleTagContent: "",
	}, nil
}

func (prs *ParserEngine) ruleCssInTag(scanner *css.Scanner) error {
	var (
		isCollectPropVal bool   = false
		propName         string = ""
		propVal          string = ""
		propLine         int    = 0
	)

	for {
		token := scanner.Next()
		if token.Type == css.TokenEOF || token.Type == css.TokenError {
			log.Printf("TagStyleProperties: %v %v %v\n", propName, propVal, propLine)
			// prs.storeTagStyleProperties(propName, propVal, propLine) // store last parts, if collected
			return nil
		}

		// log.Printf("[processCssInTag] %s (line: %d, column: %d): %v (doc line: %d)", token.Type, token.Line, token.Column, token.Value, prs.styleTagLine+token.Line-1)

		switch token.Type {
		case css.TokenCDO, css.TokenCDC, css.TokenComment, css.TokenS:
			break // continue loop
		case css.TokenChar:
			switch token.Value {
			case ":":
				isCollectPropVal = true
			case ";":
				log.Printf("TagStyleProperties: %v %v %v\n", propName, propVal, propLine)
				// prs.storeTagStyleProperties(propName, propVal, propLine)
				// reset for new prop
				isCollectPropVal = false
				propName = ""
				propVal = ""
				propLine = 0
			case "}":
				log.Printf("TagStyleProperties: %v %v %v\n", propName, propVal, propLine)
				// prs.storeTagStyleProperties(propName, propVal, propLine)
				return nil
			}
		default:
			if isCollectPropVal {
				if len(propVal) == 0 {
					propVal = token.Value
				} else {
					propVal += fmt.Sprintf(" %s", token.Value)
				}
			} else {
				if len(propName) == 0 {
					propLine = prs.styleTagLine + token.Line - 1 // initial line for prop
					propName = token.Value
				} else {
					propName += fmt.Sprintf(" %s", token.Value)
				}

			}
		}
	}
}

func (prs *ParserEngine) processCssInTag(scanner *css.Scanner) error {
	var (
		atKeywordName string = ""
		atKeywordVal  string = ""
		atKeywordLine int    = 0

		isPseudoClass   bool   = false
		pseudoClassName string = ""
		pseudoClassLine int    = 0

		cssSelector     string = ""
		cssSelectorLine int    = 0
	)

	for {
		token := scanner.Next()
		if token.Type == css.TokenEOF || token.Type == css.TokenError {
			return nil
		}

		// log.Printf("[processCssInTag] %s (line: %d, column: %d): %v (doc line: %d)", token.Type, token.Line, token.Column, token.Value, prs.styleTagLine+token.Line-1)

		switch token.Type {
		case css.TokenCDO, css.TokenCDC, css.TokenComment, css.TokenS:
			if token.Type == css.TokenS && len(cssSelector) > 0 {
				cssSelector += token.Value
			}

			break // continue loop
		case css.TokenAtKeyword:
			line := prs.styleTagLine + token.Line - 1

			switch token.Value {
			case "@media", "@supports", "@document", "@keyframes", "@-webkit-keyframes", "@font-feature-values": // nested at-rules
				atKeywordName = token.Value
				atKeywordLine = line
			default: // at-rule block, like "@charset", "@import", "@namespace"
				log.Printf("AtRuleCSSStatement: %v %v %v\n", token.Value, "", line)
				// err := prs.storeAtRuleCSSStatement(token.Value, "", line)
				// if err != nil {
				// 	return err
				// }
				atKeywordName = ""
				atKeywordLine = 0
			}
			atKeywordVal = ""
			cssSelector = ""
			cssSelectorLine = 0
		case css.TokenChar:
			switch token.Value {
			case "}": // it is not required, it is more for failsafe
				atKeywordName = ""
				atKeywordVal = ""
				atKeywordLine = 0
				pseudoClassName = ""
				pseudoClassLine = 0
				cssSelector = ""
				cssSelectorLine = 0
			case ":":
				if len(atKeywordName) == 0 {
					if isPseudoClass && len(pseudoClassName) > 0 {
						log.Printf("CSSPseudoClass: %v %v\n", pseudoClassName, pseudoClassLine)
						// err := prs.storeCSSPseudoClass(pseudoClassName, pseudoClassLine)
						// if err != nil {
						// 	return err
						// }

						pseudoClassName = ""
						pseudoClassLine = 0
					}
					isPseudoClass = true

					if len(cssSelector) == 0 {
						cssSelectorLine = prs.styleTagLine + token.Line - 1
					}
					cssSelector += token.Value
				}
			case ",", ".":
				if isPseudoClass && len(pseudoClassName) > 0 {
					log.Printf("CSSPseudoClass: %v %v\n", pseudoClassName, pseudoClassLine)
					// err := prs.storeCSSPseudoClass(pseudoClassName, pseudoClassLine)
					// if err != nil {
					// 	return err
					// }
					pseudoClassName = ""
					pseudoClassLine = 0
				}
				isPseudoClass = false

				if len(cssSelector) == 0 {
					cssSelectorLine = prs.styleTagLine + token.Line - 1
				}
				cssSelector += token.Value
			case "{":
				if len(atKeywordName) == 0 && isPseudoClass && len(pseudoClassName) > 0 {
					log.Printf("CSSPseudoClass: %v %v\n", pseudoClassName, pseudoClassLine)
					// err := prs.storeCSSPseudoClass(pseudoClassName, pseudoClassLine)
					// if err != nil {
					// 	return err
					// }
					pseudoClassName = ""
					pseudoClassLine = 0
				}
				isPseudoClass = false

				if len(atKeywordName) == 0 && len(cssSelector) > 0 {
					log.Printf("CSS SELECTOR: %v %v\n", cssSelector, cssSelectorLine)
					// err := prs.processCssSelector(cssSelector, cssSelectorLine)
					// if err != nil {
					// 	return err
					// }
					cssSelector = ""
					cssSelectorLine = 0
				}

				if len(atKeywordName) > 0 {
					log.Printf("AtRuleCSSStatement: %v %v %v\n", atKeywordName, atKeywordVal, atKeywordLine)
					// err := prs.storeAtRuleCSSStatement(atKeywordName, atKeywordVal, atKeywordLine)
					// if err != nil {
					// 	return err
					// }
					err := prs.processCssInTag(scanner)
					if err != nil {
						return err
					}
					atKeywordName = ""
					atKeywordVal = ""
					atKeywordLine = 0
				} else {
					err := prs.ruleCssInTag(scanner)
					if err != nil {
						return err
					}
				}
			default:
				if len(cssSelector) == 0 {
					cssSelectorLine = prs.styleTagLine + token.Line - 1
				}
				cssSelector += token.Value
			}
		default:
			if len(atKeywordName) > 0 {
				if len(atKeywordVal) == 0 {
					atKeywordVal = token.Value
				} else {
					atKeywordVal += fmt.Sprintf(" %s", token.Value)
				}
			} else {
				if len(cssSelector) == 0 {
					cssSelectorLine = prs.styleTagLine + token.Line - 1
				}
				cssSelector += token.Value
			}

			if isPseudoClass {
				if len(pseudoClassName) == 0 {
					pseudoClassName = token.Value
					pseudoClassLine = prs.styleTagLine + token.Line - 1
				} else {
					pseudoClassName += fmt.Sprintf(" %s", token.Value)
				}
			}
		}
	}
}

func makeInitialHtmlReport(position int, ruleTagAttrData CaniuseDBHTMLTag) HTMLTagReport {
	lines := make(map[int]int)
	lines[position] = 1

	return HTMLTagReport{
		Rules:     ruleTagAttrData,
		Lines:     lines,
		MoreLines: false,
	}
}

func (prs *ParserEngine) saveToReportHtmlTag(tagName, tagAttr string, position int, ruleTagAttrData CaniuseDBHTMLTag) {
	prs.mx.Lock()
	defer prs.mx.Unlock()

	if tagData, ok := prs.pr.HtmlTags[tagName]; ok {
		if tagAttrData, ok := tagData[tagAttr]; ok {
			if len(tagAttrData.Lines) < LIMIT_REPORT_LINES {
				tagAttrData.Lines[position] = 1
			} else {
				tagAttrData.MoreLines = true
			}
			tagData[tagAttr] = tagAttrData
			prs.pr.HtmlTags[tagName] = tagData
		} else {
			if len(prs.pr.HtmlTags[tagName]) > 0 {
				prs.pr.HtmlTags[tagName][tagAttr] = makeInitialHtmlReport(position, ruleTagAttrData)
			} else {
				rootData := make(map[string]HTMLTagReport)
				rootData[tagAttr] = makeInitialHtmlReport(position, ruleTagAttrData)
				prs.pr.HtmlTags[tagName] = rootData
			}
		}
	} else {
		rData := make(map[string]HTMLTagReport)
		rData[tagAttr] = makeInitialHtmlReport(position, ruleTagAttrData)

		if len(prs.pr.HtmlTags) > 0 {
			prs.pr.HtmlTags[tagName] = rData
		} else {
			rootData := make(map[string]map[string]HTMLTagReport)
			rootData[tagName] = rData
			prs.pr.HtmlTags = rootData
		}
	}
}

func (prs *ParserEngine) checkHtmlTags(tagName string, attrs []html.Attribute, position int) error {
	if len(tagName) == 0 {
		return nil
	}

	tagName = strings.ToLower(tagName)

	if ruleTagData, ok := rulesDB.HtmlTags[tagName]; ok {
		if ruleTagAttrData, ok := ruleTagData[""]; ok {
			prs.saveToReportHtmlTag(tagName, "", position, ruleTagAttrData)
		}
		for _, att := range attrs {
			attrKey := strings.ToLower(att.Key)
			attrVal := strings.ToLower(att.Val)

			if ruleTagAttrData, ok := ruleTagData[attrKey]; ok {
				prs.saveToReportHtmlTag(tagName, attrKey, position, ruleTagAttrData)
			}

			attrWithVal := fmt.Sprintf(TWO_KEYS_MERGE_FORMAT, attrKey, attrVal)
			if ruleTagAttrData, ok := ruleTagData[attrWithVal]; ok {
				prs.saveToReportHtmlTag(tagName, attrWithVal, position, ruleTagAttrData)
			}

			if attrKey == "style" {
				// prs.storeInlinedStyleProperties(tx, tag, attrKey, attrVal, position)
			}
		}
	}

	log.Printf("checkHtmlTags: %v %v %v\n", tagName, attrs, position)

	return nil
}

func (prs *ParserEngine) processHtmlToken(htmlTokenizer *html.Tokenizer, token html.Token, tagLine int) error {
	switch token.Type {
	case html.TextToken:
		if prs.isStyleTagOpen {
			prs.styleTagContent += strings.Replace(token.Data, "\x00", "\ufffd", -1) // replace NULL
		}
	case html.StartTagToken:
		switch token.DataAtom {
		case a.Style:
			prs.isStyleTagOpen = true
			prs.styleTagLine = tagLine
		}
		// process html tag
		err := prs.checkHtmlTags(token.Data, token.Attr, tagLine)
		if err != nil {
			return err
		}

	case html.EndTagToken:
		switch token.DataAtom {
		case a.Style:
			if prs.isStyleTagOpen && len(prs.styleTagContent) > 0 {
				scanner := css.New(prs.styleTagContent)
				err := prs.processCssInTag(scanner)
				if err != nil {
					// Ignore or should report?
				}
			}
			// reset style tag storage
			prs.isStyleTagOpen = false
			prs.styleTagContent = ""
			prs.styleTagLine = 0
		}
	case html.SelfClosingTagToken:
		// process html tag
		err := prs.checkHtmlTags(token.Data, token.Attr, tagLine)
		if err != nil {
			return err
		}
	}

	return nil
}

func (prs ParserEngine) getLineFromOffset(tagLine, offset int) int {
	// binary search used before, but looks like it waste of time - we can just "follow" bytes offset
	// return sort.Search(len(prs.bytesToLine), func(i int) bool { return prs.bytesToLine[i] > offset })
	for {
		if len(prs.bytesToLine) <= tagLine {
			return tagLine
		}

		if prs.bytesToLine[tagLine] > offset {
			return tagLine
		}

		tagLine += 1
	}
}

func (prs *ParserEngine) calulateNewlineBytePos(document []byte) {
	var (
		bytesToLine []int
		cursorPos   int = 0
	)

	lines := bytes.Split(document, []byte("\n"))

	for _, line := range lines {
		bytesToLine = append(bytesToLine, cursorPos)
		cursorPos += len(line) + 1 // "\n" provide additional byte
	}

	prs.bytesToLine = bytesToLine
}

func (prs *ParserEngine) Report(document []byte) (*ParseReport, error) {
	var (
		err             error
		htmlTokenizer   *html.Tokenizer
		htmlBytesOffset int = 0
		tagLine         int = 0
	)

	prs.calulateNewlineBytePos(document)

	htmlTokenizer = html.NewTokenizer(bytes.NewReader(document))
	for err != io.EOF {
		// CDATA sections are not alowed
		htmlTokenizer.AllowCDATA(false)
		// Read and parse the next token.
		htmlTokenizer.Next()
		tt := htmlTokenizer.Token()
		if tt.Type == html.ErrorToken {
			err = htmlTokenizer.Err()
			if err != nil && err != io.EOF {
				return nil, err
			}
		}

		tagLine = prs.getLineFromOffset(tagLine, htmlBytesOffset)
		err := prs.processHtmlToken(htmlTokenizer, tt, tagLine)
		if err != nil {
			return nil, err
		}

		htmlBytesOffset += len(htmlTokenizer.Raw())
	}

	return &prs.pr, nil
}

func init() {
	// parse caniuse.json here
	if err := json.Unmarshal(caniuseJSON, &rulesDB); err != nil {
		panic(err)
	}
}

func ReportFromHTML(document []byte) (*ParseReport, error) {
	parser, err := InitParser()
	if err != nil {
		return nil, err
	}

	report, err := parser.Report(document)
	if err != nil {
		return nil, err
	}

	return report, nil
}
