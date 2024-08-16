package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"integration-git/main/pkg/common/config"
)

var inputFile string

var inputCmd = &cobra.Command{
	Use:   "input",
	Short: "Process an input file",
	Run: func(cmd *cobra.Command, args []string) {
		if inputFile == "" {
			fmt.Println("Please specify an input file using the --input flag")
			return
		}
		setInputFileToConfig(inputFile)
	},
}

func init() {
	rootCmd.AddCommand(inputCmd)
	inputCmd.Flags().StringVarP(&inputFile, "input", "i", "./scanoss/result.json", "Path to the input file")
}

func setInputFileToConfig(inputFile string) {
	if config.Get().Scanoss.ResultFilePath != inputFile {
		config.Get().Scanoss.ResultFilePath = inputFile
	}
}
