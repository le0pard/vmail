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

type CaniuseDB struct {
	HtmlTags      map[string]map[string]interface{} `json:"html_tags"`
	CssProperties map[string]map[string]interface{} `json:"css_properties"`
}

var rulesDB CaniuseDB

// json config structs end

// result structure begin

type HTMLTagReport struct {
	Rules     interface{}  `json:"rules"`
	Lines     map[int]bool `json:"lines"`
	MoreLines bool         `json:"more_lines"`
}

type CssPropertyReport struct {
	Rules     interface{}  `json:"rules"`
	Lines     map[int]bool `json:"lines"`
	MoreLines bool         `json:"more_lines"`
}

type ParseReport struct {
	HtmlTags      map[string]map[string]HTMLTagReport     `json:"html_tags"`
	CssProperties map[string]map[string]CssPropertyReport `json:"css_properties"`
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

func (prs *ParserEngine) parseCssProps(scanner *css.Scanner, htmlTagPosition int) {
	var (
		isCollectPropVal bool     = false
		propName         []string = []string{}
		propVal          []string = []string{}
		propLine         int      = 0
	)

	processCssProps := func() {
		if len(propName) > 0 {
			log.Printf("[CSS PROPS]: %v %v %v\n", propName, propVal, propLine)
		}
	}

	collectValues := func(token *css.Token) {
		if isCollectPropVal {
			propVal = append(propVal, token.Value)
		} else {
			if len(propName) == 0 {
				propLine = htmlTagPosition + token.Line - 1 // initial line for prop
			}
			propName = append(propName, token.Value)
		}
	}

	for {
		token := scanner.Next()
		if token.Type == css.TokenEOF || token.Type == css.TokenError {
			processCssProps()
			return
		}

		// log.Printf("[parseCssProps] %s: %v (line: %d, column: %d) (doc line: %d)", token.Type, token.Value, token.Line, token.Column, htmlTagPosition+token.Line-1)

		switch token.Type {
		case css.TokenCDO, css.TokenCDC, css.TokenComment, css.TokenS:
			break // continue loop
		case css.TokenChar:
			switch token.Value {
			case ":":
				isCollectPropVal = true
			case ";":
				processCssProps()
				// reset for new prop
				isCollectPropVal = false
				propName = []string{}
				propVal = []string{}
				propLine = 0
			case "}":
				processCssProps()
				return
			default:
				collectValues(token)
			}
		default:
			collectValues(token)
		}
	}
}

func (prs *ParserEngine) processCssInStyleTag(scanner *css.Scanner, htmlTagPosition int, isNested bool) {
	var (
		isNestedAtRule     bool     = false
		isAtPropsRule      bool     = false
		isAtRule           bool     = false
		atRuleName         string   = ""
		atRuleValue        []string = []string{}
		atRuleLine         int      = 0
		cssSelector        []string = []string{}
		cssSelectorLine    int      = 0
		curlyBracesCounter int      = 0
	)

	collectAtRuleValues := func(token *css.Token) {
		if isNestedAtRule || isAtRule || isAtPropsRule {
			val := strings.Trim(token.Value, WHITESPACE)
			if len(val) > 0 {
				atRuleValue = append(atRuleValue, val)
			}
		}
	}

	collectSelectorValues := func(token *css.Token) {
		if !(isNestedAtRule || isAtRule || isAtPropsRule) {
			if len(cssSelector) == 0 {
				cssSelectorLine = htmlTagPosition + token.Line - 1
			}
			val := strings.Trim(token.Value, WHITESPACE)
			if len(val) > 0 {
				cssSelector = append(cssSelector, val)
			}
		}
	}

	resetAtRules := func() {
		isNestedAtRule = false
		isAtPropsRule = false
		isAtRule = false
		atRuleName = ""
		atRuleValue = []string{}
		atRuleLine = 0
	}

	resetSelectorValues := func() {
		cssSelector = []string{}
		cssSelectorLine = 0
	}

	for {
		token := scanner.Next()

		log.Printf("[processCssInStyleTag] %s: %v (line: %d, column: %d) (doc line: %d) (curlyBracesCounter: %d, %v)", token.Type, token.Value, token.Line, token.Column, htmlTagPosition+token.Line-1, curlyBracesCounter, isNested)

		if token.Type == css.TokenEOF || token.Type == css.TokenError {
			return
		}

		switch token.Type {
		case css.TokenCDO, css.TokenCDC, css.TokenComment, css.TokenS:
			if token.Type == css.TokenS {
				collectAtRuleValues(token)
				collectSelectorValues(token)
			}
			break
		case css.TokenAtKeyword:
			switch token.Value {
			case "@charset", "@import", "@namespace ": // regular at-rules
				isAtRule = true
			case "@media", "@supports", "@document", "@keyframes", "@-webkit-keyframes", "@font-feature-values": // nested at-rules
				isNestedAtRule = true
			default:
				isAtPropsRule = true
			}
			atRuleName = token.Value
			atRuleValue = []string{}
			atRuleLine = htmlTagPosition + token.Line - 1
		case css.TokenChar:
			switch token.Value {
			case "}":
				curlyBracesCounter -= 1
				if isNested && curlyBracesCounter < 0 {
					return // finished with nested block parsing
				}
				resetAtRules()
				resetSelectorValues()
			case "{":
				curlyBracesCounter += 1
				if isNestedAtRule {
					log.Printf("[NESTED AT RULE]: %v %v %v\n", atRuleName, atRuleValue, atRuleLine)
					prs.processCssInStyleTag(scanner, htmlTagPosition, true)
					resetAtRules()
				} else if isAtPropsRule {
					log.Printf("[PROPS AT RULE]: %v %v %v\n", atRuleName, atRuleValue, atRuleLine)
					// parse css props
					prs.parseCssProps(scanner, htmlTagPosition)
					resetAtRules()
				} else {
					log.Printf("[CSS SELECTOR]: %v %v\n", cssSelector, cssSelectorLine)
					// parse css props
					prs.parseCssProps(scanner, htmlTagPosition)
				}
			case ";":
				if isAtRule {
					log.Printf("[AT RULE]: %v %v %v\n", atRuleName, atRuleValue, atRuleLine)
					resetAtRules()
				}
			default:
				collectSelectorValues(token)
			}
		default:
			collectAtRuleValues(token)
			collectSelectorValues(token)
		}
	}
}

func makeInitialCssPropertyReport(position int, ruleCssPropData interface{}) CssPropertyReport {
	lines := make(map[int]bool)
	lines[position] = true

	return CssPropertyReport{
		Rules:     ruleCssPropData,
		Lines:     lines,
		MoreLines: false,
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
				prs.pr.CssProperties[propertyKey][propertyVal] = makeInitialCssPropertyReport(position, ruleCssPropData)
			} else {
				rootData := make(map[string]CssPropertyReport)
				rootData[propertyVal] = makeInitialCssPropertyReport(position, ruleCssPropData)
				prs.pr.CssProperties[propertyKey] = rootData
			}
		}
	} else {
		rData := make(map[string]CssPropertyReport)
		rData[propertyVal] = makeInitialCssPropertyReport(position, ruleCssPropData)

		if len(prs.pr.CssProperties) > 0 {
			prs.pr.CssProperties[propertyKey] = rData
		} else {
			rootData := make(map[string]map[string]CssPropertyReport)
			rootData[propertyKey] = rData
			prs.pr.CssProperties = rootData
		}
	}
}

