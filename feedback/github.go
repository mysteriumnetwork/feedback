/*
 * Copyright (C) 2019 The "MysteriumNetwork/node" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package feedback

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"strconv"
	"text/template"
	"time"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

// Reporter reports issues to Github
type Reporter struct {
	client        *github.Client
	owner         string
	repository    string
	issueTemplate *template.Template
}

// NewReporterOpts Reporter initialization options
type NewReporterOpts struct {
	Token      string
	Owner      string
	Repository string
}

// NewReporter creates a new Reporter
func NewReporter(opts *NewReporterOpts) *Reporter {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: opts.Token},
	)
	oauthClient := oauth2.NewClient(context.Background(), ts)
	githubClient := github.NewClient(oauthClient)
	issueTemplate := template.Must(template.New("issueTemplate").Parse(issueTemplate))
	return &Reporter{
		client:        githubClient,
		owner:         opts.Owner,
		repository:    opts.Repository,
		issueTemplate: issueTemplate,
	}
}

// Report user report data
type Report struct {
	UserId      string
	Description string
	Email       string
	LogURL      url.URL
}

const issueTemplate = `
🆔 {{.Identity}}
📅 {{.Timestamp}}
✉️ {{.Email}}

### Description

{{.Description}}

### Logs

{{.LogURL}}

`

// ReportIssue reports issue
func (rep *Reporter) ReportIssue(report *Report) (issueId string, err error) {
	templateOpts := struct {
		Description,
		Email,
		Identity,
		Timestamp,
		LogURL string
	}{
		Description: report.Description,
		Email:       report.Email,
		Identity:    report.UserId,
		Timestamp:   time.Now().String(),
		LogURL:      report.LogURL.String(),
	}
	var body bytes.Buffer
	err = rep.issueTemplate.Execute(&body, templateOpts)
	if err != nil {
		return "", fmt.Errorf("could not generate issue body with report (%+v): %w", templateOpts, err)
	}

	req := github.IssueRequest{
		Title: github.String("User report: " + report.UserId),
		Body:  github.String(body.String()),
	}
	issue, _, err := rep.client.Issues.Create(context.Background(), rep.owner, rep.repository, &req)

	if err != nil {
		return "", fmt.Errorf("could not create github issue: %w", err)
	}

	return strconv.Itoa(issue.GetNumber()), nil
}
