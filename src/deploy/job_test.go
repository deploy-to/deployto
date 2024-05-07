package deploy

import (
	"deployto/src/types"
	"reflect"
	"testing"
)

func Test_runJob(t *testing.T) {
	tests := []struct {
		name       string
		steps      []types.Step
		aliases    []string
		jobContext types.Values
		want       types.Values
		wantErr    bool
	}{
		{
			name: "all in one",
			steps: []types.Step{
				{
					Id: "step1",
					Run: `
echo "OUTPUT_VAR1=$INPUT_CONST"       >> $DEPLOYTO_OUTPUT
echo "OUTPUT_VAR2=$INPUT_TEMPLAYTING" >> $DEPLOYTO_OUTPUT
					`,
					Env: map[string]string{
						"INPUT_CONST":       "CONST1",
						"INPUT_TEMPLAYTING": "{{ .contextKey }}",
					},
				},
			},
			aliases: []string{"test", "test2"},
			jobContext: map[string]any{
				"contextKey": "contextVal",
			},
			want: map[string]any{
				"step1": map[string]any{
					"OUTPUT_VAR1": "CONST1",
					"OUTPUT_VAR2": "contextVal",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := runJob(&types.Job{Spec: &types.JobSpec{Steps: tt.steps}}, tt.aliases, tt.jobContext)
			if (err != nil) != tt.wantErr {
				t.Errorf("runJob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !MapsEqual(got, tt.want) {
				t.Errorf("runJob() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO not work for map[string]map[string]int NEED CODEREVIW
func MapsEqual(m1, m2 map[string]any) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		v2, ok := m2[k]
		if !ok {
			return false
		}

		v1m, isMap1 := v1.(map[string]any)
		v2m, isMap2 := v2.(map[string]any)
		v1ms, isMap1s := v1.(map[string]string)
		v2ms, isMap2s := v2.(map[string]string)

		if isMap1 && isMap2 {
			return MapsEqual(v1m, v2m)
		}
		if !isMap1 && isMap2 {
			if isMap1s {
				return MapsEqualAnyStr(v2m, v1ms)
			} else {
				return false
			}
		}
		if isMap1s && isMap2s {
			if isMap2s {
				return MapsEqualAnyStr(v1m, v2ms)
			} else {
				return false
			}
		}
		if isMap1 && isMap2 {
			return reflect.DeepEqual(v1ms, v2ms)
		}
	}
	return true
}

func MapsEqualAnyStr(m1 map[string]any, m2 map[string]string) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		v1s, ok1 := v1.(string)
		v2s, ok2 := m2[k]
		if !ok1 || !ok2 || v1s != v2s {
			return false
		}
	}
	return true
}
