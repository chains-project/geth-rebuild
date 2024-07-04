package build

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/chains-project/geth-rebuild/utils"
)

type UbuntuSpec struct {
	Dist      string
	ElfTarget string
	Packages  []string
}

func NewUbuntuSpec(afs ArtifactSpec, paths utils.Paths) (ub UbuntuSpec, err error) {
	elfTarget, err := getElfTarget(afs.Os, afs.Arch)
	if err != nil {
		return ub, fmt.Errorf("failed to get ELF target: %w", err)
	}

	packages, err := getUbuntuPackages(afs.Os, afs.Arch)
	if err != nil {
		return ub, fmt.Errorf("failed to get Ubuntu packages: %w", err)
	}

	dist, err := getUbuntuDist(paths.Files.Travis)
	if err != nil {
		return ub, fmt.Errorf("failed to get Ubuntu distribution: %w", err)
	}

	ub = UbuntuSpec{
		Dist:      dist,
		ElfTarget: elfTarget,
		Packages:  packages,
	}
	return ub, nil
}

func (u UbuntuSpec) ToMap() map[string]string {
	return map[string]string{
		"UBUNTU_DIST": u.Dist,
		"ELF_TARGET":  u.ElfTarget,
		"PACKAGES":    strings.Join(u.Packages, " "),
	}
}

func (u UbuntuSpec) PrintSpec() string {
	return fmt.Sprintf("UbuntuSpec: Dist=%s, ElfTarget=%s, Packages=%v",
		u.Dist, u.ElfTarget, u.Packages)
}

// -- helpers ---

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
			elfTarget = "elf32-littlearm" // TODO fix wrong target.
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
