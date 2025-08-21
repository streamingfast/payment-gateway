package auth

import (
	"fmt"
	"net/url"
	"strconv"
)

type Config struct {
	Endpoint  string // auth endpoint, ex: auth.thegraph.market
	Insecure  bool   // skip certificate verification on endpoint
	Plaintext bool   // skip encryption on endpoint

	IndexerAPIKey     string // Authentication key used to prevent rate limiting when calling /issue or /reissue endpoints on behalf of the user
	ReissueJWTAgeSecs uint64 // max age of JWT in seconds to trust its 'claims', after which we will ask for a reissue
	PubKeyURL         string // URL to fetch public keys for JWT verification
	PubKeyBase64      string // Base64-encoded public key for JWT verification
	ReissueURL        string // URL to fetch an updated JWT from a previous one
	IssueURL          string // URL to issue a new JWT from an API key
}

func newConfig(configURL string) (*Config, error) {
	c := &Config{
		Endpoint:          "",
		PubKeyURL:         "",
		ReissueJWTAgeSecs: 600,
		Insecure:          false,
		Plaintext:         false,
	}

	u, err := url.Parse(configURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse urls: %w", err)
	}

	protocol := u.Scheme
	if protocol != "paymentgateway" && protocol != "tgm" {
		return nil, fmt.Errorf("invalid protocol: %s", protocol)
	}

	hostname := u.Hostname()
	port := u.Port()

	if hostname == "" {
		return nil, fmt.Errorf("you must provide a hostname, ex: auth.thegraph.market")
	}
	// Only include port if it's non-standard
	if port != "" && port != "443" && port != "80" {
		c.Endpoint = fmt.Sprintf("%s:%s", hostname, port)
	} else {
		c.Endpoint = hostname
	}

	vals := u.Query()
	if vals.Get("insecure") == "true" {
		c.Insecure = true
	}

	if vals.Get("plaintext") == "true" {
		c.Plaintext = true
	}

	c.IndexerAPIKey = vals.Get("indexer-api-key")

	keyURL := vals.Get("pub-key-url")
	keyBase64 := vals.Get("pub-key-base64")

	scheme := "https"
	if c.Plaintext {
		scheme = "http"
	}

	switch {
	case keyURL != "":
		c.PubKeyURL = keyURL
		if keyBase64 != "" {
			return nil, fmt.Errorf("only one of pub-key-url or pub-key-base64 can be provided, not both")
		}
	case keyBase64 != "":
		c.PubKeyBase64 = keyBase64
		c.PubKeyURL = "" // remove default value here
	default:
		c.PubKeyURL = fmt.Sprintf("%s://%s/.well-known/jwks.json", scheme, c.Endpoint)
	}

	// Parse reissue-jwt-max-age-secs if provided
	if reissueAge := vals.Get("reissue-jwt-max-age-secs"); reissueAge != "" {
		age, err := strconv.ParseUint(reissueAge, 10, 64)
		if err == nil {
			c.ReissueJWTAgeSecs = age
		}
	}

	c.ReissueURL = fmt.Sprintf("%s://%s/v1/auth/reissue", scheme, c.Endpoint)
	c.IssueURL = fmt.Sprintf("%s://%s/v1/auth/issue", scheme, c.Endpoint)

	return c, nil
}
