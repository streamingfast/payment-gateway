package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/streamingfast/dauth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAuthenticator_ParseJWT(t *testing.T) {
	// Create a test key
	key, err := jwk.New([]byte("test-secret-key"))
	require.NoError(t, err)
	err = key.Set(jwk.AlgorithmKey, jwa.HS256)
	require.NoError(t, err)
	err = key.Set(jwk.KeyIDKey, "test-key-id")
	require.NoError(t, err)

	set := jwk.NewSet()
	set.Add(key)

	// Create authenticator with test key set
	auth := &authenticator{
		config: &Config{
			ReissueJWTAgeSecs: 600,
		},
		logger: zap.NewNop(),
		jwkSet: set,
	}

	// Create a test token
	token := jwt.New()
	token.Set(jwt.SubjectKey, "user:123")
	token.Set(jwt.IssuedAtKey, time.Now())
	token.Set(jwt.ExpirationKey, time.Now().Add(time.Hour))
	token.Set("api_key_id", "test-api-key")
	token.Set("user_id", "user-123")

	// Sign the token with key ID
	hdrs := jws.NewHeaders()
	hdrs.Set(jws.KeyIDKey, "test-key-id")
	signed, err := jwt.Sign(token, jwa.HS256, key, jwt.WithHeaders(hdrs))
	require.NoError(t, err)

	// Test parsing
	parsed, err := auth.ParseJWT(string(signed))
	require.NoError(t, err)
	assert.Equal(t, "user:123", parsed.Subject())

	// Verify custom claims
	claims := parsed.PrivateClaims()
	assert.Equal(t, "test-api-key", claims["api_key_id"])
	assert.Equal(t, "user-123", claims["user_id"])
}

