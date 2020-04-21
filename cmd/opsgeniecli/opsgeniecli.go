// +build go1.14

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.bus.zalan.do/SRE/adaptive-paging/pkg/opsgenie"
)

func main() {
	if err := run(); err != nil {
		printUsage()
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	opsgenieAPIToken := ""
	query := ""
	command := ""

	flag.StringVar(&query, "query", "status: open", "Opsgenie query to find alerts. e.g. 'status: open'")
	flag.StringVar(&opsgenieAPIToken, "api", "", "Provide Psgenie API token")
	flag.StringVar(&command, "cmd", "list", "Action to run on query. Available commands: list, delete. Defaul: list")
	flag.Parse()

	opsgeniesvc, err := opsgenie.New(opsgenieAPIToken)
	if err != nil {
		return err
	}

	switch command {
	case "list":
		alerts, err := opsgeniesvc.Query(query)
		if err != nil {
			return fmt.Errorf("could not get alerts with query '%s' : %w", query, err)
		}

		for _, alert := range alerts {
			fmt.Printf("- %s [%s]\n  alias: %s\n  created: %v\n  tags: %s\n", alert.Message, alert.Status, alert.Alias, alert.CreatedAt, strings.Join(alert.Tags, ","))
			notes, err := opsgeniesvc.AlertNotes(alert)
			if err != nil {
				fmt.Printf("  NOTES: Could not get: %v", err)
				continue
			}
			fmt.Println("  notes:")
			for _, note := range notes {
				fmt.Println("  - created_at: ", note.CreatedAt)
				fmt.Println("    note:       ", note.Note)
				fmt.Println()
			}
			fmt.Println("---")
			fmt.Println("")
		}
	case "delete":
		fmt.Println("WARNING: Remove records is not safe. Uncomment the code")
		// err := opsgeniesvc.CleanByQuery(query)
		// if err != nil {
		// 	return fmt.Errorf("could not delete alerts with query '%s' : %w", query, err)
		// }
	default:
		return errors.New("Unknown command: " + command)
	}

	return nil
}

func printUsage() {
	fmt.Println("Usage:")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println(`  opsgeniecli -api "xxxxx" -query "status: closed AND alias: yyy"`)
	fmt.Println()
}
