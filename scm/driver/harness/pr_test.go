// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
	"github.com/nwxleo/go-scm/scm"
	"github.com/nwxleo/go-scm/scm/transport"
)

func TestPRFind(t *testing.T) {
	if harnessPAT == "" {
		defer gock.Off()

		gock.New(gockOrigin).
			Get("/gateway/code/api/v1/repos/px7xd_BFRCi-pfWPYXVjvw/default/codeciintegration/thomas/+/pullreq/1").
			Reply(200).
			Type("plain/text").
			File("testdata/pr.json")
	}
	client, _ := New(gockOrigin, harnessOrg, harnessAccount, harnessProject)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("x-api-key", harnessPAT)
			},
		},
	}
	got, _, err := client.PullRequests.Find(context.Background(), harnessRepo, 1)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.PullRequest)
	raw, err := ioutil.ReadFile("testdata/pr.json.golden")
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal(raw, want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPRCommits(t *testing.T) {
	if harnessPAT == "" {
		defer gock.Off()

		gock.New(gockOrigin).
			Get("/gateway/code/api/v1/repos/px7xd_BFRCi-pfWPYXVjvw/default/codeciintegration/thomas/+/pullreq/1/commits").
			Reply(200).
			Type("plain/text").
			File("testdata/pr_commits.json")
	}
	client, _ := New(gockOrigin, harnessOrg, harnessAccount, harnessProject)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("x-api-key", harnessPAT)
			},
		},
	}
	got, _, err := client.PullRequests.ListCommits(context.Background(), harnessRepo, 1, scm.ListOptions{})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Commit{}
	raw, err := ioutil.ReadFile("testdata/pr_commits.json.golden")
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal(raw, &want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullCreate(t *testing.T) {
	defer gock.Off()
	gock.New(gockOrigin).
		Post("/gateway/code/api/v1/repos/px7xd_BFRCi-pfWPYXVjvw/default/codeciintegration/thomas/+/pullreq").
		Reply(200).
		Type("plain/text").
		File("testdata/pr.json")

	client, _ := New(gockOrigin, harnessOrg, harnessAccount, harnessProject)
	client.Client = &http.Client{
		Transport: &transport.Custom{
			Before: func(r *http.Request) {
				r.Header.Set("x-api-key", harnessPAT)
			},
		},
	}

	input := scm.PullRequestInput{
		Title:  "pull title",
		Body:   "pull description",
		Source: "bla",
		Target: "main",
	}

	got, _, err := client.PullRequests.Create(context.Background(), harnessRepo, &input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.PullRequest)
	raw, _ := ioutil.ReadFile("testdata/pr.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
