package buildconfig

import (
	"testing"
)

func TestGetUbuntuDist(t *testing.T) {
	tests := []struct {
		commit      string
		searchPages int
		want        string
		expectErr   bool
	}{
		{"978041feeaffc3f91afd98c8495a63bfff4b12f4", 25, "bionic", false}, // These will break as commit history grows
		{"42", 5, "", true},
		{"dad8f237ffa5da2c2471f2d9f32d2bc5b580f667", 25, "focal", false},
	}

	for _, tt := range tests {
		got, err := GetUbuntuDist(tt.commit, tt.searchPages)
		if (err != nil) != tt.expectErr {
			t.Errorf("\ngetUbuntuDist() error = %v, expectErr %v\n%v\n%v", err, tt.expectErr, tt, got)
		}
		if got != tt.want {
			t.Errorf("\ngetUbuntuDist() = `%v`, want `%v`\n%v\n", got, tt.want, tt)
		}
	}
}
