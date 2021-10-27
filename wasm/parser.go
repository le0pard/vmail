package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"

	css "github.com/gorilla/css/scanner"

	"golang.org/x/net/html"
	a "golang.org/x/net/html/atom"

	_ "embed"
)

//go:embed rules.json
var emailRules []byte

const (
	WHITESPACE           = " \t\r\n\f"
	HTML_NO_ATTRIBUTE    = ""
	CSS_SCOPE_STYLE_ATTR = "style_attribute"
	CSS_SCOPE_STYLE_TAG  = "style_tag"
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

type ParserEngine struct {
	bytesToLine []int

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

func calulateNewlineBytePos(document []byte) []int {
	var (
		bytesToLine []int
		cursorPos   int = 0
	)

	lines := bytes.Split(document, []byte("\n"))

	for _, line := range lines {
		bytesToLine = append(bytesToLine, cursorPos)
		cursorPos += len(line) + 1 // "\n" provide additional byte
	}

	return bytesToLine
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
			log.Printf("%v %v %v\n", propName, propVal, propLine)
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
				log.Printf("%v %v %v\n", propName, propVal, propLine)
				// prs.storeTagStyleProperties(propName, propVal, propLine)
				// reset for new prop
				isCollectPropVal = false
				propName = ""
				propVal = ""
				propLine = 0
			case "}":
				log.Printf("%v %v %v\n", propName, propVal, propLine)
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
				log.Printf("%v %v %v\n", token.Value, "", line)
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
						log.Printf("%v %v\n", pseudoClassName, pseudoClassLine)
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
					log.Printf("%v %v\n", pseudoClassName, pseudoClassLine)
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
					log.Printf("%v %v\n", pseudoClassName, pseudoClassLine)
					// err := prs.storeCSSPseudoClass(pseudoClassName, pseudoClassLine)
					// if err != nil {
					// 	return err
					// }
					pseudoClassName = ""
					pseudoClassLine = 0
				}
				isPseudoClass = false

				if len(atKeywordName) == 0 && len(cssSelector) > 0 {
					log.Printf("%v %v\n", cssSelector, cssSelectorLine)
					// err := prs.processCssSelector(cssSelector, cssSelectorLine)
					// if err != nil {
					// 	return err
					// }
					cssSelector = ""
					cssSelectorLine = 0
				}

				if len(atKeywordName) > 0 {
					log.Printf("%v %v %v\n", atKeywordName, atKeywordVal, atKeywordLine)
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
		if len(token.Attr) > 0 {
			log.Printf("%v %v %v\n", token.Data, token.Attr, tagLine)
			// err := prs.storeHtmlAttribute(token.Data, token.Attr, tagLine)
			// if err != nil {
			// 	return err
			// }
		} else {
			// store tag with no attributes
			log.Print(
				"%v %v %v\n",
				token.Data,
				[]html.Attribute{
					{
						Namespace: "",
						Key:       HTML_NO_ATTRIBUTE,
						Val:       HTML_NO_ATTRIBUTE,
					},
				},
				tagLine,
			)
			// err := prs.storeHtmlAttribute(
			// 	token.Data,
			// 	[]html.Attribute{
			// 		{
			// 			Namespace: "",
			// 			Key:       HTML_NO_ATTRIBUTE,
			// 			Val:       HTML_NO_ATTRIBUTE,
			// 		},
			// 	},
			// 	tagLine,
			// )
			// if err != nil {
			// 	return err
			// }
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
		log.Print("%v %v %v\n", token.Data, token.Attr, tagLine)
		// err := prs.storeHtmlAttribute(token.Data, token.Attr, tagLine)
		// if err != nil {
		// 	return err
		// }
	}

	return nil
}

func (prs *ParserEngine) ReportFromHTML(document []byte) error {
	var (
		err             error
		htmlTokenizer   *html.Tokenizer
		htmlBytesOffset int = 0
		tagLine         int = 0
	)

	prs.bytesToLine = calulateNewlineBytePos(document)

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
				return err
			}
		}

		tagLine = prs.getLineFromOffset(tagLine, htmlBytesOffset)
		err := prs.processHtmlToken(htmlTokenizer, tt, tagLine)
		if err != nil {
			return err
		}

		htmlBytesOffset += len(htmlTokenizer.Raw())
	}

	return nil
}
