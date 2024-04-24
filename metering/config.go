package metering

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type Config struct {
	Endpoint    string
	Insecure    bool
	Plaintext   bool
	Token       string
	Delay       time.Duration
	BufferSize  uint64
	PanicOnDrop bool
	Network     string
}

func newConfig(configURL string) (*Config, error) {
	c := &Config{
		Delay:       100 * time.Millisecond,
		BufferSize:  10000,
		PanicOnDrop: false,
		Insecure:    false,
		Plaintext:   false,
	}

	u, err := url.Parse(configURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse urls: %w", err)
	}

	hostname := u.Hostname()
	port := u.Port()
	if port == "" {
		port = "443"
	}

	if hostname == "" {
		return nil, fmt.Errorf("hostname not specified in %q", configURL)
	}

	c.Endpoint = fmt.Sprintf("%s:%s", hostname, port)

	vals := u.Query()
	if vals.Get("insecure") == "true" {
		c.Insecure = true
	}

	if vals.Get("plaintext") == "true" {
		c.Plaintext = true
	}

	c.Token = vals.Get("token")

	c.Network = vals.Get("network")
	if c.Network == "" {
		return nil, fmt.Errorf("network not specified (as query param)")
	}

	bufferValue := vals.Get("buffer")
	if bufferValue != "" {
		c.BufferSize, err = strconv.ParseUint(bufferValue, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid buffer value %q: %w", bufferValue, err)
		}
	}

	delayValue := vals.Get("delay")
	if delayValue != "" {
		delay, err := strconv.ParseInt(delayValue, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid delay value %q: %w", delayValue, err)
		}

		c.Delay = time.Duration(delay) * time.Millisecond
	}

	c.PanicOnDrop = vals.Get("panicOnDrop") == "true"

	return c, nil
}
