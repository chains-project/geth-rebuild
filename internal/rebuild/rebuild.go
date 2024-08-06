package rebuild

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chains-project/geth-rebuild/internal/buildconfig"
	"github.com/chains-project/geth-rebuild/internal/utils"
)

type RebuildLog struct {
	Image  string `json:"image"` // TODO name it tag?? and "result"?
	Status string `json:"status"`
	CID    string `json:"cid"`
}

// Starts a reproducing docker build for dockerfile at `dockerDir` using configured build arguments in `bi`
func RunDockerBuild(bi buildconfig.BuildInput) error {
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
		return fmt.Errorf("failed docker build: %w", err)
	}

	return nil
}

// reads and parses json log file generated by docker script
func readParseLogFile(logFile string, dockerTag string) (RebuildLog, error) {
	data, err := os.ReadFile(logFile)
	if err != nil {
		return RebuildLog{}, fmt.Errorf("failed to read log file: %w", err)
	}

	var result RebuildLog
	if err := json.Unmarshal(data, &result); err != nil {
		return RebuildLog{}, fmt.Errorf("failed to unmarshal log data: %w", err)
	}

	return result, nil
}

// TODO think about naming... verify/rebuild/reproduce/compare...
// TODO extract logic to main function...?
// Writes result of a rebuild to logs
func Verify(dockerTag string, paths utils.Paths) (RebuildLog, error) {
	_, err := utils.RunCommand(paths.Scripts.Verify, dockerTag, paths.Directories.Bin, paths.Directories.Logs)
	if err != nil {
		return RebuildLog{}, fmt.Errorf("failed docker verification: %w", err)
	}

	logFile := filepath.Join(paths.Directories.Logs, fmt.Sprintf("%s.json", dockerTag))
	rebuildLog, err := readParseLogFile(logFile, dockerTag)
	if err != nil {
		return RebuildLog{}, err
	}

	err = categorizeRebuild(rebuildLog, dockerTag, logFile, paths)
	if err != nil {
		return RebuildLog{}, err
	}

	return rebuildLog, nil
}

func categorizeRebuild(result RebuildLog, dockerTag string, logFile string, paths utils.Paths) error {
	targetDirectory, err := getTargetDir(result.Status, dockerTag, paths)
	if err != nil {
		return err
	}

	// Create target directory
	if err := os.MkdirAll(targetDirectory, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", targetDirectory, err)
	}

	categorizedPath := filepath.Join(targetDirectory, fmt.Sprintf("%s.json", dockerTag))

	if err := os.Rename(logFile, categorizedPath); err != nil {
		return fmt.Errorf("failed to move log file: %w", err)
	}

	fmt.Printf("\nLog written to %s\n", categorizedPath)
	return nil
}

// creates a directory path that is categorized (match, mismatch) and unique for each rebuild
func getTargetDir(status, dockerTag string, paths utils.Paths) (string, error) {
	switch status {
	case "match":
		return filepath.Join(paths.Directories.Match, dockerTag), nil
	case "mismatch":
		return filepath.Join(paths.Directories.Mismatch, dockerTag), nil
	default:
		return "", fmt.Errorf("error: unexpected rebuild status: %s", status)
	}
}

// Generates a Diffoscope html report for unsuccessful rebuilds identified by their docker tag
func GenerateDiffReport(dockerTag string, paths utils.Paths) error {
	binDir := filepath.Join(paths.Directories.Bin, dockerTag)
	rebuildLogDir := filepath.Join(paths.Directories.Logs, "mismatch", dockerTag)
	diffReportPath := filepath.Join(rebuildLogDir, fmt.Sprintf("%s.html", dockerTag))

	fmt.Printf("\nWriting diff report to %s...", diffReportPath)

	if _, err := utils.RunCommand(paths.Scripts.DiffReport, binDir, diffReportPath); err != nil {
		return fmt.Errorf("failed to run diff report command: %w", err)
	}

	return nil
}
