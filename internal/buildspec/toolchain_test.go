package buildspec

import (
	"path/filepath"
	"testing"
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
		{"linux", "arm7", "GOARM=7 go run build/ci.go install -dlgo -arch arm -cc arm-linux-gnueabihf-gcc", false},
		{"linux", "arm1", "", true},
		{"darwin", "amd64", "", true},
	}

	for _, tt := range tests {
		got, err := getBuildCommand(tt.ops, tt.arch, travisFile)
		if (err != nil) != tt.expectErr {
			t.Errorf("\ngetBuildCommand() error = %v, expectErr %v\n%v\n", err, tt.expectErr, tt)
		}
		if got != tt.want {
			t.Errorf("\ngetBuildCommand() = `%v`, want `%v`\n%v\n", got, tt.want, tt)
		}
	}
}

func TestGetGoVersion(t *testing.T) {
	checksumFile := filepath.Join("testdata", "checksums.txt")
	badFile := filepath.Join("testdata", "bad-checksums.txt")

	tests := []struct {
		fileContent string
		want        string
		expectErr   bool
	}{
		{checksumFile, "1.22.4", false},
		{"invalid_file", "", true},
		{badFile, "", true},
	}
	for _, tt := range tests {
		got, err := getGoVersion(tt.fileContent)
		if (err != nil) != tt.expectErr {
			t.Errorf("\ngetGoVersion() error = %v, expectErr %v\n%v\n", err, tt.expectErr, tt)
		}
		if got != tt.want {
			t.Errorf("\ngetGoVersion() = %v, want %v\n%v\n", got, tt.want, tt)
		}
	}
}

func TestGetCC(t *testing.T) {
	tests := []struct {
		ops       string
		arch      string
		want      string
		expectErr bool
	}{
		{"linux", "amd64", "gcc-multilib", false},
		{"linux", "arm64", "gcc-aarch64-linux-gnu", false},
		{"linux", "arm66", "", true},
		{"darwin", "amd64", "", true},
		{"", "amd64", "", true},
		{"linux", "", "", true},
	}

	for _, tt := range tests {
		got, err := getCC(tt.ops, tt.arch)
		if (err != nil) != tt.expectErr {
			t.Errorf("getCC() error = %v, wantErr %v\n%v", err, tt.expectErr, tt)
		}
		if got != tt.want {
			t.Errorf("getCC() = %v, want %v\n%v", got, tt.want, tt)
		}
	}
}
