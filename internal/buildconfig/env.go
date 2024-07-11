package buildconfig

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/chains-project/geth-rebuild/internal/utils"
)

type EnvSpec struct {
	UbuntuDist   string
	ElfTarget    string
	ArmVersion   string
	Dependencies []string
}

// Returns configured rebuild Environment specification
func NewEnvSpec(af ArtifactSpec, paths utils.Paths) (ub EnvSpec, err error) {
	dist, err := getUbuntuDist(paths.Files.Travis) // TODO !!!
	if err != nil {
		return ub, fmt.Errorf("failed to get Ubuntu distribution: %w", err)
	}

	elfTarget, err := getElfTarget(af.GOOS, af.GOARCH)
	if err != nil {
		return ub, fmt.Errorf("failed to get ELF target: %w", err)
	}

	armV, err := getArmVersion(af.GOOS, af.GOARCH)
	if err != nil {
		return ub, fmt.Errorf("failed to get arm version: %w", err)
	}

	ub = EnvSpec{
		UbuntuDist:   dist,
		ElfTarget:    elfTarget,
		ArmVersion:   armV,
		Dependencies: DefaultConfig.UtilDeps,
	}
	return ub, nil
}

func (ub EnvSpec) ToMap() map[string]string {
	return map[string]string{
		"UBUNTU_DIST": ub.UbuntuDist,
		"ELF_TARGET":  ub.ElfTarget,
		"GOARM":       ub.ArmVersion,
		"UB_DEPS":     strings.Join(ub.Dependencies, " "),
	}
}

func (ub EnvSpec) String() string {
	return fmt.Sprintf("UbuntuSpec: (Dist:%s, ELFTarget:%s, Packages:%v, ARMV: %s)",
		ub.UbuntuDist, ub.ElfTarget, ub.Dependencies, ub.ArmVersion)
}

// **
// HELPERS
// **

// Retrieves Ubuntu distribution as defined in `travisFile` (dist : dddd)
func getUbuntuDist(travisFile string) (dist string, err error) { // TODO this logic does not work due to travis ci issue
	fileContent, err := os.ReadFile(travisFile)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %v", travisFile, err)
	}

	re := regexp.MustCompile(`dist:\s*([a-z]+)`)
	distDefinition := re.Find(fileContent)

	if distDefinition == nil {
		return "", fmt.Errorf("no Ubuntu dist found in file `%s`", travisFile)
	}
	dist = strings.Split(string(distDefinition), ": ")[1]
	return dist, nil

}

func getElfTarget(ops utils.OS, arch utils.Arch) (string, error) {
	if targets, ok := DefaultConfig.ElfTargets[ops]; ok {
		if target, ok := targets[arch]; ok {
			return target, nil
		}
	}
	return "", fmt.Errorf("no elf version found for os `%s` or arch `%s`", ops, arch)
}

func getArmVersion(ops utils.OS, arch utils.Arch) (string, error) {
	if versions, ok := DefaultConfig.ArmVersions[ops]; ok {
		if version, ok := versions[arch]; ok {
			return version, nil
		}
	}
	return "", fmt.Errorf("no GOARM version found for os `%s` or arch `%s`", ops, arch)
}
