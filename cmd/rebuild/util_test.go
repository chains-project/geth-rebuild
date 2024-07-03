package main

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestValidateArgs(t *testing.T) {
	tests := []struct {
		ops       string
		arch      string
		version   string
		expectErr bool
		errMsg    string
	}{
		{"linux", "amd64", "1.14.4", false, ""},
		{"linux", "386", "0.1.0", false, ""},
		{"linux", "arm5", "1.14.3.d", true, "<geth version> must be in format 'major.minor.patch'\nExample: 1.14.4"},
		{"windows", "amd64", "1.14.4", true, "<os> limited to `linux` at the moment"},
		{"linux", "arm666", "1.14.4", true, "<arch> must be a valid linux target architecture (amd64|386|arm5|arm6|arm64|arm7)"},
		{"linux", "amd64", "1.1.", true, "<geth version> must be in format 'major.minor.patch'\nExample: 1.14.4"},
		{"linux", "amd64", "1..2", true, "<geth version> must be in format 'major.minor.patch'\nExample: 1.14.4"},
		{"linux", "amd64", ".1.2.2", true, "<geth version> must be in format 'major.minor.patch'\nExample: 1.14.4"},
	}

	for _, tt := range tests {
		err := validateArgs(tt.ops, tt.arch, tt.version)
		if (err != nil) != tt.expectErr {

			t.Errorf("validateArgs() error = %v, expectErr %v\n%v", err, tt.expectErr, tt)
		}
		if err != nil && err.Error() != tt.errMsg {
			t.Errorf("validateArgs() error = %v, errMsg %v\n%v", err, tt.errMsg, tt)
		}
	}
}

func TestGetBuildCommand(t *testing.T) {
	travisFile := filepath.Join("testdata", "travis.yml")

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
		got, err := getElfTarget(tt.ops, tt.arch)
		if (err != nil) != tt.expectErr {
			t.Errorf("getElfTarget() error = %v, expectErr %v\n%v", err, tt.expectErr, tt)
		}
		if got != tt.want {
			t.Errorf("getElfTarget() = %v, want %v\n%v", got, tt.want, tt)
		}
	}
}

func TestGetUbuntuPackages(t *testing.T) {
	tests := []struct {
		ops       string
		arch      string
		want      []string
		expectErr bool
	}{
		{"linux", "amd64", []string{"git", "ca-certificates", "wget", "binutils"}, false},
		{"linux", "arm64", []string{"git", "ca-certificates", "wget", "binutils", "libc6-dev-arm64-cross"}, false},
		{"linux", "arm42", nil, true},
		{"windows", "amd64", nil, true},
	}

	for _, tt := range tests {
		got, err := getUbuntuPackages(tt.ops, tt.arch)
		if (err != nil) != tt.expectErr {
			t.Errorf("getUbuntuPackages() error = %v, expectErr %v\n%v", err, tt.expectErr, tt)
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("getUbuntuPackages() = %v, want %v\n%v", got, tt.want, tt)
		}
	}
}

// func TestGetArtifactSpec(t *testing.T) { //TODO
// 	paths := setUpPaths()
// 	_, err := common.RunCommand(paths.Scripts.Clone, paths.Directories.Temp)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	_, err = common.RunCommand(paths.Scripts.Checkout, paths.Directories.Geth, "1.14.3")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	tests := []struct {
// 		ops       string
// 		arch      string
// 		version   string
// 		expectErr bool
// 	}{
// 		{"linux", "amd64", "1.14.4", false},
// 		{"linux", "amd64", ".1.2.3", true},
// 	}

// 	for _, tt := range tests {
// 		_, err := getArtifactSpec(tt.ops, tt.arch, tt.version, paths)
// 		if (err != nil) != tt.expectErr {
// 			t.Errorf("getArtifactSpec() error = %v, expectErr %v\n%v", err, tt.expectErr, tt)
// 		}
// 	}
// }
