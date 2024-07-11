package buildconfig

import (
	"path/filepath"
	"testing"

	"github.com/chains-project/geth-rebuild/internal/utils"
)

func TestGetUbuntuDist(t *testing.T) {
	tests := []struct {
		YML       string
		want      string
		expectErr bool
	}{
		{filepath.Join("testdata", ".travis.yml"), "noble", false},
		{filepath.Join("testdata", ".bad-travis.yml"), "", true},
	}

	for _, tt := range tests {
		got, err := getUbuntuDist(tt.YML)
		if (err != nil) != tt.expectErr {
			t.Errorf("\ngetUbuntuDist() error = %v, expectErr %v\n%v\n", err, tt.expectErr, tt)
		}
		if got != tt.want {
			t.Errorf("\ngetUbuntuDist() = `%v`, want `%v`\n%v\n", got, tt.want, tt)
		}
	}
}

func TestGetElfTarget(t *testing.T) {
	tests := []struct {
		ops       string
		arch      string
		want      string
		expectErr bool
	}{
		{"linux", "amd64", "elf64-x86-64", false},
		{"linux", "386", "elf32-i386", false},
		{"linux", "683", "", true},
		{"windows", "amd64", "", true},
	}

	for _, tt := range tests {
		got, err := getElfTarget(utils.OS(tt.ops), utils.Arch(tt.arch))
		if (err != nil) != tt.expectErr {
			t.Errorf("getElfTarget() error = %v, expectErr %v\n%v", err, tt.expectErr, tt)
		}
		if got != tt.want {
			t.Errorf("getElfTarget() = %v, want %v\n%v", got, tt.want, tt)
		}
	}
}

func TestGetArmVersion(t *testing.T) {
	tests := []struct {
		ops       string
		arch      string
		want      string
		expectErr bool
	}{
		{"linux", "amd64", "", false},
		{"linux", "arm5", "5", false},
		{"linux", "arm64", "", false},
		{"linux", "683", "", true},
		{"windows", "amd64", "", true},
	}

	for _, tt := range tests {
		got, err := getArmVersion(utils.OS(tt.ops), utils.Arch(tt.arch))
		if (err != nil) != tt.expectErr {
			t.Errorf("getArmVersion() error = %v, expectErr %v\n%v", err, tt.expectErr, tt)
		}
		if got != tt.want {
			t.Errorf("getArmVersion() = %v, want %v\n%v", got, tt.want, tt)
		}
	}
}
