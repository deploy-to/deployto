package gitclient

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestGetValues(t *testing.T) {
	//Create git repo
	tmpDir, err := os.MkdirTemp("", "deployto-unittests*")
	if err != nil {
		t.FailNow()
	}
	t.Logf("tmp dir: %s" + tmpDir)
	defer func() {
		os.RemoveAll(tmpDir)
	}()

	//git just init
	r, err := git.PlainInit(tmpDir, false)
	checkIfError(t, err)
	w, err := r.Worktree()
	checkIfError(t, err)

	output := GetValues(tmpDir)
	if len(output) != 2 {
		t.Errorf("git just init: the output does not contain 2 elements: %v", output)
	}
	if !strings.HasPrefix(output["Commit"].(string), "NO_GIT") {
		t.Errorf("git just init: prefix error: GetValues()[Commit] = %v, want %v", output, "NO_GIT")
	}
	if !strings.HasPrefix(output["CommitShort"].(string), "NO_GIT") {
		t.Errorf("git just init: prefix error: GetValues()[CommitShort] = %v, want %v", output, "NO_GIT")
	}
	if !strings.Contains(output["Commit"].(string), "+dirty.uuid") {
		t.Errorf("git just init: dirty mark not exists: GetValues()[CommitShort] = %v", output)
	}
	if !strings.Contains(output["CommitShort"].(string), "+dirty.uuid") {
		t.Errorf("git just init: dirty mark not exists: GetValues()[CommitShort] = %v", output)
	}

	//git first commit
	doChange(t, tmpDir)
	_, err = w.Add("data.txt")
	checkIfError(t, err)
	commit, err := w.Commit("example go-git commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "John Doe",
			Email: "john@doe.org",
			When:  time.Now(),
		},
	})
	checkIfError(t, err)
	// TODO check output["Tag"]
	_, err = setTag(r, "v1.0.0", &object.Signature{
		Name:  "John Doe",
		Email: "john@doe.org",
		When:  time.Now(),
	})
	checkIfError(t, err)
	_, err = setTag(r, "v1.0.1", &object.Signature{
		Name:  "John Doe",
		Email: "john@doe.org",
		When:  time.Now(),
	})
	checkIfError(t, err)
	output = GetValues(tmpDir)
	if len(output) != 3 {
		t.Errorf("git first commit: the output does not contain 2 elements: %v", output)
	}
	if output["Commit"] != commit.String() {
		t.Errorf("git first commit: GetValues()[Commit] = %v, want %v", output, commit.String())
	}
	if output["CommitShort"] != commit.String()[:7] {
		t.Errorf("git first commit: GetValues()[Commit] = %v, want %v", output, commit.String()[:7])
	}
	if output["Tag"] != "v1.0.1" {
		t.Errorf("git first commit: GetValues()[Tag] = %v, want %v", output, "v1.0.1")
	}

	//dirty git
	doChange(t, tmpDir)
	_, err = setTag(r, "v1.0.0", &object.Signature{
		Name:  "John Doe",
		Email: "john@doe.org",
		When:  time.Now(),
	})
	checkIfError(t, err)
	_, err = setTag(r, "v1.0.1", &object.Signature{
		Name:  "John Doe",
		Email: "john@doe.org",
		When:  time.Now(),
	})
	checkIfError(t, err)
	output = GetValues(tmpDir)
	if len(output) != 3 {
		t.Errorf("dirty git: the output does not contain 3 elements: %v", output)
	}
	if !strings.HasPrefix(output["Commit"].(string), commit.String()) {
		t.Errorf("dirty git: prefix error: GetValues()[Commit] = %v, want %v", output, commit.String())
	}
	if !strings.HasPrefix(output["CommitShort"].(string), commit.String()[:7]) {
		t.Errorf("dirty git: prefix error: GetValues()[CommitShort] = %v, want %v", output, commit.String()[:7])
	}
	if !strings.Contains(output["Commit"].(string), "+dirty.uuid") {
		t.Errorf("dirty git: dirty mark not exists: GetValues()[CommitShort] = %v", output)
	}
	if !strings.Contains(output["CommitShort"].(string), "+dirty.uuid") {
		t.Errorf("dirty git: dirty mark not exists: GetValues()[CommitShort] = %v", output)
	}
	if !strings.Contains(output["Tag"].(string), "+dirty.uuid") {
		t.Errorf("dirty git: dirty mark not exists: GetValues()[Tag] = %v", output)
	}
	if !strings.HasPrefix(output["Tag"].(string), "v1.0.1") {
		t.Errorf("dirty git: prefix error: GetValues()[Tag] = %v, want %v", output, "v1.0.1")
	}

}

func doChange(t *testing.T, tmpDir string) {
	f, err := os.OpenFile(filepath.Join(tmpDir, "data.txt"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		t.FailNow()
	}
	defer f.Close()
	_, err = f.WriteString("new line")
	if err != nil {
		t.FailNow()
	}
}

func checkIfError(t *testing.T, err error) {
	if err != nil {
		t.FailNow()
	}
}

func setTag(r *git.Repository, tag string, tagger *object.Signature) (bool, error) {
	if tagExists(tag, r) {
		return false, nil
	}
	h, err := r.Head()
	if err != nil {
		return false, err
	}
	_, err = r.CreateTag(tag, h.Hash(), &git.CreateTagOptions{
		Tagger:  tagger,
		Message: tag,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func tagExists(tag string, r *git.Repository) bool {
	tagFoundErr := "tag was found"
	tags, err := r.TagObjects()
	if err != nil {
		return false
	}
	res := false
	err = tags.ForEach(func(t *object.Tag) error {
		if t.Name == tag {
			res = true
			return fmt.Errorf(tagFoundErr)
		}
		return nil
	})
	if err != nil && err.Error() != tagFoundErr {
		return false
	}
	return res
}
