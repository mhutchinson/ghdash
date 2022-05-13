package main

import (
	"context"
	"flag"

	"github.com/golang/glog"
	"github.com/google/go-github/v44/github"
)

func main() {
	flag.Parse()
	ctx := context.Background()

	client := github.NewClient(nil)

	issues, _, err := client.Issues.ListByRepo(ctx, "google", "trillian", nil)
	if err != nil {
		glog.Exitf("Boom: %v", err)
	}
	glog.Infof("Found %d issues", len(issues))

	prs, _, err := client.PullRequests.List(ctx, "google", "trillian", nil)
	if err != nil {
		glog.Exitf("Boom: %v", err)
	}
	glog.Infof("Found %d PRs", len(prs))
}
