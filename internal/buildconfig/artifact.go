package buildconfig

import (
	"fmt"
	"os"

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
func NewArtifactSpec(pa *utils.ProgramArgs, paths utils.Paths) (af ArtifactSpec, err error) {
	if !pa.NoClone {
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

// Clones geth into /tmp, removes any existing geth directory
func cloneGethRepo(paths utils.Paths) error {
	url := "https://github.com/ethereum/go-ethereum.git"
	branch := "master"
	fmt.Printf("\nCloning go ethereum branch %s from %s\n\n", branch, url)

	_, err := utils.RunCommand("mkdir", "-p", paths.Directories.Temp)
	if err != nil {
		return err
	}
	if _, err := os.Stat(paths.Directories.Geth); !os.IsNotExist(err) {
		_, err := utils.RunCommand("rm", "-rf", paths.Directories.Geth)
		if err != nil {
			return err
		}
	}
	_, err = utils.RunCommand("git", "clone", "-v", "--branch", branch, url, paths.Directories.Geth)
	if err != nil {
		return fmt.Errorf("failed to clone geth sources: %w", err)
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
