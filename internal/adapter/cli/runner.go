package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/RubenDeBeer/Go_Do_Something/internal/core"
)

// Run parses args and invokes core use cases.
// Supported commands:
//
//	add <value>
//	list
//	delete-last
//
// Global flags:
//
//	-file <path> (defaults to ./data.txt)
func Run(svc *core.Service, args []string) int {
	fs := flag.NewFlagSet("hexacli", flag.ContinueOnError)
	//dataFile := fs.String("file", "./data.txt", "path to data file")
	fs.SetOutput(new(strings.Builder)) // silence default error output

	// We parse flags up to the command boundary; the remaining arg is the command.
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		fmt.Fprintln(os.Stderr, usage())
		return 2
	}

	// Configure is by side-effect of -file flag; the repository was constructed in main with same path.
	leftovers := fs.Args()
	if len(leftovers) == 0 {
		fmt.Fprintln(os.Stderr, usage())
		return 2
	}

	cmd := leftovers[0]
	switch cmd {
	case "add":
		if len(leftovers) < 2 {
			fmt.Fprintln(os.Stderr, "add requires a value")
			return 2
		}
		value := strings.Join(leftovers[1:], " ")
		if err := svc.AddValue(value); err != nil {
			return die(err)
		}
		fmt.Println("OK")
		return 0
	case "list":
		items, err := svc.ListValues()
		if err != nil {
			return die(err)
		}
		for i, v := range items {
			fmt.Printf("%d: %s\n", i+1, v)
		}
		return 0
	case "delete-last":
		if err := svc.DeleteLast(); err != nil {
			return die(err)
		}
		fmt.Println("Deleted last entry")
		return 0
	default:
		return die(errors.New("unknown command: " + cmd))
	}
}

func die(err error) int {
	fmt.Fprintln(os.Stderr, "error:", err)
	return 1
}

func usage() string {
	return `Usage:
  hexacli [-file ./data.txt] add <value>
  hexacli [-file ./data.txt] list
  hexacli [-file ./data.txt] delete-last`
}
