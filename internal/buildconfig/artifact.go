package buildconfig

import (
	"fmt"
	"os"

	"github.com/chains-project/geth-rebuild/internal/utils"
)

// specifies information about the artifact to rebuild
type ArtifactSpec struct {
	GOOS        utils.OS
	GOARCH      utils.Arch
	Version     string
	Commit      string
	ShortCommit string
}

func (af ArtifactSpec) ToMap() map[string]string {
	return map[string]string{
		"OS":           string(af.GOOS),
		"ARCH":         string(af.GOARCH),
		"GETH_VERSION": af.Version,
		"COMMIT":       af.Commit,
		"SHORT_COMMIT": af.ShortCommit,
	}
}

func (af ArtifactSpec) String() string {
	return fmt.Sprintf("ArtifactSpec: (Version:%s, GOOS:%s, GOARCH:%s, Commit:%s, ShortCommit:%s)",
		af.Version, af.GOOS, af.GOARCH, af.Commit, af.ShortCommit)
}

// Returns configured rebuild Artifact Specification
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
		GOOS:        pa.GOOS,
		GOARCH:      pa.GOARCH,
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

// Invokes script that checks out geth at a tagged version or commit
func checkoutGeth(versionOrCommit string, paths utils.Paths) error {
	_, err := utils.RunCommand(paths.Scripts.Checkout, paths.Directories.Geth, versionOrCommit)
	if err != nil {
		return err
	}
	return nil
}