func TestAuthenticator_NeedsReissue(t *testing.T) {
	auth := &authenticator{
		config: &Config{
			ReissueJWTAgeSecs: 60, // 1 minute
		},
		logger: zap.NewNop(),
	}

	tests := []struct {
		name     string
		issuedAt time.Time
		want     bool
	}{
		{
			name:     "token issued just now",
			issuedAt: time.Now(),
			want:     false,
		},
		{
			name:     "token issued 30 seconds ago",
			issuedAt: time.Now().Add(-30 * time.Second),
			want:     false,
		},
		{
			name:     "token issued 2 minutes ago",
			issuedAt: time.Now().Add(-2 * time.Minute),
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := jwt.New()
			token.Set(jwt.IssuedAtKey, tt.issuedAt)

			got := auth.needsReissue(token)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAuthenticator_ExtractAndParseJWT(t *testing.T) {
	// Create a test key
	key, err := jwk.New([]byte("test-secret-key"))
	require.NoError(t, err)
	err = key.Set(jwk.AlgorithmKey, jwa.HS256)
	require.NoError(t, err)
	err = key.Set(jwk.KeyIDKey, "test-key-id")
	require.NoError(t, err)

	set := jwk.NewSet()
	set.Add(key)

	auth := &authenticator{
		config: &Config{},
		logger: zap.NewNop(),
		jwkSet: set,
	}

	// Create a test token
	token := jwt.New()
	token.Set(jwt.SubjectKey, "user:123")
	token.Set(jwt.IssuedAtKey, time.Now())

	// Sign the token with key ID
	hdrs := jws.NewHeaders()
	hdrs.Set(jws.KeyIDKey, "test-key-id")
	signed, err := jwt.Sign(token, jwa.HS256, key, jwt.WithHeaders(hdrs))
	require.NoError(t, err)
	tokenString := string(signed)

	tests := []struct {
		name      string
		header    string
		wantErr   bool
		errString string
	}{
		{
			name:    "valid bearer token",
			header:  "Bearer " + tokenString,
			wantErr: false,
		},
		{
			name:    "valid token without bearer",
			header:  tokenString,
			wantErr: false,
		},
		{
			name:      "invalid bearer format",
			header:    "Basic " + tokenString,
			wantErr:   true,
			errString: "authorization header format must be Bearer {token}",
		},
		{
			name:      "too many parts",
			header:    "Bearer " + tokenString + " extra",
			wantErr:   true,
			errString: "authorization header format must be Bearer {token}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := auth.extractAndParseJWT(tt.header)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errString != "" {
					assert.Contains(t, err.Error(), tt.errString)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthenticator_IssueJWTFromAPIKey(t *testing.T) {
	// Create a test key
	key, err := jwk.New([]byte("test-secret-key"))
	require.NoError(t, err)
	err = key.Set(jwk.AlgorithmKey, jwa.HS256)
	require.NoError(t, err)
	err = key.Set(jwk.KeyIDKey, "test-key-id")
	require.NoError(t, err)

	set := jwk.NewSet()
	set.Add(key)

	// Create a valid test token for the mock response
	mockToken := jwt.New()
	mockToken.Set("api_key_id", "test-api-key")
	mockToken.Set("user_id", "user-123")
	mockToken.Set(jwt.IssuedAtKey, time.Now())
	mockToken.Set(jwt.ExpirationKey, time.Now().Add(time.Hour))

	// Sign the mock token with key ID
	hdrs := jws.NewHeaders()
	hdrs.Set(jws.KeyIDKey, "test-key-id")
	signed, err := jwt.Sign(mockToken, jwa.HS256, key, jwt.WithHeaders(hdrs))
	require.NoError(t, err)
	tokenString := string(signed)

	// Create a mock server for the issue endpoint
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/auth/issue", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Return a mock JWT response
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"token": tokenString}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	auth := &authenticator{
		config: &Config{
			IssueURL: server.URL + "/v1/auth/issue",
		},
		logger:     zap.NewNop(),
		jwkSet:     set,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	ctx := context.Background()
	token, err := auth.issueJWTFromAPIKey(ctx, "test-api-key")
	require.NoError(t, err)
	assert.NotNil(t, token)

	// Verify the token contains expected claims
	claims := token.PrivateClaims()
	assert.Equal(t, "test-api-key", claims["api_key_id"])
	assert.Equal(t, "user-123", claims["user_id"])
}

func TestAuthenticator_Authenticate(t *testing.T) {
	// Create a test key
	key, err := jwk.New([]byte("test-secret-key"))
	require.NoError(t, err)
	err = key.Set(jwk.AlgorithmKey, jwa.HS256)
	require.NoError(t, err)
	err = key.Set(jwk.KeyIDKey, "test-key-id")
	require.NoError(t, err)

	set := jwk.NewSet()
	set.Add(key)

	// Create a valid test token
	token := jwt.New()
	token.Set(jwt.SubjectKey, "uid:user-123")
	token.Set(jwt.IssuedAtKey, time.Now())
	token.Set(jwt.ExpirationKey, time.Now().Add(time.Hour))
	token.Set("api_key_id", "test-api-key")
	token.Set("user_id", "user-123")

	// Add feature configs
	token.Set("feature_configs", map[string]interface{}{
		"SUBSTREAMS_PARALLEL_JOBS": "10",
	})

	// Sign the token with key ID
	hdrs := jws.NewHeaders()
	hdrs.Set(jws.KeyIDKey, "test-key-id")
	signed, err := jwt.Sign(token, jwa.HS256, key, jwt.WithHeaders(hdrs))
	require.NoError(t, err)
	tokenString := string(signed)

	auth := &authenticator{
		config: &Config{
			ReissueJWTAgeSecs: 3600, // Don't trigger reissue
		},
		logger:     zap.NewNop(),
		jwkSet:     set,
		httpClient: &http.Client{},
	}

	tests := []struct {
		name      string
		headers   map[string][]string
		ipAddress string
		wantErr   bool
		checkCtx  func(t *testing.T, ctx context.Context)
	}{
		{
			name: "valid JWT in authorization header",
			headers: map[string][]string{
				"Authorization": {"Bearer " + tokenString},
			},
			ipAddress: "192.168.1.1",
			wantErr:   false,
			checkCtx: func(t *testing.T, ctx context.Context) {
				// Check that context has the expected values
				headers := dauth.FromContext(ctx)
				assert.Equal(t, "user-123", headers[dauth.SFHeaderUserID])
				assert.Equal(t, "test-api-key", headers[dauth.SFHeaderApiKeyID])
				assert.Equal(t, "192.168.1.1", headers[dauth.SFHeaderIP])
				assert.Equal(t, "10", headers["x-sf-substreams-parallel-jobs"])
			},
		},
		{
			name: "no auth headers",
			headers: map[string][]string{
				"Content-Type": {"application/json"},
			},
			ipAddress: "192.168.1.1",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			newCtx, err := auth.Authenticate(ctx, "/test", tt.headers, tt.ipAddress)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.checkCtx != nil {
					tt.checkCtx(t, newCtx)
				}
			}
		})
	}
}

func TestAuthenticator_AddClaimsToContext(t *testing.T) {
	auth := &authenticator{
		logger: zap.NewNop(),
	}

	token := jwt.New()
	token.Set(jwt.SubjectKey, "uid:legacy-user")
	token.Set("user_id", "user-123")
	token.Set("api_key_id", "api-key-456")
	token.Set("feature_configs", map[string]interface{}{
		"SUBSTREAMS_PARALLEL_JOBS": "5",
		"INDEXER_IDENTIFIER":       "indexer-1",
	})

	ctx := context.Background()
	newCtx := auth.addClaimsToContext(ctx, token, "10.0.0.1")

	// Get headers from context
	headers := dauth.FromContext(newCtx)

	// Check standard headers
	assert.Equal(t, "user-123", headers[dauth.SFHeaderUserID])
	assert.Equal(t, "api-key-456", headers[dauth.SFHeaderApiKeyID])
	assert.Equal(t, "10.0.0.1", headers[dauth.SFHeaderIP])

	// Check feature configs
	assert.Equal(t, "5", headers["x-sf-substreams-parallel-jobs"])
	assert.Equal(t, "indexer-1", headers["x-sf-indexer-identifier"])
}

func TestAuthenticator_AddClaimsToContext_LegacyUserID(t *testing.T) {
	auth := &authenticator{
		logger: zap.NewNop(),
	}

	// Test with legacy uid: prefix in subject
	token := jwt.New()
	token.Set(jwt.SubjectKey, "uid:legacy-user-789")
	token.Set("api_key_id", "api-key-456")

	ctx := context.Background()
	newCtx := auth.addClaimsToContext(ctx, token, "10.0.0.1")

	// Get headers from context
	headers := dauth.FromContext(newCtx)

	// Should extract user ID from subject
	assert.Equal(t, "legacy-user-789", headers[dauth.SFHeaderUserID])
}
