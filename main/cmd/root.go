package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"integration-git/main/pkg/common/config"
	"os"
	"strings"
)

var inputFile string
var configurationPath string
var scanRoot string

const ROOT_FOLDER = "."

var rootCmd = &cobra.Command{
	Use:   "integration-git",
	Short: "Lightweight UI, that presents the findings from the latest scan and prompt the user to review pending identifications.",
	Run: func(cmd *cobra.Command, args []string) {
		setConfigFile(configurationPath)
		setInputFile(inputFile)
		setScanRoot(scanRoot)
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
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Unable to read user home directory")
			os.Exit(1)
		}
		pathToConfig = homeDir + string(os.PathSeparator) + ROOT_FOLDER + config.GetDefaultGlobalFolder() + string(os.PathSeparator) + config.GetDefaultConfigFileName()
	}

	// Load the config
	if _, err := config.LoadConfig(pathToConfig); err != nil {
		fmt.Printf("Make sure you have a %s file in the root of your project", config.GetDefaultConfigFileName())
		os.Exit(1)
	}
}

func setInputFile(resultFile string) {
	input := resultFile
	if input != "" {
		config.Get().ResultFilePath = inputFile
	} else {
		resultFilePath := config.Get().ResultFilePath
		if resultFilePath != "" && strings.HasPrefix(resultFilePath, ROOT_FOLDER) {
			if currentDir, err := os.Getwd(); err == nil {
				// Workaround due to path.Join removes "." when join current dir with resultFilePath
				config.Get().ResultFilePath = currentDir + string(os.PathSeparator) + ROOT_FOLDER + resultFilePath[2:]
			}
		}
	}

	fmt.Println("CONFIG", config.Get())

}

func setScanRoot(root string) {
	if root != "" {
		config.Get().ScanRoot = root
	}
	if config.Get().ScanRoot == ROOT_FOLDER || config.Get().ScanRoot == "" {
		currentDir, _ := os.Getwd()
		config.Get().ScanRoot = currentDir
	}

}

func Execute() error {
	return rootCmd.Execute()
}
