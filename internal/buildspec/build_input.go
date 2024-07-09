package buildspec

import (
	"fmt"
	"time"

	"github.com/chains-project/geth-rebuild/internal/utils"
)

type BuildInput struct {
	Artifact  ArtifactSpec
	Toolchain ToolchainSpec
	DockerEnv DockerEnvSpec
	DockerTag string
	DockerDir string
}

func (bi BuildInput) GetBuildArgs() map[string]string {
	buildArgs := make(map[string]string)

	for k, v := range bi.Artifact.ToMap() {
		buildArgs[k] = v
	}
	for k, v := range bi.Toolchain.ToMap() {
		buildArgs[k] = v
	}
	for k, v := range bi.DockerEnv.ToMap() {
		buildArgs[k] = v
	}

	return buildArgs
}

func (bi BuildInput) String() string {
	args := bi.GetBuildArgs()
	var str string = "\n[BUILD ARGUMENTS]\n\n"
	for key, value := range args {
		str += fmt.Sprintf("%s=%s\n", key, value)
	}
	return str
}

func NewBuildInput(af ArtifactSpec, tc ToolchainSpec, de DockerEnvSpec, paths utils.Paths) BuildInput {
	return BuildInput{
		Artifact:  af,
		Toolchain: tc,
		DockerEnv: de,
		DockerTag: createDockerTag(af.Version, af.Os, af.Arch),
		DockerDir: paths.Directories.Docker,
	}
}

// Returns a tag to identify a Docker image build
func createDockerTag(version string, ops string, arch string) string {
	now := time.Now()
	timestamp := now.Format("2006-01-02-15:04")
	tag := fmt.Sprintf("rebuild-geth-v%s-%s-%s-%s", version, ops, arch, timestamp)
	return tag
}
