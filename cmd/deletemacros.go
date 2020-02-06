package cmd

import (
	"fmt"
	"log"

	zendesk "github.com/MEDIGO/go-zendesk/zendesk"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// deletemacrosCmd represents the deletemacros command
var deletemacrosCmd = &cobra.Command{
	Use:   "delete-macros",
	Short: "Delete Macros",
	Long: fmt.Sprintf(`
This tool requires a CSV of Macro IDs, no headers.

Example:
%s`, color.GreenString("zen delete-macros -f FILENAME.csv")),

	Run: func(cmd *cobra.Command, args []string) {
		deleteMacros()
	},
}

func init() {
	rootCmd.AddCommand(deletemacrosCmd)

	// Flags
	deletemacrosCmd.Flags().StringVarP(&FilePath, "Filename", "f", "", "The name of your file: filename.csv")
	deletemacrosCmd.MarkFlagRequired("Filename")
}

func deleteMacros() {

	// Create new Zendesk Client
	zd, err := zendesk.NewClient(subdomain, email, apikey)

	// Get IDs from CSV
	fmt.Println("\nReading CSV")
	ids, err := readCSV(FilePath)
	if err != nil {
		log.Fatalf("Failed to get IDs from file: %s", FilePath)
	}

	// Deleting Macros
	fmt.Printf("%d Macros to Delete.\n \n", len(ids))
	for _, id := range ids {
		fmt.Printf("Deleting: %v \n", id)
		zd.DeleteMacro(id)
	}

	fmt.Println(color.GreenString("\nFinished!\n"))
}
