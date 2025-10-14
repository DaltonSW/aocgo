package main

import (
	"errors"

	"github.com/spf13/cobra"
	"go.dalton.dog/aocgo/internal/cache"
	"go.dalton.dog/aocgo/internal/resources"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Commands for interfacing with different users.",
	Args:  cobra.NoArgs,
	Annotations: map[string]string{
		requiresAuthAnnotation: "true",
	},
}

var addUserCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new user.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("Not implemented")
	},
}

var clearUserCmd = &cobra.Command{
	Use:   "clear",
	Short: "Deletes the database file for the active user.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("Not implemented")
	},
}

var listUserCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all registered users alongside their labels and masked tokens.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("Not implemented")
	},
}

var switchUserCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch active user by choosing from a list.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("Not implemented")
	},
}

var removeUserCmd = &cobra.Command{
	Use:   "remove",
	Short: "Select a user to remove from a list.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("Not implemented")
	},
}

var viewUserCmd = &cobra.Command{
	Use:   "view",
	Short: "Load puzzle information for the active user and display a star table.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("Not implemented")
	},
}

// User will print out a table visualization of the user's star progress.
// Command: `aocli user [--clear]`
func User(user *resources.User, clearUser bool) {
	if clearUser {
		cache.ClearUserDatabase(user.SessionTok)
	} else {
		user.Display()
	}
}

func init() {
	userCmd.GroupID = "puzzles"
	userCmd.AddCommand(addUserCmd)
	userCmd.AddCommand(clearUserCmd)
	userCmd.AddCommand(listUserCmd)
	userCmd.AddCommand(switchUserCmd)
	userCmd.AddCommand(removeUserCmd)
	userCmd.AddCommand(viewUserCmd)
	rootCmd.AddCommand(userCmd)
}
