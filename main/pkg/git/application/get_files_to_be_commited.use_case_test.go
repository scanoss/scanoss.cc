package application

import (
	"integration-git/main/pkg/git/domain"
	"integration-git/main/pkg/git/infrastructure"
	"reflect"
	"testing"
)

func TestGetFilesToBeCommitedUseCase_Execute(t *testing.T) {
	M, _ := domain.NewFileStatus("M")
	D, _ := domain.NewFileStatus("D")
	Q, _ := domain.NewFileStatus("?")

	f1, _ := domain.NewFile("main.c", "src/main.c", M, D)     //Modified in staging area and removed on working directory
	f2, _ := domain.NewFile("crypto.c", "src/crypto.c", M, M) //Modified in staging area and modified on working directory
	f3, _ := domain.NewFile("crypto.h", "src/crypto.h", Q, Q) //Untracked

	gr := infrastructure.NewGitInMemoryRepository()
	gr.AddFile(*f1)
	gr.AddFile(*f2)
	gr.AddFile(*f3)

	gs := NewGitService(gr)
	got, err := NewGetFilesToBeCommitedUseCase(gs).Execute()

	if err != nil {
		t.Fatalf("error while executing NewGetFilesToBeCommitedUseCase - %v", err)
	}

	t.Run("returns files to be commited", func(t *testing.T) {
		expected := []domain.File{*f1, *f2, *f3}

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("got: %+v\nwant: %+v", got, expected)
		}
	})

}
