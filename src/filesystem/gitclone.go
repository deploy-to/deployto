package filesystem

import (
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"

	"github.com/lithammer/shortuuid/v3"
	"github.com/rs/zerolog/log"
)

func Clone2Tmp(url string) (path string, err error) {
	authMethod, err := ssh.NewSSHAgentAuth("git")
	if err != nil {
		log.Error().Err(err).Msg("new ssh agent auth error")
		return "", err
	}

	path = filepath.Join(os.TempDir(), "deployto-gitclones", shortuuid.New())
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
		log.Error().Err(err).Str("url", url).Msg("Can't clone git")
		extraError := os.RemoveAll(path)
		if extraError != nil {
			log.Error().Err(extraError).Str("path", url).Msg("Can't remove tmp path")
		}
		return "", err
	}
	return path, err
}
