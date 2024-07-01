package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	util "github.com/chains-project/geth-rebuild/util"
)

var osArchPatterns = map[string]string{
	"linux-amd64": "amd64",
	"linux-386":   "386",
	"linux-arm64": "arm64",
	"linux-arm5":  "GOARM=5",
	"linux-arm6":  "GOARM=6",
	"linux-arm7":  "GOARM=7",
}

// TODO REMOVE HARD CODING by retrieving from travis yml?
var compilers = map[string]string{
	"linux-amd64": "gcc-multilib",
	"linux-386":   "gcc-multilib",
	"linux-arm64": "gcc-aarch64-linux-gnu",
	"linux-arm5":  "gcc-arm-linux-gnueabi",
	"linux-arm6":  "gcc-arm-linux-gnueabi",
	"linux-arm7":  "gcc-arm-linux-gnueabihf",
}

// TODO REMOVE HARD CODING by retrieving from travis yml?
var archSpecificPackages = map[string]string{
	"linux-amd64": "",
	"linux-386":   "",
	"linux-arm64": "libc6-dev-arm64-cross",
	"linux-arm5":  "libc6-dev-armel-cross",
	"linux-arm6":  "libc6-dev-armel-cross",
	"linux-arm7":  "libc6-dev-armhf-cross",
}

var commonPackages = []string{"git", "ca-certificates", "wget"}

// Validates input parameters to main program.
func validParams(osArch string, gethVersion string) error {
	osArchPattern := "^linux-(amd64|386|arm5|arm6|arm64|arm7)$"
	versionPattern := "^[0-9]+.[0-9]+.[0-9]+$"

	osArchRegex := regexp.MustCompile(osArchPattern)
	versionRegex := regexp.MustCompile(versionPattern)

	if !osArchRegex.MatchString(osArch) {
		return fmt.Errorf("<os-arch> must be a valid linux target architecture\nExample: linux-amd64")
	}
	if !versionRegex.MatchString(gethVersion) {
		return fmt.Errorf("<geth version> must be in format 'major.minor.patch'\nExample: 1.14.4")
	}
	return nil
}

// Returns commit hash at latest commit in dir.
func getCommit(dir string) string {
	gitDir := fmt.Sprintf("--git-dir=%s/.git", dir)
	workTree := fmt.Sprintf("--work-tree=%s", dir)
	var commit string = util.RunCommand("git", gitDir, workTree, "log", "-1", "--format=%H")
	commit = strings.ReplaceAll(commit, "\n", "")
	return commit
}

// Returns link to download geth binary reference build.
func getDownloadURL(osArch string, gethVersion string, commit string) string {
	shortCommit := commit[0:8]
	targetPackage := "geth-" + osArch + "-" + gethVersion + "-" + shortCommit
	url := "https://gethstore.blob.core.windows.net/builds/" + targetPackage + ".tar.gz"
	return url
}

// Returns common and architecture specific package for osArc.
func getPackages(osArch string) []string {
	archSpecific := archSpecificPackages[osArch]
	if archSpecific == "" {
		return commonPackages
	} else {
		all := append(commonPackages, archSpecific)
		return all
	}
}

// Returns compiler for osArch as described by compilers map. Returns error if not found.
func getCompiler(osArch string) (string, error) {
	c := compilers[osArch]
	if c == "" {
		return "", fmt.Errorf("no compiler found for arch id %s", osArch)
	}
	return c, nil
}

// Retrieves build commands for osArch in given file. Returns error if not found.
func getBuildCommand(osArch string, file string) (string, error) {
	archID := osArchPatterns[osArch]

	if archID == "" {
		err := fmt.Errorf("architecture id not found for %s", osArch)
		return "", err
	}

	if archID == "amd64" {
		buildCmd := "go run build/ci.go install -dlgo"
		return buildCmd, nil
	}

	f, err := os.ReadFile(file)
	if err != nil {
		err = fmt.Errorf("error reading file %s: %v", file, err)
		return "", err
	}

	goarmPattern := fmt.Sprintf(`%s.go\s*run\s*build/ci\.go\s*install.*`, regexp.QuoteMeta(archID))        // pattern: GOARM=[5-7] ... go run ...
	archPattern := fmt.Sprintf(`go\s*run\s*build/ci\.go\s*install.*-arch\s%s.*`, regexp.QuoteMeta(archID)) // pattern: go run ... -arch (386|arm64) ...
	pat := goarmPattern + `|` + archPattern

	re := regexp.MustCompile(pat)
	match := re.Find(f)
	buildCmd := string(match)

	if buildCmd == "" {
		err = fmt.Errorf("no build command found for archID %s in file `%s`", archID, file)
		return "", err
	}
	return buildCmd, nil
}

// Returns build configurations for osArch retrieved from build config file (travis.yml).
func getBuildConfigs(osArch string, file string) (cc string, cmd string, packages []string, err error) {

	cc, err = getCompiler(osArch)
	if err != nil {
		return
	}

	cmd, err = getBuildCommand(osArch, file)

	if err != nil {
		return
	}

	packages = getPackages(osArch)
	return
}

// // Checks if Docker is currently running.
// func isDockerRunning() bool {
// 	cmd := exec.Command("docker", "info")
// 	err := cmd.Run()
// 	return err == nil
// }

// // Opens Docker application on system.
// func openDocker() error {
// 	cmd := exec.Command("open", "-a", "Docker") // todo... arch specific command.
// 	fmt.Println("[CMD]	", printArgs(cmd.Args))
// 	return cmd.Run()
// }

// // Ensures Docker is running on system, or returns error.
// func EnsureDocker() error {
// 	if !isDockerRunning() {
// 		fmt.Println("Docker is not running. Opening Docker...")
// 		err := openDocker()
// 		if err != nil {
// 			return fmt.Errorf("failed to start docker")
// 		}

// 		start := time.Now()
// 		timeout := start.Add(75 * time.Second)

// 		for !isDockerRunning() {
// 			fmt.Println("Waiting for Docker to start...")
// 			if time.Now().After(timeout) {
// 				return fmt.Errorf("failed to start docker")
// 			}
// 			time.Sleep(5 * time.Second)
// 		}

// 	}
// 	fmt.Println("Docker is running.")
// 	return nil
// }

func runDockerBuild(buildArgs map[string]string, dockerPath string) {
	// set docker build args
	cmdArgs := []string{"build", "-t", buildArgs["TAG"]}
	for key, value := range buildArgs {
		cmdArgs = append(cmdArgs, fmt.Sprint("--build-arg="+key+"="+value)) // TODO wordsplitting....
	}
	cmdArgs = append(cmdArgs, dockerPath)
	o := util.RunCommand("docker", cmdArgs...)
	fmt.Println(o)
}
