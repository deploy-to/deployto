package gitclient

import (
	"deployto/src/helper"
	"deployto/src/types"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/lithammer/shortuuid/v3"
	"github.com/rs/zerolog/log"
)

func GetValues(path string) types.Values {
	gitRoot, err := helper.GetProjectRoot(path, ".git")
	if err != nil {
		log.Error().Err(err).Msg("Search git root error")
		return nil
	}

	rep, err := git.PlainOpen(filepath.Dir(gitRoot))
	if err != nil {
		log.Error().Err(err).Msg("Error opening git repository")
		return nil
	}
	ref, err := rep.Head()
	if err != nil {
		log.Warn().Err(err).Msg("Error getting git Head")
		return nil
	}

	wt, err := rep.Worktree()
	if err != nil {
		log.Error().Err(err).Msg("Error getting git Worktree")
		return nil
	}
	s, err := wt.Status()
	if err != nil {
		log.Error().Err(err).Msg("Error getting git Status")
		return nil
	}

	var dirtyPostfix string
	if !s.IsClean() {
		dirtyPostfix = "+dirty.uuid" + shortuuid.New()
	}

	values := make(types.Values)
	values["Commit"] = ref.Hash().String() + dirtyPostfix
	values["CommitShort"] = ref.Hash().String()[:7] + dirtyPostfix

	//Find tag (semver or another one)
	tag, err := rep.TagObject(ref.Hash())
	if err != nil {
		log.Debug().Err(err).Msg("Git tag not found")
	}
	if tag != nil {
		values["Tag"] = tag.Name + dirtyPostfix
	}

	log.Debug().Any("values", values).Msg("gitclient.GetValues")

	return values
}
