package utils

import (
	"fmt"
	"regexp"
)

// Validates input arguments to rebuild main program.
func ValidateArgs(ops string, arch string, version string) error {
	var validArchs = []string{"amd64", "386", "arm5", "arm6", "arm64", "arm7"}
	versionRegex := regexp.MustCompile(`^(\d+\.\d+\.\d+)$`)

	if ops != "linux" {
		return fmt.Errorf("<os> limited to `linux` at the moment")
	}
	if !Contains(validArchs, arch) {
		return fmt.Errorf("<arch> must be a valid linux target architecture (amd64|386|arm5|arm6|arm64|arm7)")
	}

	if !versionRegex.MatchString(version) {
		return fmt.Errorf("<geth version> must be in format 'major.minor.patch'\nExample: 1.14.4")
	}

	return nil
}
