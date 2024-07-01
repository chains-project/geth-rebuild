package main

import (
	"fmt"
	"testing"
)

// valid params

func TestValidArgs(t *testing.T) {
	tests := []struct {
		osArch      string
		gethVersion string
		expectError bool
	}{
		{"linux-amd64", "1.14.3", false},     // valid
		{"linux-arm64", "v1.0.0", true},      // invalid gethVersion
		{"linux-arm5", "1.2", true},          // invalid gethVersion
		{"darwin-amd64", "1.14.3", true},     // invalid osArch
		{"linux-386", "1.14.3", false},       // valid
		{"linux-arm7", "2.0.0-beta", true},   // invalid gethVersion
		{"linux-arm6", "1.14.3", false},      // valid
		{"linux-arm64", "1.14.3-rc.1", true}, // invalid gethVersion
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s-%s", tc.osArch, tc.gethVersion), func(t *testing.T) {
			err := validArgs(tc.osArch, tc.gethVersion)
			if tc.expectError {
				if err == nil {
					t.Errorf("expected error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
