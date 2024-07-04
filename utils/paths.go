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
	rootDir, err := GetRootDir("geth-rebuild")
	if err != nil {
		log.Fatal(err)
	}
	paths := Paths{
		Directories: Directories{
			Root:    rootDir,
			Temp:    filepath.Join(rootDir, "tmp"),
			Geth:    filepath.Join(rootDir, "tmp", "go-ethereum"),
			Scripts: filepath.Join(rootDir, "scripts"),
			Bin:     filepath.Join(rootDir, "bin"),
		},
		Files: Files{
			Travis:    filepath.Join(rootDir, "tmp", "go-ethereum", ".travis.yml"),
			Docker:    filepath.Join(rootDir, "Dockerfile"),
			Checksums: filepath.Join(rootDir, "tmp", "go-ethereum", "build", "checksums.txt"),
		},
		Scripts: Scripts{
			Clone:           filepath.Join(rootDir, "scripts", "clone.sh"),
			Checkout:        filepath.Join(rootDir, "scripts", "checkout.sh"),
			StartDocker:     filepath.Join(rootDir, "scripts", "start_docker.sh"),
			CopyBinaries:    filepath.Join(rootDir, "scripts", "copy_bin.sh"),
			CompareBinaries: filepath.Join(rootDir, "scripts", "compare_bin.sh"),
		},
	}
	return paths
}
