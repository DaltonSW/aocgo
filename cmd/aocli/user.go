package main

import (
	"github.com/spf13/cobra"
	"go.dalton.dog/aocgo/internal/cache"
	"go.dalton.dog/aocgo/internal/resources"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Shows a visual representation of the user's puzzle progress.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		User(UserRsrc, ClearUser)
	},
}

// User will print out a table visualization of the user's star progress.
// Command: `aocli user [--clear]`
func User(user *resources.User, clearUser bool) {
	if ClearUser {
		cache.ClearUserDatabase(UserRsrc.SessionTok)
	} else {
		UserRsrc.Display()
	}
}
