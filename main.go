package main

import (
	"cleaner/internal"
	"os"
)

// Запускає програму
func main() {
	os.Exit(internal.Run(os.Args[1:], os.Stdout, os.Stderr))
}
