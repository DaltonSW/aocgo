package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"dalton.dog/aocgo/internal/api"
	"dalton.dog/aocgo/internal/cache"
	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/log"
)

var (
	HeaderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Bold(true).Align(lipgloss.Center)
	BorderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	FirstPlace  = lipgloss.NewStyle().Foreground(lipgloss.Color("#d4af37"))
	SecondPlace = lipgloss.NewStyle().Foreground(lipgloss.Color("#c0c0c0"))
	ThirdPlace  = lipgloss.NewStyle().Foreground(lipgloss.Color("#cd7f32"))
	CellStyle   = lipgloss.NewStyle().Width(12)
)

// Leaderboard is a representation of a year or year+day leaderboard
type Leaderboard struct {
	year     int
	day      int
	bucketID string
	places   []*Placing
}

// Placing represents a single placing on a leaderboard
type Placing struct {
	score       int
	displayName string
	userID      string
	href        string
	placement   int
}

func (l *Leaderboard) GetID() string                { return l.bucketID }
func (l *Leaderboard) GetBucketName() string        { return cache.LEADERBOARDS }
func (l *Leaderboard) MarshalData() ([]byte, error) { return json.Marshal(l) }
func (l *Leaderboard) SaveResource()                { cache.SaveResource(l) }

func NewLeaderboard(year, day int) *Leaderboard {
	bucketID := strconv.Itoa(day) + strconv.Itoa(year)
	data := cache.LoadResource(cache.LEADERBOARDS, bucketID)
	var lb *Leaderboard
	if data != nil {
		json.Unmarshal(data, &lb)
		return lb
	}
	lb = &Leaderboard{
		year:     year,
		day:      day,
		bucketID: strconv.Itoa(day) + strconv.Itoa(year),
		places:   []*Placing{},
	}

	lb.loadPlacements()
	lb.SaveResource()

	return lb
}

// Attempts to display a leaderboard for a given year or year+day
func (l *Leaderboard) Display() {
	// Try load from storage. If can't, load from API
	if len(l.places) == 0 {
		l.loadPlacements()
	}

	log.Debug("Leaderboard data loaded")
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		Headers("Place", "Score", "Display Name").
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				return HeaderStyle
			case row == 1:
				return FirstPlace
			case row == 2:
				return SecondPlace
			case row == 3:
				return ThirdPlace
			default:
				if col == 2 {
					return lipgloss.NewStyle().Width(40)
				} else if col == 0 {
					return lipgloss.NewStyle().Width(7)
				}
				return CellStyle
			}
		})

	for _, p := range l.places {
		t.Row(strconv.Itoa(p.placement), strconv.Itoa(p.score), p.displayName)
	}

	// fmt.Println(t.Render())
}

func (l *Leaderboard) loadPlacements() error {
	URL := fmt.Sprintf("https://adventofcode.com/%v/leaderboard", l.year)
	if l.day != 0 {
		URL += fmt.Sprintf("/day/%v", l.day)
	}

	resp, err := api.NewGetReq(URL, "")

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	placings, err := parseLeaderboardHTML(resp)
	if err != nil {
		return err
	}

	l.places = placings
	return nil
}

func parseLeaderboardHTML(resp *http.Response) ([]*Placing, error) {
	bodyDoc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	placings := make([]*Placing, 0, 100)

	var intPlace int

	bodyDoc.Find("div.leaderboard-entry").Each(func(i int, s *goquery.Selection) {
		rowText := s.Text()
		// log.Debug("Row text", "text", s.Text())

		var remainder string

		if len(rowText) > 4 && rowText[3] == ')' {
			rowText = strings.TrimSpace(s.Text())
			splitRow := strings.SplitN(rowText, " ", 2)

			placement := splitRow[0]
			placement = placement[:len(placement)-1]

			intPlace, err = strconv.Atoi(placement)
			if err != nil {
				log.Error("Error parsing placement", "err", err)
			}

			remainder = strings.TrimSpace(splitRow[1])
		} else {
			remainder = strings.TrimSpace(s.Text())
		}
		splitRemainder := strings.SplitN(remainder, " ", 2)
		totalScore, err := strconv.Atoi(splitRemainder[0])
		if err != nil {
			log.Error("Error parsing score", "err", err)
		}

		displayName := strings.TrimSpace(splitRemainder[1])
		userID, _ := s.Attr("data-user-id")

		placings = append(placings, &Placing{
			score:       totalScore,
			userID:      userID,
			displayName: displayName,
			placement:   intPlace,
		})
	})

	return placings, nil
}

func (l *Leaderboard) saveToDisk() {
	ID := strconv.Itoa(l.year) + strconv.Itoa(l.day)
	dataToSave, err := json.Marshal(l.places)
	log.Debug("Trying to save leaderboard to disk", "data", dataToSave)
	if err != nil {
		log.Error("Unable to save leaderboard to disk", "err", err)
	}
	cache.SaveGenericResource(cache.LEADERBOARDS, ID, dataToSave)
}
