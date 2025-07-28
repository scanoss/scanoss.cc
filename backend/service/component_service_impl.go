// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package service

import (
	"sort"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	purlutils "github.com/scanoss/go-purl-helper/pkg"
	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/backend/mappers"
	"github.com/scanoss/scanoss.cc/backend/repository"
	"github.com/scanoss/scanoss.cc/internal/utils"
)

type ComponentServiceImpl struct {
	repo                repository.ComponentRepository
	scanossSettingsRepo repository.ScanossSettingsRepository
	resultRepo          repository.ResultRepository
	scanossApiService   ScanossApiService
	mapper              mappers.ComponentMapper
	initialFilters      []entities.ComponentFilterDTO
	undoStack           [][]entities.ComponentFilterDTO
	redoStack           [][]entities.ComponentFilterDTO
}

func NewComponentServiceImpl(repo repository.ComponentRepository, scanossSettingsRepo repository.ScanossSettingsRepository, resultRepo repository.ResultRepository, scanossApiService ScanossApiService, mapper mappers.ComponentMapper) ComponentService {
	service := &ComponentServiceImpl{
		repo:                repo,
		scanossSettingsRepo: scanossSettingsRepo,
		resultRepo:          resultRepo,
		scanossApiService:   scanossApiService,
		initialFilters:      []entities.ComponentFilterDTO{},
		undoStack:           [][]entities.ComponentFilterDTO{},
		redoStack:           [][]entities.ComponentFilterDTO{},
		mapper:              mapper,
	}

	service.setInitialFilters()

	return service
}

func (s *ComponentServiceImpl) setInitialFilters() {
	initialFilters := s.GetInitialFilters()

	for _, include := range initialFilters.Include {
		s.initialFilters = append(s.initialFilters, entities.ComponentFilterDTO{
			Path:   include.Path,
			Purl:   include.Purl,
			Usage:  string(include.Usage),
			Action: entities.Include,
		})
	}
	for _, remove := range initialFilters.Remove {
		s.initialFilters = append(s.initialFilters, entities.ComponentFilterDTO{
			Path:   remove.Path,
			Purl:   remove.Purl,
			Usage:  string(remove.Usage),
			Action: entities.Remove,
		})
	}
	for _, replace := range initialFilters.Replace {
		s.initialFilters = append(s.initialFilters, entities.ComponentFilterDTO{
			Path:   replace.Path,
			Purl:   replace.Purl,
			Usage:  string(replace.Usage),
			Action: entities.Replace,
		})
	}
}

func (s *ComponentServiceImpl) ClearAllFilters() error {
	return s.scanossSettingsRepo.ClearAllFilters()
}

func (s *ComponentServiceImpl) GetInitialFilters() entities.InitialFilters {
	sf := s.scanossSettingsRepo.GetSettings()
	include, remove, replace := sf.Bom.Include, sf.Bom.Remove, sf.Bom.Replace

	return entities.InitialFilters{
		Include: include,
		Remove:  remove,
		Replace: replace,
	}
}

func (s *ComponentServiceImpl) GetComponentByPath(filePath string) (entities.ComponentDTO, error) {
	c, err := s.repo.FindByFilePath(filePath)
	if err != nil {
		return entities.ComponentDTO{}, err
	}
	dto := s.mapper.MapToComponentDTO(c)

	return dto, nil
}

func (s *ComponentServiceImpl) FilterComponents(dto []entities.ComponentFilterDTO) error {
	for _, filter := range dto {
		err := utils.GetValidator().Struct(filter)
		if err != nil {
			log.Error().Err(err).Msg("Validation error")
			return err
		}
	}

	s.applyFilters(dto)

	s.undoStack = append(s.undoStack, dto)
	s.redoStack = nil

	return nil
}

