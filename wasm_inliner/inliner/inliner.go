package inliner

import (
	"bytes"
	"log"
	"regexp"
	"strings"
	"sync"

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

type ExternalStylesheetsTags struct {
	Parent *html.Node
	Node   *html.Node
	url    string
}

type InlineEngine struct {
	// array for bytes -> lines
	bytesToLine []int
	// group for parallel processing
	wg sync.WaitGroup
	// lock for report
	mx sync.RWMutex
}

func InitInliner() *InlineEngine {
	return &InlineEngine{}
}

func extractStylesheets(doc *html.Node) ([]string, error) {
	var (
		externalStylesheets []ExternalStylesheetsTags
		urls                []string
		crawler             func(*html.Node)
	)

	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "link" {
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

			if isStylesheet && isMediaGood {
				externalStylesheets = append(externalStylesheets, ExternalStylesheetsTags{
					Parent: node.Parent,
					Node:   node,
					url:    foundedHref,
				})
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)

	for _, item := range externalStylesheets {
		urls = append(urls, item.url)
		item.Parent.RemoveChild(item.Node)
	}

	return urls, nil
}

func (inlr *InlineEngine) InlineCss(htmlDoc []byte) ([]byte, error) {
	var (
		doc *html.Node
		err error
	)

	if doc, err = html.Parse(bytes.NewReader(htmlDoc)); err != nil {
		return []byte{}, err
	}

	stylesheetUrls, err := extractStylesheets(doc)
	if err != nil {
		return []byte{}, err
	}

	log.Printf("[stylesheetUrls]: %v\n", stylesheetUrls)

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
