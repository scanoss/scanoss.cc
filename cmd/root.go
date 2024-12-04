package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/scanoss/scanoss.lui/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var inputFile string
var cfgFile string
var scanRoot string
var apiKey string
var apiUrl string

var rootCmd = &cobra.Command{
	Use:   "scanoss-lui",
	Short: "Lightweight UI, that presents the findings from the latest scan and prompt the user to review pending identifications.",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	cobra.OnInitialize(initConfig)

	// viper.SetKeysCaseSensitive(true)

	rootCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "Config file (optional - default: $HOME/.scanoss/.scanoss-lui-settings.json)")
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Path to scan result file (optional - default: $WORKDIR/.scanoss/results.json)")
	rootCmd.Flags().StringVarP(&scanRoot, "scan-root", "s", "", "Scanned folder root path (optional - default: $WORKDIR)")
	rootCmd.Flags().StringVarP(&apiKey, "key", "k", "", "SCANOSS API Key token (optional)")
	rootCmd.Flags().StringVarP(&apiUrl, "apiUrl", "u", "", fmt.Sprintf("SCANOSS API URL (optional - default: %s)", config.DEFAULT_API_URL))

	rootCmd.Root().CompletionOptions.HiddenDefaultCmd = true
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		configDir := filepath.Dir(config.GetDefaultConfigLocation())
		configName := strings.TrimSuffix(filepath.Base(config.GLOBAL_CONFIG_FILE_NAME), ".json")

		viper.AddConfigPath(configDir)
		viper.SetConfigName(configName)
		viper.SetConfigType("json")
	}

	// Default values
	viper.SetDefault("apiUrl", config.DEFAULT_API_URL)
	viper.SetDefault("apiToken", "")
	viper.SetDefault("resultFilePath", config.GetDefaultResultFilePath())
	viper.SetDefault("scanRoot", "")
	viper.SetDefault("scanSettingsFilePath", config.GetDefaultScanSettingsFilePath())

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Create default config
			if err := viper.SafeWriteConfig(); err != nil {
				fmt.Println("Error creating config file:", err)
				os.Exit(1)
			}
			fmt.Println("Created config file:", viper.ConfigFileUsed())
		} else {
			fmt.Println("Error reading config file:", err)
			os.Exit(1)
		}
	}

	// Override with command line flags
	if scanRoot != "" {
		viper.Set("scanRoot", scanRoot)
	}
	if apiKey != "" {
		viper.Set("apiToken", apiKey)
	}
	if apiUrl != "" {
		viper.Set("apiUrl", apiUrl)
	}
	if inputFile != "" {
		viper.Set("resultFilePath", inputFile)
	}

	if scanRoot != "" || apiKey != "" || apiUrl != "" || inputFile != "" {
		if err := viper.WriteConfig(); err != nil {
			fmt.Println("Error saving config:", err)
			os.Exit(1)
		}
	}

	cfg := config.Get()
	cfg.ApiToken = viper.GetString("apiToken")
	cfg.ApiUrl = viper.GetString("apiUrl")
	cfg.ResultFilePath = viper.GetString("resultFilePath")
	cfg.ScanRoot = viper.GetString("scanRoot")
	cfg.ScanSettingsFilePath = viper.GetString("scanSettingsFilePath")
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
