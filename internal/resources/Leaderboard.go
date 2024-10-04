package resources

import (
	"fmt"
	"strconv"
	"strings"

	"go.dalton.dog/aocgo/internal/api"
	"go.dalton.dog/aocgo/internal/styles"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/log"
)

// Placing represents a single placing on a leaderboard
type Placing struct {
	DisplayName string
	UserID      string
	Position    int
	IsSupporter bool

	// UserLink is the link for the username
	UserLink string
	// SponsorLink is the link for the user's associated sponsor
	SponsorLink string

	// FinishTime is only used for daily leaderboards
	FinishTime string

	// Score is only used for yearly leaderboards
	Score int
}

// Leaderboard is a wrapper for either a yearly or daily leaderboard
type Leaderboard struct {
	Year int
	Day  int

	FirstHundred  []*Placing
	SecondHundred []*Placing

	BucketName string
}

// NewLeaderboard will create a leaderboard object based on the parameters.
// If you want to create a leaderboard for an entire year, pass in 0 for day
func NewLeaderboard(year, day int) *Leaderboard {
	lb := &Leaderboard{
		Year:         year,
		Day:          day,
		FirstHundred: make([]*Placing, 0, 100),
	}
	if day > 0 {
		lb.SecondHundred = make([]*Placing, 0, 100)
	}

	lb.LoadPlacings()

	return lb
}

// LoadPlacings will load all of the placings for a given year or date
func (lb *Leaderboard) LoadPlacings() {
	if lb.Day == 0 {
		lb.loadYearlyLB()
	} else {
		lb.loadDailyLB()
	}

}

// GetTitle will get the appropriate viewport title for the leaderboard
func (lb *Leaderboard) GetTitle() string {
	if lb.Day == 0 {
		return fmt.Sprintf("Leaderboard -- Year: %d", lb.Year)
	}
	return fmt.Sprintf("Leaderboard -- Year: %d, Day: %d", lb.Year, lb.Day)
}

// GetContent will get the lb content in a printable format
func (lb *Leaderboard) GetContent() string {
	if lb.Day == 0 {
		return lb.getYearlyContent()
	} else {
		return lb.getDailyContent()
	}
}

func (lb *Leaderboard) loadYearlyLB() error {
	URL := fmt.Sprintf("https://adventofcode.com/%v/leaderboard", lb.Year)
	resp, err := api.NewGetReq(URL, "")

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyDoc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	placings := make([]*Placing, 0, 100)

	var intPlace int

	bodyDoc.Find("div.leaderboard-entry").Each(func(i int, s *goquery.Selection) {
		rowText := s.Text()
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
			Score:       totalScore,
			UserID:      userID,
			DisplayName: displayName,
			Position:    intPlace,
		})
	})
	lb.FirstHundred = placings

	return nil
}

func (lb *Leaderboard) getYearlyContent() string {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		Headers("Place", "Score", "Display Name").
		StyleFunc(styles.GetLeaderboardStyle)

	for _, p := range lb.FirstHundred {
		t.Row(strconv.Itoa(p.Position), strconv.Itoa(p.Score), p.DisplayName)
	}

	return t.Render()
}

// LoadPositions will get all of the daily positions
func (lb *Leaderboard) loadDailyLB() error {
	URL := fmt.Sprintf("https://adventofcode.com/%v/leaderboard/day/%v", lb.Year, lb.Day)
	resp, err := api.NewGetReq(URL, "")

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyDoc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	placings := make([]*Placing, 0, 100)

	firstPass := true
	var intPlace int

	bodyDoc.Find("div.leaderboard-entry").Each(func(i int, s *goquery.Selection) {
		rowText := s.Text()
		var remainder string

		if len(rowText) > 4 && rowText[3] == ')' {
			rowText = strings.TrimSpace(s.Text())
			splitRow := strings.SplitN(rowText, " ", 2)

			placement := splitRow[0]
			placement = placement[:len(placement)-1] // This trims the lingering parenthesis

			intPlace, err = strconv.Atoi(placement)
			if err != nil {
				log.Error("Error parsing placement", "err", err)
				return
			}

			remainder = strings.TrimSpace(splitRow[1])
		} else {
			remainder = strings.TrimSpace(s.Text())
		}

		splitRemainder := strings.SplitN(remainder, "  ", 3)
		finishTime := splitRemainder[0] + " " + splitRemainder[1]

		displayName := strings.TrimSpace(splitRemainder[2])
		userID, _ := s.Attr("data-user-id")

		if intPlace == 1 {
			if firstPass {
				firstPass = false
			} else {
				lb.FirstHundred = placings
				placings = make([]*Placing, 0)
			}
		}

		placings = append(placings, &Placing{
			FinishTime:  finishTime,
			UserID:      userID,
			DisplayName: displayName,
			Position:    intPlace,
		})
	})

	lb.SecondHundred = placings

	return nil
}

// GetContent will get the lb content in a printable format
func (lb *Leaderboard) getDailyContent() string {
	tOne := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		Headers("Place", "Time Done (EST)", "Display Name").
		StyleFunc(styles.GetLeaderboardStyle)

	for _, p := range lb.FirstHundred {
		pos := strconv.Itoa(p.Position)
		ft := p.FinishTime
		name := p.DisplayName
		tOne.Row(pos, ft, name)
	}

	sOut := "First People to Obtain Both Stars\n" + tOne.Render()

	tTwo := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		Headers("Place", "Time Done (EST)", "Display Name").
		StyleFunc(styles.GetLeaderboardStyle)

	for _, p := range lb.SecondHundred {
		pos := strconv.Itoa(p.Position)
		ft := p.FinishTime
		name := p.DisplayName
		tTwo.Row(pos, ft, name)
	}

	sOut += "\nFirst People to Obtain The First Star\n" + tTwo.Render()

	return sOut
}
