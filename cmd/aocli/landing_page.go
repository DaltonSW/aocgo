package main

import (
	"fmt"
	"math/rand"

	"go.dalton.dog/aocgo/internal/styles"
	"go.dalton.dog/aocgo/internal/utils"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

const tree = `
*
/_\
/_o_\
/o___o\
/__o____\
/_o___o___\
/__o__o__o__\
/o_________o__\
/____o___o______\
/o______o____o____\
/__o___o_____o______\
/_o___o_______o___o___\
|||
|||
`

var Colors = []lipgloss.Style{
	styles.BlueTextStyle,
	styles.RedTextStyle,
	styles.PurpleTextStyle,
	styles.CyanTextStyle,
}

func RunLandingPage() {
	var sOut string
	for i := range tree {
		c := string(tree[i])
		switch c {
		case "*":
			sOut += styles.YellowTextStyle.Render(c)
		case "o":
			sOut += Colors[rand.Intn(len(Colors))].Render(c)
		case "|":
			sOut += styles.BrownTextStyle.Render(c)
		default:
			sOut += styles.GreenTextStyle.Render(c)
		}
	}

	w, _, err := utils.GetTerminalSize()

	if err != nil {
		log.Fatal(err)
	}

	footOne := "Welcome to aocli!"

	footTwo := "Run " + styles.CodeStyle.Render("aocli help") + " to see a list of available commands."

	footThree := styles.SubtitleStyle.Render("\naocli by Dalton Williams (https://dalton.dog) ")

	footFour := styles.SubtitleStyle.Render("Advent of Code by Eric Wastl (http://was.tl)")

	outStr := lipgloss.PlaceHorizontal(w, lipgloss.Center, sOut)
	outStr += lipgloss.PlaceHorizontal(w, lipgloss.Center, footOne)
	outStr += lipgloss.PlaceHorizontal(w, lipgloss.Center, footTwo)
	outStr += lipgloss.PlaceHorizontal(w, lipgloss.Center, footThree)
	outStr += lipgloss.PlaceHorizontal(w, lipgloss.Center, footFour)

	utils.ClearTerminal()
	fmt.Println("\n" + outStr)

	// fmt.Println(lipgloss.JoinVertical(lipgloss.Center, header, sOut, footer))
}
