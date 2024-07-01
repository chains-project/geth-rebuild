package common

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var commonPackages = []string{"git", "ca-certificates", "wget"}

// Runs command and exits if encountering error.
func RunCommand(cmd string, args ...string) (out string) {
	exeCmd := exec.Command(cmd, args...)

	var outBuffer bytes.Buffer
	multiWriter := io.MultiWriter(os.Stdout, &outBuffer)
	exeCmd.Stdout = multiWriter
	exeCmd.Stderr = multiWriter

	fmt.Println("[CMD]	", printArgs(exeCmd.Args))
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

// Returns root path for basePath. E.g. Users/xxxx/geth-rebuild
func GetBaseDir(basePath string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := wd
	for !strings.HasSuffix(dir, basePath) {
		dir = filepath.Dir(dir)
		if dir == "/" {
			return "", fmt.Errorf("error. cannot find root %s in '%s'", basePath, wd)
		}
	}
	return dir, nil
}
