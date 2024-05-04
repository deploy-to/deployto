package types

import (
	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog/log"
)

type Labels map[string]string

func (ls Labels) Has(label string) bool {
	_, ok := ls[label]
	return ok
}

func (ls Labels) SimilarTo(ls2 Labels) bool {
	for k, v := range ls {
		if k == "version" {
			lsVersion, err := semver.NewConstraint(v)
			if err != nil {
				log.Error().Err(err).Str("version", v).Msg("ls.version is not semver. Look https://github.com/Masterminds/semver?tab=readme-ov-file#checking-version-constraints")
				return false
			}

			ls2VersionStr, ok := ls2["version"]
			if !ok {
				return false
			}
			ls2Version, err := semver.NewVersion(ls2VersionStr)
			if err != nil {
				log.Error().Err(err).Str("version", v).Msg("ls2.version is not semver. Look https://github.com/Masterminds/semver?tab=readme-ov-file#checking-version-constraints")
				return false
			}
			if !lsVersion.Check(ls2Version) {
				return false
			}
		} else {
			if v2, ok := ls2[k]; !ok || v2 != v {
				return false
			}
		}
	}
	return true
}
