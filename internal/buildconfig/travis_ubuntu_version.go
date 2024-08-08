package buildconfig

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chains-project/geth-rebuild/internal/utils"
)

/* This file is for calling the Travis CI API to fetch the actual ubuntu distribution used in a build.
Accessing the API is necessary to correctly retrieve the used distribution as a bug in Travis lets
one distribution be defined in .travis.yml whereas another (random) distribution will be used in practice 
See https://travis-ci.community/t/travis-uses-unexpected-ubuntu-distribution/14286/2

This is particularly bad Travis CI sometimes chooses bionic, which is known to have path issues (https://github.com/golang/go/issues/67011)

In the beginning of a Travis log, the actual distribution used is found in the line `Codename: xx`
*/

const (
	travisAPIVersion = "3"
	baseURL          = "https://api.travis-ci.com"
	repoSlug         = "ethereum%2Fgo-ethereum"
)

// JSON response
type Response struct {
	Builds []Build `json:"builds"`
}

// A Travis CI build
type Build struct {
	ID     int    `json:"id"`
	Commit Commit `json:"commit"`
	Jobs   []Job  `json:"jobs"`
}

func (b Build) String() string {
	return fmt.Sprintf("[BUILD]\nID:		%d\nCommit:		%s\nJobs:		%v", b.ID, b.Commit.SHA, b.Jobs)
}

// A job within a build
type Job struct {
	ID int `json:"id"`
}

// Git commit for a build
type Commit struct {
	SHA string `json:"sha"`
}

// A Travis CI job log
type JobLog struct {
	Content string `json:"content"`
}

// TODO: if stable build, use previous approach (faster)
// buildsURL := fmt.Sprintf("%s/repo/%s/builds?branch.name=%s", baseURL, slug, branchName)
// TODO use the latest tagged version for unstable builds to adjust pagination (forward search)
// Retrieves the **actual** Ubuntu distribution used in a Travis CI build with given git commit
func GetUbuntuDist(gitCommit string, searchPages int) (string, error) {
	var build Build
	var found bool

	// Search 50 builds at a time, searchPages # of times
	for page := range searchPages {
		resp, err := getTravisBuilds(50, page)
		if err != nil {
			return "", fmt.Errorf("Could not determine ubuntu distribution, error: %w", err)
		}
		found, build = findBuildByCommit(gitCommit, resp.Builds)
		if found {
			break
		}
	}
	if !found {
		return "", fmt.Errorf("Failed to find a Travis build with given commit %s in the first %d pages", gitCommit, searchPages)
	}

	firstJob := build.Jobs[0]
	log, err := getJobLog(firstJob.ID)
	if err != nil {
		return "", err
	}

	dist := extractCodename(log.Content)

	if dist == "" {
		return "", fmt.Errorf("no ubuntu distribution could be determined from build log")
	}

	return dist, nil
}

// Retrieves a slice with limit # of travis builds, starting from build # limit*page
func getTravisBuilds(limit int, page int) (Response, error) {
	offset := page * limit
	url := fmt.Sprintf("%s/repo/%s/builds?limit=%s&offset=%s", baseURL, repoSlug, fmt.Sprint(limit), fmt.Sprint(offset))

	travisBuilds, err := utils.HttpGetRequest(url, map[string]string{"Travis-API-Version": travisAPIVersion})
	if err != nil {
		return Response{}, fmt.Errorf("failed to fetch Travis builds: %w", err)
	}

	var resp Response
	if err := json.Unmarshal([]byte(travisBuilds), &resp); err != nil {
		return Response{}, fmt.Errorf("failed to unmarshal builds data: %w", err)
	}
	return resp, nil
}

// Searches for a build by its git commit
// This approach is necessary as the commit.sha query endpoint is not safelisted by geth
func findBuildByCommit(gitCommit string, builds []Build) (found bool, build Build) {
	for _, b := range builds {
		if b.Commit.SHA == gitCommit {
			return true, b
		}
	}
	return false, Build{}
}

// Retrieves the Travis CI job log defined by given job ID
func getJobLog(jobID int) (JobLog, error) {
	url := fmt.Sprintf("%s/job/%d/log", baseURL, jobID)

	resp, err := utils.HttpGetRequest(url, map[string]string{"Travis-API-Version": travisAPIVersion})
	if err != nil {
		return JobLog{}, fmt.Errorf("failed to fetch Travis job log for ID %d: %w", jobID, err)
	}

	var log JobLog
	if err := json.Unmarshal([]byte(resp), &log); err != nil {
		return JobLog{}, fmt.Errorf("failed to unmarshal log data: %w", err)
	}
	return log, nil
}

// Gets ubuntu dist as defined by line `Codename:` in the travis job log
func extractCodename(content string) string {
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
