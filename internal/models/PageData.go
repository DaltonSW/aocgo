package models

// import (
// 	"encoding/json"
// 	"fmt"
// 	"strconv"
// 	"strings"
//
// 	"dalton.dog/aocgo/internal/api"
// 	"dalton.dog/aocgo/internal/cache"
// 	"dalton.dog/aocgo/internal/tui"
// 	"github.com/PuerkitoBio/goquery" // Bless this package
// 	"github.com/charmbracelet/lipgloss"
// 	"github.com/charmbracelet/log"
// )
//
// // Information about the actual HTML contents of the page for a certain puzzle
// type PageData struct {
// 	day      int
// 	year     int
// 	bucketID string
//
// 	// Stuff relevant to printing
// 	PuzzleTitle     string
// 	PrettyStringArr []string
// 	AnswerOne       string
// 	AnswerTwo       string
//
// 	// The article consists of articleContents (as you might expect)
// 	articleOneSel goquery.Selection
// 	articleTwoSel goquery.Selection
// 	mainContents  goquery.Selection
// }
//
// func (p *PageData) GetID() string               { return p.bucketID }
// func (p *PageData) GetBucketName() string       { return cache.PAGE_DATA }
// func (p PageData) MarshalData() ([]byte, error) { return json.Marshal(p) }
// func (p *PageData) SaveResource()               { cache.SaveResource(p) }
//
// func LoadOrCreatePageData(year, day int, userSession, URL string) *PageData {
// 	bucketID := strconv.Itoa(year) + strconv.Itoa(day)
// 	data := cache.LoadResource(cache.PAGE_DATA, bucketID)
//
// 	if data != nil {
// 		var pageData *PageData
// 		json.Unmarshal(data, &pageData)
// 		return pageData
// 	}
//
// 	return NewPageData(userSession, URL)
//
// }
//
// func NewPageData(userSesssion, URL string) *PageData {
// 	resp, err := api.NewGetReq(URL, userSesssion)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	defer resp.Body.Close()
// 	doc, err := goquery.NewDocumentFromReader(resp.Body)
// 	if err != nil {
// 		log.Fatal("Error constructing new PageData.", "error", err)
// 	}
//
// 	// HACK: Admittedly making some assumptions on input, but should be fine. AoC is very consistent
// 	title := strings.Split(doc.Find("title").Text(), " ")
// 	day, _ := strconv.Atoi(title[1])
// 	year, _ := strconv.Atoi(title[len(title)-1])
//
// 	mainContents := doc.Find("main")
// 	puzzleTitle := mainContents.Find("h2").First().Text()
//
// 	pageData := &PageData{
// 		PuzzleTitle:  titleStyle.Render(puzzleTitle),
// 		day:          day,
// 		year:         year,
// 		bucketID:     strconv.Itoa(year) + strconv.Itoa(day),
// 		mainContents: *mainContents,
// 	}
//
// 	pageData.processPageData()
// 	pageData.PrettyStringArr = pageData.GetPageDataPrettyString()
// 	pageData.SaveResource()
// 	return pageData
// }
//
// // Stylings
// // const ParagraphWidth = 120
//
// // TODO: answerStyle
// var (
// 	titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#FFFF00"))
// 	italStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF3374"))
// 	starStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#F1FA8C"))
// 	linkStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#8BE9FD")).Underline(true)
// 	codeStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAC3D5")).Bold(true)
// 	wordWrap   = lipgloss.NewStyle().Width(tui.ViewportWidth)
// )
//
// func (p *PageData) GetPageDataPrettyString() []string {
// 	if len(p.PrettyStringArr) != 0 {
// 		return p.PrettyStringArr
// 	}
//
// 	p.processPageData()
//
// 	sOut := printArticle(p.articleOneSel)
//
// 	if p.AnswerOne != "" {
// 		sOut = append(sOut, "\n"+p.AnswerOne)
// 	}
//
// 	if len(p.articleTwoSel.Nodes) != 0 {
// 		sOut = append(sOut, "\n\n")
// 		sOut = append(sOut, "\n"+titleStyle.Render("--- Part Two ---"))
// 		sOut = append(sOut, "\n")
// 		sOut = append(sOut, printArticle(p.articleTwoSel)...)
// 		sOut = append(sOut, "\n")
//
// 		if p.AnswerTwo != "" {
// 			sOut = append(sOut, p.AnswerTwo)
// 		}
// 	}
//
// 	// wrappedText := wrapText(sOut, tui.ViewportWidth)
// 	// return wrappedText
//
// 	p.PrettyStringArr = sOut
// 	p.SaveResource()
//
// 	return sOut
// }
//
// func printArticle(article goquery.Selection) []string {
// 	var articleOut []string
//
// 	article.Contents().Each(func(i int, sel *goquery.Selection) {
// 		if goquery.NodeName(sel) == "h2" {
// 			return
// 		}
//
// 		loopContents := ""
// 		sel.Contents().Each(func(j int, s *goquery.Selection) {
// 			// TODO: Try to fix links. Maybe try "termlink" module
//
// 			// if goquery.NodeName(s) == "a" {
// 			// 	href, exists := s.Attr("href")
// 			// 	if exists {
// 			// 		// Links get made blue with an underline
// 			// 		articleOut += createLink(href, linkStyle.Render(s.Text()))
// 			// 		// articleOut += linkStyle.Render(s.Text() + "(" + href + ")")
// 			// 	}
// 			// } else
//
// 			if goquery.NodeName(s) == "em" {
// 				parent := s.Parent()
// 				if goquery.NodeName(parent) == "code" {
// 					// Emphatic code should get rendered as code and emphasis
// 					loopContents += italStyle.Render(codeStyle.Render(s.Text()))
// 				} else if s.HasClass("star") {
// 					loopContents += starStyle.Render(s.Text())
// 				} else {
// 					loopContents += italStyle.Render(s.Text())
// 				}
// 			} else if goquery.NodeName(s) == "code" {
// 				loopContents += codeStyle.Render(s.Text())
// 			} else if goquery.NodeName(s) != "h2" {
// 				loopContents += s.Text()
// 			}
// 		})
//
// 		articleOut = append(articleOut, wrapText(loopContents, tui.ViewportWidth)+"\n")
// 	})
//
// 	return articleOut
// }
//
// func (p *PageData) processPageData() {
// 	p.AnswerOne = ""
// 	p.AnswerTwo = ""
// 	p.articleOneSel = goquery.Selection{}
// 	p.articleTwoSel = goquery.Selection{}
//
// 	p.mainContents.Find("article").Each(func(i int, s *goquery.Selection) {
// 		if len(p.articleOneSel.Nodes) == 0 {
// 			p.articleOneSel = *s
// 		} else {
// 			p.articleTwoSel = *s
// 		}
// 	})
//
// 	// This should only grab "Your puzzle answer was: " tags
// 	p.mainContents.Find("article + p").Each(func(i int, s *goquery.Selection) {
// 		outStr := s.Text()
// 		if p.AnswerOne == "" {
// 			p.AnswerOne = outStr
// 		} else {
// 			p.AnswerTwo = outStr
// 		}
// 	})
// }
//
// func createLink(url string, text string) string {
// 	// return fmt.Sprintf("\033]8;;%s\a%s\033]8;;\a", url, text)
// 	return fmt.Sprintf("\x1b]8;;" + url + "\x07" + text + "\x1b]8;;\x07" + "\u001b[0m")
// }
