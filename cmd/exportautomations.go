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

// exportautomationsCmd represents the exportautomations command
var exportautomationsCmd = &cobra.Command{
	Use:   "export-automations",
	Short: "Export Automations",

	Run: func(cmd *cobra.Command, args []string) {
		exportAutomations()
	},
}

func init() {
	rootCmd.AddCommand(exportautomationsCmd)
}

func exportAutomations() {

	// Create new Zendesk Client
	zd, err := zendesk.NewClient(subdomain, email, apikey)

	// Export Automations
	fmt.Println("\nRetrieving Automations from Zendesk...")
	automations := []zendesk.Automation{}
	counter := 1
	options := &zendesk.ListOptions{Page: counter}
	page, err := zd.ListAutomations(options)
	automations = append(automations, page...)
	if err != nil {
		log.Fatalf("Failed retrieving Automations from Zendesk: %v", err)
	}

	for len(page) > 0 {
		counter++
		options := &zendesk.ListOptions{Page: counter}
		page, err = zd.ListAutomations(options)
		automations = append(automations, page...)
	}

	// Write to CSV
	fmt.Println("Writing Automations to CSV file...")
	err = WriteFileAutomations(automations)
	if err != nil {
		log.Fatalf("Failed writing automations to CSV: %v", err)
	}

	fmt.Println(color.GreenString("\nFinished!"))
}

func WriteFileAutomations(automations []zendesk.Automation) error {

	// Create a blank CSV file.
	fileName := fmt.Sprintf("zendesk_automation_export_%v.csv", time.Now().Unix())
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
	// header = append(header, "actions")
	// header = append(header, "conditions")
	header = append(header, "created_at")
	header = append(header, "updated_at")

	// Commit the header.
	writer.Write(header)

	// Now loop through each  object and populate the CSV.
	for _, automation := range automations {

		var id, title, active, createdAt, updatedAt string

		if automation.ID != nil {
			id = fmt.Sprintf("%v", *automation.ID)
		}
		if automation.Title != nil {
			title = *automation.Title
		}
		if automation.Active != nil {
			active = fmt.Sprintf("%v", *automation.Active)
		}
		// if automation.Conditions != nil {
		// 	conditions = fmt.Sprintf("%v", *automation.Conditions)
		// }
		// if automation.Actions != nil {
		// 	actions = fmt.Sprintf("%v", *automation.Actions)
		// }
		if automation.CreatedAt != nil {
			createdAt = fmt.Sprintf("%v", *automation.CreatedAt)
		}
		if automation.UpdatedAt != nil {
			updatedAt = fmt.Sprintf("%v", *automation.UpdatedAt)
		}

		var record []string
		record = append(record, id)
		record = append(record, title)
		record = append(record, active)
		// record = append(record, conditions)
		// record = append(record, actions)
		record = append(record, createdAt)
		record = append(record, updatedAt)
		writer.Write(record)
	}

	writer.Flush()
	return err
}
