package utils

import (
	"log"
	"path/filepath"
)

// path structs
type Paths struct {
	Directories Directories
	Scripts     Scripts
	Files       Files
}

type Directories struct {
	Root    string
	Rebuild string
	Temp    string
	Geth    string
	Scripts string
	Bin     string
}

type Files struct {
	Travis    string
	Docker    string
	Checksums string
}

type Scripts struct {
	Clone           string
	Checkout        string
	StartDocker     string
	CompareBinaries string
	CopyBinaries    string
}

// Sets project paths.
func SetUpPaths() Paths {
	baseDir, err := GetBaseDir("geth-rebuild")
	if err != nil {
		log.Fatal(err)
	}
	paths := Paths{
		Directories: Directories{
			Root:    baseDir,
			Temp:    filepath.Join(baseDir, "tmp"),
			Geth:    filepath.Join(baseDir, "tmp", "go-ethereum"),
			Scripts: filepath.Join(baseDir, "scripts"),
			Bin:     filepath.Join(baseDir, "bin"),
		},
		Files: Files{
			Travis:    filepath.Join(baseDir, "tmp", "go-ethereum", ".travis.yml"),
			Docker:    filepath.Join(baseDir, "Dockerfile"),
			Checksums: filepath.Join(baseDir, "tmp", "go-ethereum", "build", "checksums.txt"),
		},
		Scripts: Scripts{
			Clone:           filepath.Join(baseDir, "scripts", "clone.sh"),
			Checkout:        filepath.Join(baseDir, "scripts", "checkout.sh"),
			StartDocker:     filepath.Join(baseDir, "scripts", "start_docker.sh"),
			CopyBinaries:    filepath.Join(baseDir, "scripts", "copy_bin.sh"),
			CompareBinaries: filepath.Join(baseDir, "scripts", "compare_bin.sh"),
		},
	}
	return paths
}
