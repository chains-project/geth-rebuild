package buildconfig

import (
	"fmt"
	"strings"

	"github.com/chains-project/geth-rebuild/internal/utils"
)

type EnvSpec struct {
	UbuntuDist   string
	ArtifactSlug string
	Flags        FlagSpec
	Dependencies []string
}

type FlagSpec struct {
	CGO_ENABLED string
	ElfTarget   string
	ArmVersion  string
}

// Returns configured rebuild Environment specification
func NewEnvSpec(af ArtifactSpec, paths utils.Paths) (env EnvSpec, err error) {
	dist, err := GetUbuntuDist(af.Commit, 20) // TODO decide search length...
	if err != nil {
		return env, fmt.Errorf("failed to get Ubuntu distribution: %w", err)
	}

	envFlags, err := newFlagSpec(af)
	if err != nil {
		return env, fmt.Errorf("failed to set environment flag specification: %w", err)
	}

	env = EnvSpec{
		UbuntuDist:   dist,
		ArtifactSlug: createSlug(af),
		Flags:        envFlags,
		Dependencies: DefaultConfig.UtilDeps,
	}
	return env, nil
}

// Returns a string -> string map with the environment specific build arguments needed in the docker rebuild
func (env EnvSpec) ToMap() map[string]string {
	return map[string]string{
		"UBUNTU_DIST":   env.UbuntuDist,
		"ARTIFACT_SLUG": env.ArtifactSlug,
		"CGO_ENABLED":   env.Flags.CGO_ENABLED,
		"ELF_TARGET":    env.Flags.ElfTarget,
		"GOARM":         env.Flags.ArmVersion,
		"UTIL_DEPS":     strings.Join(env.Dependencies, " "),
	}
}

func (env EnvSpec) String() string {
	return fmt.Sprintf("Environment specification: (Ubuntu dist:%s, CGO_ENABLED:%s, ELF target:%s, ARM version: %s, Util dependencies:%v)",
		env.UbuntuDist, env.Flags.CGO_ENABLED, env.Flags.ElfTarget, env.Flags.ArmVersion, env.Dependencies)
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

func newFlagSpec(af ArtifactSpec) (FlagSpec, error) {
	var CGOEnabled string
	if af.Arch == utils.AMD64 {
		CGOEnabled = "1"
	} else {
		CGOEnabled = "0" // cross compilation - CGO will be enabled through geth source code (using /build/ci.go)
	}

	version, err := getArmVersion(af.OS, af.Arch)
	if err != nil {
		return FlagSpec{}, fmt.Errorf("failed to get arm version: %w", err)
	}

	elfTarget, err := getElfTarget(af.OS, af.Arch)
	if err != nil {
		return FlagSpec{}, fmt.Errorf("failed to get ELF target: %w", err)
	}

	return FlagSpec{CGO_ENABLED: CGOEnabled, ArmVersion: version, ElfTarget: elfTarget}, nil

}

func createSlug(af ArtifactSpec) string {
	var version = af.GethVersion
	if af.Unstable {
		version = fmt.Sprintf("%s-unstable", af.GethVersion)
	}
	return fmt.Sprintf("geth-%s-%s-%s-%s", string(af.OS), string(af.Arch), version, af.ShortCommit)
}
