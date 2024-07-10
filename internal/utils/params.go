package utils

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
)

type OS string

type Arch string

const (
	Linux   OS = "linux"
	Darwin  OS = "darwin"
	Windows OS = "windows"
)

const (
	AMD64 Arch = "amd64"
	ARM64 Arch = "arm64"
	ARM5  Arch = "arm5"
	ARM6  Arch = "arm6"
	ARM7  Arch = "arm7"
	A386  Arch = "386"
)

// Map of allowed architectures for each OS
var allowedArch = map[OS][]Arch{
	Linux: {AMD64, ARM64, ARM5, ARM6, ARM7, A386},
}

// Program Args holds parsed input arguments to main program
type ProgramArgs struct {
	GOOS        OS
	GOARCH      Arch
	GethVersion string
	GethDir     string
	Unstable    string
}

func init() {
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <os> <arch> <version> [--geth-dir <geth directory>] [--unstable <commit hash>]\nExample: %s linux amd64 1.14.3\n\n", filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "Reproduce a geth linux binary release.\n")
}

func ParseArgs() (*ProgramArgs, error) {
	if len(os.Args) < 4 {
		usage()
		return nil, fmt.Errorf("not enough arguments")
	}

	pa := &ProgramArgs{
		GOOS:        OS(os.Args[1]),
		GOARCH:      Arch(os.Args[2]),
		GethVersion: os.Args[3],
	}

	// parse optional flags into ProgramArgs struct
	optional := flag.NewFlagSet("optional", flag.ExitOnError)
	optional.StringVar(&pa.GethDir, "geth-dir", "", "Path to the geth directory")
	optional.StringVar(&pa.Unstable, "unstable", "", "Commit hash for unstable version")
	help := optional.Bool("help", false, "Show command documentation")
	optional.Parse(os.Args[4:])

	if *help == true {
		usage()
		os.Exit(1)
	}

	return pa, nil
}

func validOs(os OS) error {
	switch OS(os) {
	case Linux:
		return nil
	case Darwin, Windows:
		return fmt.Errorf("rebuilding not supported for %s", os)
	default:
		return fmt.Errorf("invalid OS `%s`", os)
	}
}

func validArch(os OS, arch Arch) error {
	allowedArchs, ok := allowedArch[os]
	if !ok {
		return fmt.Errorf("no architectures found for OS `%s`", os)
	}
	if !slices.Contains(allowedArchs, arch) {
		return fmt.Errorf("unsupported architecture `%s` for OS `%s`", arch, os)
	}
	return nil
}

// Helper function to validate version format
func validVersion(version string) error {
	versionRegex := `^\d+\.\d+\.\d+$`
	re := regexp.MustCompile(versionRegex)
	if !re.MatchString(version) {
		return fmt.Errorf("error: <version> must be in format `major.minor.patch`\nExample: 1.14.4")
	}
	return nil
}

// Validate input function
func ValidateInput(pa *ProgramArgs) error {
	if err := validOs(pa.GOOS); err != nil {
		return err
	}
	if err := validArch(pa.GOOS, pa.GOARCH); err != nil {
		return err
	}
	if err := validVersion(pa.GethVersion); err != nil {
		return err
	}
	return nil
}

// // Parses command line flags
// func ParseFlags() (ops string, arch string, version string, gethDir string, unstableCommit string) {
// 	flag.Parse()

// 	// Check mandatory positional arguments
// 	if flag.NArg() < 3 {
// 		flag.Usage()
// 		os.Exit(1)
// 	}

// 	gethDir = *gd
// 	unstableCommit = *uc

// 	// Mandatory positional arguments
// 	ops = flag.Arg(0)
// 	arch = flag.Arg(1)
// 	version = flag.Arg(2)

// 	return ops, arch, version, gethDir, unstableCommit
// }

// TODO use enums for allowed os and allowed arch?

// // Validates the program arguments
// func ValidateArgs(ops, arch, version string) error { // TODO validate optional program args.
// 	var validArchs = []string{"amd64", "386", "arm5", "arm6", "arm64", "arm7"}
// 	versionRegex := regexp.MustCompile(`^(\d+\.\d+\.\d+)$`)

// 	if ops != "linux" {
// 		flag.Usage()
// 		return fmt.Errorf("<os> limited to `linux` at the moment")
// 	}
// 	if !utils.Contains(validArchs, arch) {
// 		flag.Usage()
// 		return fmt.Errorf("<arch> must be a valid linux target architecture `amd64|386|arm5|arm6|arm64|arm7`")
// 	}

// 	if !versionRegex.MatchString(version) {
// 		flag.Usage()
// 		return fmt.Errorf("<version> must be in format `major.minor.patch`\nExample: 1.14.4")
// 	}

// 	return nil
// }
