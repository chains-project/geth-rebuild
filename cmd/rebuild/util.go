package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/chains-project/geth-rebuild/common"
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

// build inputs structs
type BuildInput struct {
	Toolchain ToolchainSpec
	Artifact  ArtifactSpec
	Ubuntu    UbuntuSpec
}

// TODO: interface Spec? methods e.g. getCongfid/SetSpec and showSpec.

type ArtifactSpec struct {
	Version     string
	Os          string
	Arch        string
	Commit      string
	ShortCommit string
}

type ToolchainSpec struct {
	GoVersion string
	CC        string
	BuildCmd  string
	CVersion  string // TODO retrieve (from binary)
}

type UbuntuSpec struct {
	Dist      string
	ElfTarget string
	Packages  []string
}

type DockerSpec struct {
	Dir      string
	BuildTag string
}

// Sets project paths.
func setUpPaths() Paths {
	baseDir, err := common.GetBaseDir("geth-rebuild")
	if err != nil {
		log.Fatal(err)
	}
	paths := Paths{
		Directories: Directories{
			Root:    baseDir,
			Rebuild: filepath.Join(baseDir, "cmd", "rebuild"),
			Temp:    filepath.Join(baseDir, "cmd", "rebuild", "tmp"),
			Geth:    filepath.Join(baseDir, "cmd", "rebuild", "tmp", "go-ethereum"),
			Scripts: filepath.Join(baseDir, "cmd", "rebuild", "scripts"),
			Bin:     filepath.Join(baseDir, "cmd", "rebuild", "bin"),
		},
		Files: Files{
			Travis:    filepath.Join(baseDir, "cmd", "rebuild", "tmp", "go-ethereum", ".travis.yml"),
			Docker:    filepath.Join(baseDir, "cmd", "rebuild", "Dockerfile"),
			Checksums: filepath.Join(baseDir, "cmd", "rebuild", "tmp", "go-ethereum", "build", "checksums.txt"),
		},
		Scripts: Scripts{
			Clone:           filepath.Join(baseDir, "cmd", "rebuild", "scripts", "clone.sh"),
			Checkout:        filepath.Join(baseDir, "cmd", "rebuild", "scripts", "checkout.sh"),
			StartDocker:     filepath.Join(baseDir, "cmd", "rebuild", "scripts", "start_docker.sh"),
			CopyBinaries:    filepath.Join(baseDir, "cmd", "rebuild", "scripts", "copy_bin.sh"),
			CompareBinaries: filepath.Join(baseDir, "cmd", "rebuild", "scripts", "compare_bin.sh"),
		},
	}
	return paths
}

// Validates input arguments to rebuild main program.
func validateArgs(ops string, arch string, version string) error {
	var validArchs = []string{"amd64", "386", "arm5", "arm6", "arm64", "arm7"}
	versionRegex := regexp.MustCompile(`^(\d+\.\d+\.\d+)$`)

	if ops != "linux" {
		return fmt.Errorf("<os> limited to `linux` at the moment")
	}
	if !common.Contains(validArchs, arch) {
		return fmt.Errorf("<arch> must be a valid linux target architecture (amd64|386|arm5|arm6|arm64|arm7)")
	}

	if !versionRegex.MatchString(version) {
		return fmt.Errorf("<geth version> must be in format 'major.minor.patch'\nExample: 1.14.4")
	}

	return nil
}

// Returns ArifactSpec.
func getArtifactSpec(ops string, arch string, version string, paths Paths) (afs ArtifactSpec, err error) {
	_, err = common.RunCommand(paths.Scripts.Clone, paths.Directories.Temp)
	if err != nil {
		return afs, err
	}
	_, err = common.RunCommand(paths.Scripts.Checkout, paths.Directories.Geth, version) // TODO what if unstable,
	if err != nil {
		return afs, err
	}
	commit, err := common.GetCommit(paths.Directories.Geth)
	if err != nil {
		return afs, err
	}

	afs = ArtifactSpec{
		Version:     version,
		Os:          ops,
		Arch:        arch,
		Commit:      commit,
		ShortCommit: commit[0:8],
	}

	return afs, nil
}

