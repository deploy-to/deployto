package gitclient

import (
	"deployto/src/helper"
	"deployto/src/types"
	"path/filepath"
	"sort"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/lithammer/shortuuid/v3"
	"github.com/rs/zerolog/log"
)

func GetValues(path string) (values types.Values) {
	values = make(types.Values)
	dirtyPostfix := "+dirty.uuid" + shortuuid.New()
	values["Commit"] = "NO_GIT" + dirtyPostfix
	values["CommitShort"] = "NO_GIT" + dirtyPostfix

	gitRoot, err := helper.GetProjectRoot(path, ".git")
	if err != nil {
		log.Error().Err(err).Msg("Search git root error")
		return
	}

	rep, err := git.PlainOpen(filepath.Dir(gitRoot))
	if err != nil {
		log.Error().Err(err).Msg("Error opening git repository")
		return
	}
	ref, err := rep.Head()
	if err != nil {
		log.Warn().Err(err).Msg("Error getting git Head")
		return
	}

	wt, err := rep.Worktree()
	if err != nil {
		log.Error().Err(err).Msg("Error getting git Worktree")
		return nil
	}
	s, err := wt.Status()
	if err != nil {
		log.Error().Err(err).Msg("Error getting git Status")
		return
	}

	if s.IsClean() {
		dirtyPostfix = ""
	}

	values["Commit"] = ref.Hash().String() + dirtyPostfix
	values["CommitShort"] = ref.Hash().String()[:7] + dirtyPostfix

	//Find tag (semver or another one)
	tagrefs, err := rep.Tags()
	if err != nil {
		log.Debug().Err(err).Msg("Git tag not found")
	}
	var Tags []string
	err = tagrefs.ForEach(func(t *plumbing.Reference) error {
		if t.Hash().String() == ref.Hash().String() {
			Tags = append(Tags, t.Name().Short())
		}
		return nil
	})
	if err != nil {
		log.Debug().Err(err).Msg("Git tag not found")
	}

	if len(Tags) > 0 {
		sort.Strings(Tags)
		values["Tag"] = Tags[len(Tags)-1] + dirtyPostfix
	}

	log.Debug().Any("values", values).Msg("gitclient.GetValues")

	return values
}
