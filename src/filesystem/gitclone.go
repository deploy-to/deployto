package filesystem

import (
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"

	"crypto/sha1"
	"encoding/base64"

	"github.com/rs/zerolog/log"
)

func Clone2Tmp(url string) (path string, err error) {
	authMethod, err := ssh.NewSSHAgentAuth("git")
	if err != nil {
		log.Error().Err(err).Msg("new ssh agent auth error")
		return "", err
	}

	hasher := sha1.New()
	hasher.Write([]byte(url))
	urlsha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	path = filepath.Join(os.TempDir(), "deployto-gitclones", urlsha)
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Error().Err(err).Msg("Can't make tmp dir for git clone")
		return "", err
	}

	_, err = git.PlainClone(path, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Auth:              authMethod,
	})
	if err != nil {
		if err.Error() != "repository already exists" {
			log.Error().Err(err).Str("url", url).Msg("Can't clone git")
			extraError := os.RemoveAll(path)
			if extraError != nil {
				log.Error().Err(extraError).Str("path", url).Msg("Can't remove tmp path")
			}
			return "", err
		}
		// git pull
		r, err := git.PlainOpen(path)
		if err != nil {
			log.Error().Err(err).Str("path", path).Msg("git.PlainOpen error")
			return "", nil
		}
		w, err := r.Worktree()
		if err != nil {
			log.Error().Err(err).Str("path", path).Msg("git get worktree error")
			return "", nil
		}
		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil {
			log.Error().Err(err).Str("path", path).Msg("git pull error")
			return "", nil
		}
	}
	return path, err
}
