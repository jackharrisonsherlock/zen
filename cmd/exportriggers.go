// FIX ACTIONS & CONDITIONS UNMARSHALLING

package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	zendesk "github.com/MEDIGO/go-zendesk/zendesk"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// exportriggersCmd represents the exportriggers command
var exportriggersCmd = &cobra.Command{
	Use:   "export-triggers",
	Short: "Export Triggers",

	Run: func(cmd *cobra.Command, args []string) {
		exportTriggers()
	},
}

func init() {
	rootCmd.AddCommand(exportriggersCmd)
}

func exportTriggers() {

	// Create new Zendesk Client
	zd, err := zendesk.NewClient(subdomain, email, apikey)

	// Export Triggers
	fmt.Println("\nRetrieving Triggers from Zendesk...")
	triggers := []zendesk.Trigger{}
	counter := 1
	options := &zendesk.ListOptions{Page: counter}
	page, err := zd.ListTriggers(options)
	triggers = append(triggers, page...)
	if err != nil {
		log.Fatalf("Failed retrieving Triggers from Zendesk: %v", err)
	}

	for len(page) > 0 {
		counter++
		options := &zendesk.ListOptions{Page: counter}
		page, err = zd.ListTriggers(options)
		triggers = append(triggers, page...)
	}

	// Write to CSV
	fmt.Println("Writing Triggers to CSV file...")
	err = WriteFileTriggers(triggers)
	if err != nil {
		log.Fatalf("Failed writing Triggers to CSV: %v", err)
	}

	fmt.Println(color.GreenString("\nFinished!"))
}

func WriteFileTriggers(triggers []zendesk.Trigger) error {

	// Create a blank CSV file.
	fileName := fmt.Sprintf("zendesk_trigger_export_%v.csv", time.Now().Unix())
	file, err := os.Create(fmt.Sprintf("./%s", fileName))
	if err != nil {
		return err
	}

	// Ensure the file is closed at the end.
	defer file.Close()
	writer := csv.NewWriter(file)

	// Write the header line.
	var header []string
	header = append(header, "id")
	header = append(header, "title")
	header = append(header, "active")
	// header = append(header, "conditionsAll")
	// header = append(header, "conditionsAny")
	// header = append(header, "actions")
	header = append(header, "description")
	header = append(header, "created_at")
	header = append(header, "updated_at")

	// Commit the header.
	writer.Write(header)

	// Now loop through each  object and populate the CSV.
	for _, trigger := range triggers {
		// for _, condition := range trigger.Conditions {

		var id, title, active, description, createdAt, updatedAt string

		if trigger.ID != nil {
			id = fmt.Sprintf("%v", *trigger.ID)
		}
		if trigger.Title != nil {
			title = *trigger.Title
		}
		if trigger.Active != nil {
			active = fmt.Sprintf("%v", *trigger.Active)
		}
		// if condition.All != nil {
		// 	conditionsAll = fmt.Sprintf("%v", condition.All)
		// }
		// if condition.Any != nil {
		// 	conditionsAny = fmt.Sprintf("%v", condition.Any)
		// }
		// if trigger.Actions != nil {
		// 	actions = fmt.Sprintf("%v", *trigger.Actions)
		// }
		if trigger.Description != nil {
			description = *trigger.Description
		}
		if trigger.CreatedAt != nil {
			createdAt = fmt.Sprintf("%v", *trigger.CreatedAt)
		}
		if trigger.UpdatedAt != nil {
			updatedAt = fmt.Sprintf("%v", *trigger.UpdatedAt)
		}

		var record []string
		record = append(record, id)
		record = append(record, title)
		record = append(record, active)
		// record = append(record, conditionsAll)
		// record = append(record, conditionsAny)
		// record = append(record, actions)
		record = append(record, description)
		record = append(record, createdAt)
		record = append(record, updatedAt)
		writer.Write(record)
	}

	writer.Flush()
	return err
}
