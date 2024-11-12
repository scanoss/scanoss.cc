package service

import (
	"sort"
	"strings"

	"github.com/labstack/gommon/log"
	"github.com/scanoss/scanoss.lui/backend/main/entities"
	"github.com/scanoss/scanoss.lui/backend/main/mappers"
	"github.com/scanoss/scanoss.lui/backend/main/repository"
	"github.com/scanoss/scanoss.lui/backend/main/utils"
)

type ResultServiceImpl struct {
	repo   repository.ResultRepository
	mapper mappers.ResultMapper
}

func NewResultServiceImpl(repo repository.ResultRepository, mapper mappers.ResultMapper) ResultService {
	return &ResultServiceImpl{
		repo:   repo,
		mapper: mapper,
	}
}

func (s *ResultServiceImpl) GetAll(dto *entities.RequestResultDTO) ([]entities.ResultDTO, error) {
	err := utils.GetValidator().Struct(dto)
	if err != nil {
		log.Errorf("Validation error: %v", err)
		return []entities.ResultDTO{}, err
	}

	filter := entities.NewResultFilterFactory().Create(dto)
	results, err := s.repo.GetResults(filter)
	if err != nil {
		return []entities.ResultDTO{}, err
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Path < results[j].Path
	})

	return s.mapper.MapToResultDTOList(results), nil
}

func (s *ResultServiceImpl) GetAllTree(dto *entities.RequestResultDTO) ([]entities.ResultTreeDTO, error) {
	err := utils.GetValidator().Struct(dto)
	if err != nil {
		log.Errorf("Validation error: %v", err)
		return []entities.ResultTreeDTO{}, err
	}

	filter := entities.NewResultFilterFactory().Create(dto)
	results, err := s.repo.GetResults(filter)
	if err != nil {
		return []entities.ResultTreeDTO{}, err
	}

	tree := s.buildResultTree(results)
	return tree, nil
}

func (s *ResultServiceImpl) buildResultTree(results []entities.Result) []entities.ResultTreeDTO {
	treeMap := make(map[string]*entities.ResultTreeDTO)

	for _, result := range results {
		parts := strings.Split(result.Path, "/")
		var currentNode *entities.ResultTreeDTO
		var parent *string

		for i, part := range parts {
			if i == 0 {
				if node, exists := treeMap[part]; exists {
					currentNode = node
				} else {
					currentNode = &entities.ResultTreeDTO{ID: part, Name: part, Parent: nil}
					treeMap[part] = currentNode
				}
			} else {
				if node, exists := currentNode.Children[part]; exists {
					currentNode = node
					parent = &currentNode.Name
				} else {
					newNode := &entities.ResultTreeDTO{ID: part, Name: part, Parent: &currentNode.Name}
					if currentNode.Children == nil {
						currentNode.Children = make(map[string]*entities.ResultTreeDTO)
					}
					currentNode.Children[part] = newNode
					currentNode = newNode
				}
			}
		}
		currentNode.Result = s.mapper.MapToResultDTO(result)
		currentNode.Parent = parent
	}

	var tree []entities.ResultTreeDTO
	for _, node := range treeMap {
		tree = append(tree, *node)
	}

	return tree
}
