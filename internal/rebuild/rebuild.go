package rebuild

import (
	"fmt"

	"github.com/chains-project/geth-rebuild/internal/buildconfig"
	"github.com/chains-project/geth-rebuild/internal/utils"
)

// Starts a reproducing docker build for dockerfile at `dockerDir` using configured build arguments in `bi`
func RunDockerBuild(bi buildconfig.BuildInput) error {
	// set docker build args
	cmdArgs := []string{"build", "-t", bi.DockerTag, "--progress=plain"}

	args := bi.GetBuildArgs()

	for key, value := range args {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--build-arg=%s=%s", key, value))
	}
	cmdArgs = append(cmdArgs, bi.DockerDir)

	// run docker build
	_, err := utils.RunCommand("docker", cmdArgs...)
	if err != nil {
		return fmt.Errorf("failed docker build: %w", err)
	}

	return nil
}

func Verify(dockerTag string, paths utils.Paths) (reproduces bool, err error) {
	_, err = utils.RunCommand(paths.Scripts.CopyBinaries, dockerTag, paths.Directories.Bin)
	if err != nil {
		return false, err
	}
	return false, err
	// run detached mode "copy script"
	// i.e. compare binaries and copy them to local bin...
	
}
