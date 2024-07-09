package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// Runs command with print and returns any outputs and errors
func RunCommand(cmd string, args ...string) (out string, err error) {
	exeCmd := exec.Command(cmd, args...)

	var outBuffer bytes.Buffer
	multiWriter := io.MultiWriter(os.Stdout, &outBuffer)
	exeCmd.Stdout = multiWriter
	exeCmd.Stderr = multiWriter

	fmt.Println("[CMD]	", printArgs(exeCmd.Args))

	if err := exeCmd.Run(); err != nil {
		return "", fmt.Errorf("command failed: %s\nerror: %w", printArgs(exeCmd.Args), err)
	}

	return outBuffer.String(), nil
}

// Pretty print args. A little copying from geth.
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

// Returns root path for `base`
func GetRootDir(base string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := wd
	for !strings.HasSuffix(dir, base) {
		dir = filepath.Dir(dir)
		if dir == "/" {
			return "", fmt.Errorf("error. cannot find root %s in '%s'", base, wd)
		}
	}
	return dir, nil
}

// Returns commit hash at latest commit in `dir`
func GetGitCommit(dir string) (string, error) {
	dirFlag := fmt.Sprintf("--git-dir=%s/.git", dir)
	treeFlag := fmt.Sprintf("--work-tree=%s", dir)
	commit, err := RunCommand("git", dirFlag, treeFlag, "log", "-1", "--format=%H")
	if err != nil {
		return "", err
	}
	if commit == "" {
		return "", fmt.Errorf("no commit found in dir %s", dir)
	}
	commit = strings.ReplaceAll(commit, "\n", "")
	return commit, nil
}

// Changes permissions of scripts to `mode`
func ChangePermissions(scripts []string, mode os.FileMode) error { // TODO test
	for _, script := range scripts {
		err := os.Chmod(script, mode)
		if err != nil {
			return fmt.Errorf("error changing permissions for %s: %v", script, err)
		}
	}
	return nil
}


func StartDocker(paths Paths) error {
	_, err := RunCommand(paths.Scripts.StartDocker)
	if err != nil {
		return err
	}
	return nil
}
