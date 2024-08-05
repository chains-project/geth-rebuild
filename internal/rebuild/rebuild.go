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

type Rebuild struct {
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
	cmdArgs = append(cmdArgs, bi.DockerDir)

	// run docker build
	_, err := utils.RunCommand("docker", cmdArgs...)
	if err != nil {
		return fmt.Errorf("failed docker build: %w", err)
	}

	return nil
}

func Verify(dockerTag string, paths utils.Paths) (reproduces bool, err error) { // TODO think about naming...
	_, err = utils.RunCommand(paths.Scripts.Verify, dockerTag, paths.Directories.Bin, paths.Directories.Logs)
	if err != nil {
		return false, fmt.Errorf("failed rebuild verification: %w", err)
	}
	// TODO split up function
	logFile := filepath.Join(paths.Directories.Logs, dockerTag)
	data, err := os.ReadFile(logFile)
	if err != nil {
		log.Fatal(err)
	}

	var result Rebuild
	if err := json.Unmarshal(data, &result); err != nil {
		log.Fatal(err)
	}

	os.Mkdir(paths.Directories.Match, 0755)
	os.Mkdir(paths.Directories.Mismatch, 0755)

	switch result.Status {
	case "match":
		logCategorized := filepath.Join(paths.Directories.Match, dockerTag)
		os.Rename(logFile, logCategorized)
		fmt.Printf("\nLog written to %s\n", logCategorized)
		return true, nil
	case "mismatch":
		logCategorized := filepath.Join(paths.Directories.Mismatch, dockerTag)
		os.Rename(logFile, logCategorized)
		fmt.Printf("\nLog written to %s\n", logCategorized)
		return false, nil
	default:
		return false, fmt.Errorf("error: unknown rebuild status: %s", result.Status)
	}
}

func diff(paths utils.Paths) {
	utils.RunCommand("")
	// TODO handle with a flag --diffoscope

}
