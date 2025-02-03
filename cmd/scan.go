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
	"strings"

	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/backend/service"
	"github.com/spf13/cobra"
)

func NewScanCmd(scanService service.ScanService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scan [scanDirPath]",
		Short: "Run a scan on the specified folder",
		Args: func(cmd *cobra.Command, args []string) error {
			filesFlag := cmd.Flag("files")
			// If no files are specified and no folder is specified, return an error
			if filesFlag == nil || !filesFlag.Changed && len(args) == 0 {
				return fmt.Errorf("you must specify a folder to scan")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := scanService.CheckDependencies(); err != nil {
				return err
			}

			scanOptions := make([]string, 0)
			scanOptions = append(scanOptions, "--quiet")

			var scanDirPath string

			// Check if the folder path is specified as an argument
			// Could be that the user specifies only --files flag (e.g scan --files file1.go file2.go)
			if len(args) == 1 && args[0] != "" {
				scanDirPath = args[0]
				scanOptions = append(scanOptions, scanDirPath)
			}

			for _, arg := range entities.ScanArguments {
				flag := cmd.Flag(arg.Name)
				if flag == nil || !flag.Changed {
					continue
				}

				switch arg.Type {
				case "string":
					scanOptions = append(scanOptions, fmt.Sprintf("--%s", arg.Name), flag.Value.String())
				case "stringSlice":
					values, err := cmd.Flags().GetStringSlice(arg.Name)
					if err != nil {
						return fmt.Errorf("an error occurred with argument %s: %w", arg.Name, err)
					}
					commaSeparatedValues := strings.Join(values, ",")
					scanOptions = append(scanOptions, fmt.Sprintf("--%s", arg.Name), commaSeparatedValues)
				case "int":
					scanOptions = append(scanOptions, fmt.Sprintf("--%s", arg.Name), flag.Value.String())
				case "bool":
					if flag.Value.String() == "true" {
						scanOptions = append(scanOptions, fmt.Sprintf("--%s", arg.Name))
					}
				}
			}

			return scanService.Scan(scanOptions)
		},
	}

	for _, arg := range entities.ScanArguments {
		switch arg.Type {
		case "string":
			cmd.Flags().StringP(arg.Name, arg.Shorthand, arg.Default.(string), arg.Usage)
		case "stringSlice":
			cmd.Flags().StringSliceP(arg.Name, arg.Shorthand, arg.Default.([]string), arg.Usage)
		case "int":
			cmd.Flags().IntP(arg.Name, arg.Shorthand, arg.Default.(int), arg.Usage)
		case "bool":
			cmd.Flags().BoolP(arg.Name, arg.Shorthand, arg.Default.(bool), arg.Usage)
		}
	}

	setupHelpCommand(cmd)
	return cmd
}

func init() {
	service := service.NewScanServicePythonImpl()
	scanCmd := NewScanCmd(service)

	// This is a workaround to prevent the scan command opening the code compare when running tests
	if os.Getenv("GO_TEST") != "true" {
		scanCmd.PostRun = func(cmd *cobra.Command, args []string) {
			os.Exit(0)
		}
	}

	rootCmd.AddCommand(scanCmd)
}
