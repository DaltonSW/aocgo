package resources

import (
	"fmt"
	"strconv"

	"go.dalton.dog/aocgo/internal/styles"
	"go.dalton.dog/aocgo/internal/utils"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/log"
)

// Message to indicate a puzzle has finished loading
// Updates status message and triggers next load
type loadDoneMsg struct {
	year int
	day  int
}

// Message to indicate that the user table is ready to display
type tableDoneMsg struct {
	table table.Table
}

// LoadUserModel is the BubbleTea model for loading and displaying
// a user's information
type LoadUserModel struct {
	user     *User
	userName string
	curYear  int
	curDate  int
	finished bool

	table   table.Table
	spinner spinner.Model
	status  string
}

func (u *User) NewModel() tea.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Spinner.FPS = 20
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(styles.UpdateSpinnerColor))

	model := LoadUserModel{
		user:    u,
		spinner: s,
		curYear: utils.FIRST_YEAR,
		curDate: 1,
		status:  "Starting up!",
	}

	return model
}

func (m LoadUserModel) Init() tea.Cmd {
	return tea.Batch(loadPuzzle(m.curYear, m.curDate, m.user.SessionTok), m.spinner.Tick)
}

func (m LoadUserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)
		}

	case loadDoneMsg:
		maxYear, _ := utils.GetCurrentMaxYearAndDay()
		year, day := msg.year, msg.day
		if day == 25 {
			if year < maxYear {
				m.curDate = 1
				m.curYear++

				m.status = fmt.Sprintf("Loading... Year %v - Day %v", m.curYear, m.curDate)
				cmds = append(cmds, loadPuzzle(m.curYear, m.curDate, m.user.SessionTok))

			} else {
				m.status = "Done loading, generating table!"
				cmds = append(cmds, generateTable(m.user.GetToken()))
			}
		} else {
			m.curDate++
			m.status = fmt.Sprintf("Loading... Year %v - Day %v", m.curYear, m.curDate)
			cmds = append(cmds, loadPuzzle(m.curYear, m.curDate, m.user.SessionTok))
		}

	case tableDoneMsg:
		m.status = "Table is done, good to go!"
		m.finished = true
		m.table = msg.table
	}

	return m, tea.Batch(cmds...)
}

func (m LoadUserModel) View() string {
	if m.finished {
		sOut := fmt.Sprintf("%v\n%v\n%v\n", styles.NormalTextStyle.Render(header(m.user.DisplayName)), m.table.Render(), styles.NormalTextStyle.Render(footer()))
		return styles.GlobalSpacingStyle.Render(sOut)
	} else {
		return styles.GlobalSpacingStyle.Render(m.spinner.View() + " " + m.status)
	}
}

func header(displayName string) string {
	outStr := fmt.Sprintf("%v's Star Breakdown\n", displayName)
	return lipgloss.PlaceHorizontal(ViewportWidth, lipgloss.Center, outStr)
}

func footer() string {
	return lipgloss.PlaceHorizontal(ViewportWidth, lipgloss.Center, "\nPress q or ctrl+c to quit\n")

}

func loadPuzzle(year, day int, userToken string) tea.Cmd {
	log.Debug("Entered loadPuzzle message")
	return func() tea.Msg {
		LoadOrCreatePuzzle(year, day, userToken)
		return loadDoneMsg{year: year, day: day}
	}
}

func generateTable(userToken string) tea.Cmd {
	return func() tea.Msg {
		maxYear, maxDay := utils.GetCurrentMaxYearAndDay()

		t := table.New().
			Headers("Year", "01", "02", "03", "04", "05", "06", "07", "08", "09",
				"10", "11", "12", "13", "14", "15", "16", "17", "18", "19",
				"20", "21", "22", "23", "24", "25", "Num").
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99")))

		y := utils.FIRST_YEAR
		for y <= maxYear {
			day := 25
			if y == maxYear {
				day = maxDay
			}
			t.Row(getRowForYear(userToken, y, day)...)
			y++
		}

		return tableDoneMsg{table: *t}
	}
}

func getRowForYear(userToken string, year, day int) []string {
	stars := make([]string, 27)
	d := 1
	numStars := 0

	for d <= day {
		p := LoadOrCreatePuzzle(year, d, userToken)
		var sOut string
		if p.AnswerTwo != "" {
			sOut = lipgloss.NewStyle().Foreground(styles.BothStarsColor).Render("*")
			numStars += 2
		} else if p.AnswerOne != "" {
			if numStars == 48 {
				sOut = lipgloss.NewStyle().Foreground(styles.BothStarsColor).Render("*")
				numStars += 2
			} else {
				sOut = lipgloss.NewStyle().Foreground(styles.FirstStarColor).Render("*")
				numStars += 1
			}
		} else {
			sOut = lipgloss.NewStyle().Foreground(styles.NoStarsColor).Render(".")
		}
		stars[d] = sOut
		d++
	}

	if day < 25 {
		for d := day + 1; d <= 25; d++ {
			stars[d] = lipgloss.NewStyle().Foreground(styles.NoStarsColor).Render("-")
		}
	}

	stars[0] = strconv.Itoa(year)
	stars[26] = strconv.Itoa(numStars)

	return stars
}
