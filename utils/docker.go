package utils

import (
	"fmt"
	"time"
)

func CreateDockerTag(gethVersion string, ops string, arch string) string {
	now := time.Now()
	timestamp := now.Format("2006-01-02-15:04")
	tag := fmt.Sprintf("rebuild-geth-v%s-%s-%s-%s", gethVersion, ops, arch, timestamp)
	return tag
}

// Starts a docker build for dockerfile at `dockerPath` with given `buildArgs`.
func RunDockerBuild(buildArgs map[string]string, dockerTag string, dockerDir string) error {
	// set docker build args
	cmdArgs := []string{"build", "-t", dockerTag, "--progress=plain"} // TODO test tty
	for key, value := range buildArgs {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--build-arg=%s=%s", key, value))
	}
	cmdArgs = append(cmdArgs, dockerDir)
	// run docker build
	_, err := RunCommand("docker", cmdArgs...)
	if err != nil {
		return err
	}
	return nil
}
