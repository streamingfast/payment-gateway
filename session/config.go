package session

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type Config struct {
	Endpoint                         string        // session endpoint, ex: session.thegraph.market
	Insecure                         bool          // skip certificate verification on endpoint
	Plaintext                        bool          // skip encryption on endpoint
	RequestKeepAliveDelay            time.Duration // delay between keep alive requests
	DefaultMaxRequestPerOrganization uint64        // default maximum requests per organization
	IndexerApiKey                    string        // indexer API key for X-Api-Key header
	MinimalWorkerLifeDuration        time.Duration // minimal worker life duration
}

func newConfig(configURL string) (*Config, error) {
	c := &Config{
		Endpoint:                         "session.thegraph.market",
		Insecure:                         false,
		Plaintext:                        false,
		RequestKeepAliveDelay:            20 * time.Second,
		DefaultMaxRequestPerOrganization: 10,
		MinimalWorkerLifeDuration:        time.Second * 5,
	}

	if configURL == "" {
		return c, nil
	}

	u, err := url.Parse(configURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	protocol := u.Scheme
	if protocol != "" && protocol != "tgm" {
		return nil, fmt.Errorf("invalid protocol: %s, expected 'tgm'", protocol)
	}

	vals := u.Query()
	for k := range vals {
		if k == "insecure" || k == "plaintext" || k == "request-keep-alive-delay" ||
			k == "default-max-request-per-user" || k == "default-max-request-per-organization" ||
			k == "indexer-api-key" || k == "minimal-worker-life-duration" {
			continue
		}
		return nil, fmt.Errorf("unknown query parameter: %s", k)
	}

	if vals.Get("insecure") == "true" {
		c.Insecure = true
	}

	if vals.Get("plaintext") == "true" {
		c.Plaintext = true
	}

	hostname := u.Hostname()
	port := u.Port()
	if port == "" {
		port = "443"
		if c.Plaintext {
			port = "80"
		}
	}

	if hostname != "" {
		c.Endpoint = fmt.Sprintf("%s:%s", hostname, port)
	}

	// Parse request-keep-alive-delay if provided
	if keepAliveDelay := vals.Get("request-keep-alive-delay"); keepAliveDelay != "" {
		delay, err := time.ParseDuration(keepAliveDelay)
		if err != nil {
			return nil, fmt.Errorf("invalid request-keep-alive-delay: %w", err)
		}
		c.RequestKeepAliveDelay = delay
	}

	if minimalWorkerLifeDuration := vals.Get("minimal-worker-life-duration"); minimalWorkerLifeDuration != "" {
		duration, err := time.ParseDuration(minimalWorkerLifeDuration)
		if err != nil {
			return nil, fmt.Errorf("invalid minimal-worker-life-duration: %w", err)
		}
		c.MinimalWorkerLifeDuration = duration
	}

	// Parse default-max-request-per-organization if provided
	maxRequests := vals.Get("default-max-request-per-organization")
	if maxRequests == "" {
		// Check for legacy parameter name
		maxRequests = vals.Get("default-max-request-per-user")
	}

	if maxRequests != "" {
		maxReq, err := strconv.ParseUint(maxRequests, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid default-max-request-per-organization: %w", err)
		}
		c.DefaultMaxRequestPerOrganization = maxReq
	}

	// Parse indexer-api-key if provided
	if apiKey := vals.Get("indexer-api-key"); apiKey != "" {
		c.IndexerApiKey = apiKey
	}

	return c, nil
}
