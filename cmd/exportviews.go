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

// exportviewsCmd represents the exportviews command
var exportviewsCmd = &cobra.Command{
	Use:   "export-views",
	Short: "Export Views",

	Run: func(cmd *cobra.Command, args []string) {
		exportViews()
	},
}

func init() {
	rootCmd.AddCommand(exportviewsCmd)
}

func exportViews() {

	// Create new Zendesk Client
	zd, err := zendesk.NewClient(subdomain, email, apikey)

	// Export Views
	fmt.Println("\nRetrieving Views from Zendesk...")
	views := []zendesk.View{}
	counter := 1
	options := &zendesk.ListOptions{Page: counter}
	page, err := zd.ListViews(options)
	views = append(views, page...)
	if err != nil {
		log.Fatalf("Failed retrieving views from Zendesk: %v", err)
	}

	for len(page) > 0 {
		counter++
		options := &zendesk.ListOptions{Page: counter}
		page, err = zd.ListViews(options)
		views = append(views, page...)
	}

	// Write to CSV
	fmt.Println("Writing views to CSV file...")
	err = WriteFileViews(views)
	if err != nil {
		log.Fatalf("Failed writing views to CSV: %v", err)
	}

	fmt.Println(color.GreenString("\nFinished!"))
}

func WriteFileViews(views []zendesk.View) error {

	// Create a blank CSV file.
	fileName := fmt.Sprintf("zendesk_view_export_%v.csv", time.Now().Unix())
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
	header = append(header, "description")
	header = append(header, "active")
	header = append(header, "created_at")
	header = append(header, "updated_at")
	// header = append(header, "execution")
	// header = append(header, "conditions")
	// header = append(header, "restriction")

	// Commit the header.
	writer.Write(header)

	// Now loop through each object and populate the CSV.
	for _, view := range views {

		var url, id, title, description, active, createdAt, updatedAt string

		if view.URL != nil {
			url = *view.URL
		}
		if view.ID != nil {
			id = fmt.Sprintf("%v", *view.ID)
		}
		if view.Title != nil {
			title = *view.Title
		}
		if view.Description != nil {
			description = *view.Description
		}
		if view.Active != nil {
			active = fmt.Sprintf("%v", *view.Active)
		}
		if view.CreatedAt != nil {
			createdAt = fmt.Sprintf("%v", *view.CreatedAt)
		}
		if view.UpdatedAt != nil {
			updatedAt = fmt.Sprintf("%v", *view.UpdatedAt)
		}

		// if view.Execution != nil {
		// 	execution = fmt.Sprintf("%v", *view.Execution)
		// }
		// if view.Conditions != nil {
		// 	conditions = fmt.Sprintf("%v", *view.Conditions)
		// }
		// if view.Restriction != nil {
		// 	restriction = fmt.Sprintf("%v", *view.Restriction)
		// }

		var record []string
		record = append(record, url)
		record = append(record, id)
		record = append(record, title)
		record = append(record, description)
		record = append(record, active)
		record = append(record, createdAt)
		record = append(record, updatedAt)
		// record = append(record, execution)
		// record = append(record, conditions)
		// record = append(record, restriction)
		writer.Write(record)
	}

	writer.Flush()
	return err
}
