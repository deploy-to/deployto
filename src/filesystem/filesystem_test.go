package filesystem

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-test/deep"
)

func TestGetFilesystem(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "deployto-testgetfilesystem*")
	checkIfError(t, err)
	defer func() {
		err := os.RemoveAll(tmpDir)
		checkIfError(t, err)
	}()
	err = os.WriteFile(filepath.Join(tmpDir, "data.txt"), []byte("test"), 0644)
	checkIfError(t, err)

	fs := Get("file://" + tmpDir)

	if fs.URI != "file://"+tmpDir {
		t.Errorf("BaseDir has changed: %s", fs.URI)
	}

	f, err := fs.FS.Open("data.txt")
	checkIfError(t, err)
	defer f.Close()
	data := make([]byte, 512)
	n, err := f.Read(data)
	data = data[:n]
	checkIfError(t, err)
	if string(data) != "test" {
		t.Errorf("the read data does not match the reference data: %s", string(data))
	}
}

func checkIfError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGetGitRootFilesystem(t *testing.T) {
	tmpFS := Get("temp")
	defer tmpFS.Destroy()

	//Check no git
	result := GetGitRootFilesystem(tmpFS, "/")
	if result != nil {
		t.Errorf("if git is not initialized, I expect the nil result")
	}

	//Check root git
	err := os.Mkdir(filepath.Join(tmpFS.LocalPath, ".git"), 0700)
	checkIfError(t, err)
	result = GetGitRootFilesystem(tmpFS, "/")
	if result == nil || result.URI != tmpFS.URI {
		t.Errorf("wait same BaseDir, get: %s", result.URI)
	}

	//Check root git path from sub sub sub path
	subsubsubPath := filepath.Join("A", "B", "C")
	err = os.MkdirAll(filepath.Join(tmpFS.LocalPath, subsubsubPath), 0700)
	checkIfError(t, err)
	result = GetGitRootFilesystem(tmpFS, subsubsubPath)
	if result == nil || result.URI != tmpFS.URI {
		t.Errorf("wait same BaseDir, get: %s", result.URI)
	}

	// Check git in git
	err = os.MkdirAll(filepath.Join(tmpFS.LocalPath, "A", ".git"), 0700)
	checkIfError(t, err)
	result = GetGitRootFilesystem(tmpFS, subsubsubPath)
	if result == nil || result.URI != "file://"+filepath.Join(tmpFS.LocalPath, "A") {
		t.Errorf("wait same BaseDir, get: %s", result.URI)
	}

	//Check root git path from sub sub sub filesystem
	subsubsubFS := Get("file://" + filepath.Join(tmpFS.LocalPath, subsubsubPath))
	result = GetGitRootFilesystem(subsubsubFS, "/")
	if result == nil || result.URI != "file://"+filepath.Join(tmpFS.LocalPath, "A") {
		t.Errorf("wait same BaseDir, get: %s", result.URI)
	}
}

func TestFilesystem_AsValues(t *testing.T) {
	tests := []struct {
		name         string
		asFilesystem *Filesystem
		asValues     map[string]any
	}{
		{
			name: "filesystem",
			asFilesystem: &Filesystem{
				URI:       "testURI",
				LocalPath: "testLocal",
				Type:      LOCAL,
				FS:        memfs.New(),
			},
			asValues: map[string]any{
				"URI":       "testURI",
				"LocalPath": "testLocal",
				"Type":      "LOCAL",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deep.NilSlicesAreEmpty = true
			deep.NilMapsAreEmpty = true

			gotValues := tt.asFilesystem.AsValues()
			if diff := deep.Equal(gotValues, tt.asValues); diff != nil {
				t.Error(strings.Join(diff, "; "))
			}
		})
	}
}
