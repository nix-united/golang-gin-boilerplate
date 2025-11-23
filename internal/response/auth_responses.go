package response

// AuthTokenResponse represents a RFC 6749 compliant token response with refresh token.
type AuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`

	// ExpiresIn represents the lifetime in seconds of the access token.
	ExpiresIn int64 `json:"expires_in"`

	RefreshToken string `json:"refresh_token"`
}
