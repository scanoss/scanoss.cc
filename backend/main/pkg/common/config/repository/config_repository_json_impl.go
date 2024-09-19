package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/config/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/utils"
)

var (
	ErrReadingFile       = errors.New("error reading file")
	ErrUnmarshallingFile = errors.New("error unmarshalling file file")
)

type ConfigJsonRepository struct {
	configPath string
}

func NewConfigJsonRepository(path string) *ConfigJsonRepository {
	return &ConfigJsonRepository{
		configPath: path,
	}
}

func (r *ConfigJsonRepository) Save(config *entities.Config) {
	selectedData := make(map[string]string)
	selectedData["apiUrl"] = config.ApiUrl
	selectedData["apiToken"] = config.ApiToken
	utils.WriteJsonFile(r.configPath, selectedData)
}

func (r *ConfigJsonRepository) Read() (entities.Config, error) {
	fileData, err := os.ReadFile(r.configPath)

	if err != nil {
		fmt.Printf("Unable to read configuration file %v\n", err)
		os.Exit(1)
	}

	// Marshal the default config into JSON
	defaultConfig, err := json.Marshal(getDefaultConfigFile())
	if err != nil {
		return entities.Config{}, err
	}

	// Merge the JSON by unmarshalling the file's content into the default JSON
	var mergedData map[string]json.RawMessage
	if err := json.Unmarshal(defaultConfig, &mergedData); err != nil {
		return entities.Config{}, err
	}
	if err := json.Unmarshal(fileData, &mergedData); err != nil {
		return entities.Config{}, err
	}

	// Marshal the merged JSON back into the config struct
	var cfg entities.Config
	mergedBytes, err := json.Marshal(mergedData)
	if err != nil {
		return entities.Config{}, err
	}
	if err := json.Unmarshal(mergedBytes, &cfg); err != nil {
		return entities.Config{}, err
	}
	return cfg, nil
}

func getDefaultConfigFile() entities.Config {
	workingDir, _ := os.Getwd()
	var defaultConfigFile entities.Config = entities.Config{
		ScanRoot:             "",
		ResultFilePath:       workingDir + string(os.PathSeparator) + ".scanoss" + string(os.PathSeparator) + "results.json",
		ApiToken:             "",
		ApiUrl:               "https://api.osskb.org",
		ScanSettingsFilePath: workingDir + string(os.PathSeparator) + "scanoss.json",
	}

	return defaultConfigFile
}

func (r *ConfigJsonRepository) Init() error {

	err := utils.FileExist(r.configPath)
	if err != nil {
		fmt.Println("Configuration file does not exists. Creating default config file...")
		err = r.createConfigFile()
		if err != nil {
			return err
		}
		if utils.IsWritableFile(r.configPath) != true {
			return errors.New(fmt.Sprintf("File %s is not writable\n", r.configPath))
		}
		return nil
	}

	if utils.IsWritableFile(r.configPath) != true {
		return errors.New(fmt.Sprintf("File %s is not writable\n", r.configPath))
	}

	return nil
}

func (r *ConfigJsonRepository) createConfigFile() error {
	// Get the directory path
	dirPath := filepath.Dir(r.configPath)
	// Check if the directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// Directory does not exist, create it
		err := os.MkdirAll(dirPath, 0700)
		if err != nil {
			return err
		}
	}

	// Create the file with read/write permissions (owner and group only)
	file, err := os.OpenFile(r.configPath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		fmt.Printf("Failed to create configuration file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	defaultConfig := getDefaultConfigFile()
	content := fmt.Sprintf(`{ "apiToken": "" , "apiUrl": "%s"}`, defaultConfig.ApiUrl)
	// Write the content to the file
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}
