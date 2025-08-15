package auth

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/lestrrat-go/jwx/jwk"
)

func jwkSetFromBase64(inBase64 string) (jwk.Set, error) {
	decoded, err := base64.StdEncoding.DecodeString(inBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	return jwk.Parse(decoded)
}

func jwkSetFromURL(jwksURL string) (jwk.Set, error) {
	ctx := context.TODO()
	return jwk.Fetch(ctx, jwksURL)
}
