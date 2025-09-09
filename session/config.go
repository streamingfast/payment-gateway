package session

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type Config struct {
	Endpoint                 string        // session endpoint, ex: session.thegraph.market
	Insecure                 bool          // skip certificate verification on endpoint
	Plaintext                bool          // skip encryption on endpoint
	RequestKeepAliveDelay    time.Duration // delay between keep alive requests
	DefaultMaxRequestPerUser uint64        // default maximum requests per user
	IndexerApiKey            string        // indexer API key for X-Api-Key header
}

func newConfig(configURL string) (*Config, error) {
	c := &Config{
		Endpoint:                 "session.thegraph.market",
		Insecure:                 false,
		Plaintext:                false,
		RequestKeepAliveDelay:    20 * time.Second,
		DefaultMaxRequestPerUser: 10,
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

	hostname := u.Hostname()
	port := u.Port()

	if hostname != "" {
		c.Endpoint = fmt.Sprintf("%s:%s", hostname, port)
	}

	vals := u.Query()
	for k := range vals {
		if k == "insecure" || k == "plaintext" || k == "request-keep-alive-delay" ||
			k == "default-max-request-per-user" || k == "indexer-api-key" {
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

	// Parse request-keep-alive-delay if provided
	if keepAliveDelay := vals.Get("request-keep-alive-delay"); keepAliveDelay != "" {
		delay, err := time.ParseDuration(keepAliveDelay)
		if err != nil {
			return nil, fmt.Errorf("invalid request-keep-alive-delay: %w", err)
		}
		c.RequestKeepAliveDelay = delay
	}

	// Parse default-max-request-per-user if provided
	if maxRequests := vals.Get("default-max-request-per-user"); maxRequests != "" {
		maxReq, err := strconv.ParseUint(maxRequests, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid default-max-request-per-user: %w", err)
		}
		c.DefaultMaxRequestPerUser = maxReq
	}

	// Parse indexer-api-key if provided
	if apiKey := vals.Get("indexer-api-key"); apiKey != "" {
		c.IndexerApiKey = apiKey
	}

	return c, nil
}
