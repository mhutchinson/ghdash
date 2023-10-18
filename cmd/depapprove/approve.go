package main

import (
	"context"
	"flag"
	"os"

	"github.com/golang/glog"
	"github.com/google/go-github/v55/github"
)

var (
	dry_run = flag.Bool("dry_run", true, "Set to false to apply the approvals")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	token := os.Getenv("GH_TOKEN")
	if token == "" {
		glog.Exitf("Create a PAT (https://github.com/settings/tokens) and set as env GH_TOKEN")
	}
	client := github.NewClient(nil).WithAuthToken(token)

	repos := []repo{
		{
			"transparency-dev",
			"armored-witness",
		}, {
			"transparency-dev",
			"armored-witness-applet",
		}, {
			"transparency-dev",
			"armored-witness-boot",
		}, {
			"transparency-dev",
			"armored-witness-common",
		}, {
			"transparency-dev",
			"armored-witness-os",
		},
		{
			"transparency-dev",
			"distributor",
		},
		{
			"transparency-dev",
			"formats",
		},
		{
			"transparency-dev",
			"merkle",
		},
		{
			"transparency-dev",
			"serverless-log",
		},
		{
			"transparency-dev",
			"witness",
		},
	}

	approve := "APPROVE"

	for _, r := range repos {
		if pss, _, err := client.PullRequests.List(ctx, r.owner, r.project, nil); err != nil {
			glog.Exitf("Boom: %v", err)
		} else {
			for _, p := range pss {
				isDependabot := *p.User.ID == 49699333
				fs, _, err := client.PullRequests.ListFiles(ctx, r.owner, r.project, p.GetNumber(), nil)
				if err != nil {
					glog.Exitf("Failed to list files: %v", err)
				}
				approval := isDependabot
				for _, f := range fs {
					if !approval {
						break
					}
					fn := f.GetFilename()
					glog.V(2).Infof("%s: %s", p.GetHTMLURL(), fn)
					approval = approval && (fn == "go.mod" || fn == "go.sum")
				}
				glog.V(1).Infof("dependabot=%t: approve: %t, %d %s %s", isDependabot, approval, *p.ID, *p.Title, p.GetHTMLURL())
				if approval {
					glog.Infof("Approval can be granted for: %s", p.GetHTMLURL())
					if !*dry_run {
						req := &github.PullRequestReviewRequest{
							Event: &approve,
						}
						_, _, err := client.PullRequests.CreateReview(ctx, r.owner, r.project, p.GetNumber(), req)
						if err != nil {
							glog.Errorf("Error approving %s: %s", p.GetHTMLURL(), err)
						}
					}
				} else if isDependabot {
					glog.Infof("Dependabot PR requiring human review: %s", p.GetHTMLURL())
				}
			}
		}
	}
}

type repo struct {
	owner, project string
}
