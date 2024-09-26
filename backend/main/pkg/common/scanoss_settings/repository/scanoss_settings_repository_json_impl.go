package repository

import (
	"fmt"

	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/config"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/utils"
)

type ScanossSettingsJsonRepository struct {
}

func NewScanossSettingsJsonRepository() *ScanossSettingsJsonRepository {
	return &ScanossSettingsJsonRepository{}
}

func (r *ScanossSettingsJsonRepository) Save() error {
	if err := utils.WriteJsonFile(config.Get().ScanSettingsFilePath, entities.ScanossSettingsJson.SettingsFile); err != nil {
		return err
	}
	return nil
}

func (r *ScanossSettingsJsonRepository) Read() (entities.SettingsFile, error) {

	if config.Get() == nil {
		return entities.SettingsFile{}, fmt.Errorf("config is nil")
	}
	scanSettingsFileBytes, err := utils.ReadFile(config.Get().ScanSettingsFilePath)
	if err != nil {
		return entities.SettingsFile{}, err
	}

	scanossSettings, err := utils.JSONParse[entities.SettingsFile](scanSettingsFileBytes)
	if err != nil {
		return entities.SettingsFile{}, err
	}

	return scanossSettings, nil
}

func (r *ScanossSettingsJsonRepository) HasUnsavedChanges() (bool, error) {
	originalBom, err := r.Read()
	if err != nil {
		return false, err
	}

	return !originalBom.Equal(entities.ScanossSettingsJson.SettingsFile), nil
}

func (r *ScanossSettingsJsonRepository) AddBomEntry(newEntry entities.ComponentFilter, filterAction string) error {
	sf := entities.ScanossSettingsJson.SettingsFile
	var targetList *[]entities.ComponentFilter

	switch filterAction {
	case "remove":
		targetList = &sf.Bom.Remove
	case "include":
		targetList = &sf.Bom.Include
	default:
		return fmt.Errorf("invalid filter action: %s", filterAction)
	}

	removeDuplicatesFromAllLists(newEntry)

	*targetList = append(*targetList, newEntry)

	return nil
}

func removeDuplicatesFromAllLists(newEntry entities.ComponentFilter) {
	sf := entities.ScanossSettingsJson.SettingsFile

	sf.Bom.Remove = removeDuplicatesFromList(sf.Bom.Remove, newEntry)
	sf.Bom.Include = removeDuplicatesFromList(sf.Bom.Include, newEntry)
}

func removeDuplicatesFromList(list []entities.ComponentFilter, newEntry entities.ComponentFilter) []entities.ComponentFilter {
	result := make([]entities.ComponentFilter, 0, len(list))
	for _, entry := range list {
		if !isDuplicate(entry, newEntry) {
			result = append(result, entry)
		}
	}
	return result
}

func isDuplicate(entry, newEntry entities.ComponentFilter) bool {
	if newEntry.Path == "" {
		return entry.Purl == newEntry.Purl
	}
	return entry.Purl == newEntry.Purl && entry.Path == newEntry.Path
}
