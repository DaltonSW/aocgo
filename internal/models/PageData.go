package models

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

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

	titleWidth := lipgloss.Width(puzzleTitle)
	titlePad := (ParagraphWidth - titleWidth) / 2
	titleStyle.PaddingLeft(titlePad).PaddingRight(titlePad)

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
const ParagraphWidth = 80

// TODO: answerStyle
var (
	titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#FFFF00"))
	italStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#E10045"))
	starStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#F1FA8C"))
	linkStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#8BE9FD")).Underline(true)
	codeStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAC3D5")).Bold(true)
	wordWrap   = lipgloss.NewStyle().Width(ParagraphWidth)
)

// TODO: Rewrite function to return a string instead of just print so it can be passed into viewport
// ... OR maybe just make a PageData model (or add a Page Data to the viewport), store the data, and then call this function in the View() method?

func (p *PageData) PrintPageData() {
	p.processPageData()

	fmt.Print("\033[H\033[2J") // Clear terminal
	fmt.Println(wordWrap.Render(printArticle(p.articleOneSel)))

	if p.answerOne != "" {
		fmt.Println(wordWrap.Render(p.answerOne))
	}

	if p.articleTwoSel != nil {
		fmt.Println(wordWrap.Render(printArticle(p.articleTwoSel)))

		if p.answerTwo != "" {
			fmt.Println(wordWrap.Render(p.answerTwo))
		}
	}

}

func (p *PageData) GetPageDataPrettyString() string {
	p.processPageData()

	sOut := printArticle(p.articleOneSel)

	if p.answerOne != "" {
		sOut += p.answerOne
	}

	if p.articleTwoSel != nil {
		titleWidth := lipgloss.Width("--- Part Two ---")
		titlePad := (ParagraphWidth - titleWidth) / 2
		titleStyle.PaddingLeft(titlePad).PaddingRight(titlePad)
		sOut += "\n" + titleStyle.Render("--- Part Two ---")
		sOut += "\n" + printArticle(p.articleTwoSel)

		if p.answerTwo != "" {
			sOut += p.answerTwo
		}
	}

	return sOut
}

func printArticle(article *goquery.Selection) string {
	articleOut := ""

	article.Contents().Each(func(i int, sel *goquery.Selection) {
		if goquery.NodeName(sel) == "h2" {
			return
		}

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
					articleOut += italStyle.Render(codeStyle.Render(s.Text()))
				} else if s.HasClass("star") {
					articleOut += starStyle.Render(s.Text())
				} else {
					articleOut += italStyle.Render(s.Text())
				}
			} else if goquery.NodeName(s) == "code" {
				articleOut += codeStyle.Render(s.Text())
			} else if goquery.NodeName(s) != "h2" {
				articleOut += s.Text()
			}
		})

		articleOut += "\n"
	})

	return "\n" + articleOut
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

func createLink(url string, text string) string {
	// return fmt.Sprintf("\033]8;;%s\a%s\033]8;;\a", url, text)
	return fmt.Sprintf("\x1b]8;;" + url + "\x07" + text + "\x1b]8;;\x07" + "\u001b[0m")
}
