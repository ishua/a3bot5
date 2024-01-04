package repo_test

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/ishua/a3bot5/fsnotes/internal/repo"
	"github.com/ishua/a3bot5/fsnotes/internal/repo/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testFile = "gitfile_test_readrows.txt"
var tempTestFile = "gitfile_test_readrows_temp.txt"

func TestGitFile_ReadRows(t *testing.T) {
	tests := []struct {
		name           string
		git            *mocks.Gitter
		expectedRows   []string
		expectedErrStr string
	}{
		{
			name: "Successful ReadRows",
			git: func() *mocks.Gitter {
				git := new(mocks.Gitter)
				git.On("Pull", mock.Anything).Return(nil)
				return git
			}(),
			expectedRows:   []string{"line1", "line2", "line3"},
			expectedErrStr: "",
		},
		{
			name: "Error on Git Pull",
			git: func() *mocks.Gitter {
				git := new(mocks.Gitter)
				git.On("Pull", mock.Anything).Return(errors.New("pull error"))
				return git
			}(),
			expectedRows:   nil,
			expectedErrStr: "gitfile ReadRows pull: pull error",
		},
		// Add more test cases as needed.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := repo.NewGitFile(testFile, tt.git)
			rows, err := g.ReadRows(context.Background())
			var errStr string
			if err != nil {
				errStr = err.Error()
			}

			assert.Equal(t, tt.expectedRows, rows)
			assert.Equal(t, tt.expectedErrStr, errStr)

			tt.git.AssertExpectations(t)
		})
	}
}

func TestGitFile_ReadRows_noFile(t *testing.T) {
	git := func() *mocks.Gitter {
		git := new(mocks.Gitter)
		git.On("Pull", mock.Anything).Return(nil)
		return git
	}()
	g := repo.NewGitFile("no_such_file.txt", git)
	rows, err := g.ReadRows(context.Background())
	assert.Empty(t, rows)
	assert.NotEmpty(t, err)

}

func setupTestCase(t *testing.T) {
	t.Log("Create copy test file")
	sourceFileStat, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("catn't setup case: %s", err.Error())
	}
	if !sourceFileStat.Mode().IsRegular() {
		t.Fatalf("is not a regular file %s", testFile)
	}

	source, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("is not a regular file %s", testFile)
	}
	defer source.Close()

	destination, err := os.Create(tempTestFile)
	if err != nil {
		t.Fatalf("can't create file %s", tempTestFile)
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		cleanTestCase(t)
		t.Fatalf("can't copy file %s", tempTestFile)
	}
}

func cleanTestCase(t *testing.T) {
	t.Log("Remove test file")
	os.Remove(tempTestFile)
}

func TestGitFile_AddRows(t *testing.T) {
	setupTestCase(t)
	defer cleanTestCase(t)

	git := func() *mocks.Gitter {
		git := new(mocks.Gitter)
		git.On("Pull", mock.Anything).Return(nil)
		git.On("CommitAndPush", mock.Anything, mock.Anything).Return(nil)
		return git
	}()
	g := repo.NewGitFile(tempTestFile, git)
	g.AddRows(context.Background(), []string{"line4", "line5"})
	rows, err := g.ReadRows(context.Background())
	if err != nil {
		t.Fatalf("%s can't read new lines in file %s", t.Name(), tempTestFile)
	}
	assert.Equal(t, rows, []string{"line1", "line2", "line3", "line4", "line5"})
}
