package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	parse "github.com/tdewolff/parse/v2"
	css "github.com/tdewolff/parse/v2/css"

	"golang.org/x/net/html"
	a "golang.org/x/net/html/atom"

	_ "embed"
)

//go:embed caniuse.json
var caniuseJSON []byte

const (
	WHITESPACE            = " \t\r\n\f"
	LIMIT_REPORT_LINES    = 50
	TWO_KEYS_MERGE_FORMAT = "%s||%s"
)

// css selectors
type CssSelectorType int

const (
	ADJACENT_SIBLING_COMBINATOR_TYPE CssSelectorType = iota // The adjacent sibling combinator (`h1 + p`) allows to target an element that is directly after another.
	ATTRIBUTE_SELECTOR_TYPE                                 // The attribute selector (`[attr]`) targets elements with this specific attribute.
	CHAINING_SELECTORS_TYPE                                 // Chaining selectors (`.foo.bar`) allows to apply styles to elements matching all the corresponding selectors.
	CHILD_COMBINATOR_TYPE                                   // The child combinator is represented by a superior sign (`>`) between two selectors and matches the second selector if it is a direct child of the first selector.
	CLASS_SELECTOR_TYPE                                     // The class selector (`.className`) allows to apply styles to a group of elements with the corresponding `class` attribute.
	DESCENDANT_COMBINATOR_TYPE                              // The descendant combinator is represented by a space (` `) between two selectors and matches the second selector if it has ancestor matching the first selector.
	GENERAL_SIBLING_COMBINATOR_TYPE                         // The general sibling combinator (`img ~ p`) allows to target any element that after another (directly or indirectly).
	GROUPING_SELECTORS_TYPE                                 // Grouping selectors (`.foo, .bar`) allows to apply the same styles to the different corresponding elements.
	ID_SELECTOR_TYPE                                        // The ID selector (`#id`) allows to apply styles to an element with the corresponding `id` attribute.
	TYPE_SELECTOR_TYPE                                      // Type selector or element selectors allow to apply styles by HTML element names.
	UNIVERSAL_SELECTOR_STAR_TYPE                            // The universal selector (`*`) allows to apply styles to every elements.
)

var (
	html5DoctypeRe = regexp.MustCompile(`(?i)<!DOCTYPE\s+html>`)
	anchorLinkRe   = regexp.MustCompile(`(?i)^#(.+)`)
	mailtoLinkRe   = regexp.MustCompile(`(?i)^mailto:(.+)`)
	dimentionsRe   = regexp.MustCompile(`(\+|-)?(\d(\.\d)?|\.\d)`) // based on https://developer.mozilla.org/en-US/docs/Web/CSS/dimension
	cssUrlRe       = regexp.MustCompile(`url\(['"\s]?(.*?)['"\s]?\)`)
)

func (d CssSelectorType) String() string {
	return []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}[d]
}

//

var (
	normalizeCssPropsMap = map[string]string{
		"margin-top":     "margin",
		"margin-bottom":  "margin",
		"margin-left":    "margin",
		"margin-right":   "margin",
		"padding-top":    "padding",
		"padding-bottom": "padding",
		"padding-left":   "padding",
		"padding-right":  "padding",
	}
)

func normalizeCssProp(prop string) string {
	if newProp, ok := normalizeCssPropsMap[prop]; ok {
		return newProp
	}
	return prop
}

// json config structs begin

type CaniuseDB struct {
	HtmlTags            map[string]map[string]interface{} `json:"html_tags"`
	HtmlAttributes      map[string]map[string]interface{} `json:"html_attributes"`
	CssProperties       map[string]map[string]interface{} `json:"css_properties"`
	AtRuleCssStatements map[string]map[string]interface{} `json:"at_rule_css_statements"`
	CssSelectorTypes    map[string]interface{}            `json:"css_selector_types"`
	CssDimentions       map[string]interface{}            `json:"css_dimentions"`
	CssFunctions        map[string]interface{}            `json:"css_functions"`
	CssPseudoSelectors  map[string]interface{}            `json:"css_pseudo_selectors"`
	ImgFormats          map[string]interface{}            `json:"img_formats"`
	LinkTypes           map[string]interface{}            `json:"link_types"`
	CssVariables        interface{}                       `json:"css_variables"`
	CssImportant        interface{}                       `json:"css_important"`
	Html5Doctype        interface{}                       `json:"html5_doctype"`
}

var rulesDB CaniuseDB

// json config structs end

// result structure begin

type ReportContainer struct {
	Rules     interface{}  `json:"rules"`
	Lines     map[int]bool `json:"lines"`
	MoreLines bool         `json:"more_lines"`
}

