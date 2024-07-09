package buildconfig

import (
	"fmt"

	utils "github.com/chains-project/geth-rebuild/internal/utils"
)

// specifies information about the artifact to rebuild
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

func (a ArtifactSpec) String() string {
	return fmt.Sprintf("ArtifactSpec: (Version:%s, Os:%s, Arch:%s, Commit:%s, ShortCommit:%s)",
		a.Version, a.Os, a.Arch, a.Commit, a.ShortCommit)
}

// Returns configured Artifact Specification
func NewArtifactSpec(pa *ProgramArgs, paths utils.Paths) (af ArtifactSpec, err error) {
	if pa.GethDir == "" {
		err := cloneGethRepo(paths)
		if err != nil {
			return af, err
		}
	}

	var commit string
	if pa.Unstable == "" { // stable release, check out version tag
		err = checkoutGeth(pa.GethVersion, paths)
		if err != nil {
			return af, err
		}
		commit, err = utils.GetGitCommit(paths.Directories.Geth)
		if err != nil {
			return af, err
		}

	} else { // unstable build, check out at commit
		err = checkoutGeth(pa.Unstable, paths)
		if err != nil {
			return af, err
		}
		commit = pa.Unstable
	}

	af = ArtifactSpec{
		Version:     pa.GethVersion,
		Os:          string(pa.GOOS),
		Arch:        string(pa.GOARCH),
		Commit:      commit,
		ShortCommit: commit[0:8],
	}

	return af, nil
}

// **
// HELPERS
// **

// Runs clone script as specified in paths struct
func cloneGethRepo(paths utils.Paths) error {
	_, err := utils.RunCommand(paths.Scripts.Clone, paths.Directories.Temp)
	if err != nil {
		return err
	}
	return nil
}

// Checks out geth at a tagged version or a commit
func checkoutGeth(versionOrCommit string, paths utils.Paths) error {
	_, err := utils.RunCommand(paths.Scripts.Checkout, paths.Directories.Geth, versionOrCommit)
	if err != nil {
		return err
	}
	return nil
}
