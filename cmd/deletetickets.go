package cmd

import (
	"fmt"
	"log"

	zendesk "github.com/MEDIGO/go-zendesk/zendesk"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// deleteticketCmd represents the deleteticket command
var deleteticketCmd = &cobra.Command{
	Use:   "delete-tickets",
	Short: "Delete Tickets",
	Long: fmt.Sprintf(`
This tool requires a CSV of Ticket IDs, no headers.

Example:
%s`, color.GreenString("zen delete-tickets -f FILENAME.csv")),

	Run: func(cmd *cobra.Command, args []string) {
		deleteTickets()
	},
}

func init() {
	rootCmd.AddCommand(deleteticketCmd)

	// Flags
	deleteticketCmd.Flags().StringVarP(&FilePath, "Filename", "f", "", "The name of your file: filename.csv")
	deleteticketCmd.MarkFlagRequired("Filename")
}

func deleteTickets() {

	// Create new Zendesk Client
	zd, err := zendesk.NewClient(subdomain, email, apikey)

	// Get IDs from CSV
	fmt.Println("\nReading CSV")
	ids, err := readCSV(FilePath)
	if err != nil {
		log.Fatalf("Failed to get IDs from file: %s", FilePath)
	}

	// Deleting Tickets
	fmt.Printf("%d Tickets to Delete.\n \n", len(ids))
	for _, id := range ids {
		fmt.Printf("Deleting: %v\n", id)
		zd.DeleteTicket(id)
	}

	fmt.Println(color.GreenString("\nFinished!\n"))
}
