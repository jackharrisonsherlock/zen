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

// exportgroupsCmd represents the exportgroups command
var exportgroupsCmd = &cobra.Command{
	Use:   "export-groups",
	Short: "Export Groups",

	Run: func(cmd *cobra.Command, args []string) {
		exportGroups()
	},
}

func init() {
	rootCmd.AddCommand(exportgroupsCmd)
}

func exportGroups() {

	// Create new Zendesk Client
	zd, err := zendesk.NewClient(subdomain, email, apikey)

	// Export groups
	fmt.Println("\nRetrieving groups from Zendesk...")
	groups := []zendesk.Group{}
	counter := 1
	options := &zendesk.ListOptions{Page: counter}
	page, err := zd.ListGroups(options)
	groups = append(groups, page...)
	if err != nil {
		log.Fatalf("Failed retrieving groups from Zendesk: %v", err)
	}

	for len(page) > 0 {
		counter++
		options := &zendesk.ListOptions{Page: counter}
		page, err = zd.ListGroups(options)
		groups = append(groups, page...)
	}

	// Write to CSV
	fmt.Println("Writing groups to CSV file...")
	err = WriteFileGroups(groups)
	if err != nil {
		log.Fatalf("Failed writing groups to CSV: %v", err)
	}

	fmt.Println(color.GreenString("\nFinished!"))
}

func WriteFileGroups(groups []zendesk.Group) error {

	// Create a blank CSV file.
	fileName := fmt.Sprintf("zendesk_groups_export_%v.csv", time.Now().Unix())
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
	header = append(header, "name")
	header = append(header, "deleted")
	header = append(header, "created_at")
	header = append(header, "updated_at")

	// Commit the header.
	writer.Write(header)

	// Now loop through each  object and populate the CSV.
	for _, group := range groups {

		var url, id, name, deleted, createdAt, updatedAt string

		if group.URL != nil {
			url = fmt.Sprintf("%v", *group.URL)
		}
		if group.ID != nil {
			id = fmt.Sprintf("%v", *group.ID)
		}
		if group.Name != nil {
			name = *group.Name
		}
		if group.Deleted != nil {
			deleted = fmt.Sprintf("%v", *group.Deleted)
		}
		if group.CreatedAt != nil {
			createdAt = fmt.Sprintf("%v", *group.CreatedAt)
		}
		if group.UpdatedAt != nil {
			updatedAt = fmt.Sprintf("%v", *group.UpdatedAt)
		}

		var record []string
		record = append(record, url)
		record = append(record, id)
		record = append(record, name)
		record = append(record, deleted)
		record = append(record, createdAt)
		record = append(record, updatedAt)
		writer.Write(record)
	}

	writer.Flush()
	return err
}
