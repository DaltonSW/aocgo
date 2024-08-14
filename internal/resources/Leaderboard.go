package resources

import (
	// "encoding/json"
	"fmt"
	"strconv"

	"dalton.dog/aocgo/internal/api"
	// "dalton.dog/aocgo/internal/cache"

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

// Placing represents a single placing on a leaderboard
type Placing struct {
	DisplayName string
	UserID      string
	Position    string
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

// YearLB is the model to handle the leaderboard for a year as a whole
type YearLB struct {
	Year      int
	Positions []*Placing
}

func (l *YearLB) GetTitle() string {
	return fmt.Sprintf("Leaderboard -- Year: %d", l.Year)
}

// func (l *YearLB) GetID() string                { return strconv.Itoa(l.Year) + "0" }
// func (l *YearLB) GetBucketName() string        { return cache.LEADERBOARDS }
// func (l *YearLB) MarshalData() ([]byte, error) { return json.Marshal(l) }
// func (l *YearLB) SaveResource()                { cache.SaveResource(l) }

func NewYearLB(year int) *YearLB {

	// data := cache.LoadResource(cache.LEADERBOARDS, strconv.Itoa(year)+"0")
	// var lb *YearLB
	// if data != nil {
	// 	json.Unmarshal(data, &lb)
	// 	return lb
	// }

	lb := &YearLB{
		Year:      year,
		Positions: make([]*Placing, 0, 100),
	}

	lb.LoadPositions()
	// lb.SaveResource()

	return lb
}

func (lb *YearLB) LoadPositions() error {
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

	bodyDoc.Find("div.leaderboard-entry").Each(func(i int, sel *goquery.Selection) {
		log.Debug(sel.Text())
		placing := &Placing{}

		placing.Position = sel.Find(".leaderboard-position").Text()

		s, _ := strconv.Atoi(sel.Find(".leaderboard-totalscore").Text())
		placing.Score = s

		sel.Find("a[target='_blank']").Each(func(j int, aTag *goquery.Selection) {
			if aTag.HasClass(".sponsor-badge") {
				link, _ := aTag.Attr("href")
				placing.SponsorLink = link
			} else {
				placing.DisplayName = aTag.Text()
				link, _ := aTag.Attr("href")
				placing.UserLink = link
			}
		})

		if placing.DisplayName == "" {
			placing.DisplayName = sel.FilterFunction(func(i int, s *goquery.Selection) bool {
				return goquery.NodeName(s) == "#text"
			}).Text()
		}

		if sel.Find(".supporter-badge").Length() > 0 {
			placing.IsSupporter = true
		}

		placings = append(placings, placing)

	})

	lb.Positions = placings

	return nil
}

func (lb *YearLB) GetContent() string {
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

	for _, p := range lb.Positions {
		t.Row(p.Position, strconv.Itoa(p.Score), p.DisplayName)
	}

	return t.Render()
}

// DayLB is the model to handle a specific day's leaderboard
type DayLB struct {
	Year int
	Day  int

	BothStars []*Placing
	FirstStar []*Placing
}

func (l *DayLB) GetTitle() string {
	return fmt.Sprintf("Leaderboard -- Year: %d, Day: %d", l.Year, l.Day)
}

// func (l *DayLB) GetID() string                { return strconv.Itoa(l.Year) + strconv.Itoa(l.Day) }
// func (l *DayLB) GetBucketName() string        { return cache.LEADERBOARDS }
// func (l *DayLB) MarshalData() ([]byte, error) { return json.Marshal(l) }
// func (l *DayLB) SaveResource()                { cache.SaveResource(l) }

func NewDayLB(year, day int) *DayLB {
	// data := cache.LoadResource(cache.LEADERBOARDS, strconv.Itoa(year)+strconv.Itoa(day))
	// var lb *DayLB
	// if data != nil {
	// 	json.Unmarshal(data, &lb)
	// 	return lb
	// }

	lb := &DayLB{
		Year:      year,
		Day:       day,
		FirstStar: make([]*Placing, 100),
		BothStars: make([]*Placing, 100),
	}

	lb.LoadPositions()
	// lb.SaveResource()

	return lb
}

func (lb *DayLB) LoadPositions() error {
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

	bodyDoc.Find("div.leaderboard-entry").Each(func(i int, sel *goquery.Selection) {
		placing := &Placing{}

		placing.Position = sel.Find(".leaderboard-position").Text()
		if placing.Position == "1)" {
			if firstPass {
				firstPass = false
			} else {
				lb.BothStars = placings
				placings = make([]*Placing, 0, 100)
			}
		}

		placing.FinishTime = sel.Find(".leaderboard-time").Text()

		sel.Find("a[target='_blank']").Each(func(j int, aTag *goquery.Selection) {
			if aTag.HasClass(".sponsor-badge") {
				link, _ := aTag.Attr("href")
				placing.SponsorLink = link
			} else {
				placing.DisplayName = aTag.Text()
				link, _ := aTag.Attr("href")
				placing.UserLink = link
			}
		})

		if placing.DisplayName == "" {
			placing.DisplayName = sel.FilterFunction(func(i int, s *goquery.Selection) bool {
				return goquery.NodeName(s) == "#text"
			}).Text()
		}

		if sel.Find(".supporter-badge").Length() > 0 {
			placing.IsSupporter = true
		}

		placings = append(placings, placing)

	})

	lb.FirstStar = placings

	return nil
}

func (lb *DayLB) GetContent() string {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		Headers("Place", "Time Finished (EST)", "Display Name").
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

	for _, p := range lb.BothStars {
		pos := p.Position
		ft := p.FinishTime
		name := p.DisplayName
		t.Row(pos, ft, name)
	}

	sOut := "First People to Obtain Both Stars\n" + t.Render()

	t.ClearRows()

	for _, p := range lb.FirstStar {
		pos := p.Position
		ft := p.FinishTime
		name := p.DisplayName
		t.Row(pos, ft, name)
	}

	sOut += "First People to Obtain The First Star\n" + t.Render()

	return sOut
}
