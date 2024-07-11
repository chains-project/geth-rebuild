package buildconfig

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/chains-project/geth-rebuild/internal/utils"
)

// TODO fix

func TestGetBuildCommand(t *testing.T) {
	travisFile := filepath.Join("testdata", ".travis.yml")

	tests := []struct {
		ops       string
		arch      string
		want      string
		expectErr bool
	}{
		{"linux", "amd64", "go run build/ci.go install -dlgo", false},
		{"linux", "386", "go run build/ci.go install -dlgo -arch 386", false},
		{"linux", "arm7", "go run build/ci.go install -dlgo -arch arm -cc arm-linux-gnueabihf-gcc", false},
		{"linux", "arm", "", true},
		{"darwin", "amd64", "", true},
	}

	for _, tt := range tests {
		got, err := getBuildCommand(utils.OS(tt.ops), utils.Arch(tt.arch), travisFile)
		if (err != nil) != tt.expectErr {
			t.Errorf("\ngetBuildCommand() error = %v, expectErr %v\n%v\n", err, tt.expectErr, tt)
		}
		if got != tt.want {
			t.Errorf("\ngetBuildCommand() = `%v`, want `%v`\n%v\n", got, tt.want, tt)
		}
	}
}

func TestGetGCVersion(t *testing.T) {
	tests := []struct {
		fileContent string
		want        string
		expectErr   bool
	}{
		{filepath.Join("testdata", "checksums.txt"), "1.22.4", false},
		{"invalid_file", "", true},
		{filepath.Join("testdata", "bad-checksums.txt"), "", true},
	}
	for _, tt := range tests {
		got, err := getGCVersion(tt.fileContent)
		if (err != nil) != tt.expectErr {
			t.Errorf("\ngetGCVersion() error = %v, expectErr %v\n%v\n", err, tt.expectErr, tt)
		}
		if got != tt.want {
			t.Errorf("\ngetGCVersion() = %v, want %v\n%v\n", got, tt.want, tt)
		}
	}
}

func TestgetToolChainDeps(t *testing.T) {
	tests := []struct {
		ops       string
		arch      string
		want      []string
		expectErr bool
	}{
		{"linux", "amd64", []string{"gcc-multilib"}, false},
		{"linux", "arm5", []string{"gcc-arm-linux-gnueabi", "libc6-dev-armel-cross"}, false},
		{"linux", "arm64", []string{"gcc-aarch64-linux-gnu", "libc6-dev-arm64-cross"}, false},
		{"linux", "683", nil, true},
		{"windows", "amd64", nil, true},
	}

	for _, tt := range tests {
		got, err := getToolChainDeps(utils.OS(tt.ops), utils.Arch(tt.arch))
		if (err != nil) != tt.expectErr {
			t.Errorf("getToolChainDeps() error = %v, expectErr %v, test case: %+v", err, tt.expectErr, tt)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("getToolChainDeps() = %v, want %v, test case: %+v", got, tt.want, tt)
		}
	}
}
