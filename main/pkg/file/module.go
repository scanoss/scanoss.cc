package file

import (
	"integration-git/main/pkg/file/adapter"
	"integration-git/main/pkg/file/application/services"
	"integration-git/main/pkg/file/infrastructure"
)

func init() {
	adapter.NewFileController(services.NewFileService(infrastructure.NewFileRepository()))
}
