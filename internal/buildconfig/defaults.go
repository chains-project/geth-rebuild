package buildconfig

import "github.com/chains-project/geth-rebuild/internal/utils"

type DefaultConfigs struct {
	ToolchainDeps map[utils.OS]map[utils.Arch][]string
	ElfTargets    map[utils.OS]map[utils.Arch]string
	ArmVersions   map[utils.OS]map[utils.Arch]string
	CC            map[utils.OS]map[utils.Arch]string
	UtilDeps      []string
}

var DefaultConfig = DefaultConfigs{
	ToolchainDeps: map[utils.OS]map[utils.Arch][]string{
		utils.Linux: {
			utils.AMD64: {"gcc-multilib"},
			utils.A386:  {"gcc-multilib"},
			utils.ARM64: {"gcc-aarch64-linux-gnu", "libc6-dev-arm64-cross"},
			utils.ARM5:  {"gcc-arm-linux-gnueabi", "libc6-dev-armel-cross", "gcc-arm-linux-gnueabihf", "libc6-dev-armhf-cross"},
			utils.ARM6:  {"gcc-arm-linux-gnueabi", "libc6-dev-armel-cross", "gcc-aarch64-linux-gnu", "libc6-dev-arm64-cross"},
			utils.ARM7:  {"gcc-arm-linux-gnueabihf", "libc6-dev-armhf-cross"},
		},
	},
	ElfTargets: map[utils.OS]map[utils.Arch]string{
		utils.Linux: {
			utils.AMD64: "elf64-x86-64",
			utils.A386:  "elf32-i386",
			utils.ARM64: "elf64-little",
			utils.ARM5:  "elf32-little",
			utils.ARM6:  "elf32-little",
			utils.ARM7:  "elf32-little",
		},
	},
	ArmVersions: map[utils.OS]map[utils.Arch]string{
		utils.Linux: {
			utils.AMD64: "", // empty goarm value sets empty flag in docker build `GOARM=`
			utils.A386:  "",
			utils.ARM64: "", // no ARM flag optimization used for aarch64
			utils.ARM5:  "5",
			utils.ARM6:  "6",
			utils.ARM7:  "7",
		},
	},
	CC: map[utils.OS]map[utils.Arch]string{
		utils.Linux: {
			utils.AMD64: "", // no cc flag needed
			utils.A386:  "",
			utils.ARM64: "aarch64-linux-gnu-gcc",
			utils.ARM5:  "arm-linux-gnueabi-gcc",
			utils.ARM6:  "arm-linux-gnueabi-gcc",
			utils.ARM7:  "arm-linux-gnueabihf-gcc",
		},
	},
	UtilDeps: []string{"git", "ca-certificates", "wget", "binutils"},
}
