package cmd

import (
	"fmt"
	"log"

	zendesk "github.com/MEDIGO/go-zendesk/zendesk"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// deletearticlesCmd represents the deleteticket command
var deletearticleCmd = &cobra.Command{
	Use:   "delete-articles",
	Short: "Delete Articles",
	Long: fmt.Sprintf(`Archives the article. You can restore the article using the Help Center user interface.
	
This tool requires a CSV of Article IDs, no headers.

Example:
%s`, color.GreenString("zen delete-articles -f FILENAME.csv")),

	Run: func(cmd *cobra.Command, args []string) {
		deleteArticles()
	},
}

func init() {
	rootCmd.AddCommand(deletearticleCmd)

	// Flags
	deletearticleCmd.Flags().StringVarP(&FilePath, "Filename", "f", "", "The name of your file: filename.csv")
	deletearticleCmd.MarkFlagRequired("Filename")
}

func deleteArticles() {

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
		zd.DeleteArticle(id)
	}

	fmt.Println(color.GreenString("\nFinished!\n"))
}
