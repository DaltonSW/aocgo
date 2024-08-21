# Welcome to `aocGo`!

`aocGo` is a two part project. 

## `aocgo`

The first encompasses the `go` module of the same name. This module contains a main package with some functions for obtaining input for a given day based on the file's directory structure. It also contains an `aocutils` sub-package, containing some helpful functions for common things you might need when solving AoC puzzles. If you want to handle all of that yourself, don't worry about importing this.

This module should be imported using `go get dalton.dog/aocgo`.

Example:
```go
// Example file/dir structure:
// 2015/1/main.go
package main

import "dalton.dog/aocgo"

func main() {
    // Get your input data as a string, exactly as the input has it
    var inLine string = aocgo.GetInputAsString()

    // Get your input data as an array of strings, created by splitting input on newline characters
    var inLnArr []string = aocgo.GetInputAsLineArray()

    // Get your input data as an array of bytes, where each element is a single byte of the input
    var inByteArr []byte = aocgo.GetInputAsByteArray()
}
```

## `aocli`

The second is a CLI application called `aocli` that can be used to interact with the Advent of Code workflow without leaving your terminal.

It can currently do all of the following:
- View puzzle data
- Submit puzzle answers
- View both yearly and daily leaderboards
- View an overview of your user

Check out the directory's specific [README](https://github.com/DaltonSW/aocgo/tree/main/cmd/aocli) for detailed documentation!

You should install the program via this repo's ['Releases' page](https://github.com/DaltonSW/aocgo/releases/latest). There are standalone binaries for Windows, Linux, and Mac, and the program has a built-in updating system. Just put the executable in a directory on your PATH and you're good to go.

## Required Setup

1. To obtain this token, log in to Advent of Code as the user you'd like to make requests and submissions as.  
2. Open the Dev Console (Ctrl + Shift + I, F12, or Right Click -> Inspect), then go to the Network tab.
3. Navigate to any puzzle's input page and inspect the GET request headers. You should see `Cookie: session=<your session token>`.
4. Place everything after the equals sign (ignore `Cookie: session=`) in the file or the environment variable.
5. Run `aocli health` to verify that everthing is loaded properly.

