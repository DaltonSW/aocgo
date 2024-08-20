package resources

import (
	"fmt"

	"dalton.dog/aocgo/internal/styles"
	"dalton.dog/aocgo/internal/utils"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/evertras/bubble-table/table"
)

// Message to indicate a puzzle has finished loading
// Updates status message and triggers next load
type loadDoneMsg struct {
	year int
	day  int
}

// Message to indicate that the user table is ready to display
type tableDoneMsg struct {
	table table.Model
}

// LoadUserModel is the BubbleTea model for loading and displaying
// a user's information
type LoadUserModel struct {
	user     *User
	curYear  int
	curDate  int
	finished bool

	table   table.Model
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
		sOut := fmt.Sprintf("%v\n%v\n%v\n", header(), m.table.View(), footer())
		return styles.GlobalSpacingStyle.Render(sOut)
	} else {
		return styles.GlobalSpacingStyle.Render(m.spinner.View() + " " + m.status)
	}
}

func header() string {
	return lipgloss.PlaceHorizontal(ViewportWidth, lipgloss.Center, "User Breakdown\n\n")
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
		var rows []table.Row
		maxYear, maxDay := utils.GetCurrentMaxYearAndDay()

		y := utils.FIRST_YEAR
		for y <= maxYear {
			day := 25
			if y == maxYear {
				day = maxDay
			}
			rows = append(rows, getRowForYear(userToken, y, day))
			y++
		}

		newTable := table.New([]table.Column{
			table.NewColumn("Year", "Year", 4),
			table.NewColumn("1", "1", 2),
			table.NewColumn("2", "2", 2),
			table.NewColumn("3", "3", 2),
			table.NewColumn("4", "4", 2),
			table.NewColumn("5", "5", 2),
			table.NewColumn("6", "6", 2),
			table.NewColumn("7", "7", 2),
			table.NewColumn("8", "8", 2),
			table.NewColumn("9", "9", 2),
			table.NewColumn("10", "10", 2),
			table.NewColumn("11", "11", 2),
			table.NewColumn("12", "12", 2),
			table.NewColumn("13", "13", 2),
			table.NewColumn("14", "14", 2),
			table.NewColumn("15", "15", 2),
			table.NewColumn("16", "16", 2),
			table.NewColumn("17", "17", 2),
			table.NewColumn("18", "18", 2),
			table.NewColumn("19", "19", 2),
			table.NewColumn("20", "20", 2),
			table.NewColumn("21", "21", 2),
			table.NewColumn("22", "22", 2),
			table.NewColumn("23", "23", 2),
			table.NewColumn("24", "24", 2),
			table.NewColumn("25", "25", 2),
		}).WithRows(rows).BorderRounded().WithBaseStyle(styles.UserTableStyle)
		return tableDoneMsg{table: newTable}
	}
}

func getRowForYear(userToken string, year, day int) table.Row {
	stars := make([]string, 26)
	d := 1

	for d <= day {
		p := LoadOrCreatePuzzle(year, d, userToken)
		var sOut string
		if p.AnswerTwo != "" {
			sOut = lipgloss.NewStyle().Foreground(styles.BothStarsColor).Render("*")
		} else if p.AnswerOne != "" {
			sOut = lipgloss.NewStyle().Foreground(styles.FirstStarColor).Render("*")
		} else {
			sOut = lipgloss.NewStyle().Foreground(styles.NoStarsColor).Render("-")
		}
		stars[d] = sOut
		d++
	}

	if day < 25 {
		for d := day; d <= 25; d++ {
			stars[d] = lipgloss.NewStyle().Foreground(styles.NoStarsColor).Render("-")
		}
	}

	return table.NewRow(table.RowData{
		"Year": year,
		"1":    stars[1],
		"2":    stars[2],
		"3":    stars[3],
		"4":    stars[4],
		"5":    stars[5],
		"6":    stars[6],
		"7":    stars[7],
		"8":    stars[8],
		"9":    stars[9],
		"10":   stars[10],
		"11":   stars[11],
		"12":   stars[12],
		"13":   stars[13],
		"14":   stars[14],
		"15":   stars[15],
		"16":   stars[16],
		"17":   stars[17],
		"18":   stars[18],
		"19":   stars[19],
		"20":   stars[20],
		"21":   stars[21],
		"22":   stars[22],
		"23":   stars[23],
		"24":   stars[24],
		"25":   stars[25],
	})
}
