package filesystem

import (
	"os"
	"path/filepath"
	"testing"
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

	fs := GetFilesystem("file://" + tmpDir)

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
	tmpDir, err := os.MkdirTemp("", "deployto-testgetfilesystem*")
	checkIfError(t, err)
	defer func() {
		err := os.RemoveAll(tmpDir)
		checkIfError(t, err)
	}()

	//Check no git
	tmpFS := GetFilesystem("file://" + tmpDir)
	result := GetGitRootFilesystem(tmpFS, "/")
	if result != nil {
		t.Errorf("if git is not initialized, I expect the nil result")
	}

	//Check root git
	err = os.Mkdir(filepath.Join(tmpDir, ".git"), 0700)
	checkIfError(t, err)
	result = GetGitRootFilesystem(tmpFS, "/")
	if result == nil || result.URI != tmpFS.URI {
		t.Errorf("wait same BaseDir, get: %s", result.URI)
	}

	//Check root git path from sub sub sub path
	subsubsubPath := filepath.Join("A", "B", "C")
	err = os.MkdirAll(filepath.Join(tmpDir, subsubsubPath), 0700)
	checkIfError(t, err)
	result = GetGitRootFilesystem(tmpFS, subsubsubPath)
	if result == nil || result.URI != tmpFS.URI {
		t.Errorf("wait same BaseDir, get: %s", result.URI)
	}

	// Check git in git
	err = os.MkdirAll(filepath.Join(tmpDir, "A", ".git"), 0700)
	checkIfError(t, err)
	result = GetGitRootFilesystem(tmpFS, subsubsubPath)
	if result == nil || result.URI != "file://"+filepath.Join(tmpDir, "A") {
		t.Errorf("wait same BaseDir, get: %s", result.URI)
	}

	//Check root git path from sub sub sub filesystem
	subsubsubFS := GetFilesystem("file://" + filepath.Join(tmpDir, subsubsubPath))
	result = GetGitRootFilesystem(subsubsubFS, "/")
	if result == nil || result.URI != "file://"+filepath.Join(tmpDir, "A") {
		t.Errorf("wait same BaseDir, get: %s", result.URI)
	}
}
