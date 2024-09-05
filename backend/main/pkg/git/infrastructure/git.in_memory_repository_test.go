package infrastructure

import (
	"integration-git/main/pkg/git/domain"
	"reflect"
	"testing"
)

func TestGitCliRepository_GetFiles(t *testing.T) {

	M, _ := domain.NewFileStatus("M")
	D, _ := domain.NewFileStatus("D")
	Q, _ := domain.NewFileStatus("?")

	f1, _ := domain.NewFile("main.c", "src/main.c", M, D)     //Modified in staging area and removed on working directory
	f2, _ := domain.NewFile("crypto.c", "src/crypto.c", M, M) //Modified in staging area and modified on working directory
	f3, _ := domain.NewFile("crypto.h", "src/crypto.h", Q, Q) //Untracked

	gr := NewGitInMemoryRepository()
	gr.AddFile(*f1)
	gr.AddFile(*f2)
	gr.AddFile(*f3)

	files, _ := gr.GetFiles()

	if !reflect.DeepEqual([]domain.File{*f1, *f2, *f3}, files) {
		t.Fatalf("expected to return same files")
	}
}
