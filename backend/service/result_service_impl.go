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
	"context"
	"sort"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/backend/mappers"
	"github.com/scanoss/scanoss.cc/backend/repository"
	"github.com/scanoss/scanoss.cc/internal/utils"
)

type ResultServiceImpl struct {
	ctx     context.Context
	watcher *fsnotify.Watcher
	repo    repository.ResultRepository
	mapper  mappers.ResultMapper
}

func NewResultServiceImpl(repo repository.ResultRepository, mapper mappers.ResultMapper) ResultService {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error().Err(err).Msg("Error creating watcher: %v")
	}

	return &ResultServiceImpl{
		repo:    repo,
		mapper:  mapper,
		watcher: watcher,
	}
}

func (s *ResultServiceImpl) GetAll(dto *entities.RequestResultDTO) ([]entities.ResultDTO, error) {
	err := utils.GetValidator().Struct(dto)
	if err != nil {
		log.Error().Err(err).Msg("Validation error: %v")
		return []entities.ResultDTO{}, err
	}

	filter := entities.NewResultFilterFactory().Create(dto)
	results, err := s.repo.GetResults(filter)
	if err != nil {
		return []entities.ResultDTO{}, err
	}

	s.sortResults(results, dto)

	return s.mapper.MapToResultDTOList(results), nil
}

func (s *ResultServiceImpl) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *ResultServiceImpl) sortResults(results []entities.Result, dto *entities.RequestResultDTO) {
	sort.Slice(results, func(i, j int) bool {
		firstResult := results[i]
		secondResult := results[j]

		// Default values if Sort is not specified
		sortOption := entities.SortByMatchPercentage
		sortOrder := entities.SortOrderDesc
		if dto.Sort.Option != "" {
			sortOption = dto.Sort.Option
		}
		if dto.Sort.Order != "" {
			sortOrder = dto.Sort.Order
		}

		switch sortOption {
		case entities.SortByMatchPercentage:
			iPercentage := firstResult.GetMatchPercentage()
			jPercentage := secondResult.GetMatchPercentage()
			if iPercentage != jPercentage {
				return sortOrder == entities.SortOrderDesc && iPercentage > jPercentage ||
					sortOrder == entities.SortOrderAsc && iPercentage < jPercentage
			}
			return firstResult.Path < secondResult.Path // Secondary sort by path
		case entities.SortByPath:
			return (sortOrder == entities.SortOrderAsc) == (firstResult.Path < secondResult.Path)
		default:
			// Default to match percentage sorting
			iPercentage := firstResult.GetMatchPercentage()
			jPercentage := secondResult.GetMatchPercentage()
			if iPercentage != jPercentage {
				return sortOrder == entities.SortOrderDesc && iPercentage > jPercentage ||
					sortOrder == entities.SortOrderAsc && iPercentage < jPercentage
			}
			return firstResult.Path < secondResult.Path
		}
	})
}
