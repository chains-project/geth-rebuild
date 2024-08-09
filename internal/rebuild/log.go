package rebuild

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chains-project/geth-rebuild/internal/buildconfig"
	"github.com/chains-project/geth-rebuild/internal/utils"
)


// func logStatus(status string) {
// 	// TODO must read whole file
// }

func createResultsLog(bi buildconfig.BuildInput, paths utils.Paths) error {
	ResultsLog = filepath.Join(paths.Directories.Logs, fmt.Sprintf("%s.json", bi.DockerTag))

	args := bi.GetBuildArgs()
	args["STATUS"] = "incomplete"

	data, err := json.MarshalIndent(args, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal build arguments to JSON: %v", err)
	}

	if err := os.WriteFile(ResultsLog, data, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %v", ResultsLog, err)
	}

	return nil
}

// Generates a Diffoscope html report for unsuccessful rebuilds identified by their docker tag
func GenerateDiffReport(dockerTag string, paths utils.Paths) error {
	binDir := filepath.Join(paths.Directories.Bin, dockerTag)
	targetDir, _ := getCategorizedPath("mismatch", dockerTag, paths)
	htmlReport := filepath.Join(targetDir, fmt.Sprintf("%s.html", dockerTag))

	fmt.Print("\nAnalyzing binary differences...")
	if _, err := utils.RunCommand(paths.Scripts.DiffReport, binDir, htmlReport); err != nil {
		return fmt.Errorf("failed to run diffoscope: %w", err)
	}
	fmt.Printf("\nHTML diff report written to %s", htmlReport)

	return nil
}
