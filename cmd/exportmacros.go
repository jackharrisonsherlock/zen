//**TODO: EXPORT ACTIONS

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

// exportmacrosCmd represents the exportmacros command
var exportmacrosCmd = &cobra.Command{
	Use:   "export-macros",
	Short: "Export Macros",

	Run: func(cmd *cobra.Command, args []string) {
		exportMacros()
	},
}

func init() {
	rootCmd.AddCommand(exportmacrosCmd)
}

func exportMacros() {

	// Create new Zendesk Client
	zd, err := zendesk.NewClient(subdomain, email, apikey)

	// Export Macros
	fmt.Println("\nRetrieving macros from Zendesk...")
	macros := []zendesk.Macro{}
	counter := 1
	options := &zendesk.ListOptions{Page: counter}
	page, err := zd.ListMacros(options)
	macros = append(macros, page...)

	for len(page) > 0 {
		counter++
		options := &zendesk.ListOptions{Page: counter}
		page, err = zd.ListMacros(options)
		macros = append(macros, page...)
	}

	// Write to CSV
	fmt.Println("Writing macros to CSV file...")
	err = WriteFileMacros(macros)
	if err != nil {
		log.Fatalf("Failed writing macros to CSV: %v", err)
	}

	fmt.Println(color.GreenString("\nFinished!"))
}

func WriteFileMacros(macros []zendesk.Macro) error {

	// Create a blank CSV file.
	fileName := fmt.Sprintf("zendesk_macros_export_%v.csv", time.Now().Unix())
	file, err := os.Create(fmt.Sprintf("./%s", fileName))
	if err != nil {
		return err
	}

	// Ensure the file is closed at the end.
	defer file.Close()
	writer := csv.NewWriter(file)

	// Write the header line.
	var header []string
	header = append(header, "url")
	header = append(header, "id")
	header = append(header, "title")
	header = append(header, "active")
	header = append(header, "updated_at")
	header = append(header, "created_at")
	header = append(header, "position")
	header = append(header, "description")
	header = append(header, "actions")
	header = append(header, "restriction")

	// Commit the header.
	writer.Write(header)
	for _, macro := range macros {

		var url, id, title, active, updated_at, created_at, position, description, actions, restriction string

		if macro.URL != nil {
			url = *macro.URL
		}
		if macro.ID != nil {
			id = fmt.Sprintf("%v", *macro.ID)
		}
		if macro.Title != nil {
			title = *macro.Title
		}
		if macro.Active != nil {
			active = fmt.Sprintf("%v", *macro.Active)
		}
		if macro.UpdatedAt != nil {
			updated_at = fmt.Sprintf("%v", *macro.UpdatedAt)
		}
		if macro.CreatedAt != nil {
			created_at = fmt.Sprintf("%v", *macro.CreatedAt)
		}
		if macro.Position != nil {
			position = fmt.Sprintf("%v", *macro.Position)
		}
		if macro.Description != nil {
			description = fmt.Sprintf("%v", *macro.Description)
		}
		// if macro.Actions != nil {
		// 	actions = *macro.Actions
		// }
		if macro.Restriction != nil {
			restriction = fmt.Sprintf("%v", *macro.Restriction)
		}

		var record []string
		record = append(record, url)
		record = append(record, id)
		record = append(record, title)
		record = append(record, active)
		record = append(record, updated_at)
		record = append(record, created_at)
		record = append(record, position)
		record = append(record, description)
		record = append(record, actions)
		record = append(record, restriction)
		writer.Write(record)
	}

	writer.Flush()
	return err
}
