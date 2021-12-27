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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/intercom/intercom-go.v2"
)

const (
	INTERCOM_BASE_API = "https://api.intercom.io"
)

// IntercomReporter reports issues to Intercom
type IntercomReporter struct {
	client          *intercom.Client
	httpClient      *http.Client
	messageTemplate *template.Template
}

// NewIntercomReporterOpts Reporter initialization options
type NewIntercomReporterOpts struct {
	Token string
}

// NewIntercomReporter creates a new IntercomReporter
func NewIntercomReporter(opts *NewIntercomReporterOpts) *IntercomReporter {
	intercomClient := intercom.NewClient(opts.Token, "")
	messageTemplate := template.Must(template.New("messageTemplate").Parse(messageTemplate))
	return &IntercomReporter{
		client:          intercomClient,
		messageTemplate: messageTemplate,
		httpClient:      &http.Client{},
	}
}

// Report user report data
type IntercomReport struct {
	UserId       string
	NodeIdentity string
	UserType     string
	NodeCountry  string
	IpType       string
	Ip           string
	Description  string
	Email        string
	LogURL       url.URL
}

const messageTemplate = `
📅 {{.Timestamp}}

Description:

{{.Description}}

Logs:

{{.LogURL}}

`

// ReportIssue creates a issue message for the user in intercom
func (rep *IntercomReporter) ReportIssue(report *IntercomReport) (issueId string, err error) {
	templateOpts := struct {
		Description,
		Timestamp,
		LogURL string
	}{
		Description: report.Description,
		Timestamp:   time.Now().String(),
		LogURL:      report.LogURL.String(),
	}
	var body bytes.Buffer
	err = rep.messageTemplate.Execute(&body, templateOpts)
	if err != nil {
		return "", fmt.Errorf("could not generate message body with report (%+v): %w", templateOpts, err)
	}

	// try update visitor
	err = rep.updateVisitor(report.UserId, &updateVisitorRequest{
		Email: report.Email,
		CustomAttributes: updateVisitorRequestCustomAttributes{
			NodeIdentity: report.NodeIdentity,
			IsConsumer:   (strings.ToLower(report.UserType) == "consumer"),
			NodeCountry:  report.NodeCountry,
			IpType:       report.IpType,
			Ip:           report.Ip,
		},
	})
	if err != nil {
		log.Warn().Msg("could not update visitor")
	}
	visitorUpdated := (err == nil)

	// try update contact
	contactUpdated := false
	if !visitorUpdated {
		contact, err := rep.client.Contacts.FindByUserID(report.UserId)
		if err != nil {
			log.Warn().Msg("could not update contact")
		}
		if err == nil {
			contact.Email = report.Email
			contact.CustomAttributes["node_identity"] = report.NodeIdentity
			contact.CustomAttributes["node_country"] = report.NodeCountry
			contact.CustomAttributes["is_consumer"] = (strings.ToLower(report.UserType) == "consumer")
			contact.CustomAttributes["ip_type"] = report.IpType
			contact.CustomAttributes["ip"] = report.Ip
			_, err := rep.client.Contacts.Update(&contact)
			if err != nil {
				return "", fmt.Errorf("could not update contact (%s): %w", contact.ID, err)
			}
			contactUpdated = true
		}
	}
	// try update user
	userUpdated := false
	if !visitorUpdated && !contactUpdated {
		user, err := rep.client.Users.FindByUserID(report.UserId)
		if err != nil {
			log.Warn().Msg("could not update user")
		}
		if err == nil {
			user.Email = report.Email
			user.CustomAttributes["node_identity"] = report.NodeIdentity
			user.CustomAttributes["node_country"] = report.NodeCountry
			user.CustomAttributes["is_consumer"] = (strings.ToLower(report.UserType) == "consumer")
			user.CustomAttributes["ip_type"] = report.IpType
			user.CustomAttributes["ip"] = report.Ip
			_, err := rep.client.Users.Save(&user)
			if err != nil {
				return "", fmt.Errorf("could not update user (%s): %w", user.ID, err)
			}
			userUpdated = true
		}
	}

	if !visitorUpdated && !contactUpdated && !userUpdated {
		return "", fmt.Errorf("could not update visitor, contact or user (%s): %w", report.UserId, err)
	}

	userType := "contact"
	if userUpdated {
		userType = "user"
	}

	err = rep.createConversation(&createConversationRequest{
		From: createConversationRequestFrom{
			UserType: userType,
			UserId:   report.UserId,
		},
		Body: body.String(),
	})
	if err != nil {
		return "", fmt.Errorf("could not create conversation for user (%s): %w", report.UserId, err)
	}

	return strconv.Itoa(123), nil
}

type updateVisitorRequestCustomAttributes struct {
	NodeIdentity string `json:"node_identity"`
	IsConsumer   bool   `json:"is_consumer"`
	NodeCountry  string `json:"node_country"`
	IpType       string `json:"ip_type"`
	Ip           string `json:"ip"`
}

type updateVisitorRequest struct {
	Email            string                               `json:"email"`
	CustomAttributes updateVisitorRequestCustomAttributes `json:"custom_attributes"`
}

//updates visitor so it becomes a lead
func (rep *IntercomReporter) updateVisitor(userId string, updateVisitorRequest *updateVisitorRequest) error {
	data, err := json.Marshal(updateVisitorRequest)
	if err != nil {
		return fmt.Errorf("marshal updateVisitorRequest failed: %w", err)
	}
	req, err := http.NewRequest(http.MethodPut, INTERCOM_BASE_API+"/visitors?user_id="+userId, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("creating updateVisitor http request failed: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+rep.client.AppID)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := rep.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("updateVisitor http request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("updateVisitor http request failed with code (%d): %w", resp.StatusCode, err)
	}
	return nil
}

type createConversationRequestFrom struct {
	UserType string `json:"type"`
	UserId   string `json:"user_id"`
}

type createConversationRequest struct {
	From createConversationRequestFrom `json:"from"`
	Body string                        `json:"body"`
}

func (rep *IntercomReporter) createConversation(createConversationRequest *createConversationRequest) error {
	data, err := json.Marshal(createConversationRequest)
	if err != nil {
		return fmt.Errorf("marshal updateVisitorRequest failed: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, INTERCOM_BASE_API+"/conversations", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("creating createConversation http request failed: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+rep.client.AppID)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := rep.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("createConversation http request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("createConversation http request failed with code (%d): %w", resp.StatusCode, err)
	}

	return nil
}
