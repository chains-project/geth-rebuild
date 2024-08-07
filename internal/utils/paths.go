package utils

import (
	"log"
	"path/filepath"
)

// Sets up paths that are needed in program

type Paths struct {
	Directories Directories
	Scripts     Scripts
	Files       Files
}

type Directories struct {
	Root     string
	Docker   string
	Temp     string
	Geth     string
	Bin      string
	Logs     string
	Match    string
	Mismatch string
}

type Files struct {
	Travis       string
	Checksums    string
	ReferenceBin string
	ReproduceBin string
}

type Scripts struct {
	Checkout        string
	CompareBinaries string
	DiffReport      string
	StartDocker     string
	Verify          string
}

// Sets up project paths according to predefined values
func SetUpPaths() Paths {
	rootDir, err := GetRootDir("geth-rebuild")
	if err != nil {
		log.Fatal(err)
	}
	paths := Paths{
		Directories: Directories{
			Root:     rootDir,
			Docker:   rootDir,
			Temp:     filepath.Join(rootDir, "tmp"),
			Geth:     filepath.Join(rootDir, "tmp", "go-ethereum"),
			Bin:      filepath.Join(rootDir, "bin"),
			Logs:     filepath.Join(rootDir, "logs"),
			Match:    filepath.Join(rootDir, "logs", "match"),
			Mismatch: filepath.Join(rootDir, "logs", "mismatch"),
		},
		Files: Files{
			Travis:       filepath.Join(rootDir, "tmp", "go-ethereum", ".travis.yml"),
			Checksums:    filepath.Join(rootDir, "tmp", "go-ethereum", "build", "checksums.txt"),
			ReferenceBin: filepath.Join(rootDir, "bin", "geth-reference"),
			ReproduceBin: filepath.Join(rootDir, "bin", "geth-reproduce"),
		},
		Scripts: Scripts{
			Checkout:        filepath.Join(rootDir, "internal", "scripts", "git_checkout.sh"),
			CompareBinaries: filepath.Join(rootDir, "internal", "scripts", "compare_SHA256.sh"),
			DiffReport:      filepath.Join(rootDir, "internal", "scripts", "html_diff_report.sh"),
			StartDocker:     filepath.Join(rootDir, "internal", "scripts", "start_docker.sh"),
			Verify:          filepath.Join(rootDir, "internal", "scripts", "docker_verify.sh"),
		},
	}
	return paths
}
