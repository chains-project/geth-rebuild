package rebuild

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/chains-project/geth-rebuild/internal/buildconfig"
	"github.com/chains-project/geth-rebuild/internal/utils"
)

type RebuildLog struct {
	Image  string `json:"image"` // TODO name it tag??
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

// TODO split up megalong function... make nice.
// TODO think about naming... verify/rebuild/reproduce/compare...
func Verify(dockerTag string, paths utils.Paths) (reproduces bool, err error) {
	_, err = utils.RunCommand(paths.Scripts.Verify, dockerTag, paths.Directories.Bin, paths.Directories.Logs)
	if err != nil {
		return false, fmt.Errorf("failed rebuild verification: %w", err)
	}

	logFile := filepath.Join(paths.Directories.Logs, fmt.Sprintf("%s.json", dockerTag))
	data, err := os.ReadFile(logFile)
	if err != nil {
		log.Fatal(err)
	}

	var result RebuildLog
	if err := json.Unmarshal(data, &result); err != nil {
		log.Fatal(err)
	}

	os.Mkdir(paths.Directories.Match, 0755)
	os.Mkdir(paths.Directories.Mismatch, 0755)

	var logDir string
	switch result.Status {
	case "match":
		reproduces = true
		logDir = filepath.Join(paths.Directories.Match, dockerTag)
	case "mismatch":
		logDir = filepath.Join(paths.Directories.Mismatch, dockerTag)
		reproduces = false
	default:
		return false, fmt.Errorf("error: unknown rebuild status: %s", result.Status)
	}

	os.Mkdir(logDir, 0755)
	logCategorized := filepath.Join(logDir, fmt.Sprintf("%s.json", dockerTag))
	os.Rename(logFile, logCategorized)
	fmt.Printf("\nLog written to %s\n", logCategorized)
	return
}

func DiffReport(dockerTag string, paths utils.Paths) {
	binDir := filepath.Join(paths.Directories.Bin, dockerTag)
	fmt.Printf("\nWriting diff report to %s...", binDir)
	diffLog := filepath.Join(paths.Directories.Logs, "mismatch", dockerTag, fmt.Sprintf("%s.html", dockerTag))
	utils.RunCommand(paths.Scripts.DiffReport, binDir, diffLog)
}
