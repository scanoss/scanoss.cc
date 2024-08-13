package infrastructure

import (
	"integration-git/main/pkg/git/domain"
	"reflect"
	"testing"
)

var tests = map[string]struct {
	input       string
	getExpected func() []domain.File
}{

	"staged file": {
		input: " M src/scanner.c\n",
		getExpected: func() []domain.File {
			M, _ := domain.NewFileStatus("M")
			S, _ := domain.NewFileStatus(" ")

			f1, _ := domain.NewFile("scanner.c", "src/scanner.c", S, M)
			return []domain.File{*f1}
		},
	},
	"modified files, staged files and untracked": {
		input: "M  src/main.c\n" +
			" M src/scanner.c\n" +
			"?? new_file.c\n",
		getExpected: func() []domain.File {
			M, _ := domain.NewFileStatus("M")
			Q, _ := domain.NewFileStatus("?")
			S, _ := domain.NewFileStatus(" ")

			f1, _ := domain.NewFile("main.c", "src/main.c", M, S)
			f2, _ := domain.NewFile("scanner.c", "src/scanner.c", S, M)
			f3, _ := domain.NewFile("new_file.c", "new_file.c", Q, Q)
			return []domain.File{*f1, *f2, *f3}
		},
	},
}

func TestGitCliRepository_GetFileStatusFromOutputGitCli(t *testing.T) {
	gitRepository := NewGitRepository()

	for name, test := range tests {
		test := test // NOTE: uncomment for Go < 1.22, see /doc/faq#closures_and_goroutines
		t.Run(name, func(t *testing.T) {

			got, err := gitRepository.GetFileStatusFromOutputGitCli(test.input)
			expected := test.getExpected()

			if err != nil {
				t.Fatalf("error while calling GetFileStatusFromOutputGitCli for %s - %v", test.input, err)
				return
			}

			if reflect.DeepEqual(got, expected) == false {

				t.Fatalf("\ngot: %+v; \nexpected %+v\n", got, expected)
			}
		})
	}

}