type ParseReport struct {
	HtmlTags            map[string]map[string]ReportContainer `json:"html_tags"`
	HtmlAttributes      map[string]map[string]ReportContainer `json:"html_attributes"`
	CssProperties       map[string]map[string]ReportContainer `json:"css_properties"`
	AtRuleCssStatements map[string]map[string]ReportContainer `json:"at_rule_css_statements"`
	CssSelectorTypes    map[string]ReportContainer            `json:"css_selector_types"`
	CssDimentions       map[string]ReportContainer            `json:"css_dimentions"`
	CssFunctions        map[string]ReportContainer            `json:"css_functions"`
	CssPseudoSelectors  map[string]ReportContainer            `json:"css_pseudo_selectors"`
	ImgFormats          map[string]ReportContainer            `json:"img_formats"`
	LinkTypes           map[string]ReportContainer            `json:"link_types"`
	CssVariables        ReportContainer                       `json:"css_variables"`
	CssImportant        ReportContainer                       `json:"css_important"`
	Html5Doctype        ReportContainer                       `json:"html5_doctype"`
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

func InitParser() *ParserEngine {
	return &ParserEngine{
		bytesToLine:     []int{},
		isStyleTagOpen:  false,
		styleTagContent: "",
	}
}

func makeInitialReportContainer(position int, ruleCssPropData interface{}) ReportContainer {
	lines := make(map[int]bool)
	lines[position] = true

	return ReportContainer{
		Rules:     ruleCssPropData,
		Lines:     lines,
		MoreLines: false,
	}
}

func (prs *ParserEngine) saveToReportHtmlAttributes(attrKey, attrVal string, position int, ruleCssPropData interface{}) {
	prs.mx.Lock()
	defer prs.mx.Unlock()

	if prKeyData, ok := prs.pr.HtmlAttributes[attrKey]; ok {
		if prValData, ok := prKeyData[attrVal]; ok {
			if len(prValData.Lines) < LIMIT_REPORT_LINES {
				prValData.Lines[position] = true
			} else {
				prValData.MoreLines = true
			}
			prKeyData[attrVal] = prValData
			prs.pr.HtmlAttributes[attrKey] = prKeyData
		} else {
			if len(prs.pr.HtmlAttributes[attrKey]) > 0 {
				prs.pr.HtmlAttributes[attrKey][attrVal] = makeInitialReportContainer(position, ruleCssPropData)
			} else {
				rootData := make(map[string]ReportContainer)
				rootData[attrVal] = makeInitialReportContainer(position, ruleCssPropData)
				prs.pr.HtmlAttributes[attrKey] = rootData
			}
		}
	} else {
		rData := make(map[string]ReportContainer)
		rData[attrVal] = makeInitialReportContainer(position, ruleCssPropData)

		if len(prs.pr.HtmlAttributes) > 0 {
			prs.pr.HtmlAttributes[attrKey] = rData
		} else {
			rootData := make(map[string]map[string]ReportContainer)
			rootData[attrKey] = rData
			prs.pr.HtmlAttributes = rootData
		}
	}
}

func (prs *ParserEngine) saveToReportCssVariables(position int) {
	if len(prs.pr.CssVariables.Lines) > 0 {
		if len(prs.pr.CssVariables.Lines) < LIMIT_REPORT_LINES {
			prs.pr.CssVariables.Lines[position] = true
		} else {
			prs.pr.CssVariables.MoreLines = true
		}
	} else {
		lines := make(map[int]bool)
		lines[position] = true

		prs.pr.CssVariables = ReportContainer{
			Rules:     rulesDB.CssVariables,
			Lines:     lines,
			MoreLines: false,
		}
	}
}

func (prs *ParserEngine) saveToReportCssImportant(position int) {
	if len(prs.pr.CssImportant.Lines) > 0 {
		if len(prs.pr.CssImportant.Lines) < LIMIT_REPORT_LINES {
			prs.pr.CssImportant.Lines[position] = true
		} else {
			prs.pr.CssImportant.MoreLines = true
		}
	} else {
		lines := make(map[int]bool)
		lines[position] = true

		prs.pr.CssImportant = ReportContainer{
			Rules:     rulesDB.CssImportant,
			Lines:     lines,
			MoreLines: false,
		}
	}
}

func (prs *ParserEngine) saveToReportHtml5Doctype(position int) {
	if len(prs.pr.Html5Doctype.Lines) > 0 {
		if len(prs.pr.Html5Doctype.Lines) < LIMIT_REPORT_LINES {
			prs.pr.Html5Doctype.Lines[position] = true
		} else {
			prs.pr.Html5Doctype.MoreLines = true
		}
	} else {
		lines := make(map[int]bool)
		lines[position] = true

		prs.pr.Html5Doctype = ReportContainer{
			Rules:     rulesDB.Html5Doctype,
			Lines:     lines,
			MoreLines: false,
		}
	}
}

func (prs *ParserEngine) checkHtmlAttribute(attrKey, attrVal string, position int) {
	attrKey = strings.ToLower(strings.Trim(attrKey, WHITESPACE))
	attrVal = strings.ToLower(strings.Trim(attrVal, WHITESPACE))

	if cssKeyData, ok := rulesDB.HtmlAttributes[attrKey]; ok {
		if cssValData, ok := cssKeyData[attrVal]; ok {
			prs.saveToReportHtmlAttributes(attrKey, attrVal, position, cssValData)
		}
	}
}

func (prs *ParserEngine) saveToReportAtRuleCssStatements(propertyKey, propertyVal string, position int, ruleCssPropData interface{}) {
	prs.mx.Lock()
	defer prs.mx.Unlock()

	if prKeyData, ok := prs.pr.AtRuleCssStatements[propertyKey]; ok {
		if prValData, ok := prKeyData[propertyVal]; ok {
			if len(prValData.Lines) < LIMIT_REPORT_LINES {
				prValData.Lines[position] = true
			} else {
				prValData.MoreLines = true
			}
			prKeyData[propertyVal] = prValData
			prs.pr.AtRuleCssStatements[propertyKey] = prKeyData
		} else {
			if len(prs.pr.AtRuleCssStatements[propertyKey]) > 0 {
				prs.pr.AtRuleCssStatements[propertyKey][propertyVal] = makeInitialReportContainer(position, ruleCssPropData)
			} else {
				rootData := make(map[string]ReportContainer)
				rootData[propertyVal] = makeInitialReportContainer(position, ruleCssPropData)
				prs.pr.AtRuleCssStatements[propertyKey] = rootData
			}
		}
	} else {
		rData := make(map[string]ReportContainer)
		rData[propertyVal] = makeInitialReportContainer(position, ruleCssPropData)

		if len(prs.pr.AtRuleCssStatements) > 0 {
			prs.pr.AtRuleCssStatements[propertyKey] = rData
		} else {
			rootData := make(map[string]map[string]ReportContainer)
			rootData[propertyKey] = rData
			prs.pr.AtRuleCssStatements = rootData
		}
	}
}

func (prs *ParserEngine) checkAtRuleCssStatements(propertyKey, propertyVal string, position int) {
	propertyKey = strings.ToLower(strings.Trim(propertyKey, WHITESPACE))
	propertyVal = strings.ToLower(strings.Trim(propertyVal, WHITESPACE))

	if cssKeyData, ok := rulesDB.AtRuleCssStatements[propertyKey]; ok {
		if cssValData, ok := cssKeyData[propertyVal]; ok {
			prs.saveToReportAtRuleCssStatements(propertyKey, propertyVal, position, cssValData)
		}
	}
}

func (prs *ParserEngine) saveToReportImgFormats(psSelectorValue string, position int, ruleCssPropData interface{}) {
	prs.mx.Lock()
	defer prs.mx.Unlock()

	if prKeyData, ok := prs.pr.ImgFormats[psSelectorValue]; ok {
		if len(prKeyData.Lines) < LIMIT_REPORT_LINES {
			prKeyData.Lines[position] = true
		} else {
			prKeyData.MoreLines = true
		}
		prs.pr.ImgFormats[psSelectorValue] = prKeyData
	} else {
		if len(prs.pr.ImgFormats) > 0 {
			prs.pr.ImgFormats[psSelectorValue] = makeInitialReportContainer(position, ruleCssPropData)
		} else {
			rootData := make(map[string]ReportContainer)
			rootData[psSelectorValue] = makeInitialReportContainer(position, ruleCssPropData)
			prs.pr.ImgFormats = rootData
		}
	}
}

func (prs *ParserEngine) checkImgFormat(imgUrl string, position int) {
	if cssUrlRe.MatchString(imgUrl) {
		imgUrl = cssUrlRe.FindStringSubmatch(imgUrl)[1] // parse url from "url(img.path)"
	}

	if strings.HasPrefix(imgUrl, "data:") && strings.Contains(imgUrl, "base64") {
		if imgFormatsData, ok := rulesDB.ImgFormats["base64"]; ok {
			prs.saveToReportImgFormats("base64", position, imgFormatsData)
		}
		return
	}

	urlData, err := url.Parse(imgUrl)
	if err != nil {
		return
	}

	format := strings.Replace(filepath.Ext(urlData.Path), ".", "", 1) // remove dot from extension
	if imgFormatsData, ok := rulesDB.ImgFormats[format]; ok {
		prs.saveToReportImgFormats(format, position, imgFormatsData)
	}
}

func (prs *ParserEngine) checkAttrImgFormat(attrKey, imgUrl string, position int) {
	imgUrl = strings.ToLower(strings.Trim(imgUrl, WHITESPACE))

	if attrKey == "srcset" && (strings.Contains(imgUrl, " ") || strings.Contains(imgUrl, ",")) {
		for _, imgWithSize := range strings.Split(imgUrl, ",") {
			imgAndOther := strings.Split(strings.Trim(imgWithSize, WHITESPACE), " ")
			if len(imgAndOther) > 0 {
				prs.checkImgFormat(strings.Trim(imgAndOther[0], WHITESPACE), position)
			}
		}
		return
	}

	prs.checkImgFormat(imgUrl, position)
}

func (prs *ParserEngine) saveToReportCssPseudoSelectors(psSelectorValue string, position int, ruleCssPropData interface{}) {
	prs.mx.Lock()
	defer prs.mx.Unlock()

	if prKeyData, ok := prs.pr.CssPseudoSelectors[psSelectorValue]; ok {
		if len(prKeyData.Lines) < LIMIT_REPORT_LINES {
			prKeyData.Lines[position] = true
		} else {
			prKeyData.MoreLines = true
		}
		prs.pr.CssPseudoSelectors[psSelectorValue] = prKeyData
	} else {
		if len(prs.pr.CssPseudoSelectors) > 0 {
			prs.pr.CssPseudoSelectors[psSelectorValue] = makeInitialReportContainer(position, ruleCssPropData)
		} else {
			rootData := make(map[string]ReportContainer)
			rootData[psSelectorValue] = makeInitialReportContainer(position, ruleCssPropData)
			prs.pr.CssPseudoSelectors = rootData
		}
	}
}

func (prs *ParserEngine) checkCssPseudoSelector(psSelectorValue string, position int) {
	psSelectorValue = strings.ToLower(strings.Trim(psSelectorValue, WHITESPACE))
	if cssFunctionsData, ok := rulesDB.CssPseudoSelectors[psSelectorValue]; ok {
		prs.saveToReportCssPseudoSelectors(psSelectorValue, position, cssFunctionsData)
	}
}

func (prs *ParserEngine) saveToReportCssFunctions(functionValue string, position int, ruleCssPropData interface{}) {
	prs.mx.Lock()
	defer prs.mx.Unlock()

	if prKeyData, ok := prs.pr.CssFunctions[functionValue]; ok {
		if len(prKeyData.Lines) < LIMIT_REPORT_LINES {
			prKeyData.Lines[position] = true
		} else {
			prKeyData.MoreLines = true
		}
		prs.pr.CssFunctions[functionValue] = prKeyData
	} else {
		if len(prs.pr.CssFunctions) > 0 {
			prs.pr.CssFunctions[functionValue] = makeInitialReportContainer(position, ruleCssPropData)
		} else {
			rootData := make(map[string]ReportContainer)
			rootData[functionValue] = makeInitialReportContainer(position, ruleCssPropData)
			prs.pr.CssFunctions = rootData
		}
	}
}

func (prs *ParserEngine) checkCssFunction(functionValue string, position int) {
	functionValue = strings.ToLower(strings.Trim(strings.ReplaceAll(functionValue, "(", ""), WHITESPACE))
	if cssFunctionsData, ok := rulesDB.CssFunctions[functionValue]; ok {
		prs.saveToReportCssFunctions(functionValue, position, cssFunctionsData)
	}
}

func (prs *ParserEngine) saveToReportCssDimention(dimentionValue string, position int, ruleCssPropData interface{}) {
	prs.mx.Lock()
	defer prs.mx.Unlock()

	if prKeyData, ok := prs.pr.CssDimentions[dimentionValue]; ok {
		if len(prKeyData.Lines) < LIMIT_REPORT_LINES {
			prKeyData.Lines[position] = true
		} else {
			prKeyData.MoreLines = true
		}
		prs.pr.CssDimentions[dimentionValue] = prKeyData
	} else {
		if len(prs.pr.CssDimentions) > 0 {
			prs.pr.CssDimentions[dimentionValue] = makeInitialReportContainer(position, ruleCssPropData)
		} else {
			rootData := make(map[string]ReportContainer)
			rootData[dimentionValue] = makeInitialReportContainer(position, ruleCssPropData)
			prs.pr.CssDimentions = rootData
		}
	}
}

func (prs *ParserEngine) checkCssDimention(dimentionValue string, position int) {
	dimentionValue = strings.ToLower(strings.Trim(dimentionsRe.ReplaceAllString(dimentionValue, ""), WHITESPACE))
	if cssDimentionsData, ok := rulesDB.CssDimentions[dimentionValue]; ok {
		prs.saveToReportCssDimention(dimentionValue, position, cssDimentionsData)
	}
}

func (prs *ParserEngine) saveToReportCssSelectorType(selectorType CssSelectorType, position int, ruleCssPropData interface{}) {
	prs.mx.Lock()
	defer prs.mx.Unlock()

	if prKeyData, ok := prs.pr.CssSelectorTypes[selectorType.String()]; ok {
		if len(prKeyData.Lines) < LIMIT_REPORT_LINES {
			prKeyData.Lines[position] = true
		} else {
			prKeyData.MoreLines = true
		}
		prs.pr.CssSelectorTypes[selectorType.String()] = prKeyData
	} else {
		if len(prs.pr.CssSelectorTypes) > 0 {
			prs.pr.CssSelectorTypes[selectorType.String()] = makeInitialReportContainer(position, ruleCssPropData)
		} else {
			rootData := make(map[string]ReportContainer)
			rootData[selectorType.String()] = makeInitialReportContainer(position, ruleCssPropData)
			prs.pr.CssSelectorTypes = rootData
		}
	}
}

func (prs *ParserEngine) checkCssSelectorType(selectorType CssSelectorType, position int) {
	// log.Printf("[checkCssSelectorType]: %v - %v\n", selectorType, position)
	if cssSelectorTypeData, ok := rulesDB.CssSelectorTypes[selectorType.String()]; ok {
		prs.saveToReportCssSelectorType(selectorType, position, cssSelectorTypeData)
	}
}

func (prs *ParserEngine) checkCssPropertyStyle(propertyKey, propertyVal string, position int) {
	propertyKey = normalizeCssProp(strings.ToLower(strings.Trim(propertyKey, WHITESPACE)))
	propertyVal = strings.Trim(strings.ReplaceAll(propertyVal, "!important", ""), WHITESPACE)

	if cssKeyData, ok := rulesDB.CssProperties[propertyKey]; ok {
		if cssValData, ok := cssKeyData[""]; ok {
			prs.saveToReportCssProperty(propertyKey, "", position, cssValData)
		}
		if cssValData, ok := cssKeyData[propertyVal]; ok {
			prs.saveToReportCssProperty(propertyKey, propertyVal, position, cssValData)
		}
	}
}

func (prs *ParserEngine) saveToReportCssProperty(propertyKey, propertyVal string, position int, ruleCssPropData interface{}) {
	prs.mx.Lock()
	defer prs.mx.Unlock()

	if prKeyData, ok := prs.pr.CssProperties[propertyKey]; ok {
		if prValData, ok := prKeyData[propertyVal]; ok {
			if len(prValData.Lines) < LIMIT_REPORT_LINES {
				prValData.Lines[position] = true
			} else {
				prValData.MoreLines = true
			}
			prKeyData[propertyVal] = prValData
			prs.pr.CssProperties[propertyKey] = prKeyData
		} else {
			if len(prs.pr.CssProperties[propertyKey]) > 0 {
				prs.pr.CssProperties[propertyKey][propertyVal] = makeInitialReportContainer(position, ruleCssPropData)
			} else {
				rootData := make(map[string]ReportContainer)
				rootData[propertyVal] = makeInitialReportContainer(position, ruleCssPropData)
				prs.pr.CssProperties[propertyKey] = rootData
			}
		}
	} else {
		rData := make(map[string]ReportContainer)
		rData[propertyVal] = makeInitialReportContainer(position, ruleCssPropData)

		if len(prs.pr.CssProperties) > 0 {
			prs.pr.CssProperties[propertyKey] = rData
		} else {
			rootData := make(map[string]map[string]ReportContainer)
			rootData[propertyKey] = rData
			prs.pr.CssProperties = rootData
		}
	}
}

func (prs *ParserEngine) checkCssParsedToken(p *css.Parser, gt css.GrammarType, data []byte, position int) {
	switch gt {
	case css.CustomPropertyGrammar:
		prs.saveToReportCssVariables(position)
	case css.AtRuleGrammar:
		prs.checkAtRuleCssStatements(string(data), "", position)
		for _, val := range p.Values() {
			prs.checkAtRuleCssStatements(string(data), string(val.Data), position)
		}
	case css.BeginAtRuleGrammar:
		prs.checkAtRuleCssStatements(string(data), "", position)
		for _, val := range p.Values() {
			prs.checkAtRuleCssStatements(string(data), string(val.Data), position)

			if val.TokenType == css.DimensionToken || val.TokenType == css.PercentageToken {
				prs.checkCssDimention(string(val.Data), position)
			}
		}
	case css.QualifiedRuleGrammar, css.BeginRulesetGrammar:
		if gt == css.QualifiedRuleGrammar {
			prs.checkCssSelectorType(GROUPING_SELECTORS_TYPE, position)
		}

		prevTokenType := css.Token{
			TokenType: css.ErrorToken,
			Data:      []byte{},
		}

		chainingSelectorsCount := 0
		typeSelector := false
		for _, val := range p.Values() {
			if val.TokenType == css.LeftBracketToken {
				prs.checkCssSelectorType(ATTRIBUTE_SELECTOR_TYPE, position)
			}
			if val.TokenType == css.WhitespaceToken {
				if chainingSelectorsCount > 1 {
					prs.checkCssSelectorType(CHAINING_SELECTORS_TYPE, position)
				}
				if typeSelector {
					prs.checkCssSelectorType(TYPE_SELECTOR_TYPE, position)
				}
				chainingSelectorsCount = 0
				typeSelector = false
				prs.checkCssSelectorType(DESCENDANT_COMBINATOR_TYPE, position)
			}
			if val.TokenType == css.HashToken {
				prs.checkCssSelectorType(ID_SELECTOR_TYPE, position)
			}

			if val.TokenType == css.DelimToken {
				typeSelector = false
				delimVal := strings.ToLower(strings.Trim(string(val.Data), WHITESPACE))
				if delimVal == "." {
					chainingSelectorsCount += 1
				} else {
					if chainingSelectorsCount > 1 {
						prs.checkCssSelectorType(CHAINING_SELECTORS_TYPE, position)
					}
					chainingSelectorsCount = 0
				}
				if delimVal == "*" {
					prs.checkCssSelectorType(UNIVERSAL_SELECTOR_STAR_TYPE, position)
				}
				if delimVal == "~" {
					prs.checkCssSelectorType(GENERAL_SIBLING_COMBINATOR_TYPE, position)
				}
				if delimVal == "+" {
					prs.checkCssSelectorType(ADJACENT_SIBLING_COMBINATOR_TYPE, position)
				}
				if delimVal == ">" {
					prs.checkCssSelectorType(CHILD_COMBINATOR_TYPE, position)
				}
			}

			if prevTokenType.TokenType == css.ColonToken && (val.TokenType == css.IdentToken || val.TokenType == css.FunctionToken) {
				pseudoName := string(val.Data)
				if val.TokenType == css.FunctionToken { // convert "has(" to "has"
					pseudoName = strings.ReplaceAll(pseudoName, "(", "")
				}
				prs.checkCssPseudoSelector(pseudoName, position)
			}
			if prevTokenType.TokenType == css.DelimToken && string(prevTokenType.Data) == "." && val.TokenType == css.IdentToken {
				prs.checkCssSelectorType(CLASS_SELECTOR_TYPE, position)
			}
			if prevTokenType.TokenType != css.ColonToken && prevTokenType.TokenType != css.DelimToken && val.TokenType == css.IdentToken {
				typeSelectorVal := strings.ToLower(strings.Trim(string(val.Data), WHITESPACE))
				if typeSelectorVal != "from" && typeSelectorVal != "to" { // ignore @keyframe values as tags
					typeSelector = true
				}
			}
			prevTokenType = val
		}

		if chainingSelectorsCount > 1 {
			prs.checkCssSelectorType(CHAINING_SELECTORS_TYPE, position)
		}
		if typeSelector {
			prs.checkCssSelectorType(TYPE_SELECTOR_TYPE, position)
		}
	case css.DeclarationGrammar:
		propVal := ""
		isPrevDelimCanBeImportant := false
		for _, val := range p.Values() {
			cssPropVal := strings.ToLower(strings.Trim(string(val.Data), WHITESPACE))

			if val.TokenType == css.DimensionToken || val.TokenType == css.PercentageToken {
				prs.checkCssDimention(string(val.Data), position)
			}
			if val.TokenType == css.IdentToken && cssPropVal == "initial" { // dimention unit "initial"
				prs.checkCssDimention(cssPropVal, position)
			}
			if val.TokenType == css.FunctionToken {
				prs.checkCssFunction(string(val.Data), position)
			}
			if val.TokenType == css.URLToken {
				prs.checkImgFormat(string(val.Data), position)
			}
			if val.TokenType == css.CustomPropertyNameToken {
				prs.saveToReportCssVariables(position)
			}
			if isPrevDelimCanBeImportant && val.TokenType == css.IdentToken && cssPropVal == "important" {
				prs.saveToReportCssImportant(position)
			}
			isPrevDelimCanBeImportant = (val.TokenType == css.DelimToken && cssPropVal == "!")
			propVal += cssPropVal
		}
		prs.checkCssPropertyStyle(string(data), propVal, position)
	}
}

func (prs *ParserEngine) checkTagInlinedStyle(inlineStyle string, position int) {
	p := css.NewParser(parse.NewInput(bytes.NewBufferString(inlineStyle)), true)
	for {
		gt, _, data := p.Next()

		// log.Printf("[checkTagInlinedStyle]: %v - %v - %v\n", gt, tt, string(data))

		if gt == css.ErrorGrammar {
			return
		}

		prs.checkCssParsedToken(p, gt, data, position)
	}
}

func (prs *ParserEngine) processCssInStyleTag(inlineStyle string, htmlTagPosition int) {
	var (
		bytesToLine []int
		cursorPos   int = 0
		cssLine     int = 0
	)

	inlineStyleBytes := []byte(inlineStyle)
	lines := bytes.Split(inlineStyleBytes, []byte("\n"))
	for _, line := range lines {
		bytesToLine = append(bytesToLine, cursorPos)
		cursorPos += len(line) + 1 // "\n" provide additional byte
	}

	correctLine := func(gt css.GrammarType, offset, cssLine int) int {
		if css.DeclarationGrammar == gt && len(inlineStyleBytes) > offset && cssLine > 1 {
			newLine := cssLine
			startOffset := bytesToLine[newLine-1]
			endOffset := offset

			for {
				if len(bytes.Trim(inlineStyleBytes[startOffset:endOffset], " \t\r\n\f}")) > 0 {
					return newLine
				}

				newLine = newLine - 1
				if newLine < 1 {
					return cssLine
				}

				endOffset = startOffset
				startOffset = bytesToLine[newLine-1]
			}
		}
		return cssLine
	}

	getLineByOffset := func(gt css.GrammarType, offset int) int {
		for {
			if len(bytesToLine) <= cssLine {
				return correctLine(gt, offset, cssLine)
			}

			if bytesToLine[cssLine] > offset {
				return correctLine(gt, offset, cssLine)
			}

			cssLine += 1
		}
	}

	p := css.NewParser(parse.NewInput(bytes.NewBufferString(inlineStyle)), false)
	for {
		gt, _, data := p.Next()

		// log.Printf("[checkTagInlinedStyle]: %v - %v - %v - %v\n", gt, string(data), p.Values(), p.Offset())

		if gt == css.ErrorGrammar {
			return
		}

		position := htmlTagPosition + getLineByOffset(gt, p.Offset()) - 1
		prs.checkCssParsedToken(p, gt, data, position)
	}
}

func (prs *ParserEngine) saveToReportHtmlTag(tagName, tagAttr string, position int, ruleTagAttrData interface{}) {
	prs.mx.Lock()
	defer prs.mx.Unlock()

	if tagData, ok := prs.pr.HtmlTags[tagName]; ok {
		if tagAttrData, ok := tagData[tagAttr]; ok {
			if len(tagAttrData.Lines) < LIMIT_REPORT_LINES {
				tagAttrData.Lines[position] = true
			} else {
				tagAttrData.MoreLines = true
			}
			tagData[tagAttr] = tagAttrData
			prs.pr.HtmlTags[tagName] = tagData
		} else {
			if len(prs.pr.HtmlTags[tagName]) > 0 {
				prs.pr.HtmlTags[tagName][tagAttr] = makeInitialReportContainer(position, ruleTagAttrData)
			} else {
				rootData := make(map[string]ReportContainer)
				rootData[tagAttr] = makeInitialReportContainer(position, ruleTagAttrData)
				prs.pr.HtmlTags[tagName] = rootData
			}
		}
	} else {
		rData := make(map[string]ReportContainer)
		rData[tagAttr] = makeInitialReportContainer(position, ruleTagAttrData)

		if len(prs.pr.HtmlTags) > 0 {
			prs.pr.HtmlTags[tagName] = rData
		} else {
			rootData := make(map[string]map[string]ReportContainer)
			rootData[tagName] = rData
			prs.pr.HtmlTags = rootData
		}
	}
}

func (prs *ParserEngine) checkHtmlTagWithAttr(attrKey, attrVal string, position int) {
	prs.checkHtmlAttribute(attrKey, "", position)
	prs.checkHtmlAttribute(attrKey, attrVal, position)

	if attrKey == "style" {
		prs.checkTagInlinedStyle(attrVal, position)
	}

	if attrKey == "src" || attrKey == "srcset" {
		prs.checkAttrImgFormat(attrKey, attrVal, position)
	}
}

func (prs *ParserEngine) checkHtmlTags(tagName string, attrs []html.Attribute, position int) {
	if len(tagName) == 0 {
		return
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

			prs.checkHtmlTagWithAttr(attrKey, attrVal, position)
		}
	} else {
		// check inline style for valid elements too
		for _, att := range attrs {
			attrKey := strings.ToLower(att.Key)
			attrVal := strings.ToLower(att.Val)

			prs.checkHtmlTagWithAttr(attrKey, attrVal, position)
		}
	}
}

