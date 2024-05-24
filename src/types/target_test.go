package types

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
)

func TestTarget_Code(t *testing.T) {
	tests := []struct {
		name     string
		asTarget *Target
		asValues Values
	}{
		{
			name:     "only name",
			asTarget: &Target{Base: Base{Meta: MetaData{Name: "TeSt"}}},
			asValues: Values{
				"metadata": Values{"name": "TeSt"},
				"spec": Values{
					"Kubeconfig": Values{
						"Namespace":  "",
						"Filename":   "",
						"UseDefault": false,
					},
					"Terraform": Values(nil),
				},
			},
		},
		{
			name:     "filesystem",
			asTarget: &Target{Base: Base{Meta: MetaData{Name: "TeSt"}}},
			asValues: Values{
				"metadata": Values{"name": "TeSt"},
				"spec": Values{
					"Kubeconfig": Values{
						"Namespace":  "",
						"Filename":   "",
						"UseDefault": false,
					},
					"Terraform": Values(nil),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deep.NilSlicesAreEmpty = true
			deep.NilMapsAreEmpty = true

			gotValues := tt.asTarget.AsValues()
			if diff := deep.Equal(gotValues, tt.asValues); diff != nil {
				t.Error(strings.Join(diff, "; "))
			}
			gotTarget := DecodeTarget(tt.asValues)
			if diff := deep.Equal(gotTarget, tt.asTarget); diff != nil {
				t.Error(diff)
			}
		})
	}
}
