package main

import (
	"fmt"
	"os"

	"github.com/taka1156/bastion/internal/app"
)

func main() {
	if err := app.NewApp().Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