// Retrieves build commands for os arch in given travis build file (travis.yml). Returns error if not found.
func getBuildCommand(ops string, arch string, travisFile string) (string, error) {
	var pattern string

	switch ops {
	case "linux":
		switch arch {
		case "amd64":
			return "go run build/ci.go install -dlgo", nil
		case "386", "arm64":
			pattern = fmt.Sprintf(`go\s*run\s*build/ci\.go\s*install.*-arch\s%s.*`, regexp.QuoteMeta(arch))
		case "arm5", "arm6", "arm7":
			v := strings.Split(arch, "arm")[1]
			pattern = fmt.Sprintf(`%s.go\s*run\s*build/ci\.go\s*install.*`, regexp.QuoteMeta(fmt.Sprintf("GOARM=%s", v)))
		default:
			return "", fmt.Errorf("no build command found for linux arch `%s`", arch)
		}
	default:
		return "", fmt.Errorf("no build command found for os `%s`", ops)
	}

	fileContent, err := os.ReadFile(travisFile)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %v", travisFile, err)
	}

	re := regexp.MustCompile(pattern)
	match := re.Find(fileContent)
	if match == nil {
		return "", fmt.Errorf("no build command found for architecture `%s` in file `%s`", arch, travisFile)
	}

	return string(match), nil
}

// Returns the Go compiler version form `major.minor.patch` as specified by geth checksumFile.
func getGoVersion(checksumFile string) (string, error) {
	fileContent, err := os.ReadFile(checksumFile)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %v", checksumFile, err)
	}

	checksumVersionRegex := regexp.MustCompile(`#\s+version:golang\s+(\d+\.\d+\.\d+)`)
	versionRegex := regexp.MustCompile(`(\d+\.\d+\.\d+)`)
	versionLine := checksumVersionRegex.Find(fileContent)

	if versionLine == nil {
		return "", fmt.Errorf("no go version found in file `%s`", checksumFile)
	}
	version := versionRegex.Find(versionLine)

	if version == nil {
		return "", fmt.Errorf("no go version found in file `%s` for line %s", checksumFile, versionLine)
	}

	return string(version), nil
}

// Returns compiler for osArch as described by compilers map. Returns error if not found.
// TODO this is hard coded
func getCC(ops string, arch string) (cc string, err error) {
	switch ops {
	case "linux":
		switch arch {
		case "amd64", "386":
			cc = "gcc-multilib"
		case "arm64":
			cc = "gcc-aarch64-linux-gnu"
		case "arm5", "arm6":
			cc = "gcc-arm-linux-gnueabi"
		case "arm7":
			cc = "gcc-arm-linux-gnueabihf"
		default:
			err = fmt.Errorf("no C compiler found for linux arch `%s`", arch)
		}
	default:
		err = fmt.Errorf("no C compiler found for os `%s`", ops)
	}
	return
}

// Returns build configurations for osArch retrieved from build config file (travis.yml).
func getToolchainSpec(afs ArtifactSpec, paths Paths) (tc ToolchainSpec, err error) {
	goVersion, err := getGoVersion(paths.Files.Checksums)
	if err != nil {
		return tc, fmt.Errorf("failed to get Go version: %w", err)
	}

	cmd, err := getBuildCommand(afs.Os, afs.Arch, paths.Files.Travis)
	if err != nil {
		return tc, fmt.Errorf("failed to get build command: %w", err)
	}

	cc, err := getCC(afs.Os, afs.Arch)
	if err != nil {
		return tc, fmt.Errorf("failed to get C compiler: %w", err)
	}

	tc = ToolchainSpec{
		GoVersion: goVersion,
		CC:        cc,
		BuildCmd:  cmd,
	}
	return tc, nil
}

