package cmd

import (
	"fmt"
	"log"

	zendesk "github.com/MEDIGO/go-zendesk/zendesk"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// deleteSectionCmd represents the deleteSection command
var deleteSectionCmd = &cobra.Command{
	Use:   "delete-sections",
	Short: "Delete Sections",
	Long: fmt.Sprintf(`Deletes the Help Centre Section.

This tool requires a CSV of Sections IDs, no headers.

Example:
%s`, color.GreenString("zen delete-sections -f FILENAME.csv")),

	Run: func(cmd *cobra.Command, args []string) {
		deleteSections()
	},
}

func init() {
	rootCmd.AddCommand(deleteSectionCmd)

	// Flags
	deleteSectionCmd.Flags().StringVarP(&FilePath, "Filename", "f", "", "The name of your file: filename.csv")
	deleteSectionCmd.MarkFlagRequired("Filename")
}

func deleteSections() {

	// Create new Zendesk Client
	zd, err := zendesk.NewClient(subdomain, email, apikey)

	// Get IDs from CSV
	fmt.Println("\nReading CSV")
	ids, err := readCSV(FilePath)
	if err != nil {
		log.Fatalf("Failed to get IDs from file: %s", FilePath)
	}

	// Deleting Sections
	fmt.Printf("%d Sections to Delete.\n \n", len(ids))
	for _, id := range ids {
		fmt.Printf("Deleting: %v\n", id)
		zd.DeleteSection(id)
	}

	fmt.Println(color.GreenString("\nFinished!\n"))
}
