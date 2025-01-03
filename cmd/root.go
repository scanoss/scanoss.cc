// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/scanoss/scanoss.cc/internal/config"
	"github.com/spf13/cobra"
)

var apiKey string
var apiUrl string
var cfgFile string
var debug bool
var inputFile string
var scanossSettingsFilePath string
var scanRoot string

var rootCmd = &cobra.Command{
	Use:   "scanoss-cc",
	Short: "Lightweight UI, that presents the findings from the latest scan and prompt the user to review pending identifications.",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "Config file (optional - default: $HOME/.scanoss/scanoss-cc-settings.json)")
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Path to scan result file (optional - default: $WORKDIR/.scanoss/results.json)")
	rootCmd.Flags().StringVarP(&scanRoot, "scan-root", "s", "", "Scanned folder root path (optional - default: $WORKDIR)")
	rootCmd.Flags().StringVar(&scanossSettingsFilePath, "settings", "", "Path to scanoss settings file (optional - default: $WORKDIR/scanoss.json)")
	rootCmd.Flags().StringVarP(&apiKey, "key", "k", "", "SCANOSS API Key token (optional)")
	rootCmd.Flags().StringVarP(&apiUrl, "apiUrl", "u", "", fmt.Sprintf("SCANOSS API URL (optional - default: %s)", config.DEFAULT_API_URL))
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")

	rootCmd.Root().CompletionOptions.HiddenDefaultCmd = true
}

func initConfig() {
	if err := config.InitializeConfig(cfgFile, scanRoot, apiKey, apiUrl, inputFile, scanossSettingsFilePath, debug); err != nil {
		log.Fatal().Err(err).Msg("Error initializing config")
	}
}

func Execute() error {
	// Workaround to exit process when help command is called
	isHelpCmd := len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h")
	if isHelpCmd {
		rootCmd.Help()
		os.Exit(0)
	}

	return rootCmd.Execute()
}
