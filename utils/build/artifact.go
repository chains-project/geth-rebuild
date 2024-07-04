package build

import (
	"fmt"

	"github.com/chains-project/geth-rebuild/utils"
)

// specifies information about reproduing artifact
type ArtifactSpec struct {
	Version     string
	Os          string
	Arch        string
	Commit      string
	ShortCommit string
	Unstable    bool
}

func cloneGethRepo(script string, destDir string) error {
	_, err := utils.RunCommand(script, destDir)
	if err != nil {
		return err
	}
	return nil
}

func checkoutGeth(script string, gethDir string, versionOrCommit string) error {
	_, err := utils.RunCommand(script, gethDir, versionOrCommit)
	if err != nil {
		return err
	}
	return nil
}

// Returns configured ArifactSpec.
func NewArtifactSpec(ops string, arch string, version string, unstableCommit string, noClone bool, paths utils.Paths) (afs ArtifactSpec, err error) {
	var commit string

	if !noClone {
		err := cloneGethRepo(paths.Scripts.Clone, paths.Directories.Temp)
		if err != nil {
			return afs, err
		}
	}

	if unstableCommit != "" {
		commit = unstableCommit
		err = checkoutGeth(paths.Scripts.Checkout, paths.Directories.Geth, commit)
		if err != nil {
			return afs, err
		}
	} else {
		err = checkoutGeth(paths.Scripts.Checkout, paths.Directories.Geth, version)
		if err != nil {
			return afs, err
		}

		commit, err = utils.GetGitCommit(paths.Directories.Geth)
		if err != nil {
			return afs, err
		}
	}

	afs = ArtifactSpec{
		Version:     version,
		Os:          ops,
		Arch:        arch,
		Commit:      commit,
		ShortCommit: commit[0:8],
	}

	return afs, nil
}

func (a ArtifactSpec) ToMap() map[string]string {
	return map[string]string{
		"GETH_VERSION": a.Version,
		"OS":           a.Os,
		"ARCH":         a.Arch,
		"COMMIT":       a.Commit,
		"SHORT_COMMIT": a.ShortCommit,
	}
}

func (a ArtifactSpec) PrintSpec() string {
	return fmt.Sprintf("ArtifactSpec: Version=%s, Os=%s, Arch=%s, Commit=%s, ShortCommit=%s",
		a.Version, a.Os, a.Arch, a.Commit, a.ShortCommit)
}
