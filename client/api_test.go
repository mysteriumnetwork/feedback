package client

import (
	"net/url"
	"os"
	"testing"

	"github.com/mysteriumnetwork/feedback/feedback"
	"github.com/mysteriumnetwork/feedback/infra"
	"github.com/mysteriumnetwork/feedback/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockReporter struct{}

func (m *mockReporter) ReportIssue(report *feedback.Report) (issueId string, err error) {
	return "12", nil
}

func (m *mockReporter) GetBugReportMessage(report *feedback.Report) (message string, err error) {
	return "test message 123", nil
}

type mockUploader struct{}

func (m *mockUploader) Upload(filepath string) (url *url.URL, err error) {
	url, err = url.Parse("http://uploadurl.com/" + filepath)
	if err != nil {
		return nil, err
	}
	return url, nil
}

func TestApi(t *testing.T) {
	endpoint := feedback.NewEndpoint(&mockReporter{}, &mockReporter{}, &mockUploader{}, infra.NewRateLimiter(99999))

	srvr := server.New(
		endpoint,
	)
	go srvr.Serve()

	tests := []struct {
		name          string
		request       *feedback.CreateIntercomIssueRequest
		filename      string
		fileContent   []byte
		fails         bool
		errorContains string
	}{
		{
			"description too short",
			&feedback.CreateIntercomIssueRequest{
				UserId:       "sdadas-424-dsfsd",
				Description:  "too short",
				NodeIdentity: "0x55345345345345345",
			},
			"filename1.txt",
			[]byte("hello test file"),
			true,
			"too short",
		},
		{
			"empty identity",
			&feedback.CreateIntercomIssueRequest{
				UserId:       "sdadas-424-dsfsd",
				Description:  "description which has enough length to be okay",
				NodeIdentity: "",
			},
			"filename1.txt",
			[]byte("hello test file"),
			true,
			"field is required: nodeIdentity",
		},
		{
			"no identity",
			&feedback.CreateIntercomIssueRequest{
				UserId:      "sdadas-424-dsfsd",
				Description: "description which has enough length to be okay",
			},
			"filename1.txt",
			[]byte("hello test file"),
			true,
			"field is required: nodeIdentity",
		},
		{
			"no userid or email",
			&feedback.CreateIntercomIssueRequest{
				Description:  "description which has enough length to be okay",
				NodeIdentity: "0x55345345345345345",
			},
			"filename1.txt",
			[]byte("hello test file"),
			true,
			"field is required: email or userId",
		},
		{
			"userid",
			&feedback.CreateIntercomIssueRequest{
				UserId:       "sdadas-424-dsfsd",
				Description:  "description which has enough length to be okay",
				NodeIdentity: "0x55345345345345345",
			},
			"filename1.txt",
			[]byte("hello test file"),
			false,
			"",
		},
		{
			"email",
			&feedback.CreateIntercomIssueRequest{
				Email:        "dfsfsdf@gmail.com",
				Description:  "description which has enough length to be okay",
				NodeIdentity: "0x55345345345345345",
			},
			"filename1.txt",
			[]byte("hello test file"),
			false,
			"",
		},
	}
	for _, test := range tests {
		apiClient, err := NewFeedbackAPI("http://localhost:8080")
		require.NoError(t, err)
		t.Run(test.name, func(t *testing.T) {
			f, err := os.CreateTemp("", test.filename)
			require.NoError(t, err)
			defer os.Remove(f.Name())

			_, err = f.Write(test.fileContent)
			require.NoError(t, err)

			resp, apierr, err := apiClient.CreateIntercomIssue(*test.request, f.Name())
			if test.fails {
				assert.NotNil(t, apierr)
				assert.Nil(t, err)
				assert.Contains(t, apierr.Errors[0].Message, test.errorContains)
			} else {
				assert.NotNil(t, resp)
				assert.Nil(t, apierr)
				assert.Nil(t, err)
			}
		})
	}
}
