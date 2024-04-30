package types

import (
	"reflect"
	"testing"
)

func TestDecodeScript(t *testing.T) {
	tests := []struct {
		name       string
		values     Values
		wantScript *Script
	}{
		{
			name: "simple",
			values: Values{
				"type": "testType",
				"root": true,
				"v1":   "k1",
				"v2":   2,
			},
			wantScript: &Script{
				Type:   "testType",
				Root:   true,
				Values: Values{"v1": "k1", "v2": 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotScript := DecodeScript(tt.values)
			if !reflect.DeepEqual(gotScript, tt.wantScript) {
				t.Errorf("DecodeScript() = %v, want %v", gotScript, tt.wantScript)
			}
		})
	}
}
