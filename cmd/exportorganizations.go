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

// exportorganizationsCmd represents the exportorganizations command
var exportorganizationsCmd = &cobra.Command{
	Use:   "export-organizations",
	Short: "Export Organizations",

	Run: func(cmd *cobra.Command, args []string) {
		exportOrganizations()
	},
}

func init() {
	rootCmd.AddCommand(exportorganizationsCmd)
}

func exportOrganizations() {

	// Create new Zendesk Client
	zd, err := zendesk.NewClient(subdomain, email, apikey)

	// Export Organizations
	fmt.Println("\nRetrieving Organizations from Zendesk...")
	organizations := []zendesk.Organization{}
	counter := 1
	options := &zendesk.ListOptions{Page: counter}
	page, err := zd.ListOrganizations(options)
	organizations = append(organizations, page...)
	if err != nil {
		log.Fatalf("Failed retrieving users from Zendesk: %v", err)
	}

	for len(page) > 0 {
		counter++
		options := &zendesk.ListOptions{Page: counter}
		page, err = zd.ListOrganizations(options)
		organizations = append(organizations, page...)
	}

	// Write to CSV
	fmt.Println("Writing organizations to CSV file...")
	err = WriteFileOrganizations(organizations)
	if err != nil {
		log.Fatalf("Failed writing organizations to CSV: %v", err)
	}

	fmt.Println(color.GreenString("\nFinished!"))
}

func WriteFileOrganizations(organizations []zendesk.Organization) error {

	// Create a blank CSV file.
	fileName := fmt.Sprintf("zendesk_organizations_export_%v.csv", time.Now().Unix())
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
	header = append(header, "url")
	header = append(header, "name")
	header = append(header, "domain_names")
	header = append(header, "details")
	header = append(header, "notes")

	// Commit the header.
	writer.Write(header)

	// Now loop through each  object and populate the CSV.
	for _, organization := range organizations {

		var id, url, name, domainNames, details, notes string

		if organization.ID != nil {
			id = fmt.Sprintf("%v", *organization.ID)
		}
		if organization.URL != nil {
			url = *organization.URL
		}
		if organization.Name != nil {
			name = *organization.Name
		}
		if organization.DomainNames != nil {
			domainNames = fmt.Sprintf("%v", *organization.DomainNames)
		}
		if organization.Details != nil {
			details = *organization.Details
		}
		if organization.Notes != nil {
			notes = *organization.Notes
		}

		var record []string
		record = append(record, id)
		record = append(record, url)
		record = append(record, name)
		record = append(record, domainNames)
		record = append(record, details)
		record = append(record, notes)
		writer.Write(record)
	}

	writer.Flush()
	return err
}
