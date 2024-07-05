package rebuild

import (
	"fmt"
	"strings"

	utils "github.com/chains-project/geth-rebuild/internal/utils"
)

type OS string

// specifies information about reproducing artifact
type ArtifactSpec struct {
	Version     string
	Os          string
	Arch        string
	ArmVersion  string
	Commit      string
	ShortCommit string
}

func (a ArtifactSpec) ToMap() map[string]string {
	return map[string]string{
		"GETH_VERSION": a.Version,
		"OS":           a.Os,
		"ARCH":         a.Arch,
		"ARM_V":        a.ArmVersion,
		"COMMIT":       a.Commit,
		"SHORT_COMMIT": a.ShortCommit,
	}
}


func (a ArtifactSpec) PrintSpec() string {
	return fmt.Sprintf("ArtifactSpec: Version=%s, Os=%s, Arch=%s, Commit=%s, ShortCommit=%s",
		a.Version, a.Os, a.Arch, a.Commit, a.ShortCommit)
}

// Returns configured ArifactSpec.
func NewArtifactSpec(ops string, arch string, version string, unstableCommit string, noClone bool, paths utils.Paths) (afs ArtifactSpec, err error) {
	var commit string

	if !noClone {
		err := cloneGethRepo(paths)
		if err != nil {
			return afs, err
		}
	}

	if unstableCommit != "" {
		commit = unstableCommit
		err = checkoutGeth(paths, commit)
		if err != nil {
			return afs, err
		}
	} else {
		err = checkoutGeth(paths, version)
		if err != nil {
			return afs, err
		}

		commit, err = utils.GetGitCommit(paths.Directories.Geth)
		if err != nil {
			return afs, err
		}
	}

	armVersion, err := getArmVersion(ops, arch)

	if err != nil {
		return afs, fmt.Errorf("failed to get GOARM: %w", err)
	}

	afs = ArtifactSpec{
		Version:     version,
		Os:          ops,
		Arch:        arch,
		Commit:      commit,
		ShortCommit: commit[0:8],
		ArmVersion:  armVersion,
	}

	return afs, nil
}

// -- helpers --

// Runs clone script as specified at `script`
func cloneGethRepo(paths utils.Paths) error {
	_, err := utils.RunCommand(paths.Scripts.Clone, paths.Directories.Temp)
	if err != nil {
		return err
	}
	return nil
}

func checkoutGeth(paths utils.Paths, versionOrCommit string) error {
	_, err := utils.RunCommand(paths.Scripts.Checkout, paths.Directories.Geth, versionOrCommit)
	if err != nil {
		return err
	}
	return nil
}

// Returns the ARM version if arch is arm5|arm6
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
