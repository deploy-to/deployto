package src

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

type ContextDump struct {
	Root    string
	Current string

	counter *int32
}

func GetContextDump(workdir, ContextDumpDir string) *ContextDump {
	if ContextDumpDir == "" {
		return nil
	}
	path := filepath.Join(workdir, ContextDumpDir, time.Now().Format("2006-01-02-150405"))
	return &ContextDump{
		Root:    path,
		Current: path,
		counter: new(int32),
	}
}

func (cd *ContextDump) Next(alias string) *ContextDump {
	if cd == nil {
		return nil
	}
	counter := atomic.AddInt32(cd.counter, 1)
	current := filepath.Join(cd.Current, fmt.Sprintf("%03d", counter)+alias)
	log.Info().Str("resultPath", current).Msg("context debug")
	err := os.MkdirAll(current, 0766)
	if err != nil {
		log.Error().Err(err).Msg("ContextDump MkdirAll error")
	}
	return &ContextDump{
		Root:    cd.Root,
		Current: current,
		counter: cd.counter,
	}
}

func (dr *ContextDump) Push(name string, values any) {
	if dr == nil {
		return
	}
	counter := atomic.AddInt32(dr.counter, 1)
	filename := filepath.Join(dr.Current, fmt.Sprintf("%03d", counter)+name)

	var data []byte
	var err error

	switch v := values.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		data, err = yaml.Marshal(v)
		if err != nil {
			log.Error().Err(err).Msg("ContextDump marshal yaml error")
			return
		}
	}
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		log.Error().Err(err).Msg("ContextDump write yaml file error")
	}
}
