package feedback

import (
	"net/url"
	"testing"

	"github.com/mysteriumnetwork/feedback/client"
	"github.com/mysteriumnetwork/feedback/e2e"
	"github.com/mysteriumnetwork/feedback/infra"
	"github.com/mysteriumnetwork/feedback/server"
	"github.com/stretchr/testify/assert"
)

type mockReporter struct{}

func (m *mockReporter) ReportIssue(report *Report) (issueId string, err error) {
	return "12", nil
}

type mockUploader struct{}

func (m *mockUploader) Upload(filepath string) (url *url.URL, err error) {
	url, err = url.Parse("http://uploadurl.com/file-id")
	if err != nil {
		return nil, err
	}
	return url, nil
}

func TestApi(t *testing.T) {
	skipFileUpload := false
	endpoint := NewEndpoint(&mockReporter{}, &mockReporter{}, &mockUploader{}, infra.NewRateLimiter(99999), &skipFileUpload)

	srvr := server.New(
		endpoint,
	)
	go srvr.Serve()

	tests := []struct {
		name          string
		request       *client.CreateIntercomIssueRequest
		fileContent   []byte
		fails         bool
		errorContains string
		errorCode     int
	}{
		{
			"description too short",
			&client.CreateIntercomIssueRequest{
				UserId:       "sdadas-424-dsfsd",
				Description:  "too short",
				NodeIdentity: "0x55345345345345345",
			},
			[]byte("hello test file"),
			true,
			"too short",
			400,
		},
		{
			"empty identity",
			&client.CreateIntercomIssueRequest{
				UserId:       "sdadas-424-dsfsd",
				Description:  "description which has enough length to be okay",
				NodeIdentity: "",
			},
			[]byte("hello test file"),
			true,
			"field is required: nodeIdentity",
			400,
		},
		{
			"no identity",
			&client.CreateIntercomIssueRequest{
				UserId:      "sdadas-424-dsfsd",
				Description: "description which has enough length to be okay",
			},
			[]byte("hello test file"),
			true,
			"field is required: nodeIdentity",
			400,
		},
		{
			"no file",
			&client.CreateIntercomIssueRequest{
				UserId:       "sdadas-424-dsfsd",
				Description:  "description which has enough length to be okay",
				NodeIdentity: "0x55345345345345345",
			},
			[]byte{},
			true,
			"field is required: file",
			400,
		},
		{
			"no userid or email",
			&client.CreateIntercomIssueRequest{
				Description:  "description which has enough length to be okay",
				NodeIdentity: "0x55345345345345345",
			},
			[]byte("hello test file"),
			true,
			"field is required: email or userId",
			400,
		},
		{
			"userid",
			&client.CreateIntercomIssueRequest{
				UserId:       "sdadas-424-dsfsd",
				Description:  "description which has enough length to be okay",
				NodeIdentity: "0x55345345345345345",
			},
			[]byte("hello test file"),
			false,
			"",
			0,
		},
		{
			"email",
			&client.CreateIntercomIssueRequest{
				Email:        "dfsfsdf@gmail.com",
				Description:  "description which has enough length to be okay",
				NodeIdentity: "0x55345345345345345",
			},
			[]byte("hello test file"),
			false,
			"",
			0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := e2e.SendReportIntercomIssueRequest(t, test.request, test.fileContent)
			if test.fails {
				assert.NotNil(t, err)
				assert.Equal(t, 400, err.Code)
				assert.Contains(t, err.Errors[0].Message, test.errorContains)
			} else {
				assert.Nil(t, err)
			}

		})
	}
}
