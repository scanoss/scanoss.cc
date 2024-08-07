package adapter

import (
	"integration-git/main/pkg/file/application/services"
	"integration-git/main/pkg/file/application/use_cases"
)

type Controller struct {
	getRemoteFileUseCase *use_cases.GetRemoteFileUseCase
}

func NewFileController(service *services.FileService) *Controller {
	return &Controller{
		getRemoteFileUseCase: use_cases.NewGetRemoteFileUseCase(service),
	}
}

func (c *Controller) GetRemoteFile(path string) (FileDTO, error) {
	//c.getRemoteFileUseCase.ReadFile(path)
	return FileDTO{
		Name:    "main.c",
		Path:    "src/main.c",
		Content: "#include<stdio.h>\n void main() {\nreturn 0;}\n",
	}, nil
}

func (c *Controller) GetLocalFile() (FileDTO, error) {
	return FileDTO{
		Name:    "main.c",
		Path:    "src/main.c",
		Content: "#include<stdio.h>\n void main() {\nreturn 0;}\n",
	}, nil
}
