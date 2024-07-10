package buildconfig

import (
	"reflect"
	"testing"
)

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

func TestGetUbuntuDeps(t *testing.T) {
	tests := []struct {
		ops       string
		arch      string
		want      []string
		expectErr bool
	}{
		{"linux", "amd64", []string{"git", "ca-certificates", "wget", "binutils", "gcc-multilib"}, false},
		{"linux", "arm64", []string{"git", "ca-certificates", "wget", "binutils", "libc6-dev-arm64-cross"}, false},
		{"linux", "arm42", nil, true},
		{"windows", "amd64", nil, true},
	}

	for _, tt := range tests {
		got, err := getUbuntuDeps(tt.ops, tt.arch)
		if (err != nil) != tt.expectErr {
			t.Errorf("getUbuntuPackages() error = %v, expectErr %v\n%v", err, tt.expectErr, tt)
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("getUbuntuPackages() = %v, want %v\n%v", got, tt.want, tt)
		}
	}
}
