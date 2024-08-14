# Welcome to aocli

This is a command line tool that lets you interact with Advent of Code without ever leaving your terminal.

It is implemented with rate limiting and local caching to ensure we're not hitting the servers more often than necessary.

## Available Commands

Syntax Example: `aocli command_name <RequiredParam> [OptionalParam] -optionOne valueOne -x`

### `get`

Allows you to get the user input for a given year and day. Can be passed in as parameters. If not passed in, will attempt to be derived from the current directory.

Syntax: `aocli get [year] [day]`

[![aocli get demo](https://asciinema.org/a/lduYJUOBrHWqwe9UieBHX9hU4.svg)](https://asciinema.org/a/lduYJUOBrHWqwe9UieBHX9hU4?autoplay=1)

### `view`

Allows you to view the puzzle page for a given year and day. Can be passed in as parameters. If not passed in, will attempt to be derived from the current directory.

Syntax: `aocli view [year] [day]`

[![aocli view demo](https://asciinema.org/a/bq5KqnHaY8ozybzxTGIAWFM8Z.svg)](https://asciinema.org/a/bq5KqnHaY8ozybzxTGIAWFM8Z?autoplay=1)

### `leaderboard`

Allows you to view the leaderboard for a given year, or given year + day. Passed in as parameters.

Syntax: `aocli leaderboard <year> [day]`

[![aocli leaderboard demo](https://asciinema.org/a/misVkiiAbGsJb0xq1iq3WXhfk.svg)](https://asciinema.org/a/misVkiiAbGsJb0xq1iq3WXhfk?autoplay=1)

### `submit`

This exists! Need to document it!

### `check-update`

Will check the internal version against the latest GitHub repo release to see if there's a new version available.

`aocli check-update`

### `update`

Will attempt to download and install the newest version of `aocli` in-place, if there is one newer.

`aocli update`

### `reload`

Will reload the contents of the puzzle page for a given year and day. Can be passed in as parameters. If not passed in, will attempt to be derived from the current directory.

Syntax: `aocli reload [year] [day]`

### `clear-user`

Will clear the stored information for the user in the session token file, or AOC_SESSION_TOKEN environment variable.

`aocli clear-user`
