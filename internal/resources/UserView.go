package resources

import (
	"fmt"
	"strconv"

	"go.dalton.dog/aocgo/internal/output"
	"go.dalton.dog/aocgo/internal/styles"
	"go.dalton.dog/aocgo/internal/utils"

	"github.com/charmbracelet/bubbles/v2/progress"
	"github.com/charmbracelet/bubbles/v2/spinner"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/table"
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

// Message to indicate the total number of puzzles needed to be loaded
type setTotalMsg struct {
	total int
}

// LoadUserModel is the BubbleTea model for loading and displaying
// a user's information
type LoadUserModel struct {
	user     *User
	userName string
	finished bool

	totalLoaded int
	totalToLoad int

	table    table.Table
	spinner  spinner.Model
	progress progress.Model
	status   string
}

func (u *User) NewModel() tea.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Spinner.FPS = 10
	s.Style = lipgloss.NewStyle().Foreground(styles.UpdateSpinnerColor)

	p := progress.New()

	model := LoadUserModel{
		user:     u,
		spinner:  s,
		progress: p,
		status:   "Starting up!",
	}

	return model
}

func (m LoadUserModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	startYear := utils.FIRST_YEAR

	maxYear, maxDay := utils.GetCurrentMaxYearAndDay()

	for loopYear := startYear; loopYear <= maxYear; loopYear++ {
		maxLoopDay := 25
		if loopYear == maxYear {
			maxLoopDay = maxDay
		}

		for loopDay := 1; loopDay <= maxLoopDay; loopDay++ {
			cmds = append(cmds, loadPuzzle(loopYear, loopDay, m.user.SessionTok))
		}
	}

	return tea.Sequence(
		func() tea.Msg { return setTotalMsg{total: len(cmds)} },
		m.spinner.Tick,
		tea.Batch(cmds...),
	)
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

	case progress.FrameMsg:
		m.progress, cmd = m.progress.Update(msg)
		return m, cmd

	case setTotalMsg:
		m.totalToLoad = msg.total
		return m, nil

	case loadDoneMsg:
		m.totalLoaded++
		if m.totalLoaded >= m.totalToLoad {
			return m, generateTable(m.user.GetToken())
		}

		return m, nil
		// maxYear, _ := utils.GetCurrentMaxYearAndDay()
		// year, day := msg.year, msg.day
		// if day == 25 {
		// 	if year < maxYear {
		// 		m.curDate = 1
		// 		m.curYear++
		//
		// 		m.status = fmt.Sprintf("Loading... Year %v - Day %v", m.curYear, m.curDate)
		// 		cmds = append(cmds, loadPuzzle(m.curYear, m.curDate, m.user.SessionTok))
		//
		// 	} else {
		// 		m.status = "Done loading, generating table!"
		// 		cmds = append(cmds, generateTable(m.user.GetToken()))
		// 	}
		// } else {
		// 	m.curDate++
		// 	m.status = fmt.Sprintf("Loading... Year %v - Day %v", m.curYear, m.curDate)
		// 	cmds = append(cmds, loadPuzzle(m.curYear, m.curDate, m.user.SessionTok))
		// }

	case tableDoneMsg:
		m.status = "Table is done, good to go!"
		m.finished = true
		m.table = msg.table

		utils.ClearTerminal()
		return m, tea.Quit
	}

	return m, tea.Batch(cmds...)
}

func (m LoadUserModel) View() string {
	if m.user.DisplayName == "" {
		m.user.LoadDisplayName()
	}
	if m.finished {
		sOut := fmt.Sprintf("\n%v\n%v\n", styles.IncorrectAnswerStyle.Bold(true).Render(header(m.user.DisplayName)), m.table.Render())
		return styles.GlobalSpacingStyle.Render(sOut)
	}

	outStr := fmt.Sprintf("Loading %s's stars... %s \n%s", m.user.DisplayName, m.spinner.View(), m.progress.ViewAs(float64(m.totalLoaded)/float64(m.totalToLoad)))
	return styles.GlobalSpacingStyle.Render(outStr)
}

func header(displayName string) string {
	outStr := fmt.Sprintf("%v's Star Breakdown\n", displayName)
	return lipgloss.PlaceHorizontal(ViewportWidth, lipgloss.Center, outStr)
}

func loadPuzzle(year, day int, userToken string) tea.Cmd {
	output.Debug("Entered loadPuzzle message")
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
			Border(lipgloss.RoundedBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(styles.TableBorderColor))

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
