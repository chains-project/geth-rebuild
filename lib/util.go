package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Runs command and exits if encountering error.
func RunCommand(dir string, cmd string, args ...string) {
	exeCmd := exec.Command(cmd, args...)
	exeCmd.Dir = dir
	fmt.Println("[CMD] ", printArgs(exeCmd.Args))
	exeCmd.Stderr = os.Stderr
	exeCmd.Stdout = os.Stdout
	if err := exeCmd.Run(); err != nil {
		log.Fatal(err)
	}
}

// Geth function copying.
func printArgs(args []string) string {
	var s strings.Builder
	for i, arg := range args {
		if i > 0 {
			s.WriteByte(' ')
		}
		if strings.IndexByte(arg, ' ') >= 0 {
			arg = strconv.QuoteToASCII(arg)
		}
		s.WriteString(arg)
	}
	return s.String()
}

// Validates input parameters to main program.
func ValidParams(osArch string, gethVersion string) error {
	osArchPattern := "^linux-(amd64|386|arm5|arm6|arm64|arm7)$"
	versionPattern := "^[0-9]+.[0-9]+.[0-9]+$"

	osArchRegex := regexp.MustCompile(osArchPattern)
	versionRegex := regexp.MustCompile(versionPattern)

	if !osArchRegex.MatchString(osArch) {
		return fmt.Errorf("<os-arch> must be a valid linux target architecture\nExample: linux-amd64")
	}
	if !versionRegex.MatchString(gethVersion) {
		return fmt.Errorf("<geth version> must be in format 'major.minor.patch'\nExample: 1.14.4")
	}
	return nil
}
