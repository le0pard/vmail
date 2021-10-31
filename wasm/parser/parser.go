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

	css "github.com/le0pard/vmail/wasm/parser/css" // patched for offset
	parse "github.com/tdewolff/parse/v2"

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

func (prs *ParserEngine) processCssInStyleTag(inlineStyle string, htmlTagPosition int) {
	var (
		bytesToLine []int
		cursorPos   int = 0
		cssLine     int = 0
	)

	lines := bytes.Split([]byte(inlineStyle), []byte("\n"))
	for _, line := range lines {
		bytesToLine = append(bytesToLine, cursorPos)
		cursorPos += len(line) + 1 // "\n" provide additional byte
	}

	getLineByOffset := func(offset int) int {
		for {
			if len(bytesToLine) <= cssLine {
				return cssLine
			}

			if bytesToLine[cssLine] > offset {
				return cssLine
			}

			cssLine += 1
		}
	}

	p := css.NewParser(parse.NewInput(bytes.NewBufferString(inlineStyle)), false)
	for {
		gt, _, data := p.Next()

		log.Printf("[checkTagInlinedStyle]: %v - %v - %v - %v\n", gt, string(data), p.Values(), htmlTagPosition+getLineByOffset(p.Offset())-1)

		if gt == css.ErrorGrammar {
			return
		}

		switch gt {
		case css.BeginRulesetGrammar:
			propVal := ""
			for _, val := range p.Values() {
				propVal += string(val.Data)
			}
			log.Printf("[CSS SELECTOR]: %v - %v - %v\n", gt, propVal, htmlTagPosition+getLineByOffset(p.Offset())-1)
		case css.DeclarationGrammar:
			propVal := ""
			for _, val := range p.Values() {
				propVal += string(val.Data)
			}
			prs.checkCssPropertyStyle(string(data), propVal, htmlTagPosition+getLineByOffset(p.Offset())-1)
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
	p := css.NewParser(parse.NewInput(bytes.NewBufferString(inlineStyle)), true)
	for {
		gt, _, data := p.Next()

		// log.Printf("[checkTagInlinedStyle]: %v - %v - %v\n", gt, tt, string(data))

		if gt == css.ErrorGrammar {
			return
		}

		switch gt {
		case css.DeclarationGrammar:
			propVal := ""
			for _, val := range p.Values() {
				propVal += string(val.Data)
			}
			prs.checkCssPropertyStyle(string(data), propVal, position)
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
