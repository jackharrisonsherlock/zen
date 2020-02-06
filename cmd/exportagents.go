package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	zendesk "github.com/MEDIGO/go-zendesk/zendesk"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// exportagentsCmd represents the exportagents command
var exportagentsCmd = &cobra.Command{
	Use:   "export-agents",
	Short: "Export Agents",

	Run: func(cmd *cobra.Command, args []string) {
		exportAgents()
	},
}

func init() {
	rootCmd.AddCommand(exportagentsCmd)
}

func exportAgents() {

	// Create new Zendesk Client
	zd, err := zendesk.NewClient(subdomain, email, apikey)

	// Export Agents
	fmt.Println("\nRetrieving Agents from Zendesk...")
	agents := []zendesk.User{}
	counter := 1
	options := &zendesk.ListUsersOptions{Page: counter}
	page, err := zd.ListAgents(options)
	agents = append(agents, page...)
	if err != nil {
		log.Fatalf("Failed retrieving agents from Zendesk: %v", err)
	}

	for len(page) > 0 {
		counter++
		options := &zendesk.ListUsersOptions{Page: counter}
		page, err = zd.ListAgents(options)
		agents = append(agents, page...)
	}

	// Write to CSV
	fmt.Println("Writing agents to CSV file...")
	err = WriteFileAgents(agents)
	if err != nil {
		log.Fatalf("Failed writing agents to CSV: %v", err)
	}

	fmt.Println(color.GreenString("\nFinished!"))
}

func WriteFileAgents(agents []zendesk.User) error {

	// Create a blank CSV file.
	fileName := fmt.Sprintf("zendesk_agent_export_%v.csv", time.Now().Unix())
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
	for _, agent := range agents {

		var url, id, name, alias, createdAt, updatedAt, active, verified, shared, timezone, lastLoginAt, phone, email, signature, details, role, restrictedAgent, suspended string

		if agent.URL != nil {
			url = *agent.URL
		}
		if agent.ID != nil {
			id = fmt.Sprintf("%v", *agent.ID)
		}
		if agent.Name != nil {
			name = *agent.Name
		}
		if agent.Alias != nil {
			alias = *agent.Alias
		}
		if agent.CreatedAt != nil {
			createdAt = fmt.Sprintf("%v", *agent.CreatedAt)
		}
		if agent.UpdatedAt != nil {
			updatedAt = fmt.Sprintf("%v", *agent.UpdatedAt)
		}
		if agent.Active != nil {
			active = fmt.Sprintf("%v", *agent.Active)
		}
		if agent.Verified != nil {
			verified = fmt.Sprintf("%v", *agent.Verified)
		}
		if agent.Shared != nil {
			shared = fmt.Sprintf("%v", *agent.Shared)
		}
		if agent.TimeZone != nil {
			timezone = fmt.Sprintf("%v", *agent.TimeZone)
		}
		if agent.Email != nil {
			email = *agent.Email
		}
		if agent.Phone != nil {
			phone = *agent.Phone
		}
		if agent.Signature != nil {
			signature = *agent.Signature
		}
		if agent.Details != nil {
			details = *agent.Details
		}
		if agent.Role != nil {
			role = *agent.Role
		}
		if agent.RestrictedAgent != nil {
			restrictedAgent = fmt.Sprintf("%v", *agent.RestrictedAgent)
		}
		if agent.Suspended != nil {
			suspended = fmt.Sprintf("%v", *agent.Suspended)
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
