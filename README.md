# Welcome to aocGo

`aocGo` is a two part project. 

## `aocgo`

The first encompasses the `go` module of the same name. This module contains a main package with some functions for obtaining input for a given day based on the file's directory structure.  

It also contains an `aocutils` sub-package, containing some helpful functions for common things you might need when solving AoC puzzles.

This module should be imported using `go get dalton.dog/aocgo`, and is intended to be used in code alone.

## `aocli`

The second is a CLI application called `aocli` that can be used to interact with the Advent of Code workflow without leaving your terminal.

I recommend that you install the tool via the GitHub Releases page. There are standalone binaries there for Windows, Linux, and Mac, and the CLI program has a built in updating system. Put the executable in a directory on your PATH and you're good to go.

## Required Setup

All of the functionality here requires a user session token to be available. It should be placed in `~/.config/aocgo/session.token`, or stored in the `AOC_SESSION_TOKEN` environment variable.

1. To obtain this token, log in to Advent of Code as the user you'd like to make requests and submissions as.  
2. Open the Dev Console (Ctrl + Shift + I, F12, or Right Click -> Inspect), then go to the Network tab.
3. Navigate to any puzzle's input page and inspect the GET request headers. You should see `Cookie: session=<your session token>`.
4. Place everything after the equals sign (ignore `Cookie: session=`) in the file or the environment variable.
5. Run `aocli health` to verify that everthing is loaded properly.

