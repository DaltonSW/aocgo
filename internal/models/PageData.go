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

	headerOne string
	headerTwo string

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

	pageData := &PageData{
		day:          day,
		year:         year,
		mainContents: mainContents,
	}

	pageData.processPageData()

	return pageData
}

// Stylings
const ParagraphWidth = 70

var (
	titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#FFFF00"))
	italStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#E10045"))
	starStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#F1FA8C"))
	linkStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#8BE9FD")).Underline(true)
	codeStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAC3D5")).Bold(true)
	wordWrap   = lipgloss.NewStyle().Width(ParagraphWidth)
)

func (p *PageData) processPageData() {
	p.answerOne = ""
	p.answerTwo = ""
	p.headerOne = ""
	p.headerTwo = ""
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

func (p *PageData) PrintPageData() {
	p.processPageData()

	titleWidth := lipgloss.Width(p.headerOne)
	titlePad := (ParagraphWidth - titleWidth) / 2
	titleStyle.PaddingLeft(titlePad).PaddingRight(titlePad)

	fmt.Print("\033[H\033[2J") // Clear terminal
	// fmt.Println(HeaderStyle.Render(p.headerOne))
	printArticle(p.articleOneSel)

	if p.answerOne != "" {
		fmt.Println(wordWrap.Render(p.answerOne))
	}

	if p.articleTwoSel != nil {
		// fmt.Println("\n" + HeaderStyle.Render(p.headerTwo))
		printArticle(p.articleTwoSel)
	}

	if p.answerTwo != "" {
		fmt.Println(wordWrap.Render(p.answerTwo))
	}

}

func printArticle(article *goquery.Selection) {
	var articleOut string

	title := article.Find("h2").Text()

	article.Contents().Each(func(i int, s *goquery.Selection) {
		if goquery.NodeName(s) == "a" {
			href, exists := s.Attr("href")
			if exists {
				// Links get made blue with an underline
				linkText := linkStyle.Render(s.Text())
				articleOut += createLink(href, linkText)
			}
		} else if goquery.NodeName(s) == "em" {
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

	titleWidth := lipgloss.Width(title)
	titlePad := (ParagraphWidth - titleWidth) / 2
	titleStyle.PaddingLeft(titlePad).PaddingRight(titlePad)

	fmt.Println(titleStyle.Render(title))

	fmt.Println(wordWrap.Render(articleOut))
}

func createLink(url string, text string) string {
	return fmt.Sprintf("\033]8;;%s\a%s\033]8;;\a", url, text)
}
