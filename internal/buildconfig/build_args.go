package buildconfig

import (
	"fmt"
	"time"

	"github.com/chains-project/geth-rebuild/internal/utils"
)

type Spec interface {
	ToMap() map[string]string
	String() string
}

type BuildArgs struct {
	Artifact  ArtifactSpec
	Toolchain ToolchainSpec
	DockerEnv EnvSpec
	DockerTag string
	DockerDir string
}

// Configures build input for Docker rebuild
func NewBuildInput(af ArtifactSpec, tc ToolchainSpec, de EnvSpec, paths utils.Paths) BuildArgs {
	return BuildArgs{
		Artifact:  af,
		Toolchain: tc,
		DockerEnv: de,
		DockerTag: createDockerTag(af.Version, af.Os, af.Arch),
		DockerDir: paths.Directories.Docker,
	}
}

func (bi BuildArgs) String() string {
	args := bi.GetBuildArgs()
	var str string = "\n[BUILD ARGUMENTS]\n\n"
	for key, value := range args {
		str += fmt.Sprintf("%s=%s\n", key, value)
	}
	return str
}

// Gathers all build args into a string -> string map
func (bi BuildArgs) GetBuildArgs() map[string]string {
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

// Returns a tag to identify a Docker image build
func createDockerTag(version string, ops string, arch string) string {
	now := time.Now()
	timestamp := now.Format("2006-01-02-15:04")
	tag := fmt.Sprintf("rebuild-geth-v%s-%s-%s-%s", version, ops, arch, timestamp)
	return tag
}
