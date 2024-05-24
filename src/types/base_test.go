package types

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
)

func TestBase_AsValues(t *testing.T) {
	tests := []struct {
		name     string
		asBase   *Base
		asValues Values
	}{
		{
			name:   "only name",
			asBase: &Base{Meta: MetaData{Name: "TeSt"}},
			asValues: Values{
				"metadata": Values{"name": "TeSt"},
			},
		},
		{
			name: "filesystem",
			asBase: &Base{
				Meta: MetaData{Name: "TeSt"},
				Status: StatusType{
					FileName: "TestFileName",
				},
			},
			asValues: Values{
				"metadata": Values{"name": "TeSt"},
				"Status": Values{
					"FileName": "TestFileName",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deep.NilSlicesAreEmpty = true
			deep.NilMapsAreEmpty = true

			gotValues := tt.asBase.AsValues()
			if diff := deep.Equal(gotValues, tt.asValues); diff != nil {
				t.Error(strings.Join(diff, "; "))
			}
		})
	}
}