func (prs *ParserEngine) saveToReportLinkTypes(linkType string, position int, ruleCssPropData interface{}) {
	prs.mx.Lock()
	defer prs.mx.Unlock()

	if prKeyData, ok := prs.pr.LinkTypes[linkType]; ok {
		if len(prKeyData.Lines) < LIMIT_REPORT_LINES {
			prKeyData.Lines[position] = true
		} else {
			prKeyData.MoreLines = true
		}
		prs.pr.LinkTypes[linkType] = prKeyData
	} else {
		if len(prs.pr.LinkTypes) > 0 {
			prs.pr.LinkTypes[linkType] = makeInitialReportContainer(position, ruleCssPropData)
		} else {
			rootData := make(map[string]ReportContainer)
			rootData[linkType] = makeInitialReportContainer(position, ruleCssPropData)
			prs.pr.LinkTypes = rootData
		}
	}
}

func (prs *ParserEngine) checkLinkTypes(attrs []html.Attribute, position int) {
	for _, att := range attrs {
		attrKey := strings.ToLower(att.Key)
		attrVal := strings.Trim(att.Val, WHITESPACE)

		if attrKey == "href" && len(attrVal) > 0 {
			if anchorLinkRe.MatchString(attrVal) {
				if ruleLinkData, ok := rulesDB.LinkTypes["anchor"]; ok {
					prs.saveToReportLinkTypes("anchor", position, ruleLinkData)
				}
			}
			if mailtoLinkRe.MatchString(attrVal) {
				if ruleLinkData, ok := rulesDB.LinkTypes["mailto"]; ok {
					prs.saveToReportLinkTypes("mailto", position, ruleLinkData)
				}
			}
		}
	}
}

