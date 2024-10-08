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

// TODO better option than package variables...? e.g. how do i know they are properly set
var ResultsLogDir string
var ResultsBinDir string
var ResultsLogPath string

// Starts a docker rebuild using build configurations in `bi`
func DockerRebuild(bc config.BuildConfig, paths utils.Paths) error {
	// log incomplete rebuild
	err := writeResultsLog(bc, paths, Incomplete, "")
	if err != nil {
		return fmt.Errorf("could not write rebuild results log: %w", err)
	}

	// set docker build args
	cmdArgs := []string{"build", "-t", bc.DockerTag, "--progress=plain"}

	args := bc.GetBuildArgs()
	for key, value := range args {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--build-arg=%s=%s", key, value))
	}

	cmdArgs = append(cmdArgs, paths.Directories.Docker)

	// run docker build
	_, err = utils.RunCommand("docker", cmdArgs...)

	if err != nil {
		_ = writeResultsLog(bc, paths, Error, err.Error()) // ignore any errors in writing logs...
		_ = ProcessLogFile(bc.DockerTag, paths, Error)
		return fmt.Errorf("failed docker build: %w", err)
	}
	return nil
}

// Runs verification script in a Docker container to retrieve rebuild results
// Manipulates the rebuild log's json key `STATUS` : match, mismatch, or error
func RunComparison(bc config.BuildConfig, paths utils.Paths) error {
	ResultsBinDir = filepath.Join(paths.Directories.Bin, bc.DockerTag)
	_, err := utils.RunCommand(paths.Scripts.GetRebuildResults, bc.DockerTag, ResultsBinDir, ResultsLogPath)

	if err != nil { // If script fails, log as error and
		_ = writeResultsLog(bc, paths, Error, err.Error())
		_ = ProcessLogFile(bc.DockerTag, paths, Error)
		return fmt.Errorf("failed docker verification: %w", err)
	}
	return nil
}

// Reads
func ReadRebuildResult() (Status, error) {
	result, err := readParseLog(ResultsLogPath)
	if err != nil {
		return Error, err
	}
	return result.Status, nil
}

// TODO when to send paths and when to not

// Moves logged results file to corresponding status dir - match/mismatch/error
func ProcessLogFile(dockerTag string, paths utils.Paths, status Status) error {
	newDirectory, err := getCategorizedPath(status, dockerTag, paths)
	if err != nil {
		return err
	}

	ResultsLogDir = newDirectory
	err = os.MkdirAll(ResultsLogDir, 0755)
	if err != nil {
		return err
	}

	newPath := filepath.Join(ResultsLogDir, fmt.Sprintf("%s.json", dockerTag))

	err = moveLog(ResultsLogPath, newPath)
	if err != nil {
		return err
	}

	ResultsLogPath = newPath
	fmt.Printf("\nLogged results to %s\n", ResultsLogPath)

	return nil
}
