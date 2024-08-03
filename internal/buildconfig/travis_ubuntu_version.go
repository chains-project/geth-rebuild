package buildconfig

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chains-project/geth-rebuild/internal/utils"
)

/* This file is for calling the Travis CI API to fetch the actual ubuntu distribution used in a build.
Accessing the API is necessary to correctly retrieve the used distribution as a bug in Travis lets
one distribution be defined in .travis.yml whereas another (random) distribution will be used in practice.

In the beginning of a Travis log, the actual distribution used is found in the line `Codename: xx`
*/

const (
	travisAPIVersion = "3"
	baseURL          = "https://api.travis-ci.com"
	slug             = "ethereum%2Fgo-ethereum"
)

type BuildResponse struct {
	Builds []Build `json:"builds"`
}

type Build struct {
	ID     int    `json:"id"`
	Number string `json:"number"`
	State  string `json:"state"`
	Jobs   []Job  `json:"jobs"`
}

type Job struct {
	ID int `json:"id"`
}

type LogResponse struct {
	Content string `json:"content"`
}

func getCodename(content string) string {
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "Codename:") {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				return strings.TrimSpace(fields[1])
			}
		}
	}
	return ""
}

func getUbuntuDist(gethVersion string) (string, error) {
	// get builds for a push to specified branch defined by version
	branchName := fmt.Sprintf("v%s", gethVersion)
	buildsURL := fmt.Sprintf("%s/repo/%s/builds?branch.name=%s", baseURL, slug, branchName)

	travisBuilds, err := utils.HttpGetRequest(buildsURL, map[string]string{"Travis-API-Version": travisAPIVersion})
	if err != nil {
		return "", fmt.Errorf("failed to fetch Travis builds: %w", err)
	}

	var buildResponse BuildResponse
	if err := json.Unmarshal([]byte(travisBuilds), &buildResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal builds data: %w", err)
	}

	if len(buildResponse.Builds) == 0 {
		return "", fmt.Errorf("no builds found for branch %s", branchName)
	}

	firstBuild := buildResponse.Builds[0]
	if len(firstBuild.Jobs) == 0 {
		return "", fmt.Errorf("no jobs found for build ID %d", firstBuild.ID)
	}

	// get a build log, e.g. first job in first build
	jobID := firstBuild.Jobs[0].ID // assumes all jobs use same ubuntu distribution
	logURL := fmt.Sprintf("%s/job/%d/log", baseURL, jobID)

	logData, err := utils.HttpGetRequest(logURL, map[string]string{"Travis-API-Version": travisAPIVersion})
	if err != nil {
		return "", fmt.Errorf("failed to fetch log: %w", err)
	}

	var logResponse LogResponse
	if err := json.Unmarshal([]byte(logData), &logResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal log data: %w", err)
	}
	dist := getCodename(logResponse.Content)

	if dist == "" {
		return "", fmt.Errorf("no ubuntu distribution could be determined from build log")
	}
	return dist, nil
}
