package models

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	// "github.com/muesli/reflow/wordwrap"
	// "github.com/muesli/termenv"
	// These might be useful
	// "github.com/muesli/reflow/indent"
	// "github.com/muesli/reflow/padding"
)

// TODO: Replace lipgloss with `reflow` and `termenv`
// https://github.com/muesli/reflow
// https://github.com/muesli/termenv

type PageData struct {
	// Both of these two can be found in the page's <head>'s <title> tag
	day  int
	year int

	// Once we get into the <main> tag, the first child should be <article class="day-desc">
	// Below that, there should be a header that has the day's title. Ex: '--- Day 1: Trebuchet?! ---'
	header string

	// The article consists of articleContents
	articleContents *goquery.Selection
}

func NewPageData(raw []byte) PageData {
	reader := bytes.NewReader(raw)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal("Error constructing new PageData.", "error", err)
	}

	title := strings.Split(doc.Find("title").Text(), " ")
	dayDesc := doc.Find(".day-desc")
	header := dayDesc.Find("h2").Text()

	day, _ := strconv.Atoi(title[1])
	year, _ := strconv.Atoi(title[len(title)-1])
	return PageData{
		header:          header,
		day:             day,
		year:            year,
		articleContents: dayDesc,
	}
}

const ParagraphWidth = 70

var (
	titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#FFFF00"))
	italStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#E10045"))
	starStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#F1FA8C"))
	linkStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#8BE9FD")).Underline(true)
	codeStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAC3D5")).Bold(true)
	wordWrap   = lipgloss.NewStyle().Width(ParagraphWidth)
)

func (p *PageData) PrintPageData() {
	titleWidth := lipgloss.Width(p.header)
	titlePad := (ParagraphWidth - titleWidth) / 2
	titleStyle.PaddingLeft(titlePad).PaddingRight(titlePad)
	sOut := titleStyle.Render(p.header) + "\n\n"

	p.articleContents.Find("p, pre").Each(func(i int, s *goquery.Selection) {
		s.Contents().Each(func(j int, sel *goquery.Selection) {
			if goquery.NodeName(sel) == "a" {
				href, exists := sel.Attr("href")
				if exists {
					linkText := linkStyle.Render(sel.Text())
					sOut += createLink(href, linkText)
				}
			} else if goquery.NodeName(sel) == "em" {
				parent := sel.Parent()
				if goquery.NodeName(parent) == "code" {
					sOut += codeStyle.Render(sel.Text())
				}
				if sel.HasClass("star") {
					sOut += starStyle.Render(sel.Text())
				} else {
					sOut += italStyle.Render(sel.Text())
				}
			} else if goquery.NodeName(sel) == "code" {
				sOut += codeStyle.Render(sel.Text())
			} else {
				sOut += sel.Text()
			}
		})
		sOut += "\n\n"
	})

	fmt.Print("\033[H\033[2J")
	fmt.Println(wordWrap.Render(sOut))
}

func createLink(url string, text string) string {
	return fmt.Sprintf("\033]8;;%s\a%s\033]8;;\a", url, text)
}
