package buildconfig

import (
	"reflect"
	"testing"
)

func TestArtifactSpec_ToMap(t *testing.T) {
	tests := []struct {
		af   ArtifactSpec
		want map[string]string
	}{
		{
			ArtifactSpec{"1.10.8", "linux", "amd64", "abcdef1234567890", "abcdef12"},
			map[string]string{
				"GETH_VERSION": "1.10.8",
				"OS":           "linux",
				"ARCH":         "amd64",
				"COMMIT":       "abcdef1234567890",
				"SHORT_COMMIT": "abcdef12",
			},
		},
	}

	for _, tt := range tests {
		got := tt.af.ToMap()
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ArtifactSpec.ToMap() = %v, want %v", got, tt.want)
		}
	}
}
