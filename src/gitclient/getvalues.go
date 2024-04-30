package gitclient

import (
	"deployto/src/helper"
	"deployto/src/types"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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
		log.Error().Err(err).Msg("Error getting git Head")
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

	//Find Tag (semver or another one)
	Tag, err := rep.TagObject(ref.Hash())
	switch err {
	case nil:
		log.Error().Err(err).Msg("Error getting git tag")

	// Tag object present
	case plumbing.ErrObjectNotFound:
		log.Error().Err(err).Msg("Error getting git tag")
		//Default tag is empty
		Tag = &object.Tag{
			Name: "",
		}
	}

	values := make(types.Values)
	if s.IsClean() {
		values["Commit"] = ref.Hash().String()
		values["CommitShort"] = ref.Hash().String()[:7]
		values["Tag"] = Tag.Name
	} else {
		uuid := shortuuid.New()
		values["Commit"] = ref.Hash().String() + "+dirty.uuid" + uuid
		values["CommitShort"] = ref.Hash().String()[:7] + "+dirty.uuid" + uuid
		values["Tag"] = Tag.Name + uuid
	}

	log.Debug().Any("values", values).Msg("gitclient.GetValues")

	return values
}
