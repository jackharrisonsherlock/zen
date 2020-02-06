package cmd

import (
	"encoding/csv"
	"fmt"
	"os"

	zendesk "github.com/MEDIGO/go-zendesk/zendesk"
	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Variables for Client authentication details and flags
var (
	ZendeskClient *zendesk.Client
	FilePath      string
	StartTime     string
	subdomain     string
	email         string
	apikey        string
	cfgFile       string
	logo          = color.GreenString(`  ____   ___   _ __  
 |_  /  / _ \ | '_ \ 
  / /  |  __/ | | | |
 /___|  \___| |_| |_|`)
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "zen",
	Short: fmt.Sprintf(`
%s`, logo),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".zen"
		viper.AddConfigPath(home)
		viper.SetConfigName(".zen")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		configure()
		// fmt.Println("Config file not found. Run 'zen configure' command.", viper.ConfigFileUsed())
		// os.Exit(1)
	}

	subdomain = viper.GetString("subdomain")
	email = viper.GetString("email")
	apikey = viper.GetString("apikey")
}

// Read passed CSV and returns the IDs
func readCSV(FilePath string) ([]string, error) {

	// Open our provided CSV file
	file, err := os.Open(FilePath)
	if err != nil {
		fmt.Println("Could not read from CSV file")
		return nil, err
	}

	defer file.Close()
	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var rowNumber int
	entities := []string{}

	// Loop through rows and assign them
	for _, row := range rows {
		rowNumber++
		entityIDs := row[0]
		entities = append(entities, entityIDs)
	}

	return entities, err
}
