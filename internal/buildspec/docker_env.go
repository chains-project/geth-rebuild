package buildspec

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/chains-project/geth-rebuild/internal/utils"
)

type DockerEnvSpec struct {
	UbuntuDist string
	ElfTarget  string
	Packages   []string
	EnvVars    []string
}

func NewDockerSpec(af ArtifactSpec, paths utils.Paths) (ub DockerEnvSpec, err error) {

	// armVersion, err := getArmVersion(string(pa.GOOS), string(pa.GOARCH))
	// if err != nil {
	// 	return af, fmt.Errorf("failed to get GOARM: %w", err)
	// }

	elfTarget, err := getElfTarget(af.Os, af.Arch)
	if err != nil {
		return ub, fmt.Errorf("failed to get ELF target: %w", err)
	}

	packages, err := getUbuntuPackages(af.Os, af.Arch)
	if err != nil {
		return ub, fmt.Errorf("failed to get Ubuntu packages: %w", err)
	}

	dist, err := getUbuntuDist(paths.Files.Travis)
	if err != nil {
		return ub, fmt.Errorf("failed to get Ubuntu distribution: %w", err)
	}

	ub = DockerEnvSpec{
		UbuntuDist: dist,
		ElfTarget:  elfTarget,
		Packages:   packages,
	}
	return ub, nil
}

func (u DockerEnvSpec) ToMap() map[string]string {
	return map[string]string{
		"UBUNTU_DIST": u.UbuntuDist,
		"ELF_TARGET":  u.ElfTarget,
		"PACKAGES":    strings.Join(u.Packages, " "),
	}
}

func (u DockerEnvSpec) String() string {
	return fmt.Sprintf("UbuntuSpec: (Dist:%s, ElfTarget:%s, Packages:%v)",
		u.UbuntuDist, u.ElfTarget, u.Packages)
}

//
// HELPERS
//

// Retrieves Ubuntu distribution as defined in `travisFile` (dist : dddd)
func getUbuntuDist(travisFile string) (dist string, err error) {
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
			elfTarget = "elf64-littleaarch64"
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
func getUbuntuPackages(ops string, arch string) (packages []string, err error) {
	packages = append(packages, "git", "ca-certificates", "wget", "binutils") // common
	switch ops {
	case "linux":
		switch arch {
		case "amd64", "386":
			return // no arch specific packages
		case "arm64":
			packages = append(packages, "libc6-dev-arm64-cross") // TODO hard coded - can use regex ?
		case "arm5", "arm6":
			packages = append(packages, "libc6-dev-armel-cross")
		case "arm7":
			packages = append(packages, "libc6-dev-armhf-cross")
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
