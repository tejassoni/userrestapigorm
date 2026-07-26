package main

import (
	"fmt"
	"os"

	"userrestapigorm/internal/bootstrap"
)

func main() {
	app, err := bootstrap.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "initialization failed: %v\n", err)
		os.Exit(1)
	}

	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n", err)
		os.Exit(1)
	}
}
