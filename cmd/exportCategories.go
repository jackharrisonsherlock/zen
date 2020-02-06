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

// exportCategoriesCmd represents the exportCategories command
var exportCategoriesCmd = &cobra.Command{
	Use:   "export-categories",
	Short: "Export Categories",
	Long:  "Exports a list of Help Centre Categories",

	Run: func(cmd *cobra.Command, args []string) {
		exportCategories()
	},
}

func init() {
	rootCmd.AddCommand(exportCategoriesCmd)
}

func exportCategories() {

	// Create new Zendesk Client
	zd, err := zendesk.NewClient(subdomain, email, apikey)

	// Export Categories
	fmt.Println("\nRetrieving Categories from Zendesk...")
	categories := []zendesk.Categorie{}
	counter := 1
	options := &zendesk.ListOptions{Page: counter}
	page, err := zd.ListCategories(options)
	categories = append(categories, page...)
	if err != nil {
		log.Fatalf("Failed retrieving Categories from Zendesk: %v", err)
	}

	for len(page) > 0 {
		counter++
		options := &zendesk.ListOptions{Page: counter}
		page, err = zd.ListCategories(options)
		categories = append(categories, page...)
	}

	// Write to CSV
	fmt.Println("Writing Categories to CSV file...")
	err = WriteFileCategories(categories)
	if err != nil {
		log.Fatalf("Failed writing Categories to CSV: %v", err)
	}

	fmt.Println(color.GreenString("\nFinished!"))
}

func WriteFileCategories(categories []zendesk.Categorie) error {

	// Create a blank CSV file.
	fileName := fmt.Sprintf("zendesk_categories_export_%v.csv", time.Now().Unix())
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
	header = append(header, "created_at")
	header = append(header, "updated_at")

	// Commit the header.
	writer.Write(header)

	// Now loop through each  object and populate the CSV.
	for _, categorie := range categories {

		var id, name, description, url, createdAt, updatedAt string

		if categorie.ID != nil {
			id = fmt.Sprintf("%v", *categorie.ID)
		}
		if categorie.Name != nil {
			name = *categorie.Name
		}
		if categorie.Description != nil {
			description = *categorie.Description
		}
		if categorie.URL != nil {
			url = *categorie.URL
		}
		if categorie.CreatedAt != nil {
			createdAt = fmt.Sprintf("%v", *categorie.CreatedAt)
		}
		if categorie.UpdatedAt != nil {
			updatedAt = fmt.Sprintf("%v", *categorie.UpdatedAt)
		}

		var record []string
		record = append(record, id)
		record = append(record, name)
		record = append(record, description)
		record = append(record, url)
		record = append(record, createdAt)
		record = append(record, updatedAt)
		writer.Write(record)
	}

	writer.Flush()
	return err
}
