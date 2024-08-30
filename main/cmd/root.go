package cmd

import (
	"fmt"
	"integration-git/main/pkg/common/config"
	"integration-git/main/pkg/common/scanoss_bom/infraestructure"
	modules "integration-git/main/pkg/common/scanoss_bom/module"
	"os"
	"strings"
	"runtime"
	"github.com/spf13/cobra"
	"regexp"
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
		pathToConfig = config.GetDefaultConfigLocation()
	}

	// Load the config
	if _, err := config.LoadConfig(config.GetDefaultConfigLocation()); err != nil {
		fmt.Printf("Make sure you have a %s file in the root of your project", config.GetDefaultConfigFileName())
		os.Exit(1)
	}

	// Scanoss bom file should be initialozed and read after config is loaded. Not before

	r := infraestructure.NewScanossBomJonRepository()
	r.Init()
	bomFile, _ := r.Read()
	// Init scanoss bom module. Set current bom file to singleton
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
			// Win OS only
		if runtime.GOOS == "windows" {
			// Create a regex pattern to match double slashes
			re := regexp.MustCompile(`\\+`)
			pathForwardSlash := re.ReplaceAllString(root, "/")
			config.Get().ScanRoot = pathForwardSlash
		} else {
			config.Get().ScanRoot = root
		}
	}
	if config.Get().ScanRoot == ROOT_FOLDER || config.Get().ScanRoot == "" {
		currentDir, _ := os.Getwd()
		config.Get().ScanRoot = currentDir
	}

}

func Execute() error {
	return rootCmd.Execute()
}