func (prs *ParserEngine) checkCssPropertyStyle(propertyKey, propertyVal string, position int) {
	propertyKey = strings.ToLower(strings.Trim(propertyKey, WHITESPACE))
	propertyVal = strings.ToLower(strings.Trim(propertyVal, WHITESPACE))

	if cssKeyData, ok := rulesDB.CssProperties[propertyKey]; ok {
		if cssValData, ok := cssKeyData[""]; ok {
			prs.saveToReportCssProperty(propertyKey, "", position, cssValData)
		}
		if cssValData, ok := cssKeyData[propertyVal]; ok {
			prs.saveToReportCssProperty(propertyKey, propertyVal, position, cssValData)
		}
	}
}

func (prs *ParserEngine) checkTagInlinedStyle(inlineStyle string, position int) {
	cssProperties := strings.Split(inlineStyle, ";")

	for _, cssProperty := range cssProperties {
		if len(cssProperty) == 0 {
			continue
		}

		propertyKeyVal := strings.Split(strings.Trim(cssProperty, WHITESPACE), ":")
		if len(propertyKeyVal) == 2 {
			prs.checkCssPropertyStyle(propertyKeyVal[0], propertyKeyVal[1], position)
		}
	}
}

func makeInitialHtmlTagReport(position int, ruleTagAttrData interface{}) HTMLTagReport {
	lines := make(map[int]bool)
	lines[position] = true

	return HTMLTagReport{
		Rules:     ruleTagAttrData,
		Lines:     lines,
		MoreLines: false,
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
				prs.pr.HtmlTags[tagName][tagAttr] = makeInitialHtmlTagReport(position, ruleTagAttrData)
			} else {
				rootData := make(map[string]HTMLTagReport)
				rootData[tagAttr] = makeInitialHtmlTagReport(position, ruleTagAttrData)
				prs.pr.HtmlTags[tagName] = rootData
			}
		}
	} else {
		rData := make(map[string]HTMLTagReport)
		rData[tagAttr] = makeInitialHtmlTagReport(position, ruleTagAttrData)

		if len(prs.pr.HtmlTags) > 0 {
			prs.pr.HtmlTags[tagName] = rData
		} else {
			rootData := make(map[string]map[string]HTMLTagReport)
			rootData[tagName] = rData
			prs.pr.HtmlTags = rootData
		}
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

			if attrKey == "style" {
				prs.checkTagInlinedStyle(attrVal, position)
			}
		}
	} else {
		// check inline style for valid elements too
		for _, att := range attrs {
			attrKey := strings.ToLower(att.Key)
			attrVal := strings.ToLower(att.Val)

			if attrKey == "style" {
				prs.checkTagInlinedStyle(attrVal, position)
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
					scanner := css.New(content)
					prs.processCssInStyleTag(scanner, line, false)
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
