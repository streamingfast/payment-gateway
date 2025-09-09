package session

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var defaultMinimalWorkerLifeDuration = time.Second * 5

func TestConfig_new(t *testing.T) {
	tests := []struct {
		name        string
		dsn         string
		expect      *Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "default config with empty DSN",
			dsn:  "",
			expect: &Config{
				Endpoint:                  "session.thegraph.market",
				Insecure:                  false,
				Plaintext:                 false,
				RequestKeepAliveDelay:     20 * time.Second,
				DefaultMaxRequestPerUser:  10,
				IndexerApiKey:             "",
				MinimalWorkerLifeDuration: defaultMinimalWorkerLifeDuration,
			},
		},
		{
			name: "default config with tgm:// only",
			dsn:  "tgm://",
			expect: &Config{
				Endpoint:                  "session.thegraph.market",
				Insecure:                  false,
				Plaintext:                 false,
				RequestKeepAliveDelay:     20 * time.Second,
				DefaultMaxRequestPerUser:  10,
				IndexerApiKey:             "",
				MinimalWorkerLifeDuration: defaultMinimalWorkerLifeDuration,
			},
		},
		{
			name: "custom endpoint default port for plaintext",
			dsn:  "tgm://custom.thegraph.market?plaintext=true",
			expect: &Config{
				Endpoint:                  "custom.thegraph.market:80",
				Insecure:                  false,
				Plaintext:                 true,
				RequestKeepAliveDelay:     20 * time.Second,
				DefaultMaxRequestPerUser:  10,
				IndexerApiKey:             "",
				MinimalWorkerLifeDuration: defaultMinimalWorkerLifeDuration,
			},
		},
		{
			name: "custom endpoint with port",
			dsn:  "tgm://session.example.com:8080",
			expect: &Config{
				Endpoint:                  "session.example.com:8080",
				Insecure:                  false,
				Plaintext:                 false,
				RequestKeepAliveDelay:     20 * time.Second,
				DefaultMaxRequestPerUser:  10,
				IndexerApiKey:             "",
				MinimalWorkerLifeDuration: defaultMinimalWorkerLifeDuration,
			},
		},
		{
			name: "insecure and plaintext flags",
			dsn:  "tgm://localhost:8080?insecure=true&plaintext=true",
			expect: &Config{
				Endpoint:                  "localhost:8080",
				Insecure:                  true,
				Plaintext:                 true,
				RequestKeepAliveDelay:     20 * time.Second,
				DefaultMaxRequestPerUser:  10,
				IndexerApiKey:             "",
				MinimalWorkerLifeDuration: defaultMinimalWorkerLifeDuration,
			},
		},
		{
			name: "custom request keep alive delay",
			dsn:  "tgm://session.thegraph.market?request-keep-alive-delay=45s",
			expect: &Config{
				Endpoint:                  "session.thegraph.market:443",
				Insecure:                  false,
				Plaintext:                 false,
				RequestKeepAliveDelay:     45 * time.Second,
				DefaultMaxRequestPerUser:  10,
				IndexerApiKey:             "",
				MinimalWorkerLifeDuration: defaultMinimalWorkerLifeDuration,
			},
		},
		{
			name: "custom max requests per user",
			dsn:  "tgm://session.thegraph.market?default-max-request-per-user=20",
			expect: &Config{
				Endpoint:                  "session.thegraph.market:443",
				Insecure:                  false,
				Plaintext:                 false,
				RequestKeepAliveDelay:     20 * time.Second,
				DefaultMaxRequestPerUser:  20,
				IndexerApiKey:             "",
				MinimalWorkerLifeDuration: defaultMinimalWorkerLifeDuration,
			},
		},
		{
			name: "all parameters with indexer API key",
			dsn:  "tgm://custom.session.com:9090?insecure=true&plaintext=true&request-keep-alive-delay=60s&default-max-request-per-user=50&indexer-api-key=server_abcdef123456&minimal-worker-life-duration=10s",
			expect: &Config{
				Endpoint:                  "custom.session.com:9090",
				Insecure:                  true,
				Plaintext:                 true,
				RequestKeepAliveDelay:     60 * time.Second,
				DefaultMaxRequestPerUser:  50,
				IndexerApiKey:             "server_abcdef123456",
				MinimalWorkerLifeDuration: 10 * time.Second,
			},
		},
		{
			name:        "invalid max request per user",
			dsn:         "tgm://session.thegraph.market?default-max-request-per-user=invalid",
			expectError: true,
			errorMsg:    "invalid default-max-request-per-user",
		},
		{
			name: "zero max requests per user",
			dsn:  "tgm://session.thegraph.market?default-max-request-per-user=0",
			expect: &Config{
				Endpoint:                  "session.thegraph.market:443",
				Insecure:                  false,
				Plaintext:                 false,
				RequestKeepAliveDelay:     20 * time.Second,
				DefaultMaxRequestPerUser:  0,
				IndexerApiKey:             "",
				MinimalWorkerLifeDuration: defaultMinimalWorkerLifeDuration,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, err := newConfig(test.dsn)
			if test.expectError {
				require.Error(t, err)
				if test.errorMsg != "" {
					assert.Contains(t, err.Error(), test.errorMsg)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expect, c)
			}
		})
	}
}

func TestConfig_DefaultValues(t *testing.T) {
	c, err := newConfig("")
	require.NoError(t, err)

	assert.Equal(t, "session.thegraph.market", c.Endpoint)
	assert.False(t, c.Insecure)
	assert.False(t, c.Plaintext)
	assert.Equal(t, 20*time.Second, c.RequestKeepAliveDelay)
	assert.Equal(t, uint64(10), c.DefaultMaxRequestPerUser)
	assert.Equal(t, "", c.IndexerApiKey)
	assert.Equal(t, defaultMinimalWorkerLifeDuration, c.MinimalWorkerLifeDuration)
}

func TestConfig_BooleanParameters(t *testing.T) {
	tests := []struct {
		name      string
		dsn       string
		insecure  bool
		plaintext bool
	}{
		{
			name:      "insecure true",
			dsn:       "tgm://session.thegraph.market?insecure=true",
			insecure:  true,
			plaintext: false,
		},
		{
			name:      "plaintext true",
			dsn:       "tgm://session.thegraph.market?plaintext=true",
			insecure:  false,
			plaintext: true,
		},
		{
			name:      "both true",
			dsn:       "tgm://session.thegraph.market?insecure=true&plaintext=true",
			insecure:  true,
			plaintext: true,
		},
		{
			name:      "insecure false explicitly",
			dsn:       "tgm://session.thegraph.market?insecure=false",
			insecure:  false,
			plaintext: false,
		},
		{
			name:      "plaintext false explicitly",
			dsn:       "tgm://session.thegraph.market?plaintext=false",
			insecure:  false,
			plaintext: false,
		},
		{
			name:      "empty values default to false",
			dsn:       "tgm://session.thegraph.market?insecure=&plaintext=",
			insecure:  false,
			plaintext: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, err := newConfig(test.dsn)
			require.NoError(t, err)
			assert.Equal(t, test.insecure, c.Insecure)
			assert.Equal(t, test.plaintext, c.Plaintext)
		})
	}
}
