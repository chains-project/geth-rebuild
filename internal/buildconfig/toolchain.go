package buildconfig

import (
	"fmt"
	"os"
	"regexp"

	"github.com/chains-project/geth-rebuild/internal/utils"
)

type ToolchainSpec struct {
	GoVersion string
	CC        string
	BuildCmd  string // TODO move ?
	//CVersion  string // TODO retrieve (from binary) (script inside docker?)
}

// Returns configured rebuild Toolchain specification
func NewToolchainSpec(af ArtifactSpec, paths utils.Paths) (tc ToolchainSpec, err error) {
	goVersion, err := getGoVersion(paths.Files.Checksums)
	if err != nil {
		return tc, fmt.Errorf("failed to get Go version: %w", err)
	}

	cc, err := getCC(af.Os, af.Arch)
	if err != nil {
		return tc, fmt.Errorf("failed to get C compiler: %w", err)
	}

	cmd, err := createBuildCommand(af.Arch, cc)
	if err != nil {
		return tc, fmt.Errorf("failed to get build command: %w", err)
	}

	tc = ToolchainSpec{
		GoVersion: goVersion,
		CC:        cc,
		BuildCmd:  cmd,
	}
	return tc, nil
}

func (tc ToolchainSpec) ToMap() map[string]string {
	return map[string]string{
		"GO_VERSION": tc.GoVersion,
		"C_COMPILER": tc.CC,
		"BUILD_CMD":  tc.BuildCmd,
		//"CVersion":   t.CVersion,
	}
}

func (tc ToolchainSpec) String() string {
	return fmt.Sprintf("ToolchainSpec: (GoVersion:%s, CC:%s, BuildCmd:%s)",
		tc.GoVersion, tc.CC, tc.BuildCmd)
}

// **
// HELPERS
// **

func getArchType(arch string) string {
	switch arch {
	case "arm5", "arm6", "arm7":
		return "arm"
	default:
		return arch
	}
}

// Creates the build command based on target architecture and required c compiler
func createBuildCommand(arch string, cc string) (string, error) {
	cmd := fmt.Sprintf("go run build/ci.go install -dlgo -arch %s -cc %s", getArchType(arch), cc)
	return cmd, nil
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

func getCC(ops string, arch string) (cc string, err error) {
	switch ops {
	case "linux": // TODO this is hard coded
		switch arch {
		case "amd64", "386":
			cc = "gcc-multilib"
		case "arm64":
			cc = "gcc-aarch64-linux-gnu"
		case "arm5", "arm6":
			cc = "gcc-arm-linux-gnueabi"
		case "arm7":
			cc = "gcc-arm-linux-gnueabihf"
		default:
			err = fmt.Errorf("no C compiler found for linux arch `%s`", arch)
		}
	default:
		err = fmt.Errorf("no C compiler found for os `%s`", ops)
	}
	return
}
