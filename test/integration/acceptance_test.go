package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/nix-united/golang-gin-boilerplate/internal/provider"
	"github.com/nix-united/golang-gin-boilerplate/internal/request"
	"github.com/nix-united/golang-gin-boilerplate/internal/response"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAcceptance(t *testing.T) {
	registerRequest := request.RegisterRequest{
		BasicAuthRequest: &request.BasicAuthRequest{
			Email:    "example@email.com",
			Password: "some-password",
		},
		FullName: "full-name",
	}
	rawRegisterRequest, err := json.Marshal(registerRequest)
	require.NoError(t, err)

	loginRequest := registerRequest.BasicAuthRequest
	rawLoginRequest, err := json.Marshal(loginRequest)
	require.NoError(t, err)

	createPostRequest := request.CreatePostRequest{
		BasicPost: &request.BasicPost{
			Title:   "Title",
			Content: "Content",
		},
	}
	rawCreatePostRequest, err := json.Marshal(createPostRequest)
	require.NoError(t, err)

	var (
		createdPost response.CreatePostResponse
		accessToken string
	)

	t.Run("It should register an user", func(t *testing.T) {
		httpResponse, err := http.Post(
			applicationURL.JoinPath("/users").String(),
			"application/json",
			bytes.NewReader(rawRegisterRequest),
		)
		require.NoError(t, err)
		defer func() {
			assert.NoError(t, httpResponse.Body.Close())
		}()

		require.Equal(t, http.StatusOK, httpResponse.StatusCode)
	})

	t.Run("It should login", func(t *testing.T) {
		httpResponse, err := http.Post(
			applicationURL.JoinPath("/login").String(),
			"application/json",
			bytes.NewReader(rawLoginRequest),
		)
		require.NoError(t, err)
		defer func() {
			assert.NoError(t, httpResponse.Body.Close())
		}()

		require.Equal(t, http.StatusOK, httpResponse.StatusCode)

		rawResponse, err := io.ReadAll(httpResponse.Body)
		require.NoError(t, err)

		var loginResponse provider.Success
		err = json.Unmarshal(rawResponse, &loginResponse)
		require.NoError(t, err)

		require.NotEmpty(t, loginResponse.Token)
		require.NotEmpty(t, loginResponse.Expire)

		accessToken = loginResponse.Token
	})

	t.Run("It should create a post", func(t *testing.T) {
		httpRequest, err := http.NewRequest(
			http.MethodPost,
			applicationURL.JoinPath("/posts").String(),
			bytes.NewReader(rawCreatePostRequest),
		)
		require.NoError(t, err)

		httpRequest.Header.Set("Content-Type", "application/json")
		httpRequest.Header.Set("Authorization", "Bearer "+accessToken)

		httpResponse, err := http.DefaultClient.Do(httpRequest)
		require.NoError(t, err)
		defer func() {
			assert.NoError(t, httpResponse.Body.Close())
		}()

		require.Equal(t, http.StatusOK, httpResponse.StatusCode)

		rawResponse, err := io.ReadAll(httpResponse.Body)
		require.NoError(t, err)

		var createPostResponse response.CreatePostResponse
		err = json.Unmarshal(rawResponse, &createPostResponse)
		require.NoError(t, err)

		assert.NotEmpty(t, createPostResponse.ID)
		assert.Equal(t, createPostRequest.Title, createPostResponse.Title)
		assert.Equal(t, createPostRequest.Content, createPostResponse.Content)

		createdPost = createPostResponse
	})

	t.Run("It should fetch a newly created post", func(t *testing.T) {
		httpRequest, err := http.NewRequest(
			http.MethodGet,
			applicationURL.JoinPath(fmt.Sprintf("/post/%d", createdPost.ID)).String(),
			http.NoBody,
		)
		require.NoError(t, err)

		httpRequest.Header.Set("Content-Type", "application/json")
		httpRequest.Header.Set("Authorization", "Bearer "+accessToken)

		httpResponse, err := http.DefaultClient.Do(httpRequest)
		require.NoError(t, err)
		defer func() {
			assert.NoError(t, httpResponse.Body.Close())
		}()

		require.Equal(t, http.StatusOK, httpResponse.StatusCode)

		rawResponse, err := io.ReadAll(httpResponse.Body)
		require.NoError(t, err)

		var getPostResponse response.GetPostResponse
		err = json.Unmarshal(rawResponse, &getPostResponse)
		require.NoError(t, err)

		assert.Equal(t, createdPost.ID, getPostResponse.ID)
		assert.Equal(t, createdPost.Title, getPostResponse.Title)
		assert.Equal(t, createdPost.Content, getPostResponse.Content)
	})
}
