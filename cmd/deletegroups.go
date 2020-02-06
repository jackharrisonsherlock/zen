package cmd

import (
	"fmt"
	"log"

	zendesk "github.com/MEDIGO/go-zendesk/zendesk"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// deletegroupsCmd represents the deletegroups command
var deletegroupsCmd = &cobra.Command{
	Use:   "delete-groups",
	Short: "Delete Groups",
	Long: fmt.Sprintf(`
This tool requires a CSV of Group IDs, no headers.

Example:
%s`, color.GreenString("zen delete-groups -f FILENAME.csv")),

	Run: func(cmd *cobra.Command, args []string) {
		deleteGroups()
	},
}

func init() {
	rootCmd.AddCommand(deletegroupsCmd)

	// Flags
	deletegroupsCmd.Flags().StringVarP(&FilePath, "Filename", "f", "", "The name of your file: filename.csv")
	deletegroupsCmd.MarkFlagRequired("Filename")
}

func deleteGroups() {

	// Create new Zendesk Client
	zd, err := zendesk.NewClient(subdomain, email, apikey)

	// Get IDs from CSV
	fmt.Println("\nReading CSV")
	ids, err := readCSV(FilePath)
	if err != nil {
		log.Fatalf("Failed to get IDs from file: %s", FilePath)
	}

	// Deleting Groups
	fmt.Printf("%d Groups to Delete.\n \n", len(ids))
	for _, id := range ids {
		fmt.Printf("Deleting: %v\n", id)
		zd.DeleteGroup(id)
	}

	fmt.Println(color.GreenString("\nFinished!\n"))
}
