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
	Root   string
	Docker string
	Temp   string
	Geth   string
	Bin    string
}

type Files struct {
	Travis    string
	Checksums string
}

type Scripts struct {
	Clone           string
	Checkout        string
	StartDocker     string
	CopyBinaries    string
	CompareBinaries string
}

// Sets project paths.
func SetUpPaths() Paths {
	rootDir, err := GetRootDir("geth-rebuild")
	if err != nil {
		log.Fatal(err)
	}
	paths := Paths{
		Directories: Directories{
			Root:   rootDir,
			Docker: rootDir,
			Temp:   filepath.Join(rootDir, "tmp"),
			Geth:   filepath.Join(rootDir, "tmp", "go-ethereum"),
			Bin:    filepath.Join(rootDir, "bin"),
		},
		Files: Files{
			Travis:    filepath.Join(rootDir, "tmp", "go-ethereum", ".travis.yml"),
			Checksums: filepath.Join(rootDir, "tmp", "go-ethereum", "build", "checksums.txt"),
		},
		Scripts: Scripts{
			Clone:           filepath.Join(rootDir, "internal", "scripts", "clone.sh"),
			Checkout:        filepath.Join(rootDir, "internal", "scripts", "checkout.sh"),
			StartDocker:     filepath.Join(rootDir, "internal", "scripts", "start_docker.sh"),
			CopyBinaries:    filepath.Join(rootDir, "internal", "scripts", "copy_bin.sh"),
			CompareBinaries: filepath.Join(rootDir, "internal", "scripts", "compare_bin.sh"),
		},
	}
	return paths
}
