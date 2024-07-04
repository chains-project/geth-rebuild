package build

import (
	"fmt"

	"github.com/chains-project/geth-rebuild/utils"
)

type BuildInput struct {
	Toolchain ToolchainSpec
	Artifact  ArtifactSpec
	Ubuntu    UbuntuSpec
	DockerTag string
}

type Spec interface {
	ToMap() map[string]string
	PrintSpec() string
}

func (bi BuildInput) getBuildArgs() map[string]string {
	buildArgs := make(map[string]string)

	for k, v := range bi.Artifact.ToMap() {
		buildArgs[k] = v
	}
	for k, v := range bi.Toolchain.ToMap() {
		buildArgs[k] = v
	}
	for k, v := range bi.Ubuntu.ToMap() {
		buildArgs[k] = v
	}

	return buildArgs
}

func (bi BuildInput) printArgs(args ...map[string]string) string {
	var str string
	for _, arg := range args {
		for key, value := range arg {
			str += fmt.Sprintf("%s=%s\n", key, value)
		}
	}
	return str
}

// Starts a reproduing docker build for dockerfile at `dockerDir` using configured build argument in `bi`
func RunDockerBuild(bi BuildInput, dockerDir string) error {
	// set docker build args
	cmdArgs := []string{"build", "-t", bi.DockerTag, "--progress=plain"} // TODO test tty

	args := bi.getBuildArgs()
	bi.printArgs(args)

	for key, value := range args {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--build-arg=%s=%s", key, value))
	}
	cmdArgs = append(cmdArgs, dockerDir)
	// run docker build
	_, err := utils.RunCommand("docker", cmdArgs...)
	if err != nil {
		return err
	}
	return nil
}
