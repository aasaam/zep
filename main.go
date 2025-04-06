package main

import (
	"fmt"
	"os"
)

func main() {
	output, err := Run(os.Args, os.Environ())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, output)
}
