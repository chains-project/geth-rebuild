package util

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var archIds = map[string]string{
	"linux-amd64": "amd64",
	"linux-386":   "386",
	"linux-arm64": "arm64",
	"linux-arm5":  "GOARM=5",
	"linux-arm6":  "GOARM=6",
	"linux-arm7":  "GOARM=7",
}

// todo? retrieve from travis yml?
var compilers = map[string]string{
	"linux-amd64": "gcc-multilib",
	"linux-386":   "gcc-multilib",
	"linux-arm64": "gcc-aarch64-linux-gnu",
	"linux-arm5":  "gcc-arm-linux-gnueabi",
	"linux-arm6":  "gcc-arm-linux-gnueabi",
	"linux-arm7":  "gcc-arm-linux-gnueabihf",
}

var archSpecificPackages = map[string]string{
	"linux-amd64": "",
	"linux-386":   "",
	"linux-arm64": "libc6-dev-arm64-cross",
	"linux-arm5":  "libc6-dev-armel-cross",
	"linux-arm6":  "libc6-dev-armel-cross",
	"linux-arm7":  "libc6-dev-armhf-cross",
}

var commonPackages = []string{"git", "ca-certificates", "wget"}

// Runs command and exits if encountering error.
func RunCommand(dir string, cmd string, args ...string) (out string) {
	exeCmd := exec.Command(cmd, args...)
	exeCmd.Dir = dir // run command in dir

	// catch out in buffer
	var outBuffer bytes.Buffer
	exeCmd.Stdout = &outBuffer

	fmt.Println("[CMD]	", printArgs(exeCmd.Args))
	if err := exeCmd.Run(); err != nil {
		log.Fatal(err)
	}
	return outBuffer.String()
}

// Geth function copying.
func printArgs(args []string) string {
	var s strings.Builder
	for i, arg := range args {
		if i > 0 {
			s.WriteByte(' ')
		}
		if strings.IndexByte(arg, ' ') >= 0 {
			arg = strconv.QuoteToASCII(arg)
		}
		s.WriteString(arg)
	}
	return s.String()
}

// Validates input parameters to main program.
func ValidParams(osArch string, gethVersion string) error {
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

func GetRootDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := wd
	for !strings.HasSuffix(dir, "geth-rebuild") {
		dir = filepath.Dir(dir)
		if dir == "/" {
			return "", fmt.Errorf("error. cannot find root geth-rebuild in '%s'", wd)
		}
	}

	return dir, nil
}

// Returns commit hash at latest commit in dir.
func GetCommit(dir string) string {
	var commit string = RunCommand(dir, "git", "log", "-1", "--format=%H")
	commit = strings.ReplaceAll(commit, "\n", "")
	return commit
}

// Returns link to download geth binary reference build.
func GetDownloadURL(osArch string, gethVersion string, commit string) string {
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
	archID := archIds[osArch]

	if archID == "" {
		return "", fmt.Errorf("architecture id not found for %s", osArch)
	}

	if archID == "amd64" {
		buildCmd := "go run build/ci.go install -dlgo"
		return buildCmd, nil
	}

	f, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %v", file, err)
	}

	reGoARM := fmt.Sprintf(`%s.go\s*run\s*build/ci\.go\s*install.*`, regexp.QuoteMeta(archID))        // GOARM=[5-7] ... go run ...
	reArch := fmt.Sprintf(`go\s*run\s*build/ci\.go\s*install.*-arch\s%s.*`, regexp.QuoteMeta(archID)) // go run ... -arch (386|arm64) ...
	pat := reGoARM + `|` + reArch

	re := regexp.MustCompile(pat)
	match := re.Find(f)
	buildCmd := string(match)

	if buildCmd == "" {
		return "", fmt.Errorf("no build command found for archID %s in file `%s`", archID, file)
	}

	return buildCmd, nil
}

// Returns build configurations for osArch retrieved from build config file (travis.yml).
func GetBuildConfigs(osArch string, file string) (compiler string, cmd string, packages []string, err error) {

	compiler, err = getCompiler(osArch)
	if err != nil {
		return "", "", nil, err
	}

	cmd, err = getBuildCommand(osArch, file)

	if err != nil {
		return "", "", nil, err
	}

	packages = getPackages(osArch)

	return compiler, cmd, packages, nil
}

// Checks if Docker is currently running.
func isDockerRunning() bool {
	cmd := exec.Command("docker", "info")
	err := cmd.Run()
	return err == nil
}

// Opens Docker application on system.
func openDocker() error {
	cmd := exec.Command("open", "-a", "Docker") // todo... arch specific command.
	fmt.Println("[CMD]	", printArgs(cmd.Args))
	return cmd.Run()
}

// Ensures Docker is running on system, or returns error.
func EnsureDocker() error {
	if !isDockerRunning() {
		fmt.Println("Docker is not running. Opening Docker...")
		err := openDocker()
		if err != nil {
			return fmt.Errorf("failed to start docker")
		}

		start := time.Now()
		timeout := start.Add(75 * time.Second)

		for !isDockerRunning() {
			fmt.Println("Waiting for Docker to start...")
			if time.Now().After(timeout) {
				return fmt.Errorf("failed to start docker")
			}
			time.Sleep(5 * time.Second)
		}

	}
	fmt.Println("Docker is running.")
	return nil
}

func RunDockerBuild(buildArgs map[string]string, dockerPath string) {
	dockerFile := dockerPath + "/Dockerfile"

	// set docker build args
	cmdArgs := []string{"build", "-t", "test-tag"}
	for key, value := range buildArgs {
		cmdArgs = append(cmdArgs, fmt.Sprintf(`--build-arg="%s=%s"`, key, value))
	}
	cmdArgs = append(cmdArgs, "-", "<", dockerFile)
	// for _, k := range cmdArgs {
	// 	fmt.Printf("%q\n", k)
	// }

	o := RunCommand(dockerPath, "docker", cmdArgs...)
	fmt.Println(o)
}
