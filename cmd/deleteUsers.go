package cmd

import (
	"fmt"
	"log"

	zendesk "github.com/MEDIGO/go-zendesk/zendesk"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// deleteUsersCmd represents the deleteUsers command
var deleteUsersCmd = &cobra.Command{
	Use:   "delete-users",
	Short: "Delete Users",
	Long: fmt.Sprintf(`
This tool requires a CSV of User IDs, no headers.

Example:
%s`, color.GreenString("zen delete-users -f FILENAME.csv")),

	Run: func(cmd *cobra.Command, args []string) {
		deleteUsers()
	},
}

func init() {
	rootCmd.AddCommand(deleteUsersCmd)

	// Flags
	deleteUsersCmd.Flags().StringVarP(&FilePath, "Filename", "f", "", "The name of your file: filename.csv")
	deleteUsersCmd.MarkFlagRequired("Filename")
}

func deleteUsers() {

	// Create new Zendesk Client
	zd, err := zendesk.NewClient(subdomain, email, apikey)

	// Get IDs from CSV
	fmt.Println("\nReading CSV")
	ids, err := readCSV(FilePath)
	if err != nil {
		log.Fatalf("Failed to get IDs from file: %s", FilePath)
	}

	// Deleting Users
	fmt.Printf("%d Users to Delete.\n \n", len(ids))
	for _, id := range ids {
		fmt.Printf("Deleting: %v\n", id)
		zd.DeleteUser(id)
	}

	fmt.Println(color.GreenString("\nFinished!\n"))
}
