package buildconfig

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chains-project/geth-rebuild/internal/utils"
)

var p = utils.SetUpPaths()

func TestGethRepoExists(t *testing.T) {
	existingDir := filepath.Join(p.Directories.Root, "tmp", "existing-directory")
	tests := []struct {
		paths    utils.Paths
		want     bool
		wantErr  bool
		testName string
	}{
		{
			paths: utils.Paths{
				Directories: utils.Directories{
					Geth: existingDir,
				},
			},
			want:     true,
			wantErr:  false,
			testName: "Existing directory",
		},
		{
			paths: utils.Paths{
				Directories: utils.Directories{
					Geth: filepath.Join(p.Directories.Root, "tmp", "non-existing-directory"),
				},
			},
			want:     false,
			wantErr:  false,
			testName: "Non-existing directory",
		},
	}

	err := os.MkdirAll(existingDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory %s: %v", existingDir, err)
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			exists, err := gethRepoExists(tt.paths)
			if (err != nil) != tt.wantErr {
				t.Errorf("gethRepoExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if exists != tt.want {
				t.Errorf("gethRepoExists() = %v, want %v", exists, tt.want)
			}
		})
	}
	err = os.RemoveAll(existingDir)
	if err != nil {
		t.Errorf("gethRepoExists() error removing dir `%s`: %v", existingDir, err)
		return
	}
}
