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

// Runs command and exits if encountering error.
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

// Returns commit hash at latest commit in dir.
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

/* // Runs command chmod `permission` script for each script in `scripts`. Exits if error occurs.
func ChangePermissions(scripts []string, permission string) {
	for _, script := range scripts {
		//util.RunCommand("cat", script)
		RunCommand("chmod", permission, script)
		fmt.Printf("\nPermissions changed (%s) for file %s", permission, script)
	}
}
*/
// Utility function to change permissions of scripts
func ChangePermissions(scripts []string, mode os.FileMode) error { // TODO test
	for _, script := range scripts {
		err := os.Chmod(script, mode)
		if err != nil {
			return fmt.Errorf("error changing permissions for %s: %v", script, err)
		}
	}
	return nil
}

// Contains for string slices
func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
