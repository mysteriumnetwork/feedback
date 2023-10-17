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
	"path"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

// GithubReporter reports issues to Github
type GithubReporter struct {
	client          *github.Client
	owner           string
	repository      string
	issueTemplate   *template.Template
	logProxyBaseUrl string
}

// NewGithubReporterOpts GithubReporter initialization options
type NewGithubReporterOpts struct {
	Token           string
	Owner           string
	Repository      string
	LogProxyBaseUrl string
}

// NewGithubReporter creates a new GithubReporter
func NewGithubReporter(opts *NewGithubReporterOpts) *GithubReporter {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: opts.Token},
	)
	oauthClient := oauth2.NewClient(context.Background(), ts)
	githubClient := github.NewClient(oauthClient)
	issueTemplate := template.Must(template.New("issueTemplate").Parse(issueTemplate))
	return &GithubReporter{
		client:          githubClient,
		owner:           opts.Owner,
		repository:      opts.Repository,
		issueTemplate:   issueTemplate,
		logProxyBaseUrl: strings.TrimSuffix(opts.LogProxyBaseUrl, "/"),
	}
}

const issueTemplate = `
🆔 {{.Identity}}
📅 {{.Timestamp}}
✉️ {{.Email}}

### Description

{{.Description}}

### Logs

{{.LogProxyBaseUrl}}/{{.LogKey}}

`

// ReportIssue reports issue
func (rep *GithubReporter) ReportIssue(report *Report) (issueId string, err error) {
	key := path.Base(report.LogURL.String())
	templateOpts := struct {
		Description,
		Email,
		Identity,
		Timestamp,
		LogKey,
		LogProxyBaseUrl string
	}{
		Description:     report.Description,
		Email:           report.Email,
		Identity:        report.NodeIdentity,
		Timestamp:       time.Now().Format("2006-01-02 15:04:05"),
		LogKey:          key,
		LogProxyBaseUrl: rep.logProxyBaseUrl,
	}
	var body bytes.Buffer
	err = rep.issueTemplate.Execute(&body, templateOpts)
	if err != nil {
		return "", fmt.Errorf("could not generate issue body with report (%+v): %w", templateOpts, err)
	}

	req := github.IssueRequest{
		Title: github.String("User report: " + report.NodeIdentity),
		Body:  github.String(body.String()),
	}
	issue, _, err := rep.client.Issues.Create(context.Background(), rep.owner, rep.repository, &req)

	if err != nil {
		return "", fmt.Errorf("could not create github issue: %w", err)
	}

	return strconv.Itoa(issue.GetNumber()), nil
}
