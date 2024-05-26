package filesystem

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/structs"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/lithammer/shortuuid/v3"
	"github.com/rs/zerolog/log"
)

const (
	DESTROYED string = "DESTROYED"
	LOCAL     string = "LOCAL"
	GIT       string = "GIT"
	TEMP      string = "TEMP"
)

type Filesystem struct {
	URI       string           `structs:",omitempty"`
	LocalPath string           `structs:",omitempty"`
	Type      string           `structs:",omitempty"`
	FS        billy.Filesystem `structs:"-"`
}

func Supported(uri string) bool {
	return strings.HasPrefix(uri, "file://") ||
		strings.HasSuffix(uri, ".git")
}

func Get(uri string) *Filesystem {
	//TODO Howto remove temp dir (and git). Save list and remove on application exit?
	if uri == "temp" {
		localPath := filepath.Join(os.TempDir(), "deployto-tempfs", shortuuid.New())
		err := os.MkdirAll(localPath, os.ModePerm)
		if err != nil {
			log.Error().Err(err).Msg("Can't make tmp dir for git clone")
			return nil
		}
		return &Filesystem{
			Type:      TEMP,
			URI:       uri,
			LocalPath: localPath,
			FS:        osfs.New(localPath),
		}
	}
	if localPath, isLOCAL := strings.CutPrefix(uri, "file://"); isLOCAL {
		return &Filesystem{
			Type:      LOCAL,
			URI:       uri,
			LocalPath: localPath,
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
			LocalPath: localPath,
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
	// loop guard
	// если за 50 итерация не нашли директорию
	for i := 1; i < 50; i++ {
		if IsDirExists(filepath.Join(currentPath, dirName)) {
			log.Trace().Str("dirName", dirName).Str("path", currentPath).Msg("searchLocalRoot - ok")
			return Get("file://" + currentPath)
		}
		log.Trace().Str("dirName", dirName).Str("path", currentPath).Msg("searchLocalRoot - not found - go to parent")

		if strings.HasSuffix(currentPath, string(os.PathSeparator)) {
			return nil // root dir
		}
		currentPath = filepath.Dir(currentPath)
	}
	return nil
}

const DeploytoDirName = ".deployto"

func IsDirExists(localPath string) bool {
	fi, err := os.Stat(localPath)
	if err == nil {
		return fi.IsDir()
	}
	if !os.IsNotExist(err) {
		log.Trace().Err(err).Str("path", localPath).Msg("Check dir exists")
	}
	return false
}

func (fs *Filesystem) Destroy() {
	if fs.Type == LOCAL {
		log.Trace().Msg("don’t destroy the local file system")
		return
	}
	if fs.Type == GIT || fs.Type == TEMP {
		if !strings.HasPrefix(fs.LocalPath, os.TempDir()) {
			log.Error().Str("path", fs.LocalPath).Msg("I'm afraid to delete something important. I don't delete it.")
			return
		}
		err := os.RemoveAll(fs.LocalPath)
		if err != nil {
			log.Error().Err(err).Str("path", fs.LocalPath).Msg("destroy error")
		}
		fs.LocalPath = "FILESYSTEM WAS DESTROYED"
		fs.URI = "FILESYSTEM WAS DESTROYED"
		fs.Type = DESTROYED
		fs = nil
		return
	}
	log.Error().Str("path", fs.LocalPath).Msg("destroy not implemented")
}

func (fs *Filesystem) AsValues() map[string]any { // can't import types - loop
	return structs.Map(fs)
}

func (fs *Filesystem) Get(fileName string) template.HTML {
	log.Debug().Str("fileName", fileName).Msg("Filesystem.Get")
	//TODO security load from parent dir
	//TODO как быть, если ресурс запущен во вложенной папке?
	bytes, err := os.ReadFile(filepath.Join(fs.LocalPath, fileName))
	if err != nil {
		log.Error().Err(err).Str("fs", fs.URI).Str("fileName", fileName).Msg("read file error")
	}
	return template.HTML(string(bytes))
}
