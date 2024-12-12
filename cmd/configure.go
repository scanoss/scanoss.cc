package cmd

import (
	"os"
	"strings"

	"github.com/rs/zerolog/log"
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
				log.Error().Err(err).Msg("Error saving API token")
				os.Exit(1)
			}
		}

		if apiUrl != "" {
			if err := cfg.SetApiUrl(apiUrl); err != nil {
				log.Error().Err(err).Msg("Error saving API URL")
				os.Exit(1)
			}
		}

		log.Info().Msgf("API URL: %s", cfg.ApiUrl)
		log.Info().Msgf("KEY: %s", strings.Repeat("*", len(cfg.ApiToken)))
		log.Info().Msg("Configuration saved successfully!")
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)

	configureCmd.Flags().StringVarP(&apiKey, "key", "k", "", "SCANOSS API Key token (optional - not required for default OSSKB URL)")
	configureCmd.Flags().StringVarP(&apiUrl, "apiUrl", "u", "", "SCANOSS API URL (optional - default: https://api.osskb.org)")
}
