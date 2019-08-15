package main

//

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	githubToken    = flag.String("token", "", "Github access token")
	githubProj     = flag.String("p", "", "Github project/repository owner name")
	githubRepo     = flag.String("r", "", "Github repository name")
	pullRequestNum = flag.Int("pr", 0, "Github pull request number")
	attemptTimeout = flag.Int("t", 60, "attempts timeout in seconds")
)

func main() {
	flag.Parse()
	token := mustStringFromEnvIfNotSet("token", "GITHUB_AUTH_TOKEN", *githubToken)
	proj := mustStringFromEnvIfNotSet("project", "GITHUB_PROJECT", *githubProj)
	repo := mustStringFromEnvIfNotSet("r", "GITHUB_REPO", *githubRepo)
	timeout := time.Duration(int64(*attemptTimeout) * int64(time.Second.Nanoseconds()))

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	for i := 1; ; i++ {
		pr, _, err := client.PullRequests.Get(ctx, proj, repo, *pullRequestNum)
		failedOnErr(err)
		ms := pr.GetMergeableState()

		log.Printf("PR#%v: merged=%v, state=%s, mergable=%v (attemt %v)\n",
			pr.GetNumber(),
			pr.GetMerged(),
			ms,
			pr.GetMergeable(),
			i,
		)
		if pr.GetMerged() {
			finishWithMessage("\nAlready merged, nothing to do.")
		}
		if !pr.GetMergeable() {
			failedWithMessage("\nNot in mergable state.")
		}
		if ms == "blocked" || ms == "unknown" {
			time.Sleep(timeout)
			continue
		}

		prMerged, _, err := client.PullRequests.Merge(ctx, proj, repo, *pullRequestNum, "automerge", nil)
		failedOnErr(err)
		log.Printf("PR merged: %#v.", prMerged)
		break
	}
}

func failedWithMessage(msg string) {
	fmt.Println("Failed.", msg)
	os.Exit(1)
}

func finishWithMessage(msg string) {
	fmt.Println(msg)
	os.Exit(0)
}

func failedOnErr(e error) {
	if e == nil {
		return
	}
	extra := ""
	if ge, ok := e.(*github.ErrorResponse); ok {
		b, _ := json.MarshalIndent(ge.Block, "", "  ")
		if string(b) != "null" {
			extra = fmt.Sprintf(" \n%s", string(b))
		}
	}
	failedWithMessage(fmt.Sprintf("Error: %s", e.Error()) + extra)
}

func mustStringFromEnvIfNotSet(name, envVarName, value string) string {
	if value != "" {
		return value
	}
	value = os.Getenv(envVarName)
	if value == "" {
		failedWithMessage(fmt.Sprintf("-%v flag or %v env variable should be set.", name, envVarName))
	}
	return value
}
