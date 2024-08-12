package utils

import (
	"log"
	"os"
	"path/filepath"
)

// Sets up paths that are needed in program

type Paths struct {
	Directories Directories
	Scripts     Scripts
	Files       Files
}

type Directories struct {
	Root         string
	Docker       string
	Temp         string
	Geth         string
	Bin          string
	Logs         string
	MatchLogs    string
	MismatchLogs string
	ErrorLogs    string
}

type Files struct {
	Travis       string
	Checksums    string
	ReferenceBin string
	ReproduceBin string
}

type Scripts struct {
	GitCheckout        string
	CompareBinaries    string
	GenerateDiffReport string
	StartDocker        string
	GetRebuildResults  string
}

type RebuildPaths struct {
	JSONLog string
	BinDir  string
	LogDir  string
}

// Sets up project paths according to predefined values
func SetUpPaths() Paths {
	rootDir, err := GetRootDir("geth-rebuild")
	if err != nil {
		log.Fatal(err)
	}
	paths := Paths{
		Directories: Directories{
			Root:         rootDir,
			Docker:       rootDir,
			Temp:         filepath.Join(rootDir, "tmp"),
			Geth:         filepath.Join(rootDir, "tmp", "go-ethereum"),
			Bin:          filepath.Join(rootDir, "bin"),
			Logs:         filepath.Join(rootDir, "logs"),
			MatchLogs:    filepath.Join(rootDir, "logs", "match"),
			MismatchLogs: filepath.Join(rootDir, "logs", "mismatch"),
			ErrorLogs:    filepath.Join(rootDir, "logs", "error"),
		},
		Files: Files{
			Travis:       filepath.Join(rootDir, "tmp", "go-ethereum", ".travis.yml"),
			Checksums:    filepath.Join(rootDir, "tmp", "go-ethereum", "build", "checksums.txt"),
			ReferenceBin: filepath.Join(rootDir, "bin", "geth-reference"),
			ReproduceBin: filepath.Join(rootDir, "bin", "geth-reproduce"),
		},
		Scripts: Scripts{
			GitCheckout:        filepath.Join(rootDir, "internal", "scripts", "git_checkout.sh"),
			CompareBinaries:    filepath.Join(rootDir, "internal", "scripts", "compare_binary_SHA.sh"),
			GenerateDiffReport: filepath.Join(rootDir, "internal", "scripts", "gen_diff_report.sh"),
			StartDocker:        filepath.Join(rootDir, "internal", "scripts", "start_docker.sh"),
			GetRebuildResults:  filepath.Join(rootDir, "internal", "scripts", "get_rebuild_results.sh"),
		},
	}
	os.MkdirAll(paths.Directories.Temp, 0755)
	os.MkdirAll(paths.Directories.Bin, 0755)
	os.MkdirAll(paths.Directories.Logs, 0755)
	os.MkdirAll(paths.Directories.MatchLogs, 0755)
	os.MkdirAll(paths.Directories.MismatchLogs, 0755)
	os.MkdirAll(paths.Directories.ErrorLogs, 0755)

	return paths
}
