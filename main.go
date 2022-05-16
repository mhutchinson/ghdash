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

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d issues:\n", len(issues)))
	for _, issue := range issues {
		sb.WriteString(fmt.Sprintf(" %s\n", formatItem(workItem(issue))))
	}

	prs, _, err := client.PullRequests.List(ctx, owner, project, nil)
	if err != nil {
		glog.Exitf("Boom: %v", err)
	}
	sb.WriteString(fmt.Sprintf("%d PRs:\n", len(prs)))
	for _, pr := range prs {
		reviewers := make([]string, len(pr.RequestedReviewers))
		for i, v := range pr.RequestedReviewers {
			reviewers[i] = *v.Login
		}
		sb.WriteString(fmt.Sprintf(" %s %s\n", formatItem(workItem(pr)), reviewers))
	}

	fmt.Println(sb.String())
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
	return fmt.Sprintf("#%d %q: %s -> %s (%s, %s)", i.GetNumber(), i.GetTitle(), *i.GetUser().Login, assignee, i.GetCreatedAt().Format("2006-01-02"), i.GetUpdatedAt().Format("2006-01-02"))
}
