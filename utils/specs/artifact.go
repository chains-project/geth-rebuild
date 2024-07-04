package specs

import (
	"fmt"

	"github.com/chains-project/geth-rebuild/utils"
)

type ArtifactSpec struct {
	Version     string
	Os          string
	Arch        string
	Commit      string
	ShortCommit string
}

// Returns configured ArifactSpec.
func NewArtifactSpec(ops string, arch string, version string, paths utils.Paths) (afs ArtifactSpec, err error) {
	_, err = utils.RunCommand(paths.Scripts.Clone, paths.Directories.Temp)
	if err != nil {
		return afs, err
	}
	_, err = utils.RunCommand(paths.Scripts.Checkout, paths.Directories.Geth, version) // TODO what if unstable,
	if err != nil {
		return afs, err
	}
	commit, err := utils.GetGitCommit(paths.Directories.Geth)
	if err != nil {
		return afs, err
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
