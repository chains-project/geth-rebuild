package util

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Runs command and exits if encountering error.
func RunCommand(dir string, cmd string, args ...string) (out string) {
	exeCmd := exec.Command(cmd, args...)
	exeCmd.Dir = dir // run command in dir

	// catch out in buffer
	var outBuffer bytes.Buffer
	exeCmd.Stdout = &outBuffer

	fmt.Println("[CMD] ", printArgs(exeCmd.Args))
	if err := exeCmd.Run(); err != nil {
		log.Fatal(err)
	}
	return outBuffer.String()
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

func GetRootDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := wd
	for !strings.HasSuffix(dir, "geth-rebuild") {
		dir = filepath.Dir(dir)
		if dir == "/" {
			return "", fmt.Errorf("error. cannot find root geth-rebuild in '%s'", wd)
		}
	}

	return dir, nil
}

func GetArchId(osArch string) string {
	arch := strings.Split(osArch, "-")[1]

	if strings.HasPrefix(arch, "arm") && arch != "arm64" {
		armVersion := strings.TrimLeft(arch, "arm")
		return fmt.Sprintf("ARM=%s", armVersion)
	} else {
		return arch
	}
}
