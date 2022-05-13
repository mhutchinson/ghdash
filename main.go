package main

import (
	"context"
	"flag"
	"github.com/google/go-github/v44/github"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	glog.Info("hello")

	client := github.NewClient(nil)

	// list all organizations for user "willnorris"
	orgs, _, err := client.Organizations.List(context.Background(), "willnorris", nil)
	if err != nil {
		glog.Exitf("Boom: %v", err)
	}
	glog.Infof("Found %d orgs", len(orgs))
}

