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

	Key               string // key to authenticate when asking for 'reissue', used to prevent being rate-limited
	ReissueJWTAgeSecs uint64 // max age of JWT in seconds to trust its 'claims', after which we will ask for a reissue
	PubKeyURL         string // URL to fetch public keys for JWT verification
	PubKeyBase64      string // Base64-encoded public key for JWT verification
	ReissueURL        string // URL to fetch an updated JWT from a previous one
	IssueURL          string // URL to issue a new JWT from an API key
}

func newConfig(configURL string) (*Config, error) {
	c := &Config{
		Endpoint:          "auth.thegraph.market",
		PubKeyURL:         "https://auth.thegraph.market/.well-known/jwks.json",
		ReissueJWTAgeSecs: 600,
		Insecure:          false,
		Plaintext:         false,
	}

	u, err := url.Parse(configURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse urls: %w", err)
	}

	protocol := u.Scheme
	if protocol != "paymentgateway" {
		return nil, fmt.Errorf("invalid protocol: %s", protocol)
	}

	hostname := u.Hostname()
	port := u.Port()

	if hostname != "" {
		// Only include port if it's non-standard
		if port != "" && port != "443" && port != "80" {
			c.Endpoint = fmt.Sprintf("%s:%s", hostname, port)
		} else {
			c.Endpoint = hostname
		}
	}

	vals := u.Query()
	if vals.Get("insecure") == "true" {
		c.Insecure = true
	}

	if vals.Get("plaintext") == "true" {
		c.Plaintext = true
	}

	c.Key = vals.Get("key")

	keyURL := vals.Get("pubkeyurl")
	keyBase64 := vals.Get("pubkeybase64")

	if keyURL != "" {
		c.PubKeyURL = keyURL
		if keyBase64 != "" {
			return nil, fmt.Errorf("only one of pubkeyurl or pubkeybase64 can be provided, not both")
		}
	}
	if keyBase64 != "" {
		c.PubKeyBase64 = keyBase64
	}

	// Parse reissue-jwt-max-age-secs if provided
	if reissueAge := vals.Get("reissue-jwt-max-age-secs"); reissueAge != "" {
		age, err := strconv.ParseUint(reissueAge, 10, 64)
		if err == nil {
			c.ReissueJWTAgeSecs = age
		}
	}

	if c.Plaintext {
		c.ReissueURL = fmt.Sprintf("http://%s/v1/auth/reissue", c.Endpoint)
		c.IssueURL = fmt.Sprintf("http://%s/v1/auth/issue", c.Endpoint)
	} else {
		c.ReissueURL = fmt.Sprintf("https://%s/v1/auth/reissue", c.Endpoint)
		c.IssueURL = fmt.Sprintf("https://%s/v1/auth/issue", c.Endpoint)
	}

	return c, nil
}
