package utils

import (
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
		err := ValidateArgs(tt.ops, tt.arch, tt.version)
		if (err != nil) != tt.expectErr {

			t.Errorf("validateArgs() error = %v, expectErr %v\n%v", err, tt.expectErr, tt)
		}
		if err != nil && err.Error() != tt.errMsg {
			t.Errorf("validateArgs() error = %v, errMsg %v\n%v", err, tt.errMsg, tt)
		}
	}
}
