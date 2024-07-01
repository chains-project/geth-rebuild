package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	util "github.com/chains-project/geth-rebuild/util"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: <os-arch> <geth version>")
		fmt.Println("Example: linux-amd64 1.14.3") // TODO should change input params to os arch ?
		os.Exit(1)
	}

	// 1. Validate input parameters
	osArch := os.Args[1]
	gethVersion := os.Args[2]
	if err := validParams(osArch, gethVersion); err != nil {
		log.Fatal(err)
	}
	// 2. Set directory parameters
	rootDir, err := util.GetBaseDir("geth-rebuild")

	if err != nil {
		log.Fatal(err)
	}
	rebuildDir := rootDir + "/cmd/rebuild" // TODO fix path setting/retrieval more neat/logical.
	tmpDir := rebuildDir + "/tmp"
	gethDir := tmpDir + "/go-ethereum"
	travisPath := tmpDir + "/.travis.yml"

	// 3. CHMOD scripts
	cloneGeth := filepath.Join(rebuildDir, "scripts/clone_geth.sh")
	startDocker := filepath.Join(rebuildDir, "scripts/start_docker.sh")
	rebuild := filepath.Join(rebuildDir, "scripts/rebuild.sh")
	var scripts = []string{cloneGeth, startDocker, rebuild}
	fmt.Printf("\n[CHANGING FILE PERMISSIONS FOR EXECUTABLES]\n%q\n\n", scripts)

	// TODO CLI - would like to cat files?
	for _, script := range scripts {
		//util.RunCommand("cat", script)
		util.RunCommand("chmod", "+x", script)
	}

	// 3. clone geth & checkout at version
	fmt.Printf("\n[CLONING GO ETHEREUM SOURCES]\nos-arch		%s\ngeth version	%s\n\n", osArch, gethVersion)
	//util.RunCommand(cloneGeth, tmpDir, gethVersion)

	// 4. retrieve all necessary parameters for rebuilding in docker.
	fmt.Printf("\n[RETRIEVING DOCKER BUILD PARAMETERS]\n")
	gethCommit := getCommit(gethDir) // TODO need integrity of commit retrieval?..
	referenceURL := getDownloadURL(osArch, gethVersion, gethCommit)
	cc, buildCmd, packages, err := getBuildConfigs(osArch, travisPath)

	// TODO retrieve go version (major vs minor?)
	// TODO ubuntu distribution

	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	timestamp := now.Format("2006-01-02-15:04")

	var dockerArgs = map[string]string{
		"UBUNTU_VERSION": "focal",  // default
		"GO_VERSION":     "1.22.0", // default
		"C_COMPILER":     cc,
		"GETH_VERSION":   gethVersion,
		"OS_ARCH":        osArch,
		"TAG":            fmt.Sprintf("rebuild-geth-v%s-%s-%s", gethVersion, osArch, timestamp),
		"GETH_COMMIT":    gethCommit,
		"SHORT_COMMIT":   gethCommit[0:8],
		"PACKAGES":       strings.Join(packages, " "),
		"REFERENCE_URL":  referenceURL,
		"BUILD_CMD":      buildCmd,
		"ELF_TARGET":     "elf64-x86-64", // TODO fix arch specific elfer.
		"DOCKER_FILE":    rebuildDir + "/Dockerfile",
	}

	fmt.Print("\n")
	for k, v := range dockerArgs {
		fmt.Println(k + ":	" + v)
	}

	// 5. start verification in docker container
	fmt.Printf("\n[STARTING DOCKER REBUILD]\n")
	util.RunCommand(startDocker)
	runDockerBuild(dockerArgs, rebuildDir)

	// for k, v := range dockerArgs {
	// 	if err := os.Setenv(k, v); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	// util.RunCommand(rebuild)
}
