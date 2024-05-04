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
	BaseDir string
	Type    FilesystemType
	FS      billy.Filesystem
}

func GetFilesystem(baseDir string) *Filesystem {
	if localDir, isLOCAL := strings.CutPrefix(baseDir, "file://"); isLOCAL {
		return &Filesystem{
			BaseDir: baseDir,
			Type:    LOCAL,
			FS:      osfs.New(localDir),
		}
	}
	if strings.HasPrefix(baseDir, "https://") && strings.HasSuffix(baseDir, ".git") {
		//TODO git in GetFilesystem
		// 1) clone (in memory?) (in tmp? HOWTO remove?)
		// 2) set &Filesystem.FS
		log.Fatal().Msg("GIT NOT IMPLIMENTED")
		return &Filesystem{
			BaseDir: baseDir,
			Type:    GIT,
			FS:      nil,
		}
	}
	log.Fatal().Str("basedir", baseDir).Msg("filesystem not implimented")
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
	currentPath, wasCut := strings.CutPrefix(fs.BaseDir, "file://")
	if !wasCut {
		log.Error().Str("baseDir", fs.BaseDir).Msg("wait file:// prefix")
		return nil
	}
	currentPath = filepath.Clean(filepath.Join(currentPath, path))
	for {
		if IsDirExists(filepath.Join(currentPath, dirName)) {
			log.Debug().Str("searchDir", dirName).Str("path", currentPath).Msg("searchDir found")
			return GetFilesystem("file://" + currentPath)
		}
		log.Debug().Str("searchDir", dirName).Str("path", currentPath).Msg("searchDir not found - go to parent")

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
