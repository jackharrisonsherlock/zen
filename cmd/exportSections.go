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

// exportSectionsCmd represents the exportSections command
var exportSectionsCmd = &cobra.Command{
	Use:   "export-sections",
	Short: "Export Sections",

	Run: func(cmd *cobra.Command, args []string) {
		exportSections()
	},
}

func init() {
	rootCmd.AddCommand(exportSectionsCmd)
}

func exportSections() {

	// Create new Zendesk Client
	zd, err := zendesk.NewClient(subdomain, email, apikey)

	// Export Sections
	fmt.Println("\nRetrieving Sections from Zendesk...")
	sections := []zendesk.Section{}
	counter := 1
	options := &zendesk.ListOptions{Page: counter}
	page, err := zd.ListSections(options)
	sections = append(sections, page...)
	if err != nil {
		log.Fatalf("Failed retrieving Sections from Zendesk: %v", err)
	}

	for len(page) > 0 {
		counter++
		options := &zendesk.ListOptions{Page: counter}
		page, err = zd.ListSections(options)
		sections = append(sections, page...)
	}

	// Write to CSV
	fmt.Println("Writing Sections to CSV file...")
	err = WriteFileSections(sections)
	if err != nil {
		log.Fatalf("Failed writing Sections to CSV: %v", err)
	}

	fmt.Println(color.GreenString("\nFinished!"))
}

func WriteFileSections(sections []zendesk.Section) error {

	// Create a blank CSV file.
	fileName := fmt.Sprintf("zendesk_sections_export_%v.csv", time.Now().Unix())
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
	header = append(header, "name")
	header = append(header, "description")
	header = append(header, "url")
	header = append(header, "category_id")
	header = append(header, "parent_section_id")
	header = append(header, "created_at")
	header = append(header, "updated_at")

	// Commit the header.
	writer.Write(header)

	// Now loop through each  object and populate the CSV.
	for _, section := range sections {

		var id, name, description, url, categoryID, parentSectionID, createdAt, updatedAt string

		if section.ID != nil {
			id = fmt.Sprintf("%v", *section.ID)
		}
		if section.Name != nil {
			name = *section.Name
		}
		if section.Description != nil {
			description = *section.Description
		}
		if section.URL != nil {
			url = *section.URL
		}
		if section.CategoryID != nil {
			categoryID = fmt.Sprintf("%v", *section.CategoryID)
		}
		if section.ParentSectionID != nil {
			parentSectionID = fmt.Sprintf("%v", *section.ParentSectionID)
		}
		if section.CreatedAt != nil {
			createdAt = fmt.Sprintf("%v", *section.CreatedAt)
		}
		if section.UpdatedAt != nil {
			updatedAt = fmt.Sprintf("%v", *section.UpdatedAt)
		}

		var record []string
		record = append(record, id)
		record = append(record, name)
		record = append(record, description)
		record = append(record, url)
		record = append(record, categoryID)
		record = append(record, parentSectionID)
		record = append(record, createdAt)
		record = append(record, updatedAt)
		writer.Write(record)
	}

	writer.Flush()
	return err
}
