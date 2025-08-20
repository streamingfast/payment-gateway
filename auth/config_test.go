package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_new(t *testing.T) {
	tests := []struct {
		name        string
		dsn         string
		expect      *Config
		expectError bool
	}{
		{
			name: "complete config with all parameters",
			dsn:  "tgm://auth.thegraph.market?pubkeyurl=https://auth.thegraph.market/.well-known/jwks.json&reissue-jwt-max-age-secs=600&key=server_23bc9f9ccbd62e9f23d92513cdc76b55",
			expect: &Config{
				Endpoint:          "auth.thegraph.market",
				Insecure:          false,
				Plaintext:         false,
				Key:               "server_23bc9f9ccbd62e9f23d92513cdc76b55",
				ReissueJWTAgeSecs: 600,
				PubKeyURL:         "https://auth.thegraph.market/.well-known/jwks.json",
				PubKeyBase64:      "",
				ReissueURL:        "https://auth.thegraph.market/v1/auth/reissue",
				IssueURL:          "https://auth.thegraph.market/v1/auth/issue",
			},
		},
		{
			name: "config with custom port",
			dsn:  "tgm://auth.example.com:8443?key=test_key",
			expect: &Config{
				Endpoint:          "auth.example.com:8443",
				Insecure:          false,
				Plaintext:         false,
				Key:               "test_key",
				ReissueJWTAgeSecs: 600, // default value
				PubKeyURL:         "https://auth.thegraph.market/.well-known/jwks.json",
				PubKeyBase64:      "",
				ReissueURL:        "https://auth.example.com:8443/v1/auth/reissue",
				IssueURL:          "https://auth.example.com:8443/v1/auth/issue",
			},
		},
		{
			name: "config with insecure and plaintext flags",
			dsn:  "tgm://localhost:8080?insecure=true&plaintext=true",
			expect: &Config{
				Endpoint:          "localhost:8080",
				Insecure:          true,
				Plaintext:         true,
				Key:               "",
				ReissueJWTAgeSecs: 600, // default value
				PubKeyURL:         "https://auth.thegraph.market/.well-known/jwks.json",
				PubKeyBase64:      "",
				ReissueURL:        "http://localhost:8080/v1/auth/reissue",
				IssueURL:          "http://localhost:8080/v1/auth/issue",
			},
		},
		{
			name: "config with standard https port (should not be included)",
			dsn:  "tgm://auth.example.com:443?key=mykey",
			expect: &Config{
				Endpoint:          "auth.example.com",
				Insecure:          false,
				Plaintext:         false,
				Key:               "mykey",
				ReissueJWTAgeSecs: 600,
				PubKeyURL:         "https://auth.thegraph.market/.well-known/jwks.json",
				PubKeyBase64:      "",
				ReissueURL:        "https://auth.example.com/v1/auth/reissue",
				IssueURL:          "https://auth.example.com/v1/auth/issue",
			},
		},
		{
			name: "config with standard http port (should not be included)",
			dsn:  "tgm://auth.example.com:80?plaintext=true",
			expect: &Config{
				Endpoint:          "auth.example.com",
				Insecure:          false,
				Plaintext:         true,
				Key:               "",
				ReissueJWTAgeSecs: 600,
				PubKeyURL:         "https://auth.thegraph.market/.well-known/jwks.json",
				PubKeyBase64:      "",
				ReissueURL:        "http://auth.example.com/v1/auth/reissue",
				IssueURL:          "http://auth.example.com/v1/auth/issue",
			},
		},
		{
			name: "config with encoded URL in pubkeyurl",
			dsn:  "tgm://auth.example.com?pubkeyurl=https%3A%2F%2Fauth.example.com%2F.well-known%2Fjwks.json",
			expect: &Config{
				Endpoint:          "auth.example.com",
				Insecure:          false,
				Plaintext:         false,
				Key:               "",
				ReissueJWTAgeSecs: 600,
				PubKeyURL:         "https://auth.example.com/.well-known/jwks.json",
				PubKeyBase64:      "",
				ReissueURL:        "https://auth.example.com/v1/auth/reissue",
				IssueURL:          "https://auth.example.com/v1/auth/issue",
			},
		},
		{
			name: "config with base64 public key",
			dsn:  "tgm://auth.example.com?pubkeybase64=LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUF4SUpKczQ3YnBZcjRLCi0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ==",
			expect: &Config{
				Endpoint:          "auth.example.com",
				Insecure:          false,
				Plaintext:         false,
				Key:               "",
				ReissueJWTAgeSecs: 600,
				PubKeyURL:         "",
				PubKeyBase64:      "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUF4SUpKczQ3YnBZcjRLCi0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ==",
				ReissueURL:        "https://auth.example.com/v1/auth/reissue",
				IssueURL:          "https://auth.example.com/v1/auth/issue",
			},
		},
		{
			name: "complete config with base64 public key",
			dsn:  "tgm://auth.secure.com:9443?pubkeybase64=LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUF4SUpKczQ3YnBZcjRLCi0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ==&reissue-jwt-max-age-secs=300&key=api_key_123&insecure=true",
			expect: &Config{
				Endpoint:          "auth.secure.com:9443",
				Insecure:          true,
				Plaintext:         false,
				Key:               "api_key_123",
				ReissueJWTAgeSecs: 300,
				PubKeyURL:         "",
				PubKeyBase64:      "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUF4SUpKczQ3YnBZcjRLCi0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ==",
				ReissueURL:        "https://auth.secure.com:9443/v1/auth/reissue",
				IssueURL:          "https://auth.secure.com:9443/v1/auth/issue",
			},
		},
		{
			name:        "error when both pubkeyurl and pubkeybase64 are provided",
			dsn:         "tgm://auth.example.com?pubkeyurl=https://auth.example.com/.well-known/jwks.json&pubkeybase64=LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0K",
			expectError: true,
		},
		{
			name: "minimal hostname",
			dsn:  "tgm://",
			expect: &Config{
				Endpoint:          "auth.thegraph.market",
				Insecure:          false,
				Plaintext:         false,
				Key:               "",
				ReissueJWTAgeSecs: 600,
				PubKeyURL:         "https://auth.thegraph.market/.well-known/jwks.json",
				PubKeyBase64:      "",
				ReissueURL:        "https://auth.thegraph.market/v1/auth/reissue",
				IssueURL:          "https://auth.thegraph.market/v1/auth/issue",
			},
		},
		{
			name:        "empty DSN",
			dsn:         "",
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, err := newConfig(test.dsn)
			if test.expectError {
				require.Error(t, err)
				// Special check for mutual exclusivity error
				if test.name == "error when both pubkeyurl and pubkeybase64 are provided" {
					assert.Contains(t, err.Error(), "only one of pubkeyurl or pubkeybase64 can be provided")
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expect, c)
			}
		})
	}
}