// Returns ELF input target for os and arch. Used e.g. for binutils `strip` command.
func getElfTarget(ops string, arch string) (elfTarget string, err error) {
	switch ops {
	case "linux":
		switch arch {
		case "amd64":
			elfTarget = "elf64-x86-64"
		case "386":
			elfTarget = "elf32-i386"
		case "arm64":
			elfTarget = "elf64-littleaarch64"
		case "arm5", "arm6", "arm7":
			elfTarget = "elf32-littlearm" // TODO wrong target.
		default:
			err = fmt.Errorf("no elf version found for linux arch `%s`", arch)
		}
	default:
		err = fmt.Errorf("no elf version found for os `%s`", ops)
	}
	return
}

// Returns common and architecture specific package for osArc.
func getUbuntuPackages(ops string, arch string) (packages []string, err error) {
	packages = append(packages, "git", "ca-certificates", "wget", "binutils") // common
	switch ops {
	case "linux":
		switch arch {
		case "amd64", "386":
			return // no arch specific packages
		case "arm64":
			packages = append(packages, "libc6-dev-arm64-cross") // TODO hard coded - can use regex?
		case "arm5", "arm6":
			packages = append(packages, "libc6-dev-armel-cross")
		case "arm7":
			packages = append(packages, "libc6-dev-armhf-cross")
		default:
			return nil, fmt.Errorf("no packages found for linux arch `%s`", arch)
		}
	default:
		return nil, fmt.Errorf("no packages found for os `%s`", ops)
	}
	return
}

func getUbuntuSpec(afs ArtifactSpec) (ub UbuntuSpec, err error) {
	elfTarget, err := getElfTarget(afs.Os, afs.Arch)
	if err != nil {
		return ub, fmt.Errorf("failed to get ELF target: %w", err)
	}

	packages, err := getUbuntuPackages(afs.Os, afs.Arch)
	if err != nil {
		return ub, fmt.Errorf("failed to get Ubuntu packages: %w", err)
	}

	ub = UbuntuSpec{
		Dist:      "focal", //TODO: get dist from Travis file
		ElfTarget: elfTarget,
		Packages:  packages,
	}
	return ub, nil
}

func createDockerTag(gethVersion string, ops string, arch string) string {
	now := time.Now()
	timestamp := now.Format("2006-01-02-15:04")
	tag := fmt.Sprintf("rebuild-geth-v%s-%s-%s-%s", gethVersion, ops, arch, timestamp)
	return tag
}

func (bi *BuildInput) ToMap() map[string]string {
	buildArgs := make(map[string]string)
	buildArgs["GO_VERSION"] = bi.Toolchain.GoVersion
	buildArgs["CC"] = bi.Toolchain.CC
	buildArgs["C_VERSION"] = bi.Toolchain.CVersion
	buildArgs["GETH_VERSION"] = bi.Artifact.Version
	buildArgs["OS"] = bi.Artifact.Os
	buildArgs["ARCH"] = bi.Artifact.Arch
	buildArgs["COMMIT"] = bi.Artifact.Commit
	buildArgs["SHORT_COMMIT"] = bi.Artifact.ShortCommit
	buildArgs["UBUNTU_VERSION"] = bi.Ubuntu.Dist
	buildArgs["ELF_TARGET"] = bi.Ubuntu.ElfTarget
	buildArgs["PACKAGES"] = strings.Join(bi.Ubuntu.Packages, " ")
	buildArgs["BUILD_CMD"] = bi.Toolchain.BuildCmd
	return buildArgs
}

// Starts a docker build for dockerfile at `dockerPath` with given `buildArgs`.
func runDockerBuild(buildArgs map[string]string, dockerTag string, dockerDir string) error {
	// set docker build args
	cmdArgs := []string{"build", "-t", dockerTag, "--progress=plain"} // TODO test tty
	for key, value := range buildArgs {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--build-arg=%s=%s", key, value))
	}
	cmdArgs = append(cmdArgs, dockerDir)
	// run docker build
	_, err := common.RunCommand("docker", cmdArgs...)
	if err != nil {
		return err
	}
	return nil
}
