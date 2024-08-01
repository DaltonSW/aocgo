# Welcome to aocGo

`aocGo` is a two part project. 

The first encompasses the `go` module of the same name. This module contains a main package with some functions for obtaining input for a given day based on the file's directory structure. 
It also contains an `aocutils` sub-package, containing some helpful functions for common things you might need when solving AoC puzzles.

This module should be imported using `go get dalton.dog/aocgo`, and is intended to be used in code alone

The second is a CLI application called `aocli` that can be used to interact with the Advent of Code workflow without leaving your terminal. It is a separate `go` module

I recommend that you install the tool via the GitHub Releases page. There are standalone binaries there for Windows, Linux, and Mac, and the CLI program has a built in updating system. Put the executable in a directory on your PATH and you're good to go.

You are also able to install it via `go install dalton.dog/aocgo/cmd/aocli@latest`, if you'd prefer.
