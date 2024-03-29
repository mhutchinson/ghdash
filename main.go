package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/golang/glog"
	"github.com/google/go-github/v55/github"
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
			if !i.IsPullRequest() {
				issues = append(issues, issueWorkItem{i})
			}
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
	GetLabels() []string
	GetCreatedAt() *github.Timestamp
	GetUpdatedAt() *github.Timestamp
}

func formatItem(i workItem) string {
	attentionSet, _ := json.Marshal(i.GetAttentionSet())
	labels, _ := json.Marshal(i.GetLabels())
	created, updated := i.GetCreatedAt().Format("2006-01-02"), i.GetUpdatedAt().Format("2006-01-02")
	return fmt.Sprintf("#%d %q: %s -> %s (%s, %s) %s", i.GetNumber(), i.GetTitle(), *i.GetUser().Login, string(attentionSet), created, updated, string(labels))
}

type issueWorkItem struct {
	*github.Issue
}

func (i issueWorkItem) GetCreatedAt() *github.Timestamp {
	return i.CreatedAt
}

func (i issueWorkItem) GetUpdatedAt() *github.Timestamp {
	return i.UpdatedAt
}

func (i issueWorkItem) GetAttentionSet() []string {
	as := make([]string, 0)
	if i.Assignee != nil {
		as = append(as, *i.Assignee.Login)
	}
	return as
}

func (i issueWorkItem) GetLabels() []string {
	labels := make([]string, len(i.Labels))
	for i, l := range i.Labels {
		labels[i] = *l.Name
	}
	return labels
}

type prWorkItem struct {
	*github.PullRequest
}

func (p prWorkItem) GetCreatedAt() *github.Timestamp {
	return p.CreatedAt
}

func (p prWorkItem) GetUpdatedAt() *github.Timestamp {
	return p.UpdatedAt
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

func (p prWorkItem) GetLabels() []string {
	labels := make([]string, len(p.Labels))
	for i, l := range p.Labels {
		labels[i] = *l.Name
	}
	if *p.Draft {
		labels = append(labels, "draft")
	}
	return labels
}
