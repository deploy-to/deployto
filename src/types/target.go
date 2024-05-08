package types

import (
	"io"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
	"github.com/rs/zerolog/log"
)

type Target struct {
	Base `json:",inline" yaml:",inline"`
	Spec *TargetSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type TargetSpec struct {
	Namespace  string     `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Kubeconfig Kubeconfig `json:"kubeconfig,omitempty" yaml:"kubeconfig,omitempty"`
	Script     *Script    `json:"script,omitempty" yaml:"script,omitempty"`
}

type Kubeconfig struct {
	Filename   string `json:"filename,omitempty" yaml:"filename,omitempty"`
	UseDefault bool   `json:"usedefault,omitempty" yaml:"usedefault,omitempty"`
}

var SystemKubeconfig string // set from Flags

func (t *Target) LoadKubeconfig() []byte {
	if t.Spec.Kubeconfig.Filename != "" {
		workdir := filepath.Dir(t.Status.FileName)
		kubeconfigFilename := filepath.Join(workdir, t.Spec.Kubeconfig.Filename)
		kubeconfig, err := t.Base.Status.Filesystem.FS.Open(kubeconfigFilename)
		if err != nil {
			log.Error().Err(err).Str("filename", kubeconfigFilename).Msg("Open kubeconfig error")
			return nil
		}
		result, err := ReadFile(kubeconfig)
		if err != nil {
			log.Error().Err(err).Str("filename", kubeconfigFilename).Msg("ReadFile kubeconfig error (billy)")
			return nil
		}
		return result
	} else {
		if t.Spec.Kubeconfig.UseDefault {
			var kubeconfigFilename string
			userHomeDir, err := os.UserHomeDir()
			if err == nil {
				kubeconfigFilename = filepath.Join(userHomeDir, "/.kube/config")
			}
			if SystemKubeconfig != "" {
				kubeconfigFilename = SystemKubeconfig
			}
			result, err := os.ReadFile(kubeconfigFilename)
			if err != nil {
				log.Error().Err(err).Str("filename", kubeconfigFilename).Msg("ReadFile kubeconfig error (os)")
				return nil
			}
			return result
		}
	}
	return nil
}

func ReadFile(f billy.File) ([]byte, error) {
	defer f.Close()
	size := 512
	data := make([]byte, 0, size)
	for {
		n, err := f.Read(data[len(data):cap(data)])
		data = data[:len(data)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return data, err
		}

		if len(data) >= cap(data) {
			d := append(data[:cap(data)], 0)
			data = d[:len(data)]
		}
	}
}
