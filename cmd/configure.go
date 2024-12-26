package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/scanoss/scanoss.lui/internal/config"
	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure the SCANOSS LUI application",
	PostRun: func(cmd *cobra.Command, args []string) {
		os.Exit(0)
	},
	Run: func(cmd *cobra.Command, args []string) {
		if apiKey == "" && apiUrl == "" {
			os.Exit(0)
		}

		cfg := config.GetInstance()

		if apiKey != "" {
			if err := cfg.SetApiToken(apiKey); err != nil {
				fmt.Println("Error saving API token:", err)
				os.Exit(1)
			}
		}

		if apiUrl != "" {
			if err := cfg.SetApiUrl(apiUrl); err != nil {
				fmt.Println("Error saving API URL:", err)
				os.Exit(1)
			}
		}

		fmt.Println("API URL: ", cfg.ApiUrl)
		fmt.Println("KEY: ", strings.Repeat("*", len(cfg.ApiToken)))
		fmt.Println("Configuration saved successfully!")
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)

	configureCmd.Flags().StringVarP(&apiKey, "key", "k", "", "SCANOSS API Key token (optional - not required for default OSSKB URL)")
	configureCmd.Flags().StringVarP(&apiUrl, "apiUrl", "u", "", "SCANOSS API URL (optional - default: https://api.osskb.org)")
}
