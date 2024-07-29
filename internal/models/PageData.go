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

	// Once we get into the <main> tag, the first child should be <article class="day-desc">
	// Below that, there should be a header that has the day's title. Ex: '--- Day 1: Trebuchet?! ---'
	headerOne string
	headerTwo string

	answerOne string
	answerTwo string

	// The article consists of articleContents (as you might expect)
	articleOne      string
	articleTwo      string
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
	header := mainContents.Find("h2").Text()

	return &PageData{
		headerOne:    header,
		day:          day,
		year:         year,
		mainContents: mainContents,
	}
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
	processArticleFn := func(article *goquery.Selection) string {
		title := article.Find("h2").Text()
		articleOut := title

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
			} else {
				articleOut += s.Text()
			}
		})

		return articleOut
	}

	p.mainContents.Find("article").Each(func(i int, s *goquery.Selection) {
		outStr := processArticleFn(s)
		if p.articleOne == "" {
			p.articleOne = outStr
		} else {
			p.articleTwo = outStr
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
	titleWidth := lipgloss.Width(p.headerOne)
	titlePad := (ParagraphWidth - titleWidth) / 2
	titleStyle.PaddingLeft(titlePad).PaddingRight(titlePad)

	sOut := titleStyle.Render(p.headerOne) + "\n\n"

	fmt.Print("\033[H\033[2J")
	fmt.Println(wordWrap.Render(sOut))
}

func createLink(url string, text string) string {
	return fmt.Sprintf("\033]8;;%s\a%s\033]8;;\a", url, text)
}
