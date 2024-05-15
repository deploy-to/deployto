package types

import (
	"cmp"
	"slices"

	"github.com/rs/zerolog/log"
)

type Component struct {
	Base `json:",inline" yaml:",inline"`
	Spec map[string]Values `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type TheOrderOfResource = []sameOrder
type sameOrder = []*Script

func GetTheOrderOfResource(spec map[string]Values) TheOrderOfResource {
	result := make(TheOrderOfResource, 0, 2)
	for k, v := range spec {
		script := DecodeScript(k, v)
		if script == nil {
			log.Error().Str("scriptKey", k).Msg("error decode script")
			return nil
		}

		var newOrder = true
		for i := 0; i < len(result); i++ {
			if result[i][0].Order == script.Order {
				result[i] = append(result[i], script)
				newOrder = false
				break
			}
		}
		if newOrder {
			result = append(result, sameOrder{script})
		}
	}
	slices.SortStableFunc(result, func(a, b sameOrder) int { return cmp.Compare(a[0].Order, b[0].Order) })
	return result
}
