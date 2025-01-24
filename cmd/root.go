// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/scanoss/scanoss.cc/backend/entities"
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
var version bool
var originalWorkDir string

var rootCmd = &cobra.Command{
	Use:   "scanoss-cc",
	Short: "Lightweight UI, that presents the findings from the latest scan and prompt the user to review pending identifications.",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	// We need the original wd to set the default values for the config.
	// Otherwise, when running from a symlink, the default values will be incorrect.
	var err error
	originalWorkDir, err = os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting original working directory")
	}

	cobra.OnInitialize(initConfig)

	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "Show application version")
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
	if err := config.InitializeConfig(cfgFile, scanRoot, apiKey, apiUrl, inputFile, scanossSettingsFilePath, originalWorkDir, debug); err != nil {
		log.Fatal().Err(err).Msg("Error initializing config")
	}
}

func Execute() error {
	isVersionCmd := len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v")
	if isVersionCmd {
		fmt.Printf("scanoss-cc %s\n", entities.AppVersion)
		os.Exit(0)
	}

	// Workaround to exit process when help command is called
	isHelpCmd := len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h")
	if isHelpCmd {
		rootCmd.Help()
		os.Exit(0)
	}

	return rootCmd.Execute()
}
