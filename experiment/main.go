package main

import (
	"html/template"
	"log"
	"os"
	"os/exec"
)

// command execution results
type ReproduceDetails struct {
	OsArch    string
	Version   string
	Success   bool
	Binaries  string
	ErrorInfo string
}

func main() {
	// todo input validation

	testCases := []struct {
		OsArch  string
		Version string
	}{
		{os.Args[1], os.Args[2]},
		{os.Args[3], os.Args[4]},
	}

	for _, testCase := range testCases {
		osArch := testCase.OsArch
		version := testCase.Version

		// Run your command here
		cmd := exec.Command("TODO", osArch, version) //TODO binary? or go run.
		output, _ := cmd.CombinedOutput()

		success := cmd.ProcessState.Success()
		var result ReproduceDetails

		if success {
			result = ReproduceDetails{
				OsArch:   osArch,
				Version:  version,
				Success:  true,
				Binaries: string(output),
			}
		} else {
			result = ReproduceDetails{
				OsArch:    osArch,
				Version:   version,
				Success:   false,
				ErrorInfo: string(output),
			}
		}
		generateHTML(result)
	}
}

func generateHTML(result ReproduceDetails) {
	// Define the HTML template
	const htmlTemplate = `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Reproducibility Test</title>
		<style>
			body { font-family: Arial, sans-serif; }
			.success { color: green; }
			.failure { color: red; }
		</style>
	</head>
	<body>
		<h1>Reproducibility Test Result</h1>
		<h2>OS/Arch: {{ .OsArch }}</h2>
		<h2>Version: {{ .Version }}</h2>
		{{ if .Success }}
			<p class="success">Result: Reproducible</p>
			<p>Binaries: {{ .Binaries }}</p>
		{{ else }}
			<p class="failure">Result: Unreproducible</p>
			<p>Error Info: {{ .ErrorInfo }}</p>
		{{ end }}
	</body>
	</html>
	`

	tmpl, err := template.New("reproducibility").Parse(htmlTemplate)
	if err != nil {
		log.Fatalf("Error parsing HTML template: %v", err)
	}

	err = tmpl.Execute(os.Stdout, result)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
}
