package adapter

import "integration-git/main/pkg/component/domain"

func adaptToComponentDTO(componentEntity domain.Component) ComponentDTO {

	dto := ComponentDTO{
		ID:          componentEntity.ID,
		URL:         componentEntity.URL,
		Status:      componentEntity.Status,
		File:        componentEntity.File,
		Component:   componentEntity.Component,
		FileHash:    componentEntity.FileHash,
		Lines:       componentEntity.Lines,
		Purl:        componentEntity.Purl,
		FileURL:     componentEntity.FileURL,
		OssLines:    componentEntity.OssLines,
		Matched:     componentEntity.Matched,
		URLStats:    componentEntity.URLStats,
		Server:      componentEntity.Server,
		SourceHash:  componentEntity.SourceHash,
		URLHash:     componentEntity.URLHash,
		Vendor:      componentEntity.Vendor,
		Version:     componentEntity.Version,
		ReleaseDate: componentEntity.ReleaseDate,
		Provenance:  componentEntity.Provenance,
		Latest:      componentEntity.Latest,
	}

	var licenses []License
	if len(componentEntity.Licenses) > 0 {
		for _, l := range componentEntity.Licenses {
			// patentHints ,copyleft0 checklistURL string, osadlUpdated string, source string, url string, incompatibleWith string
			licenses = append(licenses, NewLicenseDTO(l.Name, l.PatentHints, l.Copyleft, l.ChecklistURL, l.OsadlUpdated, l.Source, l.URL, l.IncompatibleWith))
		}
		dto.Licenses = licenses
	}

	return dto

}
