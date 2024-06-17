package main

import (
	"fmt"
	"log"
	"os"

	util "github.com/chains-project/geth-rebuild/lib"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: <os-arch> <geth version>")
		fmt.Println("Example: linux-amd64 1.14.3")
		os.Exit(1)
	}

	// TODO runtime.GOOS affects commands?

	osArch := os.Args[1]
	gethVersion := os.Args[2]
	// repoURL := "https://github.com/ethereum/go-ethereum.git"
	// branch := "master"
	gethDir := "./tmp/go-ethereum"
	//rootDir := "~/geth-rebuild" // TODO
	checkoutVersion := "v" + gethVersion

	if err := util.ValidParams(osArch, gethVersion); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n[CLONING GO ETHEREUM SOURCES]\nos-arch		%s\ngeth version	%s\n\n", osArch, gethVersion)
	//util.RunCommand("rm", "-r", "-f", gethDir)
	//util.RunCommand("git", "clone", "--branch", branch, repoURL, gethDir) // TODO: shallow copy. Decide proper --depth OR use --single-branch
	util.RunCommand(gethDir, "git", "fetch")
	util.RunCommand(gethDir, "git", "checkout", checkoutVersion)

}
