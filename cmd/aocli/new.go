package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.dalton.dog/aocgo/internal/output"
)

var newCmd = &cobra.Command{
	Use:   "new [-b base.go]",
	Short: "Copies the given file into the ./<year>/<day>/main.go",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		New(Year, Day, BaseFilename, OutFilename)
	},
}

// New will make a copy of the given base file into the given out file in a ./yearIn/dayIn/ directory.
// Command `aocli new [-y yyyy -d dd -b base.go -o main.go]`
func New(yearIn, dayIn, baseFile, outFile string) {
	if yearIn == "0" {
		output.Fatal("Must provide a year with the -y option.")
	}

	if dayIn == "0" {
		output.Fatal("Must provide a day with the -d option.")
	}

	// Build the directory path, e.g. "./2025/15"
	dirPath := filepath.Join(".", yearIn, dayIn)

	// Create the directory (with parents, if needed).
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		output.Fatalf("Failed to create directory %s: %v", dirPath, err)
	}

	// Open the source/base file.
	src, err := os.Open(baseFile)
	if err != nil {
		output.Fatalf("Failed to open base file %s: %v", baseFile, err)
	}
	defer src.Close()

	// Create the destination file in the new directory.
	dstPath := filepath.Join(dirPath, outFile)
	dst, err := os.Create(dstPath)
	if err != nil {
		output.Fatalf("Failed to create output file %s: %v", dstPath, err)
	}
	defer dst.Close()

	// Copy the entire contents from baseFile to outFile.
	if _, err := io.Copy(dst, src); err != nil {
		output.Fatalf("Failed to copy data from %s to %s: %v", baseFile, dstPath, err)
	}

	output.Successf("Successfully copied %s to %s", baseFile, dstPath)

}
