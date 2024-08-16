package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"integration-git/main/pkg/common/config"
	"os"
	"path"
	"strings"
)

var inputFile string
var configurationPath string
var scanRoot string

var rootCmd = &cobra.Command{
	Use:   "integration-git",
	Short: "Lightweight UI, that presents the findings from the latest scan and prompt the user to review pending identifications.",
	Run: func(cmd *cobra.Command, args []string) {
		setConfigFile(configurationPath)
		setInputFile(inputFile)
		setScanRoot(scanRoot)
		return
	},
}

func init() {
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Path to the input file")
	rootCmd.Flags().StringVarP(&configurationPath, "configuration", "c", "", "Path to the configuration file")
	rootCmd.Flags().StringVarP(&scanRoot, "scanRoot", "s", "", "Path to scanned project")

}

func setConfigFile(configFile string) {
	pathToConfig := configFile
	if pathToConfig == "" {
		root, _ := os.Getwd()
		pathToConfig = path.Join(root, "config.json")
	}

	// Load the config
	if _, err := config.LoadConfig(pathToConfig); err != nil {
		fmt.Println("Error reading configuration file:", pathToConfig)
		os.Exit(1)
	}
}

func setInputFile(resultFile string) {
	input := resultFile
	if input != "" {
		config.Get().Scanoss.ResultFilePath = inputFile
	} else {
		resultFilePath := config.Get().Scanoss.ResultFilePath
		if resultFilePath != "" && strings.HasPrefix(resultFilePath, ".") {
			if currentDir, err := os.Getwd(); err == nil {
				// Workaround due to path.Join removes "." when join current dir with resultFilePath
				config.Get().Scanoss.ResultFilePath = currentDir + string(os.PathSeparator) + "." + resultFilePath[2:]
			}
		}
	}

}

func setScanRoot(root string) {
	if root != "" {
		config.Get().Scanoss.ScanRoot = root
	}
	if config.Get().Scanoss.ScanRoot == "." || config.Get().Scanoss.ScanRoot == "" {
		currentDir, _ := os.Getwd()
		config.Get().Scanoss.ScanRoot = currentDir
	}

}

func Execute() error {
	return rootCmd.Execute()
}
