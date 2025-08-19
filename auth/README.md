# Payment Gateway Authenticator

This package provides a JWT-based authentication system that integrates with the StreamingFast dauth framework. It has been migrated from an external library to be self-contained within the payment gateway.

## Features

- **JWT Token Validation**: Validates JWT tokens using JWK sets fetched from a URL or provided as base64-encoded keys
- **API Key Exchange**: Automatically exchanges API keys for JWT tokens via an issue endpoint
- **Token Reissue**: Automatically reissues JWT tokens that are older than a configured threshold
- **Feature Configurations**: Extracts and propagates feature configurations from JWT claims
- **Legacy Support**: Supports legacy JWT formats with user IDs in the subject field

## Configuration

The authenticator is configured using a URL-style connection string with the following format:

```
paymentgateway://[host[:port]]?[parameters]
```

### Parameters

- `pubkeyurl`: URL to fetch the JWK set for JWT verification (e.g., `https://auth.example.com/.well-known/jwks.json`)
- `pubkeybase64`: Base64-encoded JWK set for offline JWT verification (mutually exclusive with `pubkeyurl`)
- `reissue-jwt-max-age-secs`: Maximum age in seconds before a JWT token is reissued (default: 600)
- `key`: Authentication key for the reissue endpoint to prevent rate limiting
- `insecure`: Skip certificate verification (default: false)
- `plaintext`: Use HTTP instead of HTTPS (default: false)

### Example Configurations

**Production:**
```
paymentgateway://auth.thegraph.market?pubkeyurl=https://auth.thegraph.market/.well-known/jwks.json&reissue-jwt-max-age-secs=600&key=server_key
```

**Development:**
```
paymentgateway://localhost:8080?plaintext=true&insecure=true&pubkeyurl=http://localhost:8080/.well-known/jwks.json
```

**Offline Verification:**
```
paymentgateway://auth.example.com?pubkeybase64=LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0K...
```

## Usage

### Registration

Register the authenticator with dauth at application startup:

```go
import "github.com/streamingfast/payment-gateway/auth"

func init() {
    auth.Register()
}
```

### Creating an Authenticator

```go
import (
    "github.com/streamingfast/dauth"
    "go.uber.org/zap"
)

config := "paymentgateway://"
logger := zap.NewLogger()

authenticator, err := dauth.New(config, logger)
if err != nil {
    log.Fatal(err)
}
```

### Authentication

The authenticator supports two authentication methods:

#### 1. JWT Token Authentication

Pass a JWT token via the `Authorization` header:

```go
headers := map[string][]string{
    "Authorization": {"Bearer eyJhbGciOiJIUzI1NiIs..."},
}

ctx, err := authenticator.Authenticate(context.Background(), "/api/endpoint", headers, "192.168.1.1")
if err != nil {
    // Handle authentication failure
}

// Extract trusted headers from context
trustedHeaders := dauth.FromContext(ctx)
userID := trustedHeaders[dauth.SFHeaderUserID]
apiKeyID := trustedHeaders[dauth.SFHeaderApiKeyID]
```

#### 2. API Key Authentication

Pass an API key via the `X-API-Key` header. The authenticator will automatically exchange it for a JWT token:

```go
headers := map[string][]string{
    "X-API-Key": {"server_1234567890abcdef"},
}

ctx, err := authenticator.Authenticate(context.Background(), "/api/endpoint", headers, "192.168.1.1")
```

## Authentication Flow

1. **Header Parsing**: The authenticator checks for either an `Authorization` header with a JWT token or an `X-API-Key` header
2. **API Key Exchange**: If an API key is provided, it's sent to the issue endpoint to obtain a JWT token
3. **JWT Validation**: The JWT token is validated against the configured JWK set
4. **Token Reissue**: If the token is older than the configured threshold, it's automatically reissued (to ensure we have the latest parameters/features in it)
5. **Context Enrichment**: Claims from the JWT are extracted and added to the context as trusted headers

## Trusted Headers

The following headers are automatically extracted from JWT claims and added to the context:

- `x-sf-user-id`: User identifier
- `x-sf-api-key-id`: API key identifier
- `x-real-ip`: Client IP address
- `x-sf-plan-tier`: one of "FREE", "SCALING", "PRO", "ENTERPRISE"
- Feature configuration headers (e.g., `x-sf-substreams-parallel-jobs`)

## Issuance Endpoints

The authenticator communicates with two external endpoints:

### Issue Endpoint

- **URL**: `https://[host]/v1/auth/issue`
- **Method**: POST
- **Content-Type**: `application/json`
- **Request Body**: `{"api_key": "the_api_key"}`
- **Response**: `{"token": "jwt_token_string"}`

Used to exchange an API key for a JWT token.

### Reissue Endpoint

- **URL**: `https://[host]/v1/auth/reissue`
- **Method**: POST
- **Content-Type**: `application/json`
- **Request Body**: `{"jwt": "existing_jwt_token"}`
- **Response**: `{"token": "new_jwt_token_string"}`
- **Optional Header**: `Authorization: Bearer [key]` (if configured)

Used to obtain a new JWT token from an existing one that's approaching expiration.
