package gitclient

import (
	"deployto/src/filesystem"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestGetValues_GeneralChecks(t *testing.T) {
	tmpDir, fs, r, w := prepareGit(t)
	defer func() {
		os.RemoveAll(tmpDir)
	}()

	//just init
	output := GetValues(fs, "/")
	if len(output) != 2 {
		t.Errorf("git just init: the output does not contain 2 elements: %v", output)
	}
	if !strings.HasPrefix(output["Commit"].(string), "NO_GIT") {
		t.Errorf("git just init: prefix error: GetValues()[Commit] = %v, want %v", output, "NO_GIT")
	}
	if !strings.HasPrefix(output["CommitShort"].(string), "NO_GIT") {
		t.Errorf("git just init: prefix error: GetValues()[CommitShort] = %v, want %v", output, "NO_GIT")
	}
	if !strings.Contains(output["Commit"].(string), "-dirty.uuid") {
		t.Errorf("git just init: dirty mark not exists: GetValues()[CommitShort] = %v", output)
	}
	if !strings.Contains(output["CommitShort"].(string), "-dirty.uuid") {
		t.Errorf("git just init: dirty mark not exists: GetValues()[CommitShort] = %v", output)
	}

	//git first commit
	doChange(t, tmpDir)
	commit := doCommit(t, w)
	setTag(t, r, "v1.0.1")

	output = GetValues(fs, "/")
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

	output = GetValues(fs, "/")
	if len(output) != 3 {
		t.Errorf("dirty git: the output does not contain 3 elements: %v", output)
	}
	if !strings.HasPrefix(output["Commit"].(string), commit.String()) {
		t.Errorf("dirty git: prefix error: GetValues()[Commit] = %v, want %v", output, commit.String())
	}
	if !strings.HasPrefix(output["CommitShort"].(string), commit.String()[:7]) {
		t.Errorf("dirty git: prefix error: GetValues()[CommitShort] = %v, want %v", output, commit.String()[:7])
	}
	if !strings.HasPrefix(output["Tag"].(string), "v1.0.1") {
		t.Errorf("dirty git: prefix error: GetValues()[Tag] = %v, want %v", output, "v1.0.1")
	}
	if !strings.Contains(output["Commit"].(string), "-dirty.uuid") {
		t.Errorf("dirty git: dirty mark not exists: GetValues()[CommitShort] = %v", output)
	}
	if !strings.Contains(output["CommitShort"].(string), "-dirty.uuid") {
		t.Errorf("dirty git: dirty mark not exists: GetValues()[CommitShort] = %v", output)
	}
	if !strings.Contains(output["Tag"].(string), "-dirty.uuid") {
		t.Errorf("dirty git: dirty mark not exists: GetValues()[Tag] = %v", output)
	}

}

func TestGetValues_GetCurrentTag(t *testing.T) {
	tmpDir, fs, r, w := prepareGit(t)
	defer func() {
		os.RemoveAll(tmpDir)
	}()

	//oldest commit
	doChange(t, tmpDir)
	doCommit(t, w)
	setTag(t, r, "v6.6.6") // Previously, the project used a different versioning. But that was a long time ago and shouldn't affect current commits.

	//corrent release
	doChange(t, tmpDir)
	commit := doCommit(t, w)
	setTag(t, r, "v1.1.1")

	// // Делаем потом, не в текущей итерации
	// // TODO если у текущего коммита нет тега, а у родительского есть, то надо возвращать родительский, отмечая его +dirtyCommitXXXX
	// // corrent release
	// doChange(t, tmpDir)
	// commit := doCommit(t, w)

	//next release (maybe in another branch?)
	doChange(t, tmpDir)
	doCommit(t, w)
	setTag(t, r, "v2.2.2-RC2")

	// checkout to corrent release
	err := w.Checkout(&git.CheckoutOptions{Hash: commit})
	checkIfError(t, err)

	output := GetValues(fs, tmpDir)
	if !strings.HasPrefix(output["Tag"].(string), "v1.1.1") {
		t.Errorf("dirty git: prefix error: GetValues() = %v, want Tag: v1.1.1", output)
	}

	// many tag on the same commit
	doChange(t, tmpDir)
	doCommit(t, w)

	setTag(t, r, "v1.0.2-rc")
	//add tag after test done
	setTag(t, r, "v1.0.2")

	output = GetValues(fs, tmpDir)
	if !strings.HasPrefix(output["Tag"].(string), "v1.0.2") {
		t.Errorf("dirty git: prefix error: GetValues() = %v, want Tag: v1.0.2", output)
	}
}

func prepareGit(t *testing.T) (string, *filesystem.Filesystem, *git.Repository, *git.Worktree) {
	//Create git repo
	tmpDir, err := os.MkdirTemp("", "deployto-testgetvalues*")
	checkIfError(t, err)
	t.Logf("tmp dir: %s", tmpDir)

	fs := filesystem.Get("file://" + tmpDir)

	//git init
	r, err := git.PlainInit(tmpDir, false)
	checkIfError(t, err)
	w, err := r.Worktree()
	checkIfError(t, err)

	return tmpDir, fs, r, w
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

func doCommit(t *testing.T, w *git.Worktree) plumbing.Hash {
	_, err := w.Add("data.txt")
	checkIfError(t, err)
	commit, err := w.Commit("example go-git commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "John Doe",
			Email: "john@doe.org",
			When:  time.Now(),
		},
	})
	checkIfError(t, err)
	return commit
}

func checkIfError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func setTag(t *testing.T, r *git.Repository, tag string) {
	h, err := r.Head()
	checkIfError(t, err)

	_, err = r.CreateTag(tag, h.Hash(), &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  "John Doe",
			Email: "john@doe.org",
			When:  time.Now(),
		},
		Message: tag,
	})
	checkIfError(t, err)
}
