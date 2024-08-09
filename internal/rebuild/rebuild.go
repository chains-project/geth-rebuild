package rebuild

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chains-project/geth-rebuild/internal/buildconfig"
	"github.com/chains-project/geth-rebuild/internal/utils"
)

type RebuildResult struct {
	Status string `json:"STATUS"`
}

var ResultsDir string
var ResultsLog string

// Starts a reproducing docker build using configured build arguments in `bi`
func RunDockerBuild(bi buildconfig.BuildInput, paths utils.Paths) error {
	// create a results log and fill with build inputs
	err := createResultsLog(bi, paths)
	if err != nil {
		return fmt.Errorf("could not write rebuild results log: %w", err)
	}

	// set docker build args
	cmdArgs := []string{"build", "-t", bi.DockerTag, "--progress=plain"}

	args := bi.GetBuildArgs()
	for key, value := range args {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--build-arg=%s=%s", key, value))
	}

	cmdArgs = append(cmdArgs, bi.DockerfileDir)

	// run docker build
	_, err = utils.RunCommand("docker", cmdArgs...)
	if err != nil {
		return fmt.Errorf("failed docker build: %w", err)
	}

	return nil
}

// Runs verification script in a Docker container to retrieve rebuild results
// Manipulates the rebuild log's json key `STATUS` : match, mismatch, or error
func RunVerification(dockerTag string, paths utils.Paths) error {
	targetBinDir := filepath.Join(paths.Directories.Bin, dockerTag)
	_, err := utils.RunCommand(paths.Scripts.VerifyResult, dockerTag, targetBinDir, ResultsLog) // TODO handle bindir i.e. remove or something
	if err != nil {
		return fmt.Errorf("failed docker verification: %w", err) // TODO should return error rebuild results?
	}
	return nil
}

// Gets logged rebuild result for the docker tag
func GetRebuildResult(dockerTag string, paths utils.Paths) (RebuildResult, error) {
	result, err := readParseLog(ResultsLog)
	if err != nil {
		return RebuildResult{}, err
	}

	ResultsDir, err = getCategorizedPath(result.Status, dockerTag, paths)
	if err != nil {
		return RebuildResult{}, err
	}

	os.MkdirAll(ResultsDir, 0755)

	NewResultsLog := filepath.Join(ResultsDir, fmt.Sprintf("%s.json", dockerTag))

	err = moveLog(ResultsLog, NewResultsLog)
	if err != nil {
		return RebuildResult{}, err
	}

	ResultsLog = NewResultsLog
	fmt.Printf("\nLogged results to %s\n", ResultsLog)

	return result, nil
}
