package main

import (
	"fmt"
	"os"

	"github.com/RubenDeBeer/Go_Do_Something/internal/adapter/cli"
	"github.com/RubenDeBeer/Go_Do_Something/internal/adapter/filedb"
	"github.com/RubenDeBeer/Go_Do_Something/internal/core"
)

func main() {
	// Determine data file path: env HEXACLI_FILE overrides default
	dataFile := os.Getenv("HEXACLI_FILE")
	if dataFile == "" {
		dataFile = "./data.txt"
	}

	repo := filedb.New(dataFile)
	svc := core.NewService(repo)

	// Pass through the same args to CLI runner
	exit := cli.Run(svc, os.Args[1:])
	if exit != 0 {
		os.Exit(exit)
	}

	// Optional success message is printed by CLI runner
	fmt.Print("")
}
