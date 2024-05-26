package deploy

import (
	"bytes"
	"deployto/src/types"
	"html/template"

	"github.com/rs/zerolog/log"
)

type Templater struct {
	deploy *Deploy
	funcs  template.FuncMap
}

func NewTemplater(deploy *Deploy) *Templater {
	return &Templater{
		deploy: deploy,
		funcs: template.FuncMap{
			"inc": func(i int) int { return i + 1 },
			"dec": func(i int) int { return i + 1 },
			"add": func(i1 int, i2 int) int { return i1 + i2 },
			"sub": func(i1 int, i2 int) int { return i1 - i2 },
		},
	}
}

func (templater *Templater) TemplatingString(templ string, context types.Values) (string, error) {
	t, err := template.New("letter").Funcs(templater.funcs).Parse(templ)
	if err != nil {
		log.Error().Err(err).Str("template", templ).Msg("Template parse error")
		return "nil", err
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, context)
	if err != nil {
		log.Error().Err(err).Str("template", templ).Msg("Template execute with scriptContext error")
		return "", err
	}
	return buf.String(), nil
}

func (templater *Templater) Templating(values, context types.Values) (types.Values, error) {
	enrichedContext := types.MergeValues(
		types.Values{
			"Files": templater.deploy.FS,
		},
		context,
	)
	result := make(types.Values)
	for k, v := range values {
		switch vTyped := v.(type) {
		case types.Values:
			subResult, err := templater.Templating(vTyped, enrichedContext)
			if err != nil {
				log.Error().Err(err).Str("key", k).Msg("Template subValues execute with scriptContext error")
				return nil, err
			}
			result[k] = subResult
		case string:
			res, err := templater.TemplatingString(vTyped, enrichedContext)
			if err != nil {
				log.Error().Err(err).Str("key", k).Str("template", vTyped).Msg("Templating error")
				return nil, err
			}
			result[k] = res
		default:
			result[k] = v
		}
	}
	return result, nil
}
