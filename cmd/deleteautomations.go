package cmd

import (
	"fmt"
	"log"

	zendesk "github.com/MEDIGO/go-zendesk/zendesk"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// deleteautomationsCmd represents the deleteautomations command
var deleteautomationsCmd = &cobra.Command{
	Use:   "delete-automations",
	Short: "Delete Automations",
	Long: fmt.Sprintf(`
This tool requires a CSV of Automation IDs, no headers.

Example:
%s`, color.GreenString("zen delete-automations -f FILENAME.csv")),

	Run: func(cmd *cobra.Command, args []string) {
		deleteAutomations()
	},
}

func init() {
	rootCmd.AddCommand(deleteautomationsCmd)

	// Flags
	deleteautomationsCmd.Flags().StringVarP(&FilePath, "Filename", "f", "", "The name of your file: filename.csv")
	deleteautomationsCmd.MarkFlagRequired("Filename")
}

func deleteAutomations() {

	// Create new Zendesk Client
	zd, err := zendesk.NewClient(subdomain, email, apikey)

	// Get IDs from CSV
	fmt.Println("\nReading CSV")
	ids, err := readCSV(FilePath)
	if err != nil {
		log.Fatalf("Failed to get IDs from file: %s", FilePath)
	}

	// Deleting Automations
	fmt.Printf("%d Automations to Delete.\n \n", len(ids))
	for _, id := range ids {
		fmt.Printf("Deleting: %v\n", id)
		zd.DeleteAutomation(id)
	}

	fmt.Println(color.GreenString("\nFinished!\n"))
}
