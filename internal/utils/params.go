package utils

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// TODO params

var (
	// Optional flags
	// TODO error, these are not parsed...
	gd = flag.String("geth-dir", "", "Skip cloning the geth source repository and use specified directory instead")
	uc = flag.String("unstable", "", "Specify an unstable version hash")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <os> <arch> <version> [--geth-dir <geth directory>] [--unstable <commit hash>]\nExample: linux amd64 1.14.3\n\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nReproduce a geth linux binary release.\n")
	}
}

// Parses command line flags
func ParseFlags() (ops string, arch string, version string, gethDir string, unstableCommit string) {
	flag.Parse()

	// Check mandatory positional arguments
	if flag.NArg() < 3 {
		flag.Usage()
		os.Exit(1)
	}

	gethDir = *gd
	unstableCommit = *uc

	// Mandatory positional arguments
	ops = flag.Arg(0)
	arch = flag.Arg(1)
	version = flag.Arg(2)

	return ops, arch, version, gethDir, unstableCommit
}

// Validates the program arguments
func ValidateArgs(ops, arch, version string) error { // TODO validate optional program args.
	var validArchs = []string{"amd64", "386", "arm5", "arm6", "arm64", "arm7"}
	versionRegex := regexp.MustCompile(`^(\d+\.\d+\.\d+)$`)

	if ops != "linux" {
		flag.Usage()
		return fmt.Errorf("<os> limited to `linux` at the moment")
	}
	if !Contains(validArchs, arch) {
		flag.Usage()
		return fmt.Errorf("<arch> must be a valid linux target architecture `amd64|386|arm5|arm6|arm64|arm7`")
	}

	if !versionRegex.MatchString(version) {
		flag.Usage()
		return fmt.Errorf("<version> must be in format `major.minor.patch`\nExample: 1.14.4")
	}

	return nil
}
