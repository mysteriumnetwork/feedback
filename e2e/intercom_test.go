package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mysteriumnetwork/feedback/feedback"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/walkerus/go-wiremock"
)

func TestIntercomReporting(t *testing.T) {
	wiremockClient := wiremock.NewClient("http://localhost:8090")
	s3Downloader, err := newS3Downloader("node-user-reports")
	assert.NoError(t, err)

	if apiClient == nil {
		t.Error("apiClient is nil")
	}

	t.Run("report an issue on intercom with userId (visitor)", func(t *testing.T) {
		userId := "dfcte009-73cc-4638-be3d-f4tjd22a22a2"
		fileContent := []byte("hello")
		defer wiremockClient.Reset()
		err := wiremockClient.StubFor(wiremock.Get(wiremock.URLPathMatching("/visitors")).WithQueryParam("user_id", wiremock.EqualTo(userId)).WillReturn(
			`{"type":"visitor","id":"61cdbfc18b3349339c0b626c","user_id":"dfcte009-73cc-4638-be3d-f4tjd22a22a2","anonymous":true,"email":""}`,
			map[string]string{"Content-Type": "application/json"},
			200,
		))
		assert.Nil(t, err)
		err = wiremockClient.StubFor(wiremock.Put(wiremock.URLPathMatching("/visitors")).WithQueryParam("user_id", wiremock.EqualTo(userId)).WillReturn(
			`{"type":"visitor","id":"61cdbfc18b3349339c0b626c","user_id":"dfcde009-73cc-4658-be8d-f446922a22a2","anonymous":true,"email":""}`,
			map[string]string{"Content-Type": "application/json"},
			200,
		))
		assert.Nil(t, err)
		err = wiremockClient.StubFor(wiremock.Post(wiremock.URLPathMatching("/conversations")).
			WithBodyPattern(wiremock.MatchingJsonPath("$.from[?(@.user_id == '"+userId+"')]")).
			WillReturn(
				`{}`,
				map[string]string{"Content-Type": "application/json"},
				200,
			))
		assert.Nil(t, err)

		req := &feedback.CreateIntercomIssueRequest{
			UserId:       userId,
			NodeIdentity: "0x5345765675656",
			Description:  "long description for an issue that I have in my node",
		}

		f, err := os.CreateTemp("", "test_file")
		require.NoError(t, err)
		defer os.Remove(f.Name())

		_, err = f.Write(fileContent)
		require.NoError(t, err)

		resp, apiErr, err := apiClient.CreateIntercomIssue(*req, f.Name())
		assert.Nil(t, err)
		assert.Nil(t, apiErr)
		assert.NotNil(t, resp)

		content, err := s3Downloader.getFileContent(t, filepath.Base(f.Name()))
		assert.NoError(t, err)
		assert.Equal(t, fileContent, content)
	})

	t.Run("report an issue on intercom with userId (lead)", func(t *testing.T) {
		userId := "dfyye009-73cc-4888-be3d-f4t6666662a2"
		id := "61cd88818b3349339c44426c"
		fileContent := []byte("hello-2")
		defer wiremockClient.Reset()
		err := wiremockClient.StubFor(wiremock.Get(
			wiremock.URLPathMatching("/contacts")).
			WithQueryParam("user_id", wiremock.EqualTo(userId)).
			WillReturn(
				`{"type":"lead","id":"`+id+`","user_id":"`+userId+`","anonymous":true,"email":""}`,
				map[string]string{"Content-Type": "application/json"},
				200,
			))
		assert.Nil(t, err)
		err = wiremockClient.StubFor(wiremock.Post(
			wiremock.URLPathMatching("/contacts")).
			WithBodyPattern(wiremock.MatchingJsonPath("$[?(@.id == '"+id+"')]")).
			WillReturn(
				`{"type":"lead","id":"`+id+`","user_id":"`+userId+`","anonymous":true,"email":""}`,
				map[string]string{"Content-Type": "application/json"},
				200,
			))
		assert.Nil(t, err)
		err = wiremockClient.StubFor(
			wiremock.Post(wiremock.URLPathMatching("/conversations")).
				WithBodyPattern(wiremock.MatchingJsonPath("$.from[?(@.user_id == '"+userId+"')]")).
				WillReturn(
					`{}`,
					map[string]string{"Content-Type": "application/json"},
					200,
				))
		assert.Nil(t, err)

		req := &feedback.CreateIntercomIssueRequest{
			UserId:       userId,
			NodeIdentity: "0x5345765675656",
			Description:  "long description for an issue that I have in my node",
		}
		f, err := os.CreateTemp("", "test_file")
		require.NoError(t, err)
		defer os.Remove(f.Name())

		_, err = f.Write(fileContent)
		require.NoError(t, err)

		resp, apiErr, err := apiClient.CreateIntercomIssue(*req, f.Name())
		assert.Nil(t, err)
		assert.Nil(t, apiErr)
		assert.NotNil(t, resp)
		count, err := wiremockClient.GetCountRequests(wiremock.Get(
			wiremock.URLPathMatching("/contacts")).
			WithQueryParam("user_id", wiremock.EqualTo(userId)).Request())
		assert.Nil(t, err)
		assert.Equal(t, 1, int(count))
		count, err = wiremockClient.GetCountRequests(
			wiremock.Post(wiremock.URLPathMatching("/contacts")).
				WithBodyPattern(wiremock.MatchingJsonPath("$[?(@.id == '" + id + "')]")).Request())
		assert.Nil(t, err)
		assert.Equal(t, 1, int(count))
		content, err := s3Downloader.getFileContent(t, filepath.Base(f.Name()))
		assert.NoError(t, err)
		assert.Equal(t, fileContent, content)
	})

	t.Run("report an issue on intercom with userId (user)", func(t *testing.T) {
		userId := "d77709-73cc-4888-be3d-f4t6aaaa62a2"
		id := "61cd8887777779339c44426c"
		fileContent := []byte("hello-3")
		defer wiremockClient.Reset()
		err := wiremockClient.StubFor(wiremock.Get(
			wiremock.URLPathMatching("/users")).
			WithQueryParam("user_id", wiremock.EqualTo(userId)).WillReturn(
			`{"type":"user","id":"`+id+`","user_id":"`+userId+`","anonymous":true,"email":""}`,
			map[string]string{"Content-Type": "application/json"},
			200,
		))
		assert.Nil(t, err)
		err = wiremockClient.StubFor(wiremock.Post(
			wiremock.URLPathMatching("/users")).
			WithBodyPattern(wiremock.MatchingJsonPath("$[?(@.id == '"+id+"')]")).
			WillReturn(
				`{"type":"user","id":"`+id+`","user_id":"`+userId+`","anonymous":true,"email":""}`,
				map[string]string{"Content-Type": "application/json"},
				200,
			))
		assert.Nil(t, err)
		err = wiremockClient.StubFor(wiremock.Post(
			wiremock.URLPathMatching("/conversations")).
			WithBodyPattern(wiremock.MatchingJsonPath("$.from[?(@.user_id == '"+userId+"')]")).
			WillReturn(
				`{}`,
				map[string]string{"Content-Type": "application/json"},
				200,
			))
		assert.Nil(t, err)

		req := &feedback.CreateIntercomIssueRequest{
			UserId:       userId,
			NodeIdentity: "0x5345765675656",
			Description:  "long description for an issue that I have in my node",
		}
		f, err := os.CreateTemp("", "test_file")
		require.NoError(t, err)
		defer os.Remove(f.Name())

		_, err = f.Write(fileContent)
		require.NoError(t, err)

		resp, apiErr, err := apiClient.CreateIntercomIssue(*req, f.Name())
		assert.Nil(t, err)
		assert.Nil(t, apiErr)
		assert.NotNil(t, resp)
		count, err := wiremockClient.GetCountRequests(
			wiremock.Get(wiremock.URLPathMatching("/users")).
				WithQueryParam("user_id", wiremock.EqualTo(userId)).Request())
		assert.Nil(t, err)
		assert.Equal(t, 1, int(count))
		count, err = wiremockClient.GetCountRequests(wiremock.Post(
			wiremock.URLPathMatching("/users")).
			WithBodyPattern(wiremock.MatchingJsonPath("$[?(@.id == '" + id + "')]")).Request())
		assert.Nil(t, err)
		assert.Equal(t, 1, int(count))
		content, err := s3Downloader.getFileContent(t, filepath.Base(f.Name()))
		assert.NoError(t, err)
		assert.Equal(t, fileContent, content)
	})

	t.Run("user not found", func(t *testing.T) {
		userId := "dfcte009-73cc-4638-be3d-f4tjd22a22a2"
		fileContent := []byte("hello-4")
		defer wiremockClient.Reset()

		err := wiremockClient.StubFor(wiremock.Post(
			wiremock.URLPathMatching("/conversations")).
			WithBodyPattern(wiremock.MatchingJsonPath("$.from[?(@.user_id == '"+userId+"')]")).
			WillReturn(
				`{}`,
				map[string]string{"Content-Type": "application/json"},
				200,
			))
		assert.Nil(t, err)
		err = wiremockClient.StubFor(wiremock.Post(wiremock.URLPathMatching("/conversations")).
			WithBodyPattern(wiremock.MatchingJsonPath("$.from[?(@.user_id == '"+userId+"')]")).
			WillReturn(
				`{}`,
				map[string]string{"Content-Type": "application/json"},
				404,
			))
		assert.Nil(t, err)

		req := &feedback.CreateIntercomIssueRequest{
			UserId:       userId,
			NodeIdentity: "0x5345765675656",
			Description:  "long description for an issue that I have in my node",
		}

		f, err := os.CreateTemp("", "test_file")
		require.NoError(t, err)
		defer os.Remove(f.Name())

		_, err = f.Write(fileContent)
		require.NoError(t, err)

		_, apiErr, err := apiClient.CreateIntercomIssue(*req, f.Name())
		assert.Nil(t, err)
		assert.NotNil(t, apiErr)
		assert.Contains(t, apiErr.Message(), "service unavailable")
	})

	t.Run("report an issue on intercom with email", func(t *testing.T) {
		email := "email@something.com"
		fileContent := []byte("hello-5")
		defer wiremockClient.Reset()
		err := wiremockClient.StubFor(wiremock.Post(
			wiremock.URLPathMatching("/contacts")).
			WithBodyPattern(wiremock.MatchingJsonPath("$[?(@.email == '"+email+"')]")).
			WillReturn(
				`{"type":"contact","id":"61cdca90cbdc920d239ffa88","role":"user","email":"`+email+`","phone":null,"name":null,"avatar":null}`,
				map[string]string{"Content-Type": "application/json"},
				200,
			))
		assert.Nil(t, err)
		err = wiremockClient.StubFor(wiremock.Post(
			wiremock.URLPathMatching("/conversations")).
			WithBodyPattern(wiremock.MatchingJsonPath("$.from[?(@.id == '61cdca90cbdc920d239ffa88')]")).
			WillReturn(
				`{}`,
				map[string]string{"Content-Type": "application/json"},
				200,
			))
		assert.Nil(t, err)

		req := &feedback.CreateIntercomIssueRequest{
			Email:        email,
			NodeIdentity: "0x5345765675656",
			Description:  "long description for an issue that I have in my node",
		}

		f, err := os.CreateTemp("", "test_file")
		require.NoError(t, err)
		defer os.Remove(f.Name())

		_, err = f.Write(fileContent)
		require.NoError(t, err)

		resp, apiErr, err := apiClient.CreateIntercomIssue(*req, f.Name())
		assert.Nil(t, err)
		assert.Nil(t, apiErr)
		assert.NotNil(t, resp)
		content, err := s3Downloader.getFileContent(t, filepath.Base(f.Name()))
		assert.NoError(t, err)
		assert.Equal(t, fileContent, content)
	})
}
