package main

import (
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth [token]",
	Short: "Authenticate your Advent of Code session token",
	Args:  cobra.RangeArgs(0, 1),

	Run: func(cmd *cobra.Command, args []string) {
		Auth()
	},
}

func Auth() {

}
