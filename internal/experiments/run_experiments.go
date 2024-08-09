package experiments

import (
	"fmt"
	"log"

	"github.com/chains-project/geth-rebuild/internal/utils"
)

func RunExperiments(exps []ExperimentInput, exePath string) {
	fmt.Printf("[EXPERIMENT RUN FOR OS `%s` ARCH `%s`]\n", exps[0].OS, exps[0].Arch)

	// Sequential execution of experiments
	for _, exp := range exps {
		runExperiment(exp, exePath)
	}
}

func runExperiment(exp ExperimentInput, exePath string) {
	cmdArgs := []string{string(exp.OS), string(exp.Arch), exp.Version, "--diff"}
	if exp.Unstable != "" {
		cmdArgs = append(cmdArgs, "--unstable", exp.Unstable)
	}
	out, err := utils.RunCommand(exePath, cmdArgs...)
	if err != nil {
		log.Printf("Error running gethrebuild: %v\nOutput: %s", err, string(out))
	} else {
		fmt.Printf("Experiment successful: %s\n", exp)
	}
}
