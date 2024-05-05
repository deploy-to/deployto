package filesystem

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/rs/zerolog/log"
)

type FilesystemType int

const (
	LOCAL FilesystemType = iota
	GIT
)

type Filesystem struct {
	URI       string
	localPath string
	Type      FilesystemType
	FS        billy.Filesystem
}

func Supported(uri string) bool {
	return strings.HasPrefix(uri, "file://") ||
		strings.HasSuffix(uri, ".git")
}

func Get(uri string) *Filesystem {
	if localPath, isLOCAL := strings.CutPrefix(uri, "file://"); isLOCAL {
		return &Filesystem{
			Type:      LOCAL,
			URI:       uri,
			localPath: localPath,
			FS:        osfs.New(localPath),
		}
	}
	if strings.HasSuffix(uri, ".git") {
		localPath, err := Clone2Tmp(uri)
		//TODO HOW REMOVE TEMPORY DIRECTORY ?
		if err != nil {
			log.Error().Err(err).Msg("git clone error")
			return nil
		}
		return &Filesystem{
			Type:      GIT,
			URI:       uri,
			localPath: localPath,
			FS:        osfs.New(localPath),
		}
	}
	log.Fatal().Str("basedir", uri).Msg("filesystem not implimented")
	return nil
}

func GetGitRootFilesystem(fs *Filesystem, path string) *Filesystem {
	if fs == nil {
		return nil
	}
	if fs.Type == GIT {
		return fs
	}
	if fs.Type == LOCAL {
		return searchLocalRoot(fs, path, git.GitDirName)
	}
	return nil
}

func GetDeploytoRootFilesystem(fs *Filesystem, path string) *Filesystem {
	if fs == nil {
		return nil
	}
	if fs.Type == GIT {
		log.Error().Msg("GetDeploytoRootFilesystem for git not implimented")
		return nil
	}
	if fs.Type == LOCAL {
		return searchLocalRoot(fs, path, DeploytoDirName)
	}
	return nil
}

func searchLocalRoot(fs *Filesystem, path string, dirName string) *Filesystem {
	currentPath, wasCut := strings.CutPrefix(fs.URI, "file://")
	if !wasCut {
		log.Error().Str("baseDir", fs.URI).Msg("wait file:// prefix")
		return nil
	}
	currentPath = filepath.Clean(filepath.Join(currentPath, path))
	for {
		if IsDirExists(filepath.Join(currentPath, dirName)) {
			log.Debug().Str("searchDir", dirName).Str("path", currentPath).Msg("searchDir found")
			return Get("file://" + currentPath)
		}
		log.Trace().Str("searchDir", dirName).Str("path", currentPath).Msg("searchDir not found - go to parent")

		if strings.HasSuffix(currentPath, string(os.PathSeparator)) {
			return nil // root dir
		}
		currentPath = filepath.Dir(currentPath)
	}
}

const DeploytoDirName = ".deployto"

func IsDirExists(localPath string) bool {
	fi, err := os.Stat(localPath)
	if err == nil {
		return fi.IsDir()
	}
	if !os.IsNotExist(err) {
		log.Info().Err(err).Str("path", localPath).Msg("Check dir exists")
	}
	return false
}
