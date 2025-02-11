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
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/scanoss/scanoss.cc/internal/config"
	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure the SCANOSS Code Compare application",
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

		log.Info().Msgf("API URL: %s", cfg.GetApiUrl())
		log.Info().Msgf("KEY: %s", strings.Repeat("*", len(cfg.GetApiToken())))
		log.Info().Msg("Configuration saved successfully!")
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)

	configureCmd.Flags().StringVarP(&apiKey, "key", "k", "", "SCANOSS API Key token (optional - not required for default OSSKB URL)")
	configureCmd.Flags().StringVarP(&apiUrl, "apiUrl", "u", "", "SCANOSS API URL (optional - default: https://api.osskb.org)")

	setupHelpCommand(configureCmd)
}
