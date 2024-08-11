# Welcome to aocli

This is a command line tool that lets you interact with Advent of Code without ever leaving your terminal.

It is implemented with rate limiting and local caching to ensure we're not hitting the servers more often than necessary.

## Available Commands

### `get`

Allows you to get the user input for a given year and day. Either passed in as parameters, or attempted to be derived from the current directory.

[![aocli get demo](https://asciinema.org/a/lduYJUOBrHWqwe9UieBHX9hU4.svg)](https://asciinema.org/a/lduYJUOBrHWqwe9UieBHX9hU4?autoplay=1)

### `view`

Allows you to view the puzzle page for a given year and day. Either passed in as parameters, or attempted to be derived from the current directory.

[![aocli view demo](https://asciinema.org/a/bq5KqnHaY8ozybzxTGIAWFM8Z.svg)](https://asciinema.org/a/bq5KqnHaY8ozybzxTGIAWFM8Z?autoplay=1)
