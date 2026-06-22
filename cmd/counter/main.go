package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/zandgaim/anagram-counter/internal/facade"
)

func main() {
	// CLI flags
	dirPtr := flag.String("dir", "", "Path to the directory containing text files")
	workersPtr := flag.Int("workers", runtime.NumCPU(), "Number of concurrent workers")
	flag.Parse()

	if *dirPtr == "" {
		fmt.Println("Error: --dir flag is required.")
		fmt.Println("Usage: ./counter --dir=/path/to/files")
		os.Exit(1)
	}

	if *workersPtr <= 0 {
		fmt.Println("Error: --workers must be greater than 0")
		os.Exit(1)
	}

	if err := facade.Run(*dirPtr, *workersPtr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
