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

// exportarticlesCmd represents the exportarticles command
var exportarticlesCmd = &cobra.Command{
	Use:   "export-articles",
	Short: "Export Articles",
	Long:  "Exports a list of Help Centre Articles",

	Run: func(cmd *cobra.Command, args []string) {
		exportArticles()
	},
}

func init() {
	rootCmd.AddCommand(exportarticlesCmd)
}

func exportArticles() {

	// Create new Zendesk Client
	zd, err := zendesk.NewClient(subdomain, email, apikey)

	// Export Articles
	fmt.Println("\nRetrieving Articles from Zendesk...")
	articles := []zendesk.Article{}
	counter := 1
	options := &zendesk.ListOptions{Page: counter}
	page, err := zd.ListArticles(options)
	articles = append(articles, page...)
	if err != nil {
		log.Fatalf("Failed retrieving Articles from Zendesk: %v", err)
	}

	for len(page) > 0 {
		counter++
		options := &zendesk.ListOptions{Page: counter}
		page, err = zd.ListArticles(options)
		articles = append(articles, page...)
	}

	// Write to CSV
	fmt.Println("Writing Articles to CSV file...")
	err = WriteFileArticles(articles)
	if err != nil {
		log.Fatalf("Failed writing Articles to CSV: %v", err)
	}

	fmt.Println(color.GreenString("\nFinished!"))
}

func WriteFileArticles(articles []zendesk.Article) error {

	// Create a blank CSV file.
	fileName := fmt.Sprintf("zendesk_article_export_%v.csv", time.Now().Unix())
	file, err := os.Create(fmt.Sprintf("./%s", fileName))
	if err != nil {
		return err
	}

	// Ensure the file is closed at the end.
	defer file.Close()
	writer := csv.NewWriter(file)

	// Write the header line.
	var header []string
	header = append(header, "id")
	header = append(header, "url")
	header = append(header, "author_id")
	header = append(header, "section_id")
	header = append(header, "name")
	header = append(header, "created_at")
	header = append(header, "updated_at")
	header = append(header, "vote_sum")
	header = append(header, "vote_count")
	// header = append(header, "body")

	// Commit the header.
	writer.Write(header)

	// Now loop through each  object and populate the CSV.
	for _, article := range articles {

		var id, url, author_id, section_id, name, created_at, updated_at, vote_sum, vote_count string

		if article.ID != nil {
			id = fmt.Sprintf("%v", *article.ID)
		}
		if article.URL != nil {
			url = *article.URL
		}
		if article.AuthorID != nil {
			author_id = fmt.Sprintf("%v", *article.AuthorID)
		}
		if article.SectionID != nil {
			section_id = fmt.Sprintf("%v", *article.SectionID)
		}
		if article.Name != nil {
			name = *article.Name
		}
		if article.CreatedAt != nil {
			created_at = fmt.Sprintf("%v", *article.CreatedAt)
		}
		if article.UpdatedAt != nil {
			updated_at = fmt.Sprintf("%v", *article.UpdatedAt)
		}
		if article.VoteSum != nil {
			vote_sum = fmt.Sprintf("%v", *article.VoteSum)
		}
		if article.VoteCount != nil {
			vote_count = fmt.Sprintf("%v", *article.VoteCount)
		}
		// if article.Body != nil {
		// 	body = *article.Body
		// }

		var record []string
		record = append(record, id)
		record = append(record, url)
		record = append(record, author_id)
		record = append(record, section_id)
		record = append(record, name)
		record = append(record, created_at)
		record = append(record, updated_at)
		record = append(record, vote_sum)
		record = append(record, vote_count)
		// record = append(record, body)

		writer.Write(record)
	}

	writer.Flush()
	return err
}
