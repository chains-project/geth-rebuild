package rebuild

import (
	"fmt"
	"os"
	"path/filepath"

	config "github.com/chains-project/geth-rebuild/internal/buildconfig"
	"github.com/chains-project/geth-rebuild/internal/utils"
)

type RebuildResult struct {
	Status Status `json:"STATUS"`
}

var TargetLogDir string
var TargetBinDir string
var ResultsLogPath string

// Starts a reproducing docker build using configured build arguments in `bi`
func RunDockerBuild(bi config.BuildConfig, paths utils.Paths) error {
	// set docker build args
	cmdArgs := []string{"build", "-t", bi.DockerTag, "--progress=plain"}

	args := bi.GetBuildArgs()
	for key, value := range args {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--build-arg=%s=%s", key, value))
	}

	cmdArgs = append(cmdArgs, bi.DockerfileDir)

	// run docker build
	_, err := utils.RunCommand("docker", cmdArgs...)
	if err != nil {
		_ = logResults(bi, Error, paths) // ignore any errors here
		_, _ = CategorizeRebuild(bi.DockerTag, paths)
		return fmt.Errorf("failed docker build: %w", err)
	}
	return nil
}

// Runs verification script in a Docker container to retrieve rebuild results
// Manipulates the rebuild log's json key `STATUS` : match, mismatch, or error
func RunVerification(bi config.BuildConfig, dockerTag string, paths utils.Paths) error {
	err := logResults(bi, Incomplete, paths)
	if err != nil {
		return fmt.Errorf("could not write rebuild results log: %w", err)
	}

	TargetBinDir = filepath.Join(paths.Directories.Bin, dockerTag)

	_, err = utils.RunCommand(paths.Scripts.VerifyResult, dockerTag, TargetBinDir, ResultsLogPath)
	if err != nil {
		_ = logResults(bi, Error, paths) // TODO specify which error in the log, optional arg
		_, _ = CategorizeRebuild(bi.DockerTag, paths)
		return fmt.Errorf("failed docker verification: %w", err)
	}
	return nil
}

// Gets logged rebuild result for the docker tag
func CategorizeRebuild(dockerTag string, paths utils.Paths) (RebuildResult, error) {
	result, err := readParseLog(ResultsLogPath)
	if err != nil {
		return RebuildResult{Status: Error}, err
	}

	TargetLogDir, err = getCategorizedPath(result.Status, dockerTag, paths)
	if err != nil {
		return result, err
	}

	os.MkdirAll(TargetLogDir, 0755)
	newPath := filepath.Join(TargetLogDir, fmt.Sprintf("%s.json", dockerTag))

	err = moveLog(ResultsLogPath, newPath)
	if err != nil {
		return result, err
	}

	ResultsLogPath = newPath
	fmt.Printf("\nLogged results to %s\n", ResultsLogPath)

	return result, nil
}
