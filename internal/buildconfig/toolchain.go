package buildconfig

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/chains-project/geth-rebuild/internal/utils"
)

type ToolchainSpec struct {
	GCVersion    string
	BuildCmd     string
	Dependencies []string
	//CVersion  string // TODO retrieve (from binary) (script inside docker?)
}

// Returns configured rebuild Toolchain specification
func NewToolchainSpec(af ArtifactSpec, paths utils.Paths) (tc ToolchainSpec, err error) {
	goVersion, err := getGCVersion(paths.Files.Checksums)
	if err != nil {
		return tc, fmt.Errorf("failed to get Go version: %w", err)
	}

	deps, err := getToolChainDeps(af.GOOS, af.GOARCH)
	if err != nil {
		return tc, fmt.Errorf("failed to get C compiler: %w", err)
	}

	cmd, err := getBuildCommand(af, paths.Files.Travis)
	if err != nil {
		return tc, fmt.Errorf("failed to get build command: %w", err)
	}

	tc = ToolchainSpec{
		GCVersion:    goVersion,
		Dependencies: deps,
		BuildCmd:     cmd,
	}
	return tc, nil
}

func (tc ToolchainSpec) ToMap() map[string]string {
	return map[string]string{
		"GO_VERSION": tc.GCVersion,
		"TC_DEPS":    strings.Join(tc.Dependencies, " "),
		"BUILD_CMD":  tc.BuildCmd,
		//"CVersion":   t.CVersion,
	}
}

func (tc ToolchainSpec) String() string {
	return fmt.Sprintf("ToolchainSpec: (GoVersion:%s, Dependencies:%s, BuildCmd:%s)",
		tc.GCVersion, tc.Dependencies, tc.BuildCmd)
}

// **
// HELPERS
// **

// Retrieves build command for artifact from travis file
func getBuildCommand(af ArtifactSpec, travisFile string) (string, error) {
	switch af.GOOS {
	case utils.Linux:
		return getLinuxBuildCmd(af, travisFile)
	default:
		return "", fmt.Errorf("no build command retrievable for unsupported os `%s`", string(af.GOOS))
	}
}

// Regexp matches linux build commands for given architecture
func getLinuxBuildCmd(af ArtifactSpec, travisFile string) (string, error) {
	var pattern string

	switch af.GOARCH {
	case utils.AMD64:
		return "go run build/ci.go install -dlgo", nil
	case utils.A386, utils.ARM64:
		pattern = fmt.Sprintf(`go\s*run\s*build/ci\.go\s*install.*-arch\s%s.*`, regexp.QuoteMeta(string(af.GOARCH)))
	case utils.ARM5, utils.ARM6, utils.ARM7:
		v, err := getArmVersion(af.GOOS, af.GOARCH)
		if err != nil {
			return "", err
		}
		pattern = fmt.Sprintf(`%s.go\s*run\s*build/ci\.go\s*install.*`, regexp.QuoteMeta(fmt.Sprintf("GOARM=%s", v)))
	default:
		return "", fmt.Errorf("no build command found for `%s` arch `%s`", af.GOOS, af.GOARCH)
	}
	return findBuildCmdInFile(pattern, travisFile)
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

// Returns the Go `gc` compiler version on form `major.minor.patch` as specified by geth checksum file
func getGCVersion(checksumFile string) (string, error) {
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
func getToolChainDeps(ops utils.OS, arch utils.Arch) ([]string, error) {
	if archDeps, ok := DefaultConfig.ToolchainDeps[ops]; ok {
		if deps, ok := archDeps[arch]; ok {
			return deps, nil
		}
	}
	return nil, fmt.Errorf("no toolchain dependencies found for `%s` `%s`", ops, arch)
}