func (s *ComponentServiceImpl) applyFilters(dto []entities.ComponentFilterDTO) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(dto))

	for _, item := range dto {
		wg.Add(1)
		go func(item entities.ComponentFilterDTO) {
			defer wg.Done()
			newFilter := entities.ComponentFilter{
				Path:        item.Path,
				Purl:        item.Purl,
				Usage:       entities.ComponentFilterUsage(item.Usage),
				Comment:     item.Comment,
				ReplaceWith: item.ReplaceWith,
				License:     item.License,
			}
			if err := s.scanossSettingsRepo.AddBomEntry(newFilter, string(item.Action)); err != nil {
				log.Error().Err(err).Msg("Error adding bom entry")
				errChan <- err
			}
		}(item)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *ComponentServiceImpl) GetDeclaredComponents() ([]entities.DeclaredComponent, error) {
	results, err := s.resultRepo.GetResults(nil)
	if err != nil {
		return []entities.DeclaredComponent{}, err
	}

	scanossSettingsDeclaredPurls := s.scanossSettingsRepo.GetDeclaredPurls()

	declaredComponents := make([]entities.DeclaredComponent, 0)
	addedPurls := make(map[string]struct{})

	for _, result := range results {
		if result.Purl == nil || len(*result.Purl) == 0 {
			continue
		}
		purl := (*result.Purl)[0]
		if _, found := addedPurls[purl]; !found {
			declaredComponents = append(declaredComponents, entities.DeclaredComponent{
				Purl: purl,
				Name: result.ComponentName,
			})
			addedPurls[purl] = struct{}{}
		}
	}

	for _, purl := range scanossSettingsDeclaredPurls {
		if _, found := addedPurls[purl]; !found {
			purlName, err := purlutils.PurlNameFromString(purl)
			if err != nil {
				log.Error().Err(err).Msg("Error getting component name from purl")
				continue
			}
			nameParts := strings.Split(purlName, "/")
			name := nameParts[0]
			if len(nameParts) > 1 {
				name = nameParts[1]
			}
			declaredComponents = append(declaredComponents, entities.DeclaredComponent{
				Purl: purl,
				Name: name,
			})
		}
	}

	sort.Slice(declaredComponents, func(i, j int) bool {
		return strings.ToLower(declaredComponents[i].Name) < strings.ToLower(declaredComponents[j].Name)
	})

	return declaredComponents, nil
}

func (s *ComponentServiceImpl) Undo() error {
	if !s.CanUndo() {
		return nil
	}

	lastAction := s.undoStack[len(s.undoStack)-1]
	s.undoStack = s.undoStack[:len(s.undoStack)-1]
	s.redoStack = append(s.redoStack, lastAction)

	return s.reapplyActions()
}

func (s *ComponentServiceImpl) Redo() error {
	if !s.CanRedo() {
		return nil
	}

	nextAction := s.redoStack[len(s.redoStack)-1]
	s.redoStack = s.redoStack[:len(s.redoStack)-1]
	s.undoStack = append(s.undoStack, nextAction)

	err := s.applyFilters(nextAction)
	if err != nil {
		return err
	}
	return nil
}

func (s *ComponentServiceImpl) CanUndo() bool {
	return len(s.undoStack) > 0
}

func (s *ComponentServiceImpl) CanRedo() bool {
	return len(s.redoStack) > 0
}

func (s *ComponentServiceImpl) reapplyActions() error {
	// We need to clear bom filters so workflow state is properly resetted
	err := s.ClearAllFilters()
	if err != nil {
		return err
	}

	for _, initialFilterDto := range s.initialFilters {
		err := s.applyFilters([]entities.ComponentFilterDTO{initialFilterDto})
		if err != nil {
			return err
		}
	}

	for _, filterDto := range s.undoStack {
		err := s.applyFilters(filterDto)
		if err != nil {
			return err
		}
	}

	return nil
}

// Searches for components using the SCANOSS API
func (s *ComponentServiceImpl) SearchComponents(request entities.ComponentSearchRequest) (entities.ComponentSearchResponse, error) {
	log.Debug().
		Str("search", request.Search).
		Str("vendor", request.Vendor).
		Str("component", request.Component).
		Msg("Delegating component search to SCANOSS API service")

	return s.scanossApiService.SearchComponents(request)
}
