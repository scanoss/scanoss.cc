package repository

import (
	"errors"
	"fmt"
	"os"

	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/config"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/utils"
)

type ScanossSettingsJsonRepository struct {
	fr utils.FileReader
}

func NewScanossSettingsJsonRepository(fr utils.FileReader) ScanossSettingsRepository {
	return &ScanossSettingsJsonRepository{
		fr: fr,
	}
}

func (r *ScanossSettingsJsonRepository) Save() error {
	sf := r.GetSettingsFileContent()
	if err := utils.WriteJsonFile(config.Get().ScanSettingsFilePath, sf); err != nil {
		return err
	}
	return nil
}

func (r *ScanossSettingsJsonRepository) Read() (entities.SettingsFile, error) {
	if config.Get() == nil {
		return entities.SettingsFile{}, fmt.Errorf("config is not initialized")
	}
	scanSettingsFileBytes, err := r.fr.ReadFile(config.Get().ScanSettingsFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return entities.SettingsFile{}, nil
		}
		return entities.SettingsFile{}, err
	}

	scanossSettings, err := utils.JSONParse[entities.SettingsFile](scanSettingsFileBytes)
	if err != nil {
		return entities.SettingsFile{}, err
	}

	return scanossSettings, nil
}

func (r *ScanossSettingsJsonRepository) GetSettingsFileContent() *entities.SettingsFile {
	return entities.ScanossSettingsJson.SettingsFile
}

func (r *ScanossSettingsJsonRepository) HasUnsavedChanges() (bool, error) {
	originalBom, err := r.Read()
	if err != nil {
		return false, err
	}

	return !originalBom.Equal(r.GetSettingsFileContent()), nil
}

func (r *ScanossSettingsJsonRepository) AddBomEntry(newEntry entities.ComponentFilter, filterAction string) error {
	sf := r.GetSettingsFileContent()
	var targetList *[]entities.ComponentFilter

	switch filterAction {
	case "remove":
		targetList = &sf.Bom.Remove
	case "include":
		targetList = &sf.Bom.Include
	default:
		return fmt.Errorf("invalid filter action: %s", filterAction)
	}

	r.removeDuplicatesFromAllLists(newEntry)

	*targetList = append(*targetList, newEntry)

	return nil
}

func (r *ScanossSettingsJsonRepository) removeDuplicatesFromAllLists(newEntry entities.ComponentFilter) {
	sf := r.GetSettingsFileContent()

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

func (r *ScanossSettingsJsonRepository) ClearAllFilters() error {
	sf := r.GetSettingsFileContent()
	sf.Bom.Include = []entities.ComponentFilter{}
	sf.Bom.Remove = []entities.ComponentFilter{}
	return nil
}
