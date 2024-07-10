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

	elfTarget, err := getElfTarget(af.Os, af.Arch)
	if err != nil {
		return ub, fmt.Errorf("failed to get ELF target: %w", err)
	}

	deps, err := getUbuntuDeps(af.Os, af.Arch)
	if err != nil {
		return ub, fmt.Errorf("failed to get Ubuntu dependencies: %w", err)
	}

	armV, err := getArmVersion(af.Os, af.Arch)
	if err != nil {
		return ub, fmt.Errorf("failed to get arm version: %w", err)
	}

	ub = EnvSpec{
		UbuntuDist:   dist,
		ElfTarget:    elfTarget,
		Dependencies: deps,
		ArmVersion:   armV,
	}
	return ub, nil
}

func (ub EnvSpec) ToMap() map[string]string {
	return map[string]string{
		"UBUNTU_DIST": ub.UbuntuDist,
		"ELF_TARGET":  ub.ElfTarget,
		"UB_DEPS":     strings.Join(ub.Dependencies, " "),
		"GOARM":       ub.ArmVersion,
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

	reDistDef := regexp.MustCompile(`dist:\s*([a-z]+)`)
	distLine := reDistDef.Find(fileContent)

	if distLine == nil {
		return "", fmt.Errorf("no Ubuntu dist found in file `%s`", travisFile)
	}
	dist = strings.Split(string(distLine), ": ")[1]
	return dist, nil

}

// Returns ELF input target for os and arch. Used e.g. for binutils `strip` command.
func getElfTarget(ops string, arch string) (elfTarget string, err error) {
	switch ops {
	case "linux":
		switch arch {
		case "amd64":
			elfTarget = "elf64-x86-64"
		case "386":
			elfTarget = "elf32-i386"
		case "arm64":
			elfTarget = "elf64-littleaarch64" //"elf64-littleaarch64"
		case "arm5", "arm6", "arm7":
			elfTarget = "elf32-little"
		default:
			err = fmt.Errorf("no elf version found for linux arch `%s`", arch)
		}
	default:
		err = fmt.Errorf("no elf version found for os `%s`", ops)
	}
	return
}

// Returns common and architecture specific package for osArc.
func getUbuntuDeps(ops string, arch string) (packages []string, err error) {
	packages = append(packages, "git", "ca-certificates", "wget", "binutils") // common
	switch ops {
	case "linux":
		switch arch {
		case "amd64", "386":
			packages = append(packages, "gcc-multilib")
		case "arm64":
			packages = append(packages, "libc6-dev-arm64-cross", "gcc-aarch64-linux-gnu") // TODO hard coded - can use regex ?
		case "arm5", "arm6":
			packages = append(packages, "libc6-dev-armel-cross", "gcc-arm-linux-gnueabi")
		case "arm7":
			packages = append(packages, "libc6-dev-armhf-cross", "gcc-arm-linux-gnueabihf")
		default:
			return nil, fmt.Errorf("no packages found for linux arch `%s`", arch)
		}
	default:
		return nil, fmt.Errorf("no packages found for os `%s`", ops)
	}
	return
}

// Returns the ARM version if arch is arm5|arm6|arm7
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
			return "", fmt.Errorf("no GOARM command found. Invalid linux arch `%s`", arch)
		}
	default:
		return "", fmt.Errorf("no GOARM command found. Invalid os `%s`", ops)
	}
}
