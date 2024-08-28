package cmd

import (
	"fmt"
	"integration-git/main/pkg/common/config"
	"integration-git/main/pkg/common/scanoss_bom/infraestructure"
	"integration-git/main/pkg/common/scanoss_bom/module"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var inputFile string
var configurationPath string
var scanRoot string

const ROOT_FOLDER = "."

var rootCmd = &cobra.Command{
	Use:   "scanoss-lui",
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

	// Scanoss bom file should be read after config is loaded. Not before
	// Read current scanoss.json file and init ScanossBom module to use as singleton
	bomFile, _ := infraestructure.NewScanossBomJonRepository().Read()
	modules.NewScanossBomModule().Init(&bomFile)
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
