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
	Root string
	Temp string
	Geth string
	Bin  string
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
			Root: rootDir,
			Temp: filepath.Join(rootDir, "tmp"),
			Geth: filepath.Join(rootDir, "tmp", "go-ethereum"),
			Bin:  filepath.Join(rootDir, "bin"),
		},
		Files: Files{
			Travis:    filepath.Join(rootDir, "tmp", "go-ethereum", ".travis.yml"),
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
