package rebuild

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	config "github.com/chains-project/geth-rebuild/internal/buildconfig"
	"github.com/chains-project/geth-rebuild/internal/utils"
)

type Status string

const (
	Match      Status = "match"
	Mismatch   Status = "mismatch"
	Error      Status = "error"
	Incomplete Status = "incomplete"
)

// Writes the build args for a rebuild and the corresponding results status and any error messages
func writeResultsLog(bc config.BuildConfig, paths utils.Paths, status Status, errMsg string) error {
	ResultsLogPath = filepath.Join(paths.Directories.Logs, fmt.Sprintf("%s.json", bc.DockerTag))

	args := bc.GetBuildArgs()
	args["STATUS"] = string(status)

	if status == Error {
		args["ERROR"] = errMsg
	}

	data, err := json.MarshalIndent(args, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal build arguments to JSON: %v", err)
	}

	if err := os.WriteFile(ResultsLogPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %v", ResultsLogPath, err)
	}

	return nil
}


// Generates a Diffoscope html report for unsuccessful rebuilds identified by their docker tag
func GenerateDiffReport(dockerTag string, paths utils.Paths) error {
	htmlPath := filepath.Join(ResultsLogDir, fmt.Sprintf("%s.html", dockerTag))

	fmt.Print("\nAnalyzing binary differences...")
	if _, err := utils.RunCommand(paths.Scripts.GenerateDiffReport, ResultsBinDir, htmlPath); err != nil {
		return fmt.Errorf("failed to run diffoscope: %w", err)
	}
	fmt.Printf("\nHTML diff report written to %s", htmlPath)

	return nil
}
