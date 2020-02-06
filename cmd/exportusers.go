package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MEDIGO/go-zendesk/zendesk"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// exportusersCmd represents the exportusers command
var exportusersCmd = &cobra.Command{
	Use:   "export-users",
	Short: "Export Users",
	Long:  "Exports a list of end users",

	Run: func(cmd *cobra.Command, args []string) {
		exportUsers()
	},
}

func init() {
	rootCmd.AddCommand(exportusersCmd)
}

func exportUsers() {

	// Create new Zendesk Client
	zd, err := zendesk.NewClient(subdomain, email, apikey)

	// Export Users
	fmt.Println("\nRetrieving Users from Zendesk...")
	users := []zendesk.User{}
	counter := 1
	options := &zendesk.ListUsersOptions{Page: counter}
	page, err := zd.ListUsers(options)
	users = append(users, page...)
	if err != nil {
		log.Fatalf("Failed retrieving users from Zendesk: %v", err)
	}

	for len(page) > 0 {
		counter++
		options := &zendesk.ListUsersOptions{Page: counter}
		page, err = zd.ListUsers(options)
		users = append(users, page...)
	}

	// Write to CSV
	fmt.Println("Writing users to CSV file...")
	err = WriteFileUsers(users)
	if err != nil {
		log.Fatalf("Failed writing users to CSV: %v", err)
	}

	fmt.Println(color.GreenString("\nFinished!"))
}

func WriteFileUsers(users []zendesk.User) error {

	// Create a blank CSV file.
	fileName := fmt.Sprintf("zendesk_user_export_%v.csv", time.Now().Unix())
	file, err := os.Create(fmt.Sprintf("./%s", fileName))
	if err != nil {
		return err
	}

	// Ensure the file is closed at the end.
	defer file.Close()
	writer := csv.NewWriter(file)

	// Write the header line.
	var header []string
	header = append(header, "url")
	header = append(header, "id")
	header = append(header, "name")
	header = append(header, "alias")
	header = append(header, "created_at")
	header = append(header, "updated_at")
	header = append(header, "active")
	header = append(header, "verified")
	header = append(header, "shared")
	header = append(header, "timezone")
	header = append(header, "last_login_at")
	header = append(header, "email")
	header = append(header, "phone")
	header = append(header, "signature")
	header = append(header, "details")
	header = append(header, "role")
	header = append(header, "restricted_agent")
	header = append(header, "suspended")

	// Commit the header.
	writer.Write(header)

	// Now loop through each  object and populate the CSV.
	for _, user := range users {

		var url, id, name, alias, createdAt, updatedAt, active, verified, shared, timezone, lastLoginAt, phone, email, signature, details, role, restrictedAgent, suspended string

		if user.URL != nil {
			url = *user.URL
		}
		if user.ID != nil {
			id = fmt.Sprintf("%v", *user.ID)
		}
		if user.Name != nil {
			name = *user.Name
		}
		if user.Alias != nil {
			alias = *user.Alias
		}
		if user.CreatedAt != nil {
			createdAt = fmt.Sprintf("%v", *user.CreatedAt)
		}
		if user.UpdatedAt != nil {
			updatedAt = fmt.Sprintf("%v", *user.UpdatedAt)
		}
		if user.Active != nil {
			active = fmt.Sprintf("%v", *user.Active)
		}
		if user.Verified != nil {
			verified = fmt.Sprintf("%v", *user.Verified)
		}
		if user.Shared != nil {
			shared = fmt.Sprintf("%v", *user.Shared)
		}
		if user.TimeZone != nil {
			timezone = fmt.Sprintf("%v", *user.TimeZone)
		}
		if user.LastLoginAt != nil {
			lastLoginAt = fmt.Sprintf("%v", *user.LastLoginAt)
		}
		if user.Email != nil {
			email = *user.Email
		}
		if user.Phone != nil {
			phone = *user.Phone
		}
		if user.Signature != nil {
			signature = *user.Signature
		}
		if user.Details != nil {
			details = *user.Details
		}
		if user.Role != nil {
			role = *user.Role
		}
		if user.RestrictedAgent != nil {
			restrictedAgent = fmt.Sprintf("%v", *user.RestrictedAgent)
		}
		if user.Suspended != nil {
			suspended = fmt.Sprintf("%v", *user.Suspended)
		}

		var record []string
		record = append(record, url)
		record = append(record, id)
		record = append(record, name)
		record = append(record, alias)
		record = append(record, createdAt)
		record = append(record, updatedAt)
		record = append(record, active)
		record = append(record, verified)
		record = append(record, shared)
		record = append(record, timezone)
		record = append(record, lastLoginAt)
		record = append(record, email)
		record = append(record, phone)
		record = append(record, signature)
		record = append(record, details)
		record = append(record, role)
		record = append(record, restrictedAgent)
		record = append(record, suspended)
		writer.Write(record)
	}

	writer.Flush()
	return err
}
