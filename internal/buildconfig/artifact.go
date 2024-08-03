package buildconfig

import (
	"fmt"
	"os"

	"github.com/chains-project/geth-rebuild/internal/utils"
)

// specifies artifact info
type ArtifactSpec struct {
	OS          utils.OS
	Arch        utils.Arch
	GethVersion string
	Commit      string
	ShortCommit string
}

func (af ArtifactSpec) ToMap() map[string]string {
	return map[string]string{
		"OS":           string(af.OS),
		"ARCH":         string(af.Arch),
		"GETH_VERSION": af.GethVersion,
		"COMMIT":       af.Commit,
		"SHORT_COMMIT": af.ShortCommit,
	}
}

func (af ArtifactSpec) String() string {
	return fmt.Sprintf("ArtifactSpec: (Version:%s, GOOS:%s, GOARCH:%s, Commit:%s, ShortCommit:%s)",
		af.GethVersion, af.OS, af.Arch, af.Commit, af.ShortCommit)
}

// Returns configured rebuild Artifact Specification
func NewArtifactSpec(pa *utils.ProgramArgs, paths utils.Paths) (af ArtifactSpec, err error) {
	var commit string

	exists, err := gethRepoExists(paths)
	if err != nil {
		return af, err
	}

	if !exists || pa.ForceClone {
		err := cloneGethRepo(paths)
		if err != nil {
			return af, err
		}
	}

	if pa.Unstable == "" { // stable release, check out version tag
		err = checkoutGeth(fmt.Sprintf("v%s", pa.GethVersion), paths)
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
		GethVersion: pa.GethVersion,
		OS:          pa.OS,
		Arch:        pa.Arch,
		Commit:      commit,
		ShortCommit: commit[0:8],
	}

	return af, nil
}

// Indicates if directory /tmp/go-ethereum exists
func gethRepoExists(paths utils.Paths) (bool, error) {
	_, err := os.Stat(paths.Directories.Geth)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Clones geth into /tmp, removes any existing geth directory
func cloneGethRepo(paths utils.Paths) error {
	url := "https://github.com/ethereum/go-ethereum.git"
	branch := "master"
	fmt.Printf("\nCloning go ethereum branch %s from %s\n\n", branch, url)

	// create /tmp if not existing
	err := os.MkdirAll(paths.Directories.Temp, 0755)
	if err != nil {
		return err
	}

	// remove any existing repo
	exists, err := gethRepoExists(paths)
	if err != nil {
		return err
	}

	if exists {
		err = os.RemoveAll(paths.Directories.Geth)
		if err != nil {
			return err
		}
	}

	// clone
	_, err = utils.RunCommand("git", "clone", "-v", "--branch", branch, url, paths.Directories.Geth)
	if err != nil {
		return fmt.Errorf("failed to clone geth sources: %w", err)
	}
	return nil
}

// Invokes script that checks out geth at a tagged version or commit
func checkoutGeth(versionOrCommit string, paths utils.Paths) error {
	// TODO remove script and run cmd...
	_, err := utils.RunCommand(paths.Scripts.Checkout, paths.Directories.Geth, versionOrCommit)
	if err != nil {
		return err
	}
	return nil
}
