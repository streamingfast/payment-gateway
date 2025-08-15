package auth

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWKSetFromURL(t *testing.T) {
	tests := []struct {
		name        string
		jwksURL     string
		expectError bool
		minKeys     int
	}{
		{
			name:        "valid Google JWKS URL",
			jwksURL:     "https://storage.googleapis.com/dfuseio.com/.well-known/jwks.json",
			expectError: false,
			minKeys:     1,
		},
		{
			name:        "invalid URL",
			jwksURL:     "https://invalid-url-that-does-not-exist.example.com/jwks.json",
			expectError: true,
			minKeys:     0,
		},
		{
			name:        "empty URL",
			jwksURL:     "",
			expectError: true,
			minKeys:     0,
		},
		{
			name:        "malformed URL",
			jwksURL:     "not-a-valid-url",
			expectError: true,
			minKeys:     0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			set, err := jwkSetFromURL(test.jwksURL)

			if test.expectError {
				require.Error(t, err)
				assert.Nil(t, set)
			} else {
				require.NoError(t, err)
				require.NotNil(t, set)

				// Check that we have at least the expected number of keys
				assert.GreaterOrEqual(t, set.Len(), test.minKeys)

				// Verify that we can iterate over the keys
				for iter := set.Iterate(context.Background()); iter.Next(context.Background()); {
					pair := iter.Pair()
					assert.NotNil(t, pair.Value)
				}
			}
		})
	}
}

func TestJWKSetFromBase64(t *testing.T) {
	tests := []struct {
		name        string
		base64Input string
		expectError bool
		expectedLen int
	}{
		{
			name:        "valid EC P-256 JWK Set",
			base64Input: "ewogICJrZXlzIjogWwogICAgewogICAgICAiY3J2IjogIlAtMjU2IiwKICAgICAgImt0eSI6ICJFQyIsCiAgICAgICJ1c2UiOiAic2lnIiwKICAgICAgIngiOiAiei1WOG0xM0JuYVhuQ0lBeE5Tc0hwS01wMVVuMTJTdk92YTl0TS15Y3RYYyIsCiAgICAgICJ5IjogImw4VFJGdkYzbjlJQnUzbDQ2dHVUelZJUkRBcXBxWG9tWXI4R2w3Nk14bTgiCiAgICB9CiAgXQp9Cg==",
			expectError: false,
			expectedLen: 1,
		},
		{
			name:        "invalid base64",
			base64Input: "not-valid-base64!@#$",
			expectError: true,
			expectedLen: 0,
		},
		{
			name:        "valid base64 but invalid JSON",
			base64Input: "dGhpcyBpcyBub3QgdmFsaWQgSlNPTg==", // "this is not valid JSON" in base64
			expectError: true,
			expectedLen: 0,
		},
		{
			name:        "empty string",
			base64Input: "",
			expectError: true,
			expectedLen: 0,
		},
		{
			name:        "valid base64 of empty JWK set",
			base64Input: "eyJrZXlzIjpbXX0=", // {"keys":[]} in base64
			expectError: false,
			expectedLen: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			set, err := jwkSetFromBase64(test.base64Input)

			if test.expectError {
				require.Error(t, err)
				assert.Nil(t, set)
			} else {
				require.NoError(t, err)
				require.NotNil(t, set)
				assert.Equal(t, test.expectedLen, set.Len())

				// For the valid EC key test, verify key properties
				if test.name == "valid EC P-256 JWK Set" && set.Len() > 0 {
					// Get the first key
					key, exists := set.Get(0)
					require.True(t, exists)
					require.NotNil(t, key)

					// Verify key type
					assert.Equal(t, "EC", string(key.KeyType()))

					// Verify key usage
					usage, ok := key.Get("use")
					if ok {
						assert.Equal(t, "sig", usage)
					}

					// Verify curve
					crv, ok := key.Get("crv")
					if ok {
						// The crv field might be returned as a typed value
						assert.Equal(t, "P-256", string(crv.(fmt.Stringer).String()))
					}
				}
			}
		})
	}
}

func TestJWKSetFromBase64_DecodedContent(t *testing.T) {
	// Test that the decoded content matches what we expect
	base64Input := "ewogICJrZXlzIjogWwogICAgewogICAgICAiY3J2IjogIlAtMjU2IiwKICAgICAgImt0eSI6ICJFQyIsCiAgICAgICJ1c2UiOiAic2lnIiwKICAgICAgIngiOiAiei1WOG0xM0JuYVhuQ0lBeE5Tc0hwS01wMVVuMTJTdk92YTl0TS15Y3RYYyIsCiAgICAgICJ5IjogImw4VFJGdkYzbjlJQnUzbDQ2dHVUelZJUkRBcXBxWG9tWXI4R2w3Nk14bTgiCiAgICB9CiAgXQp9Cg=="

	set, err := jwkSetFromBase64(base64Input)
	require.NoError(t, err)
	require.NotNil(t, set)

	// Should have exactly one key
	assert.Equal(t, 1, set.Len())

	// Get the key and verify its properties
	key, exists := set.Get(0)
	require.True(t, exists)

	// Verify all the expected fields
	assert.Equal(t, "EC", string(key.KeyType()))

	x, ok := key.Get("x")
	assert.True(t, ok)
	// x and y are returned as byte arrays, not base64 strings
	// Just verify they exist and are not nil
	assert.NotNil(t, x)

	y, ok := key.Get("y")
	assert.True(t, ok)
	assert.NotNil(t, y)

	crv, ok := key.Get("crv")
	assert.True(t, ok)
	// The crv field might be returned as a typed value
	assert.Equal(t, "P-256", string(crv.(fmt.Stringer).String()))

	use, ok := key.Get("use")
	assert.True(t, ok)
	assert.Equal(t, "sig", use)
}

func TestJWKSetRoundTrip(t *testing.T) {
	// Fetch from URL
	urlSet, err := jwkSetFromURL("https://storage.googleapis.com/dfuseio.com/.well-known/jwks.json")
	if err != nil {
		t.Skip("Skipping round-trip test: could not fetch from URL (network might be unavailable)")
	}

	// Parse from base64
	base64Set, err := jwkSetFromBase64("ewogICJrZXlzIjogWwogICAgewogICAgICAiY3J2IjogIlAtMjU2IiwKICAgICAgImt0eSI6ICJFQyIsCiAgICAgICJ1c2UiOiAic2lnIiwKICAgICAgIngiOiAiei1WOG0xM0JuYVhuQ0lBeE5Tc0hwS01wMVVuMTJTdk92YTl0TS15Y3RYYyIsCiAgICAgICJ5IjogImw4VFJGdkYzbjlJQnUzbDQ2dHVUelZJUkRBcXBxWG9tWXI4R2w3Nk14bTgiCiAgICB9CiAgXQp9Cg==")
	require.NoError(t, err)

	// Both should be valid sets
	assert.NotNil(t, urlSet)
	assert.NotNil(t, base64Set)

	// Both should have at least one key
	assert.GreaterOrEqual(t, urlSet.Len(), 1)
	assert.Equal(t, 1, base64Set.Len())
}
