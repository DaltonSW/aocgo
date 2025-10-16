package main

import (
	"fmt"

	"github.com/charmbracelet/huh/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/table"
	"github.com/spf13/cobra"

	"go.dalton.dog/aocgo/internal/api"
	"go.dalton.dog/aocgo/internal/cache"
	"go.dalton.dog/aocgo/internal/output"
	"go.dalton.dog/aocgo/internal/resources"
	"go.dalton.dog/aocgo/internal/session"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage Advent of Code users and session tokens.",
	Args:  cobra.NoArgs,
}

var addUserCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new user and optionally make them active.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var label, token string
		setActive := true

		err := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("Token Label").Value(&label),
				huh.NewInput().Title("Session Token").EchoMode(huh.EchoModePassword).Value(&token),
				huh.NewConfirm().Title("Set as active user?").Value(&setActive),
			),
		).Run()

		if err != nil {
			return err
		}

		active, err := session.AddUser(label, token, setActive)
		if err != nil {
			return err
		}

		output.Success("User added.", "label", label)
		if active != nil && active.Label == label {
			output.Success("Active user updated.", "label", active.Label)
		}

		return nil
	},
}

var clearUserCmd = &cobra.Command{
	Use:   "clear",
	Short: "Delete cached data for the active user.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		active, err := session.GetActiveUser(false)
		if err != nil {
			return err
		}

		var confirm bool
		huh.NewConfirm().Title("Delete locally cached data for this user?").Value(&confirm).Run()
		if !confirm {
			return nil
		}

		cache.ClearUserDatabase(active.Label)
		output.Success("Cleared user cache.", "label", active.Label)
		return nil
	},
}

var listUserCmd = &cobra.Command{
	Use:   "list",
	Short: "List stored users with their labels and masked tokens.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		users, err := session.ListUsers()
		if err != nil {
			return err
		}

		if len(users) == 0 {
			output.Info("No users registered yet. Run `aocli user add` to create one.")
			return nil
		}

		userTable := table.New().Headers(" ", "User", "Token").Border(lipgloss.HiddenBorder())

		for _, user := range users {
			marker := " "
			if user.Active {
				marker = "*"
			}
			userTable.Row(marker, user.Label, api.MaskSecret(user.Token))
		}

		lipgloss.Println(userTable.Render())
		return nil
	},
}

var switchUserCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch the active user from a list of stored users.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		users, err := session.ListUsers()
		if err != nil {
			return err
		}
		if len(users) == 0 {
			return fmt.Errorf("no users available; add one with `aocli user add`")
		}

		var selected string
		options := make([]huh.Option[string], 0, len(users))
		for _, user := range users {
			label := user.Label
			if user.Active {
				label += " (active)"
			}
			options = append(options, huh.NewOption(label, user.Label))
		}

		if err := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().Title("Choose a user").Options(options...).Value(&selected),
			),
		).Run(); err != nil {
			return err
		}

		if selected == "" {
			output.Info("No user selected; active user unchanged.")
			return nil
		}

		active, err := session.SetActiveUser(selected)
		if err != nil {
			return err
		}

		output.Success("Active user set.", "label", active.Label)
		return nil
	},
}

var updateUserCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the active user's session token.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		active, err := session.GetActiveUser(false)
		if err != nil {
			return err
		}

		var token string
		if err := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title(fmt.Sprintf("New session token for %s", active.Label)).
					EchoMode(huh.EchoModePassword).
					Value(&token),
			),
		).Run(); err != nil {
			return err
		}

		updated, err := session.UpdateActiveUserToken(token)
		if err != nil {
			return err
		}

		output.Success("Session token updated.", "label", updated.Label)
		return nil
	},
}

var removeUserCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a stored user.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		users, err := session.ListUsers()
		if err != nil {
			return err
		}
		if len(users) == 0 {
			return fmt.Errorf("no users available to remove")
		}

		var selected string
		options := make([]huh.Option[string], 0, len(users))
		for _, user := range users {
			options = append(options, huh.NewOption(user.Label, user.Label))
		}

		if err := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().Title("Remove which user?").Options(options...).Value(&selected),
			),
		).Run(); err != nil {
			return err
		}

		if selected == "" {
			output.Info("No user removed.")
			return nil
		}

		active, err := session.RemoveUser(selected)
		if err != nil {
			return err
		}

		cache.ClearUserDatabase(selected)
		output.Success("Removed user.", "label", selected)
		if active != nil {
			output.Info("Active user is now", "label", active.Label)
		} else {
			output.Warn("No active user remains. Add a new one with `aocli user add`.")
		}
		return nil
	},
}

var viewUserCmd = &cobra.Command{
	Use:   "view",
	Short: "Display the active user's Advent of Code progress.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		active, err := session.GetActiveUser(false)
		if err != nil {
			return err
		}

		if err := cache.StartupDBM(active.Label); err != nil {
			return err
		}
		defer cache.ShutdownDBM()

		user, err := resources.NewUser(active.Label, active.Token)
		if err != nil {
			return err
		}

		// user.LoadUser()
		user.Display()
		return nil
	},
}

func init() {
	userCmd.AddCommand(addUserCmd)
	userCmd.AddCommand(clearUserCmd)
	userCmd.AddCommand(listUserCmd)
	userCmd.AddCommand(switchUserCmd)
	userCmd.AddCommand(updateUserCmd)
	userCmd.AddCommand(removeUserCmd)
	userCmd.AddCommand(viewUserCmd)
	rootCmd.AddCommand(userCmd)
}
