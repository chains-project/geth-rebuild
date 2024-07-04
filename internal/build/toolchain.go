package build

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	utils "github.com/chains-project/geth-rebuild/internal/utils"
)

type ToolchainSpec struct {
	GoVersion  string
	CC         string
	BuildCmd   string
	ArmVersion string
	CVersion   string // TODO retrieve (from binary) (script inside docker?)
}

// Returns build configurations for osArch retrieved from build config file (travis.yml).
func NewToolchainSpec(afs ArtifactSpec, paths utils.Paths) (tc ToolchainSpec, err error) {
	goVersion, err := getGoVersion(paths.Files.Checksums)
	if err != nil {
		return tc, fmt.Errorf("failed to get Go version: %w", err)
	}

	cmd, err := getBuildCommand(afs.Os, afs.Arch, paths.Files.Travis)
	if err != nil {
		return tc, fmt.Errorf("failed to get build command: %w", err)
	}

	cc, err := getCC(afs.Os, afs.Arch)
	if err != nil {
		return tc, fmt.Errorf("failed to get C compiler: %w", err)
	}

	env, err := getArmVersion(afs.Os, afs.Arch)
	if err != nil {
		return tc, fmt.Errorf("failed to get GOARM: %w", err)
	}

	tc = ToolchainSpec{
		GoVersion:  goVersion,
		CC:         cc,
		BuildCmd:   cmd,
		ArmVersion: env,
	}
	return tc, nil
}

func (t ToolchainSpec) ToMap() map[string]string {
	return map[string]string{
		"GO_VERSION": t.GoVersion,
		"C_COMPILER": t.CC,
		"BUILD_CMD":  t.BuildCmd,
		"ARM_V":      t.ArmVersion,
		//"CVersion":   t.CVersion,
	}
}

func (t ToolchainSpec) PrintSpec() string {
	return fmt.Sprintf("ToolchainSpec: GoVersion=%s, CC=%s, BuildCmd=%s, CVersion=%s",
		t.GoVersion, t.CC, t.BuildCmd, t.CVersion)
}

// --- helpers ---

func getArmVersion(ops string, arch string) (string, error) {
	switch ops {
	case "linux":
		switch arch {
		case "amd64", "386", "arm64":
			return "", nil
		case "arm5", "arm6", "arm7":
			v := strings.Split(arch, "arm")[1]
			return strings.TrimSpace(v), nil
		default:
			return "", fmt.Errorf("no GOARM command found for linux arch `%s`", arch)
		}
	default:
		return "", fmt.Errorf("no GOARM command found for os `%s`", ops)
	}

}

// Retrieves build commands for os arch in given travis build file (travis.yml). Returns error if not found.
func getBuildCommand(ops string, arch string, travisFile string) (string, error) {
	var pattern string

	switch ops {
	case "linux":
		switch arch {
		case "amd64":
			return "go run build/ci.go install -dlgo", nil
		case "386", "arm64":
			pattern = fmt.Sprintf(`go\s*run\s*build/ci\.go\s*install.*-arch\s%s.*`, regexp.QuoteMeta(arch))
		case "arm5", "arm6", "arm7":
			v, err := getArmVersion(ops, arch)
			if err != nil {
				return "", err
			}
			pattern = fmt.Sprintf(`%s.go\s*run\s*build/ci\.go\s*install.*`, regexp.QuoteMeta(fmt.Sprintf("GOARM=%s", v)))
		default:
			return "", fmt.Errorf("no build command found for linux arch `%s`", arch)
		}
	default:
		return "", fmt.Errorf("no build command found for os `%s`", ops)
	}

	fileContent, err := os.ReadFile(travisFile)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %v", travisFile, err)
	}

	re := regexp.MustCompile(pattern)
	line := re.Find(fileContent)
	if line == nil {
		return "", fmt.Errorf("no build command found for architecture `%s` in file `%s`", arch, travisFile)
	}

	reArm := regexp.MustCompile(`go run\s+(.*)`)
	cmd := reArm.Find(line)

	if cmd == nil {
		return "", fmt.Errorf("no build command found for architecture `%s` from line %s`", arch, line)
	}

	return string(cmd), nil
}

// Returns the Go compiler version form `major.minor.patch` as specified by geth checksumFile.
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

// Returns compiler for osArch as described by compilers map. Returns error if not found.
// TODO this is hard coded
func getCC(ops string, arch string) (cc string, err error) {
	switch ops {
	case "linux":
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