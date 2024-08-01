package models

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"dalton.dog/aocgo/internal/tui"
	"github.com/PuerkitoBio/goquery" // Bless this package
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type PageData struct {
	// Both of these two can be found in the page's <head>'s <title> tag
	day  int
	year int

	answerOne string
	answerTwo string

	PuzzleTitle string

	// The article consists of articleContents (as you might expect)
	articleOne      string
	articleOneSel   *goquery.Selection
	articleTwo      string
	articleTwoSel   *goquery.Selection
	articleContents *goquery.Selection

	mainContents *goquery.Selection
}

func NewPageData(raw []byte) *PageData {
	reader := bytes.NewReader(raw)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal("Error constructing new PageData.", "error", err)
	}

	// HACK: Admittedly making some assumptions on input, but should be fine. AoC is very consistent
	title := strings.Split(doc.Find("title").Text(), " ")
	day, _ := strconv.Atoi(title[1])
	year, _ := strconv.Atoi(title[len(title)-1])

	mainContents := doc.Find("main")
	puzzleTitle := mainContents.Find("h2").First().Text()

	pageData := &PageData{
		PuzzleTitle:  titleStyle.Render(puzzleTitle),
		day:          day,
		year:         year,
		mainContents: mainContents,
	}

	pageData.processPageData()

	return pageData
}

// Stylings
// const ParagraphWidth = 120

// TODO: answerStyle
var (
	titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#FFFF00"))
	italStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF3374"))
	starStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#F1FA8C"))
	linkStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#8BE9FD")).Underline(true)
	codeStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAC3D5")).Bold(true)
	wordWrap   = lipgloss.NewStyle().Width(tui.ViewportWidth)
)

func (p *PageData) GetPageDataPrettyString() []string {
	p.processPageData()

	sOut := printArticle(p.articleOneSel)

	if p.answerOne != "" {
		sOut = append(sOut, p.answerOne)
	}

	if p.articleTwoSel != nil {
		sOut = append(sOut, "\n"+titleStyle.Render("--- Part Two ---"))
		sOut = append(sOut, printArticle(p.articleTwoSel)...)

		if p.answerTwo != "" {
			sOut = append(sOut, p.answerTwo)
		}
	}

	// wrappedText := wrapText(sOut, tui.ViewportWidth)
	// return wrappedText

	return sOut
}

func printArticle(article *goquery.Selection) []string {
	var articleOut []string

	article.Contents().Each(func(i int, sel *goquery.Selection) {
		if goquery.NodeName(sel) == "h2" {
			return
		}

		loopContents := ""
		sel.Contents().Each(func(j int, s *goquery.Selection) {
			// TODO: Try to fix links. Maybe try "termlink" module

			// if goquery.NodeName(s) == "a" {
			// 	href, exists := s.Attr("href")
			// 	if exists {
			// 		// Links get made blue with an underline
			// 		articleOut += createLink(href, linkStyle.Render(s.Text()))
			// 		// articleOut += linkStyle.Render(s.Text() + "(" + href + ")")
			// 	}
			// } else

			if goquery.NodeName(s) == "em" {
				parent := s.Parent()
				if goquery.NodeName(parent) == "code" {
					// Emphatic code should get rendered as code and emphasis
					loopContents += italStyle.Render(codeStyle.Render(s.Text()))
				} else if s.HasClass("star") {
					loopContents += starStyle.Render(s.Text())
				} else {
					loopContents += italStyle.Render(s.Text())
				}
			} else if goquery.NodeName(s) == "code" {
				loopContents += codeStyle.Render(s.Text())
			} else if goquery.NodeName(s) != "h2" {
				loopContents += s.Text()
			}
		})

		articleOut = append(articleOut, wrapText(loopContents, tui.ViewportWidth)+"\n")
	})

	return articleOut
}

func (p *PageData) processPageData() {
	p.answerOne = ""
	p.answerTwo = ""
	p.articleOneSel = nil
	p.articleTwoSel = nil

	p.mainContents.Find("article").Each(func(i int, s *goquery.Selection) {
		if p.articleOneSel == nil {
			p.articleOneSel = s
		} else {
			p.articleTwoSel = s
		}
	})

	// This should only grab "Your puzzle answer was: " tags
	p.mainContents.Find("article + p").Each(func(i int, s *goquery.Selection) {
		outStr := s.Text()
		if p.answerOne == "" {
			p.answerOne = outStr
		} else {
			p.answerTwo = outStr
		}
	})
}

func wrapText(line string, width int) string {
	var result string
	words := strings.Fields(line)
	lineLength := 0

	for _, word := range words {
		if lineLength+len(word)+1 > width {
			result += "\n"
			lineLength = 0
		}
		if lineLength > 0 {
			result += " "
			lineLength++
		}

		result += word
		lineLength += len(word)
	}

	return result
}

func createLink(url string, text string) string {
	// return fmt.Sprintf("\033]8;;%s\a%s\033]8;;\a", url, text)
	return fmt.Sprintf("\x1b]8;;" + url + "\x07" + text + "\x1b]8;;\x07" + "\u001b[0m")
}
