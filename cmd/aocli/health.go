package main

import (
	"github.com/spf13/cobra"
	"go.dalton.dog/aocgo/internal/api"
	"go.dalton.dog/aocgo/internal/output"
	"go.dalton.dog/aocgo/internal/resources"
	"go.dalton.dog/aocgo/internal/session"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Checks if aocli and aocgo have proper config to run.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		Health()
	},
}

// Health will check if a session key is available so that the program can run.
// Command: `aocli health`
func Health() {
	sessionToken, err := session.GetSessionToken(true)

	if err != nil {
		output.Fatal("Test failed! Couldn't properly load a session token.", "err", err)
	}

	output.Success("Session token check success!")

	api.InitClient(sessionToken)

	output.Success("API Client initialization check success!")

	// user, err := resources.NewUser(sessionToken)

	resources.LoadOrCreatePuzzle(2016, 1, sessionToken)

	output.Success("Session token appears to be valid, happy solving!")
}
