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

	var issues []issueWorkItem
	if iss, _, err := client.Issues.ListByRepo(ctx, owner, project, nil); err != nil {
		glog.Exitf("Boom: %v", err)
	} else {
		for _, i := range iss {
			issues = append(issues, issueWorkItem{i})
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d issues:\n", len(issues)))
	for _, issue := range issues {
		sb.WriteString(fmt.Sprintf(" %s\n", formatItem(workItem(issue))))
	}

	var prs []prWorkItem
	if pss, _, err := client.PullRequests.List(ctx, owner, project, nil); err != nil {
		glog.Exitf("Boom: %v", err)
	} else {
		for _, p := range pss {
			prs = append(prs, prWorkItem{p})
		}
	}
	sb.WriteString(fmt.Sprintf("%d PRs:\n", len(prs)))
	for _, pr := range prs {
		sb.WriteString(fmt.Sprintf(" %s\n", formatItem(workItem(pr))))
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
	GetAttentionSet() []string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}

func formatItem(i workItem) string {
	return fmt.Sprintf("#%d %q: %s -> %s (%s, %s)", i.GetNumber(), i.GetTitle(), *i.GetUser().Login, i.GetAttentionSet(), i.GetCreatedAt().Format("2006-01-02"), i.GetUpdatedAt().Format("2006-01-02"))
}

type issueWorkItem struct {
	*github.Issue
}

func (i issueWorkItem) GetAttentionSet() []string {
	as := make([]string, 0)
	if i.Assignee != nil {
		as = append(as, *i.Assignee.Login)
	}
	return as
}

type prWorkItem struct {
	*github.PullRequest
}

func (p prWorkItem) GetAttentionSet() []string {
	as := make([]string, 0)
	if p.Assignee != nil {
		as = append(as, *p.Assignee.Login)
	}
	for _, r := range p.RequestedReviewers {
		as = append(as, *r.Login)
	}
	return as
}
