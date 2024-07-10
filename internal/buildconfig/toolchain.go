package buildconfig

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/chains-project/geth-rebuild/internal/utils"
)

type ToolchainSpec struct {
	GoVersion    string
	Dependencies []string
	BuildCmd     string
	//CVersion  string // TODO retrieve (from binary) (script inside docker?)
}

// Returns configured rebuild Toolchain specification
func NewToolchainSpec(af ArtifactSpec, paths utils.Paths) (tc ToolchainSpec, err error) {
	goVersion, err := getGoVersion(paths.Files.Checksums)
	if err != nil {
		return tc, fmt.Errorf("failed to get Go version: %w", err)
	}

	deps, err := getToolChainDeps(af)
	if err != nil {
		return tc, fmt.Errorf("failed to get C compiler: %w", err)
	}

	cmd, err := getBuildCommand(af, paths.Files.Travis)
	if err != nil {
		return tc, fmt.Errorf("failed to get build command: %w", err)
	}

	tc = ToolchainSpec{
		GoVersion:    goVersion,
		Dependencies: deps,
		BuildCmd:     cmd,
	}
	return tc, nil
}

func (tc ToolchainSpec) ToMap() map[string]string {
	return map[string]string{
		"GO_VERSION": tc.GoVersion,
		"TC_DEPS":    strings.Join(tc.Dependencies, " "),
		"BUILD_CMD":  tc.BuildCmd,
		//"CVersion":   t.CVersion,
	}
}

func (tc ToolchainSpec) String() string {
	return fmt.Sprintf("ToolchainSpec: (GoVersion:%s, Dependencies:%s, BuildCmd:%s)",
		tc.GoVersion, tc.Dependencies, tc.BuildCmd)
}

// **
// HELPERS
// **

// Retrieves build command for artifact from travis file
func getBuildCommand(af ArtifactSpec, travisFile string) (string, error) {
	switch af.Os {
	case "linux":
		return getLinuxBuildCmd(af, travisFile)
	default:
		return "", fmt.Errorf("no build command retrievable for unsupported os `%s`", af.Os)
	}
}

// Regexp matches linux build commands for given architecture
func getLinuxBuildCmd(af ArtifactSpec, travisFile string) (string, error) {
	switch af.Arch {
	case "amd64":
		return "go run build/ci.go install -dlgo", nil
	case "386", "arm64":
		pattern := fmt.Sprintf(`go\s*run\s*build/ci\.go\s*install.*-arch\s%s.*`, regexp.QuoteMeta(af.Arch))
		return findBuildCmdInFile(pattern, travisFile)
	case "arm5", "arm6", "arm7":
		v, err := getArmVersion(af.Os, af.Arch)
		if err != nil {
			return "", err
		}
		pattern := fmt.Sprintf(`%s.go\s*run\s*build/ci\.go\s*install.*`, regexp.QuoteMeta(fmt.Sprintf("GOARM=%s", v)))
		return findBuildCmdInFile(pattern, travisFile)
	default:
		return "", fmt.Errorf("no build command found for `%s` arch `%s`", af.Os, af.Arch)
	}
}

// Finds build command from travis file for a the given pattern
func findBuildCmdInFile(pattern string, travisFile string) (string, error) {
	fileContent, err := os.ReadFile(travisFile)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %v", travisFile, err)
	}

	re := regexp.MustCompile(pattern)
	line := re.Find(fileContent)
	if line == nil {
		return "", fmt.Errorf("no build command found in file `%s` for pattern `%s`", travisFile, pattern)
	}

	reCmd := regexp.MustCompile(`go run\s+(.*)`)
	cmd := reCmd.Find(line)

	if cmd == nil {
		return "", fmt.Errorf("no build command found in file `%s` from line %s`", travisFile, line)
	}
	return string(cmd), nil
}

// Returns the Go compiler version on form `major.minor.patch` as specified by geth checksumFile.
func getGoVersion(checksumFile string) (string, error) {
	fileContent, err := os.ReadFile(checksumFile)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %v", checksumFile, err)
	}

	reTar := regexp.MustCompile(`go(\d+.\d+.(\d+)?).src.tar.gz`)
	goTar := reTar.Find(fileContent)

	if goTar == nil {
		return "", fmt.Errorf("no go version found in file `%s`", checksumFile)
	}

	reVersion := regexp.MustCompile(`(\d+.\d+.(\d+)?)`)
	match := reVersion.Find(goTar)

	if match == nil {

		return "", fmt.Errorf("no go version derivable form line %s in file `%s`", goTar, checksumFile)
	}
	return string(match), nil
}

// Returns required gcc and libc packages for C cross compilation
func getToolChainDeps(af ArtifactSpec) ([]string, error) {
	switch af.Os {
	case "linux":
		switch af.Arch {
		case "amd64", "386":
			return []string{"gcc-multilib"}, nil
		case "arm64":
			return []string{"libc6-dev-arm64-cross", "gcc-aarch64-linux-gnu"}, nil // TODO hard coded - can use regex ?
		case "arm5", "arm6":
			return []string{"libc6-dev-armel-cross", "gcc-arm-linux-gnueabi"}, nil
		case "arm7":
			return []string{"libc6-dev-armhf-cross", "gcc-arm-linux-gnueabihf"}, nil
		default:
			return nil, fmt.Errorf("no packages found for linux arch `%s`", af.Arch)
		}
	default:
		return nil, fmt.Errorf("no packages found for os `%s`", af.Os)
	}
}
