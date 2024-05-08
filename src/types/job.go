package types

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

type Job struct {
	Base `json:",inline" yaml:",inline"`
	Spec *JobSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type JobSpec struct {
	Steps []Step `json:"steps,omitempty" yaml:"steps,omitempty"`
}

type Step struct {
	Id              string            `json:"id,omitempty" yaml:"id,omitempty"`
	Name            string            `json:"name,omitempty" yaml:"name,omitempty"`
	ContinueOnError bool              `json:"continue-on-error,omitempty" yaml:"continue-on-error,omitempty"`
	Run             string            `json:"run,omitempty" yaml:"run,omitempty"`
	Env             map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
}

func DecodeJob(values any) (job *Job) {
	if values == nil {
		log.Info().Msg("DecodeScript - input values is nil")
		return nil
	}
	job = &Job{}
	err := mapstructure.Decode(values, job)
	if err != nil {
		log.Error().Err(err).Msg("DecodeScript error")
		return nil
	}
	return job
}
