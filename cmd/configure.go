package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/scanoss/scanoss.lui/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		if apiKey != "" {
			viper.Set("apiToken", apiKey)
			config.Get().ApiToken = apiKey
		}

		if apiUrl != "" {
			viper.Set("apiUrl", apiUrl)
			config.Get().ApiUrl = apiUrl
		}

		if err := viper.WriteConfig(); err != nil {
			fmt.Println("Error saving configuration:", err)
			os.Exit(1)
		}

		fmt.Println("API URL: ", viper.GetString("apiUrl"))
		fmt.Println("KEY: ", strings.Repeat("*", len(viper.GetString("apiToken"))))
		fmt.Println("Configuration saved successfully!")
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)

	configureCmd.Flags().StringVarP(&apiKey, "key", "k", "", "SCANOSS API Key token (optional - not required for default OSSKB URL)")
	configureCmd.Flags().StringVarP(&apiUrl, "apiUrl", "u", "", "SCANOSS API URL (optional - default: https://api.osskb.org/scan/direct)")
}
