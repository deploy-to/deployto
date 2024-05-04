package gitclient

import (
	"deployto/src/filesystem"
	"deployto/src/types"
	"sort"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	gogitfilesystem "github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/lithammer/shortuuid/v3"
	"github.com/rs/zerolog/log"
)

func GetValues(fs *filesystem.Filesystem, path string) (values types.Values) {
	values = make(types.Values)

	gitRoot := filesystem.GetGitRootFilesystem(fs, path)
	if gitRoot == nil {
		log.Error().Msg("git not found")
		return
	}

	// set default result
	dirtyPostfix := "+dirty.uuid" + shortuuid.New()
	values["Commit"] = "NO_GIT" + dirtyPostfix
	values["CommitShort"] = "NO_GIT" + dirtyPostfix

	storer, err := gitRoot.FS.Chroot(git.GitDirName)
	if err != nil {
		log.Error().Err(err).Msg("git storer not fount")
		return
	}
	rep, err := git.Open(gogitfilesystem.NewStorage(storer, cache.NewObjectLRUDefault()), gitRoot.FS)
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
		Tags = append(Tags, t.Name().Short())
		return nil
	})
	if err != nil {
		log.Debug().Err(err).Msg("Git tag not found")
	}

	if len(Tags) > 0 {
		sort.Strings(Tags)
		values["Tag"] = Tags[0] + dirtyPostfix
	}

	log.Debug().Any("values", values).Msg("gitclient.GetValues")

	return values
}