func (prs *ParserEngine) processHtmlToken(htmlTokenizer *html.Tokenizer, token html.Token, tagLine int) {
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
		case a.A:
			// check link
			prs.checkLinkTypes(token.Attr, tagLine)
		}
		// process html tag
		prs.checkHtmlTags(token.Data, token.Attr, tagLine)
	case html.EndTagToken:
		switch token.DataAtom {
		case a.Style:
			if prs.isStyleTagOpen && len(prs.styleTagContent) > 0 {
				prs.wg.Add(1)
				go func(content string, line int) {
					defer prs.wg.Done()
					prs.processCssInStyleTag(content, line)
				}(prs.styleTagContent, prs.styleTagLine)
				// reset style tag storage
				prs.isStyleTagOpen = false
				prs.styleTagContent = ""
				prs.styleTagLine = 0
			}
		}
	case html.SelfClosingTagToken:
		// process html tag
		prs.checkHtmlTags(token.Data, token.Attr, tagLine)
	case html.DoctypeToken:
		// check doctype
		if html5DoctypeRe.MatchString(token.String()) {
			prs.saveToReportHtml5Doctype(tagLine)
		}
	}
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

		// log.Printf("[htmlTokenizer]: info: %v ; data: %v ; type: %v ; atom: %v; attr - %v \n", tt, tt.Data, tt.Type, tt.DataAtom, tt.Attr)

		tagLine = prs.getLineFromOffset(tagLine, htmlBytesOffset)
		prs.processHtmlToken(htmlTokenizer, tt, tagLine)

		htmlBytesOffset += len(htmlTokenizer.Raw())
	}

	prs.wg.Wait() // wait all jobs

	return &prs.pr, nil
}

func init() {
	// parse caniuse.json here
	if err := json.Unmarshal(caniuseJSON, &rulesDB); err != nil {
		panic(err)
	}
}

func ReportFromHTML(document []byte) (*ParseReport, error) {
	parser := InitParser()
	report, err := parser.Report(document)
	if err != nil {
		return nil, err
	}

	return report, nil
}
