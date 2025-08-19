package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/streamingfast/dauth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestEndToEndAuthentication(t *testing.T) {
	// Generate test keys
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	key, err := jwk.Import(&privateKey.PublicKey)
	require.NoError(t, err)
	key.Set(jwk.AlgorithmKey, jwa.ES256())
	key.Set(jwk.KeyIDKey, "test-key-id")

	jwkSet := jwk.NewSet()
	jwkSet.AddKey(key)

	// Create test JWT
	token := jwt.New()
	token.Set(jwt.IssuedAtKey, time.Now())
	token.Set(jwt.ExpirationKey, time.Now().Add(time.Hour))
	token.Set("user_id", "test-user-123")
	token.Set("api_key_id", "api-key-456")
	token.Set("plan_tier", "premium")
	token.Set("feature_configs", map[string]interface{}{
		"max_requests":   "1000",
		"enable_feature": "true",
	})

	signKey, _ := jwk.Import(privateKey)
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.ES256(), signKey))
	require.NoError(t, err)

	// Create mock server for JWK endpoint
	jwkServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jwkSet)
	}))
	defer jwkServer.Close()

	// Create mock server for issue/reissue endpoints
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/auth/issue":
			var body map[string]string
			json.NewDecoder(r.Body).Decode(&body)

			if body["api_key"] == "valid-api-key" {
				response := map[string]string{
					"token": string(signed),
				}
				json.NewEncoder(w).Encode(response)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Invalid API key"))
			}

		case "/v1/auth/reissue":
			authHeader := r.Header.Get("Authorization")
			if authHeader == "Bearer reissue-key" {
				// Create new token with updated issued time
				newToken := jwt.New()
				newToken.Set(jwt.IssuedAtKey, time.Now())
				newToken.Set(jwt.ExpirationKey, time.Now().Add(time.Hour))
				newToken.Set("user_id", "reissued-user")

				newSigned, _ := jwt.Sign(newToken, jwt.WithKey(jwa.ES256(), signKey))
				response := map[string]string{
					"token": string(newSigned),
				}
				json.NewEncoder(w).Encode(response)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
			}

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer authServer.Close()

	// Register and create authenticator
	Register()
	logger := zap.NewNop()

	config := "paymentgateway://" + authServer.Listener.Addr().String() +
		"?pubkeyurl=" + jwkServer.URL +
		"&plaintext=true" +
		"&key=reissue-key" +
		"&reissue-jwt-max-age-secs=300"

	auth, err := dauth.New(config, logger)
	require.NoError(t, err)
	require.NotNil(t, auth)

	ctx := context.Background()

	// Test 1: Authenticate with JWT
	t.Run("JWT authentication", func(t *testing.T) {
		headers := map[string][]string{
			"Authorization": {"Bearer " + string(signed)},
		}

		newCtx, err := auth.Authenticate(ctx, "/test/path", headers, "192.168.1.1")
		assert.NoError(t, err)

		trustedHeaders := dauth.FromContext(newCtx)
		assert.Equal(t, "test-user-123", trustedHeaders[dauth.SFHeaderUserID])
		assert.Equal(t, "api-key-456", trustedHeaders[dauth.SFHeaderApiKeyID])
		assert.Equal(t, "premium", trustedHeaders[dauth.SFHeaderPlanTier])
		assert.Equal(t, "192.168.1.1", trustedHeaders[dauth.SFHeaderIP])
		assert.Equal(t, "1000", trustedHeaders["x-sf-max-requests"])
		assert.Equal(t, "true", trustedHeaders["x-sf-enable-feature"])
	})

	// Test 2: Authenticate with API key
	t.Run("API key authentication", func(t *testing.T) {
		headers := map[string][]string{
			"X-API-Key": {"valid-api-key"},
		}

		newCtx, err := auth.Authenticate(ctx, "/test/path", headers, "10.0.0.1")
		assert.NoError(t, err)

		trustedHeaders := dauth.FromContext(newCtx)
		assert.Equal(t, "test-user-123", trustedHeaders[dauth.SFHeaderUserID])
		assert.Equal(t, "10.0.0.1", trustedHeaders[dauth.SFHeaderIP])
	})

	// Test 3: Invalid API key
	t.Run("Invalid API key", func(t *testing.T) {
		headers := map[string][]string{
			"X-API-Key": {"invalid-api-key"},
		}

		_, err := auth.Authenticate(ctx, "/test/path", headers, "10.0.0.1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to authenticate with API key")
	})

	// Test 4: JWT reissue for old token
	t.Run("JWT reissue", func(t *testing.T) {
		// Create an old token
		oldToken := jwt.New()
		oldToken.Set(jwt.IssuedAtKey, time.Now().Add(-10*time.Minute))
		oldToken.Set(jwt.ExpirationKey, time.Now().Add(time.Hour))
		oldToken.Set("user_id", "old-user")

		oldSigned, err := jwt.Sign(oldToken, jwt.WithKey(jwa.ES256(), signKey))
		require.NoError(t, err)

		headers := map[string][]string{
			"Authorization": {"Bearer " + string(oldSigned)},
		}

		newCtx, err := auth.Authenticate(ctx, "/test/path", headers, "192.168.1.1")
		assert.NoError(t, err)

		trustedHeaders := dauth.FromContext(newCtx)
		// Should have the reissued user ID
		assert.Equal(t, "reissued-user", trustedHeaders[dauth.SFHeaderUserID])
	})

	// Test 5: Missing authentication
	t.Run("Missing authentication", func(t *testing.T) {
		headers := map[string][]string{}

		_, err := auth.Authenticate(ctx, "/test/path", headers, "192.168.1.1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "required authorization token not found")
	})

	// Test 6: Invalid JWT
	t.Run("Invalid JWT", func(t *testing.T) {
		headers := map[string][]string{
			"Authorization": {"Bearer invalid.jwt.token"},
		}

		_, err := auth.Authenticate(ctx, "/test/path", headers, "192.168.1.1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid JWT token")
	})

	// Test 7: Case-insensitive headers
	t.Run("Case-insensitive headers", func(t *testing.T) {
		headers := map[string][]string{
			"AUTHORIZATION": {"Bearer " + string(signed)},
		}

		newCtx, err := auth.Authenticate(ctx, "/test/path", headers, "192.168.1.1")
		assert.NoError(t, err)

		trustedHeaders := dauth.FromContext(newCtx)
		assert.Equal(t, "test-user-123", trustedHeaders[dauth.SFHeaderUserID])
	})

	// Test 8: Expired JWT
	t.Run("Expired JWT", func(t *testing.T) {
		// Create an expired token
		expiredToken := jwt.New()
		expiredToken.Set(jwt.IssuedAtKey, time.Now().Add(-2*time.Hour))
		expiredToken.Set(jwt.ExpirationKey, time.Now().Add(-1*time.Hour))
		expiredToken.Set("user_id", "expired-user")

		expiredSigned, err := jwt.Sign(expiredToken, jwt.WithKey(jwa.ES256(), signKey))
		require.NoError(t, err)

		headers := map[string][]string{
			"Authorization": {"Bearer " + string(expiredSigned)},
		}

		_, err = auth.Authenticate(ctx, "/test/path", headers, "192.168.1.1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token is expired")
	})
}
