package utils

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"
)

/* Handles program argument parsing and validation */

type ProgramArgs struct {
	OS          OS
	Arch        Arch
	GethVersion string
	ForceClone  bool
	Unstable    string
	Diff        bool
}

var (
	pa       = &ProgramArgs{}
	optional = flag.NewFlagSet("optional", flag.ExitOnError)
)

func init() {
	// Set up optional program arguments
	optional.Usage = usage
	optional.BoolVar(&pa.ForceClone, "force-clone", false, "Forces a fresh clone of geth repo and removes any existing repo in ./tmp")
	optional.StringVar(&pa.Unstable, "unstable", "", "Rebuilds an unstable build specified by given commit hash\nNote: version number must be correct")
	optional.BoolVar(&pa.Diff, "diff", false, "Write diff report in case of binary mismatch")
}

func usage() {
	msg := "Usage:      gethrebuild OS ARCH VERSION [OPTIONS]\n\nExample:    gethrebuild linux amd64 1.14.3\n\nOptions:\n"
	fmt.Fprint(os.Stderr, msg)
	optional.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n")
}

// Parses program arguments from command line
func ParseArgs() (*ProgramArgs, error) {
	if len(os.Args) < 4 {
		usage()
		return nil, fmt.Errorf("not enough arguments")
	}

	pa.OS = OS(os.Args[1])
	pa.Arch = Arch(os.Args[2])
	pa.GethVersion = os.Args[3]

	optional.Parse(os.Args[4:])

	return pa, nil
}

// Validates mandatory program arguments: os, arch, version
func ValidArgs(pa *ProgramArgs) error {
	if err := validPlatform(pa.OS, pa.Arch); err != nil {
		return err
	}
	if err := validVersion(pa.GethVersion); err != nil {
		return err
	}
	return nil
}

func validPlatform(os OS, arch Arch) error {
	OSAllows, ok := validArchitectures[os]
	if !ok {
		return fmt.Errorf("rebuilding not supported for OS %s", os)
	}
	if !slices.Contains(OSAllows, arch) {
		return fmt.Errorf("unsupported architecture `%s` for OS `%s`", arch, os)
	}
	return nil
}

func validVersion(version string) error {
	versionRegex := `^\d+\.\d+\.\d+$`
	re := regexp.MustCompile(versionRegex)
	if !re.MatchString(version) {
		return fmt.Errorf("error: <version> must be in format `major.minor.patch`\nExample: 1.14.4")
	}
	return nil
}
