package utils

/* This file holds supported target platforms for rebuilding
Extend by adding valid architectures for an OS in validArchitectures map
*/

type OS string

type Arch string

const (
	Linux   OS = "linux"
	Darwin  OS = "darwin"
	Windows OS = "windows"
)

const (
	AMD64 Arch = "amd64"
	ARM64 Arch = "arm64"
	ARM5  Arch = "arm5"
	ARM6  Arch = "arm6"
	ARM7  Arch = "arm7"
	A386  Arch = "386"
)

var validArchitectures = map[OS][]Arch{
	Linux: {AMD64, ARM64, ARM5, ARM6, ARM7, A386},
}
