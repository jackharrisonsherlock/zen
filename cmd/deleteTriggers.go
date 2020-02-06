package cmd

import (
	"fmt"
	"log"

	zendesk "github.com/MEDIGO/go-zendesk/zendesk"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// deleteTriggersCmd represents the deleteTriggers command
var deleteTriggersCmd = &cobra.Command{
	Use:   "delete-triggers",
	Short: "Delete Triggers",
	Long: fmt.Sprintf(`
This tool requires a CSV of Triggers IDs, no headers.

Example:
%s`, color.GreenString("zen delete-triggers -f FILENAME.csv")),

	Run: func(cmd *cobra.Command, args []string) {
		deleteTriggers()
	},
}

func init() {
	rootCmd.AddCommand(deleteTriggersCmd)

	// Flags
	deleteTriggersCmd.Flags().StringVarP(&FilePath, "Filename", "f", "", "The name of your file: filename.csv")
	deleteTriggersCmd.MarkFlagRequired("Filename")
}

func deleteTriggers() {

	// Create new Zendesk Client
	zd, err := zendesk.NewClient(subdomain, email, apikey)

	// Get IDs from CSV
	fmt.Println("\nReading CSV")
	ids, err := readCSV(FilePath)
	if err != nil {
		log.Fatalf("Failed to get IDs from file: %s", FilePath)
	}

	// Deleting Triggers
	fmt.Printf("%d Triggers to Delete.\n \n", len(ids))
	for _, id := range ids {
		fmt.Printf("Deleting: %v\n", id)
		zd.DeleteTrigger(id)
	}

	fmt.Println(color.GreenString("\nFinished!\n"))
}
