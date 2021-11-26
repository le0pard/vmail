package inliner

import (
	"bytes"
	"errors"
	"io"
	"log"
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
)

const (
	WHITESPACE = " \t\r\n\f"
)

var (
	mediaSplitInlineRe = regexp.MustCompile(`(?i)[\s]+|,`)
	mediaInlineRe      = regexp.MustCompile(`(?i)(screen|handheld|all)`)
)

type StylesheetsTags struct {
	Parent  *html.Node
	Node    *html.Node
	Content string
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

func extractStylesheets(doc *html.Node) (string, error) {
	var (
		externalStylesheets []StylesheetsTags
		contents            []string
		crawler             func(*html.Node)
		netClient           = &http.Client{
			Timeout: 5 * time.Second,
		}
	)

	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode {
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

	return strings.Join(contents, "\n"), nil
}

func (inlr *InlineEngine) inlineStyleSheetContent(doc *html.Node, sheetContent string) error {
	// log.Printf("[sheetContent]: %v\n", sheetContent)

	p := css.NewParser(parse.NewInput(bytes.NewBufferString(sheetContent)), false)
	for {
		gt, _, data := p.Next()

		log.Printf("[inlineStyleSheetContent]: %v - %v - %v\n", gt, string(data), p.Values())

		if gt == css.ErrorGrammar {
			if p.Err() == io.EOF {
				return nil
			}
			return errors.New("Error to parse CSS")
		}

		// prs.checkCssParsedToken(p, gt, data, position)
	}
	return nil
}

func (inlr *InlineEngine) InlineCss(htmlDoc []byte) ([]byte, error) {
	var (
		doc *html.Node
		err error
	)

	if doc, err = html.Parse(bytes.NewReader(htmlDoc)); err != nil {
		return []byte{}, err
	}

	stylesheetContents, err := extractStylesheets(doc)
	if err != nil {
		return []byte{}, err
	}

	err = inlr.inlineStyleSheetContent(doc, stylesheetContents)
	if err != nil {
		return []byte{}, err
	}

	_, err = cascadia.ParseGroup(".body .title")
	if err != nil {
		log.Printf("invalid selector %s", ".body .title")
		return []byte{}, err
	}

	inlr.wg.Wait() // wait all jobs

	return []byte{}, nil
}

func InlineCssInHTML(htmlDoc []byte) ([]byte, error) {
	inliner := InitInliner()
	newHtmlDoc, err := inliner.InlineCss(htmlDoc)
	if err != nil {
		return nil, err
	}

	return newHtmlDoc, nil
}
