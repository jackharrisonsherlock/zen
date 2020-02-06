package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	input "github.com/tcnksm/go-input"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure your Zendesk Account",

	Run: func(cmd *cobra.Command, args []string) {
		configure()
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
}

func configure() {

	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	// Create config file
	home, err := homedir.Dir()
	file := filepath.Join(home, ".zen.yaml")
	f, err := os.Create(file)
	if err != nil {
		log.Fatal("Cannot create config file", err)
	}
	defer f.Close()

	// Subdomain
	q := fmt.Sprintf(color.GreenString("What is your Zendesk Subdomain (subdomain.zendesk.com)?"))
	subdomain, err := ui.Ask(q, &input.Options{
		Default:  "",
		Required: true,
		Loop:     true,
	})
	if err != nil {
		log.Fatal(err)
	}
	viper.Set("subdomain", subdomain)

	// Email
	q = fmt.Sprintf(color.GreenString("What is your Zendesk Email?"))
	email, err := ui.Ask(q, &input.Options{
		Default:  "",
		Required: true,
		Loop:     true,
	})
	if err != nil {
		log.Fatal(err)
	}
	viper.Set("email", email)

	// API Key
	q = fmt.Sprintf(color.GreenString("What is your Zendesk API Key?"))
	apikey, err := ui.Ask(q, &input.Options{
		Default:     "",
		Required:    true,
		Mask:        true,
		MaskDefault: true,
		Loop:        true,
	})
	if err != nil {
		log.Fatal(err)
	}
	viper.Set("apikey", apikey)
	viper.WriteConfig()
}
