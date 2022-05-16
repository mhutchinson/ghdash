package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/google/go-github/v44/github"
)

var (
	repo = flag.String("repo", "google/trillian", "'owner/project' components of the github repo to get stats for")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	owner, project := parseRepo()
	client := github.NewClient(nil)

	issues, _, err := client.Issues.ListByRepo(ctx, owner, project, nil)
	if err != nil {
		glog.Exitf("Boom: %v", err)
	}
	glog.Infof("Found %d issues", len(issues))
	for _, issue := range issues {
		glog.Infof(" %s", formatItem(workItem(issue)))
	}

	prs, _, err := client.PullRequests.List(ctx, owner, project, nil)
	if err != nil {
		glog.Exitf("Boom: %v", err)
	}
	glog.Infof("Found %d PRs", len(prs))
	for _, pr := range prs {
		reviewers := make([]string, len(pr.RequestedReviewers))
		for i, v := range pr.RequestedReviewers {
			reviewers[i] = *v.Login
		}
		glog.Infof(" %s %s", formatItem(workItem(pr)), reviewers)
	}
}

func parseRepo() (string, string) {
	ss := strings.Split(*repo, "/")
	if l := len(ss); l > 2 {
		glog.Exitf("Expected owner/project, but got %d components", l)
	}
	return ss[0], ss[1]
}

type workItem interface {
	GetNumber() int
	GetTitle() string
	GetUser() *github.User
	GetAssignee() *github.User
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}

func formatItem(i workItem) string {
	assignee := "UNASSIGNED"
	if ass := i.GetAssignee(); ass != nil {
		assignee = *ass.Login
	}
	return fmt.Sprintf("#%d %q: %s -> %s (%s, %s)", i.GetNumber(), i.GetTitle(), *i.GetUser().Login, assignee, i.GetCreatedAt().Format("2006-02-01"), i.GetUpdatedAt().Format("2006-02-01"))
}
