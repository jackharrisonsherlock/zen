package cmd

import (
	"fmt"
	"log"

	zendesk "github.com/MEDIGO/go-zendesk/zendesk"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// deleteviewsCmd represents the deleteviews command
var deleteviewsCmd = &cobra.Command{
	Use:   "delete-views",
	Short: "Delete Views",
	Long: fmt.Sprintf(`
This tool requires a CSV of View IDs, no headers.

Example:
%s`, color.GreenString("zen delete-views -f FILENAME.csv")),

	Run: func(cmd *cobra.Command, args []string) {
		deleteViews()
	},
}

func init() {
	rootCmd.AddCommand(deleteviewsCmd)

	// Flags
	deleteviewsCmd.Flags().StringVarP(&FilePath, "Filename", "f", "", "The name of your file: filename.csv")
	deleteviewsCmd.MarkFlagRequired("Filename")
}

func deleteViews() {

	// Create new Zendesk Client
	zd, err := zendesk.NewClient(subdomain, email, apikey)

	// Get IDs from CSV
	fmt.Println("\nReading CSV")
	ids, err := readCSV(FilePath)
	if err != nil {
		log.Fatalf("Failed to get IDs from file: %s", FilePath)
	}

	// Deleting Views
	fmt.Printf("%d Views to Delete.\n \n", len(ids))
	for _, id := range ids {
		fmt.Printf("Deleting: %v\n", id)
		zd.DeleteView(id)
	}

	fmt.Println(color.GreenString("\nFinished!\n"))
}
