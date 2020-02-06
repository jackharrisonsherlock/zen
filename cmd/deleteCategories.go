package cmd

import (
	"fmt"
	"log"

	zendesk "github.com/MEDIGO/go-zendesk/zendesk"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// deleteCategoriesCmd represents the deleteCategories command
var deleteCategoriesCmd = &cobra.Command{
	Use:   "delete-categories",
	Short: "Delete Categories",
	Long: fmt.Sprintf(`Deletes the Help Centre Category. 
	
This tool requires a CSV of Category IDs, no headers.

Example:
%s`, color.GreenString("zen delete-categories -f FILENAME.csv")),

	Run: func(cmd *cobra.Command, args []string) {
		deleteCategories()
	},
}

func init() {
	rootCmd.AddCommand(deleteCategoriesCmd)

	// Flags
	deleteCategoriesCmd.Flags().StringVarP(&FilePath, "Filename", "f", "", "The name of your file: filename.csv")
	deleteCategoriesCmd.MarkFlagRequired("Filename")
}

func deleteCategories() {

	// Create new Zendesk Client
	zd, err := zendesk.NewClient(subdomain, email, apikey)

	// Get IDs from CSV
	fmt.Println("\nReading CSV")
	ids, err := readCSV(FilePath)
	if err != nil {
		log.Fatalf("Failed to get IDs from file: %s", FilePath)
	}

	// Deleting Categories
	fmt.Printf("%d Categories to Delete.\n \n", len(ids))
	for _, id := range ids {
		fmt.Printf("Deleting: %v\n", id)
		zd.DeleteCategorie(id)
	}

	fmt.Println(color.GreenString("\nFinished!\n"))
}
