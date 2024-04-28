package types

import (
	"reflect"
	"testing"
)

// TODO https://stackoverflow.com/questions/70980279/table-testing-go-generics
func TestGetString(t *testing.T) {
	tests := []struct {
		name   string
		values Values
		def    string
		path   string
		want   string
	}{
		{
			name:   "empty values",
			values: nil,
			def:    "test default value",
			path:   "aaa",
			want:   "test default value",
		},
		{
			name:   "empty path",
			values: Values{"key": "value"},
			def:    "test default value",
			path:   "",
			want:   "test default value",
		},
		{
			name:   "single key",
			values: Values{"key": "value"},
			def:    "test default value",
			path:   "key",
			want:   "value",
		},
		{
			name:   "second level",
			values: Values{"key-level1": Values{"key-level2": "value-level2"}},
			def:    "test default value",
			path:   "key-level1.key-level2",
			want:   "value-level2",
		},
		{
			name:   "third level not found",
			values: Values{"key-level1": Values{"key-level2": "value-level2"}},
			def:    "test default value",
			path:   "key-level1.key-level2.not-found",
			want:   "test default value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Get(tt.values, tt.def, tt.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExists(t *testing.T) {
	tests := []struct {
		name   string
		values Values
		path   string
		want   bool
	}{
		{
			name:   "empty values",
			values: nil,
			path:   "aaa",
			want:   false,
		},
		{
			name:   "empty path",
			values: Values{"key": "value"},
			path:   "",
			want:   false,
		},
		{
			name:   "single key",
			values: Values{"key": "value"},
			path:   "key",
			want:   true,
		},
		{
			name:   "second level",
			values: Values{"key-level1": Values{"key-level2": "value-level2"}},
			path:   "key-level1.key-level2",
			want:   true,
		},
		{
			name:   "third level not found",
			values: Values{"key-level1": Values{"key-level2": "value-level2"}},
			path:   "key-level1.key-level2.not-found",
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Exists(tt.values, tt.path); got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}
