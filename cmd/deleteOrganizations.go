package cmd

import (
	"fmt"
	"log"

	zendesk "github.com/MEDIGO/go-zendesk/zendesk"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// deleteOrganizationsCmd represents the deleteOrganizations command
var deleteOrganizationsCmd = &cobra.Command{
	Use:   "delete-organizations",
	Short: "Delete Organizations",
	Long: fmt.Sprintf(`
This tool requires a CSV of Organizations IDs, no headers.

Example:
%s`, color.GreenString("zen delete-organizations -f FILENAME.csv")),

	Run: func(cmd *cobra.Command, args []string) {
		deleteOrganizations()
	},
}

func init() {
	rootCmd.AddCommand(deleteOrganizationsCmd)

	// Flags
	deleteOrganizationsCmd.Flags().StringVarP(&FilePath, "Filename", "f", "", "The name of your file: filename.csv")
	deleteOrganizationsCmd.MarkFlagRequired("Filename")
}

func deleteOrganizations() {

	// Create new Zendesk Client
	zd, err := zendesk.NewClient(subdomain, email, apikey)

	// Get IDs from CSV
	fmt.Println("\nReading CSV")
	ids, err := readCSV(FilePath)
	if err != nil {
		log.Fatalf("Failed to get IDs from file: %s", FilePath)
	}

	// Deleting Organizations
	fmt.Printf("%d Organizations to Delete.\n \n", len(ids))
	for _, id := range ids {
		fmt.Printf("Deleting: %v\n", id)
		zd.DeleteOrganization(id)
	}

	fmt.Println(color.GreenString("\nFinished!\n"))
}
