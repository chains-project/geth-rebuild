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

type BuildInput struct {
	Artifact      ArtifactSpec
	Toolchain     ToolchainSpec
	Environment   EnvSpec
	DockerTag     string
	DockerfileDir string
}

// Configures build input for Docker rebuild
func NewBuildInput(af ArtifactSpec, tc ToolchainSpec, de EnvSpec, paths utils.Paths) BuildInput {
	return BuildInput{
		Artifact:      af,
		Toolchain:     tc,
		Environment:   de,
		DockerTag:     createDockerTag(af),
		DockerfileDir: paths.Directories.Docker,
	}
}

func (bi BuildInput) String() string {
	args := bi.GetBuildArgs()
	var str string = "\n[BUILD ARGUMENTS]\n\n"
	for key, value := range args {
		str += fmt.Sprintf("%s=%s\n", key, value)
	}
	return str
}

// Gathers all relevant build arguments into a string -> string map
func (bi BuildInput) GetBuildArgs() map[string]string {
	buildArgs := make(map[string]string)

	for k, v := range bi.Artifact.ToMap() {
		buildArgs[k] = v
	}
	for k, v := range bi.Toolchain.ToMap() {
		buildArgs[k] = v
	}
	for k, v := range bi.Environment.ToMap() {
		buildArgs[k] = v
	}

	return buildArgs
}

// Creates unique timestamped tag for a Docker rebuild
func createDockerTag(af ArtifactSpec) string {
	version := af.GethVersion

	if af.Unstable {
		version = fmt.Sprintf("%s-%s", af.GethVersion, af.Commit)
	}

	now := time.Now()
	timestamp := now.Format("2006-01-02-15.04")
 	tag := fmt.Sprintf("rebuild-geth-v%s-%s-%s-%s", version, string(af.OS), string(af.Arch), timestamp)
	return tag
}
