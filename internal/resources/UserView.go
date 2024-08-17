package resources

import (
	"dalton.dog/aocgo/internal/styles"
	"dalton.dog/aocgo/internal/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

type Model struct {
	simpleTable table.Model
}

func getRowForYear(user User, year, day int) table.Row {
	stars := make([]string, 26)
	d := 1

	for d <= day {
		p := LoadOrCreatePuzzle(year, d, user.SessionTok)
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

func (u User) NewModel() Model {
	var rows []table.Row
	maxYear, maxDay := utils.GetCurrentMaxYearAndDay()

	y := utils.FIRST_YEAR
	for y <= maxYear {
		day := 25
		if y == maxYear {
			day = maxDay
		}
		rows = append(rows, getRowForYear(u, y, day))
		y++
	}

	return Model{
		simpleTable: table.New([]table.Column{
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
		}).WithRows(rows),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.simpleTable, cmd = m.simpleTable.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)
		}
	}

	return m, tea.Batch(cmds...)
}

// TODO: Properly style this
func (m Model) View() string {
	sOut := "User Breakdown"

	sOut += "\nPress q or ctrl+c to quit\n\n"

	sOut += m.simpleTable.View()

	return sOut
}
