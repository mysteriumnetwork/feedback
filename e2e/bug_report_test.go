package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mysteriumnetwork/feedback/feedback"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/walkerus/go-wiremock"
)

func TestBugReport(t *testing.T) {
	wiremockClient := wiremock.NewClient("http://localhost:8090")
	s3Downloader, err := newS3Downloader("node-user-reports")
	assert.NoError(t, err)

	if apiClient == nil {
		t.Error("apiClient is nil")
	}

	t.Run("simple bug report", func(t *testing.T) {
		fileContent := []byte("hello")
		defer wiremockClient.Reset()
		err := wiremockClient.StubFor(wiremock.Post(wiremock.URLPathMatching("/github/repos/github-owner/github-repo/issues")).WillReturn(
			`{"id":123,"number":456}`,
			map[string]string{"Content-Type": "application/json"},
			200,
		))
		assert.Nil(t, err)

		req := &feedback.CreateBugReportRequest{
			NodeIdentity: common.HexToAddress("0x12345").String(),
			Description:  "long description for an issue that I have in my node",
			Email:        "test@gmail.com",
		}

		f, err := os.CreateTemp("", "test_file")
		require.NoError(t, err)
		defer os.Remove(f.Name())

		_, err = f.Write(fileContent)
		require.NoError(t, err)

		resp, apiErr, err := apiClient.CreateBugReport(*req, f.Name())
		assert.Nil(t, err)
		assert.Nil(t, apiErr)
		assert.NotNil(t, resp)
		assert.Equal(t, "456", resp.IssueId)
		assert.Equal(t, req.Email, resp.Email)
		assert.Equal(t, req.NodeIdentity, resp.NodeIdentity)
		assert.True(t, strings.Contains(resp.Message, "long description for an issue that I have in my node"), fmt.Sprintf("%s does not contain %s", resp.Message, "long description for an issue that I have in my node"))
		assert.True(t, strings.Contains(resp.Message, "Logs: http://someweb.com/"), fmt.Sprintf("%s does not contain %s", resp.Message, "Logs: http://someweb.com/"))

		content, err := s3Downloader.getFileContent(t, filepath.Base(f.Name()))
		assert.NoError(t, err)
		assert.Equal(t, fileContent, content)
	})
}
