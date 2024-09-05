package cmd

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/config"
	"github.com/spf13/cobra"
)

var inputFile string
var configurationPath string
var scanRoot string
var apiKey string
var apiUrl string

const ROOT_FOLDER = "."

var rootCmd = &cobra.Command{
	Use:   "scanoss-lui",
	Short: "Lightweight UI, that presents the findings from the latest scan and prompt the user to review pending identifications.",
	Run: func(cmd *cobra.Command, args []string) {
		setConfigFile(configurationPath)
		setInputFile(inputFile)
		setScanRoot(scanRoot)
		setApiKey(apiKey)
		setApiUrl(apiUrl)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Path to the input file")
	rootCmd.Flags().StringVarP(&configurationPath, "configuration", "c", "", "Path to the configuration file")
	rootCmd.Flags().StringVarP(&scanRoot, "scanRoot", "s", "", "Path to scanned project")
	rootCmd.Flags().StringVarP(&apiKey, "key", "k", "", "API KEY")
	rootCmd.Flags().StringVarP(&apiUrl, "apiUrl", "u", "", "API URL")

}

func setConfigFile(configFile string) {
	pathToConfig := configFile
	if pathToConfig == "" {
		pathToConfig = config.GetDefaultConfigLocation()
	}

	c := config.NewConfigModule(pathToConfig)
	c.Init()
	// Load the config
	if err := c.LoadConfig(); err != nil {
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

func setApiKey(apiKey string) {
	if apiKey == "" {
		return
	}
	config.Get().ApiToken = apiKey
	// Sets Scanoss premium URL
	config.Get().ApiUrl = config.SCNOSS_PREMIUM_API_URL
}

func setApiUrl(apiUrl string) {
	if apiUrl == "" {
		return
	}
	config.Get().ApiUrl = apiUrl
}

func Execute() error {
	return rootCmd.Execute()
}
