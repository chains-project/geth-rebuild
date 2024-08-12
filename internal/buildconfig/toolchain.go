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
	BuildCmd     string
	Dependencies []string
}

// Returns configured rebuild Toolchain specification
func NewToolchainSpec(af ArtifactSpec, paths utils.Paths) (tc ToolchainSpec, err error) {
	goVersion, err := getGoVersion(paths.Files.Checksums)
	if err != nil {
		return tc, fmt.Errorf("failed to get Go version: %w", err)
	}

	deps, err := getToolChainDeps(af.OS, af.Arch)
	if err != nil {
		return tc, fmt.Errorf("failed to get C compiler: %w", err)
	}

	cmd, err := getBuildCommand(af.OS, af.Arch, paths.Files.Travis)
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

// Returns a string -> string map with the toolchain specific build arguments needed in the docker rebuild
func (tc ToolchainSpec) ToMap() map[string]string {
	return map[string]string{
		"GO_VERSION":     tc.GoVersion,
		"TOOLCHAIN_DEPS": strings.Join(tc.Dependencies, " "),
		"BUILD_CMD":      tc.BuildCmd,
	}
}

func (tc ToolchainSpec) String() string {
	return fmt.Sprintf("ToolchainSpec: (GoVersion:%s, Dependencies:%s, BuildCmd:%s)",
		tc.GoVersion, tc.Dependencies, tc.BuildCmd)
}

// Retrieves build command for artifact from travis file
func getBuildCommand(ops utils.OS, arch utils.Arch, travisYML string) (string, error) {
	switch ops {
	case utils.Linux:
		return getLinuxBuildCmd(ops, arch, travisYML)
	default:
		return "", fmt.Errorf("no build command retrievable for unsupported os `%s`", ops)
	}
}

// Regexp matches linux build commands for given architecture
func getLinuxBuildCmd(ops utils.OS, arch utils.Arch, travisYML string) (string, error) { // TODO messy functions ahead...
	var pattern string

	switch arch {
	case utils.AMD64:
		return "go run build/ci.go install -dlgo", nil
	case utils.A386, utils.ARM64:
		pattern = fmt.Sprintf(`go\s*run\s*build/ci\.go\s*install.*-arch\s%s.*`, regexp.QuoteMeta(string(arch)))
	case utils.ARM5, utils.ARM6, utils.ARM7:
		v, err := getArmVersion(ops, arch)
		if err != nil {
			return "", err
		}
		pattern = fmt.Sprintf(`%s.go\s*run\s*build/ci\.go\s*install.*`, regexp.QuoteMeta(fmt.Sprintf("GOARM=%s", v)))
	default:
		return "", fmt.Errorf("no build command found for `%s` arch `%s`", ops, arch)
	}
	return findBuildCmdInFile(pattern, travisYML)
}

// Finds build command from travis file for a the given pattern
func findBuildCmdInFile(pattern string, travisYML string) (string, error) {
	fileContent, err := os.ReadFile(travisYML)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %v", travisYML, err)
	}

	re := regexp.MustCompile(pattern)
	line := re.Find(fileContent)
	if line == nil {
		return "", fmt.Errorf("no build command found in file `%s` for pattern `%s`", travisYML, pattern)
	}

	reCmd := regexp.MustCompile(`go run\s+(.*)`)
	cmd := reCmd.Find(line)

	if cmd == nil {
		return "", fmt.Errorf("no build command found in file `%s` from line %s`", travisYML, line)
	}
	return string(cmd), nil
}

// Returns the Go `gc` compiler version on form `major.minor.patch` as specified by geth checksum file
func getGoVersion(checksumFile string) (string, error) {
	fileContent, err := os.ReadFile(checksumFile)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %v", checksumFile, err)
	}

	reTar := regexp.MustCompile(`go(\d+.\d+.(\d+)?).src.tar.gz`) // TODO this matches also ppa-builder version
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

// Returns required gcc and libc packages for C compilation
func getToolChainDeps(ops utils.OS, arch utils.Arch) ([]string, error) {
	if archDeps, ok := DefaultConfig.ToolchainDeps[ops]; ok {
		if deps, ok := archDeps[arch]; ok {
			return deps, nil
		}
	}
	return nil, fmt.Errorf("no toolchain dependencies found for `%s` `%s`", ops, arch)
}
