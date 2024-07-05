package rebuild

import (
	"fmt"
	"time"

	utils "github.com/chains-project/geth-rebuild/internal/utils"
)

type Spec interface {
	ToMap() map[string]string
	PrintSpec() string
}

type BuildInput struct {
	Toolchain ToolchainSpec
	Artifact  ArtifactSpec
	Ubuntu    DockerSpec
	DockerTag string
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

func (bi BuildInput) PrintArgs() {
	args := bi.getBuildArgs()
	var str string = "\n[BUILD ARGUMENTS]\n\n"
	for key, value := range args {
		str += fmt.Sprintf("%s=%s\n", key, value)
	}
	fmt.Println(str)
}

// Starts a reproducing docker build for dockerfile at `dockerDir` using configured build arguments in `bi`
func (bi BuildInput) RunDockerBuild(dockerDir string) error {
	// set docker build args
	cmdArgs := []string{"build", "-t", bi.DockerTag, "--progress=plain"}

	args := bi.getBuildArgs()

	for key, value := range args {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--build-arg=%s=%s", key, value))
	}
	cmdArgs = append(cmdArgs, dockerDir)
	// run docker build
	o, err := utils.RunCommand("docker", cmdArgs...)
	if err != nil {
		return err
	}

	fmt.Println("\nout is:\n\n%s", o)
	return nil
}

// Returns a tag to identify a Docker image build
func CreateDockerTag(version string, ops string, arch string) string {
	now := time.Now()
	timestamp := now.Format("2006-01-02-15:04")
	tag := fmt.Sprintf("rebuild-geth-v%s-%s-%s-%s", version, ops, arch, timestamp)
	return tag
}
