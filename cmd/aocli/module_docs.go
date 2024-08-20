// TODO: Update all of this too

// Copyright 2024 Dalton Williams
// Use of this source code is governed by GPL v2 license,
// which can be found in the repository's LICENSE file.

// Please keep the source code available for anything distributed standalone
// Though consider contributing back to the repo!

/*
`Aocli` allows you to interface with Advent of Code workflows without leaving your terminal.
This command is *NOT* Go-specific and can be used with any language for solving.

This command is *NOT* intended to be installed or updated as part of the parent module.
Though part of the parent module, that is for central code sharing purposes. Aocli has self-updating capabilities that can/should be used to keep it up-to-date. Otherwise, you can find and install the latest version by checking the GitHub Releases

Usage:

	aocli <command> [flags...]

Available commands are:

	get --------- Get the user input for a given year and day and save it to a local file
	health ------ Checks to see if the system has valid configuration in place to successfully run the program
	help -------- Shows the help information for a specific command
	leaderboard - Shows the leaderboard for the given year, or given year and day
	reload ------ Refresh the page data for the puzzle on a given year and day
	submit ------ Submit a puzzle answer for a given year and day
	user -------- View the stars obtained for the current user
	view -------- Pretty print the puzzle for a given day

Any commands that rely on a year and day will attempt to derive those values from the names of the current directory and the parent directory. If those can't be properly derived, or you wish to run the command for another date, you can pass those in manually.

Run aocli help `<command>` for more information on a specific command.

# Get puzzle input for default session token

Usage:

	aocli get [year] [day]

This will automatically download the input into the current working directory to a file called 'input.txt'.

# Check that enough config is in place to run the program

Usage:

	aocli health

This will check that there is a valid session token in the environment to use for AoC requests.
The AOC_SESSION_TOKEN environment variable will be checked, as will the ~/.config/aocgo/session.token file.

# Display the leaderboard for a given year

Usage:

	aocli leaderboard <year> <day>

This will display the leaderboard for the given year, or a given year and day, in a scrollable table.

# View the puzzle information for a given day and year

Usage:

	aocli view [year] [day]

This will display the page information about a puzzle. This includes the title, the description, any correct answers for the current user, as well as successfully loading Part B once it's unlocked.

# Refresh the puzzle information for a given day and year

Usage:

	aocli refresh [year] [day]

# Submit an answer for a day and year, SOLELY determined by the CWD directory structure

Usage:

	aocli submit <answer>
*/
package main
