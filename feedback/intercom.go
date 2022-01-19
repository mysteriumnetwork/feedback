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
	"text/template"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/intercom/intercom-go.v2"
)

// keys of intercom custom attributes
const (
	USER_ROLE_KEY             = "user_role"
	NODE_IDENTITY_KEY         = "node_identity"
	NODE_COUNTRY_KEY          = "node_country"
	IP_TYPE_KEY               = "ip_type"
	IP_KEY                    = "ip"
	CONTACT_TYPE              = "contact"
	USER_TYPE                 = "user"
	DEFAULT_INTERCOM_BASE_URL = "https://api.intercom.io"
)

// IntercomReporter reports issues to Intercom
type IntercomReporter struct {
	client          *intercom.Client
	httpClient      *http.Client
	messageTemplate *template.Template
	intercomBaseURL string
}

// NewIntercomReporterOpts Reporter initialization options
type NewIntercomReporterOpts struct {
	Token           string
	IntercomBaseURL string
}

// NewIntercomReporter creates a new IntercomReporter
func NewIntercomReporter(opts *NewIntercomReporterOpts) *IntercomReporter {
	intercomClient := intercom.NewClient(opts.Token, "")
	if opts.IntercomBaseURL != "" {
		setBaseURI := intercom.BaseURI(opts.IntercomBaseURL)
		intercomClient.Option(setBaseURI)
	}
	messageTemplate := template.Must(template.New("messageTemplate").Parse(messageTemplate))
	rep := &IntercomReporter{
		client:          intercomClient,
		messageTemplate: messageTemplate,
		httpClient:      &http.Client{},
		intercomBaseURL: DEFAULT_INTERCOM_BASE_URL,
	}
	if opts.IntercomBaseURL != "" {
		rep.intercomBaseURL = opts.IntercomBaseURL
	}
	return rep
}

const messageTemplate = `
📅 {{.Timestamp}}

Description:

{{.Description}}

Logs:

{{.LogURL}}

`

// ReportIssue creates a issue message for the user in intercom
func (rep *IntercomReporter) ReportIssue(report *Report) (conversationId string, err error) {
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

	if report.UserId != "" {
		// try update visitor (will become lead)
		err = rep.updateVisitor(report.UserId, &updateVisitorRequest{
			Email: report.Email,
			CustomAttributes: updateVisitorRequestCustomAttributes{
				NodeIdentity: report.NodeIdentity,
				UserRole:     report.UserType,
				NodeCountry:  report.NodeCountry,
				IpType:       report.IpType,
				Ip:           report.Ip,
			},
		})
		if err != nil {
			log.Warn().Msgf("could not update visitor %s\n", report.UserId)
		}
		visitorUpdated := (err == nil)

		// try update contact
		contactUpdated := false
		if !visitorUpdated {
			contact, err := rep.client.Contacts.FindByUserID(report.UserId)
			if err != nil {
				log.Warn().Msgf("could not update contact %s\n", report.UserId)
			}
			if err == nil {
				contact.Email = report.Email
				contact.CustomAttributes = map[string]interface{}{
					NODE_IDENTITY_KEY: report.NodeIdentity,
					USER_ROLE_KEY:     report.UserType,
					NODE_COUNTRY_KEY:  report.NodeCountry,
					IP_TYPE_KEY:       report.IpType,
					IP_KEY:            report.Ip,
				}
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
				log.Warn().Msgf("could not update user %s\n", report.UserId)
			}
			if err == nil {
				user.Email = report.Email
				user.CustomAttributes = map[string]interface{}{
					NODE_IDENTITY_KEY: report.NodeIdentity,
					USER_ROLE_KEY:     report.UserType,
					NODE_COUNTRY_KEY:  report.NodeCountry,
					IP_TYPE_KEY:       report.IpType,
					IP_KEY:            report.Ip,
				}
				_, err := rep.client.Users.Save(&user)
				if err != nil {
					return "", fmt.Errorf("saving user failed (%s): %w", user.ID, err)
				}
				userUpdated = true
			}
		}

		if !visitorUpdated && !contactUpdated && !userUpdated {
			return "", fmt.Errorf("could not update visitor, contact or user (%s): %w", report.UserId, err)
		}

		userType := CONTACT_TYPE
		if userUpdated {
			userType = USER_TYPE
		}

		conversationId, err = rep.createConversation(&createConversationRequest{
			From: createConversationRequestFrom{
				UserType: userType,
				UserId:   &report.UserId,
			},
			Body: body.String(),
		})
		if err != nil {
			return "", fmt.Errorf("could not create conversation for user (%s): %w", report.UserId, err)
		}
		return conversationId, nil
	}

	contact, err := rep.client.Contacts.Create(&intercom.Contact{
		Email: report.Email,
		CustomAttributes: map[string]interface{}{
			NODE_IDENTITY_KEY: report.NodeIdentity,
			NODE_COUNTRY_KEY:  report.NodeCountry,
			USER_ROLE_KEY:     report.UserType,
			IP_TYPE_KEY:       report.IpType,
			IP_KEY:            report.Ip,
		},
	})
	if err != nil {
		return "", fmt.Errorf("could not create contact: %w", err)
	}

	conversationId, err = rep.createConversation(&createConversationRequest{
		From: createConversationRequestFrom{
			UserType: CONTACT_TYPE,
			Id:       &contact.ID,
		},
		Body: body.String(),
	})
	if err != nil {
		return "", fmt.Errorf("could not create conversation for user with id (%s): %w", contact.ID, err)
	}
	return conversationId, nil
}

type updateVisitorRequestCustomAttributes struct {
	NodeIdentity string `json:"node_identity"`
	UserRole     string `json:"user_role"`
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
	req, err := http.NewRequest(http.MethodPut, rep.intercomBaseURL+"/visitors?user_id="+userId, bytes.NewBuffer(data))
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
	UserType string  `json:"type"`
	UserId   *string `json:"user_id"`
	Id       *string `json:"id"`
}

type createConversationRequest struct {
	From createConversationRequestFrom `json:"from"`
	Body string                        `json:"body"`
}

type createConversationResponse struct {
	Id string `json:"id"`
}

func (rep *IntercomReporter) createConversation(createConversationRequest *createConversationRequest) (string, error) {
	data, err := json.Marshal(createConversationRequest)
	if err != nil {
		return "", fmt.Errorf("marshal createConversationRequest failed: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, rep.intercomBaseURL+"/conversations", bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("creating createConversation http request failed: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+rep.client.AppID)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := rep.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("createConversation http request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("createConversation http request failed with code (%d): %w", resp.StatusCode, err)
	}

	var result createConversationResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("createConversationResponse parsing failed: %w", err)
	}
	return result.Id, nil
}
