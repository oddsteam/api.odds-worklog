package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// KeycloakClaims represents the standard claims in a Keycloak JWT token
type KeycloakClaims struct {
	jwt.RegisteredClaims
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
	RealmAccess       struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	ResourceAccess map[string]struct {
		Roles []string `json:"roles"`
	} `json:"resource_access"`
	AuthorizedParty string `json:"azp"`
}

// VerifyAudience implements the jwt.Claims interface
func (c *KeycloakClaims) VerifyAudience(claim string, required bool) bool {
	if !required {
		return true
	}
	// Check both standard audience and authorized party
	if c.AuthorizedParty == claim {
		return true
	}
	for _, aud := range c.Audience {
		if aud == claim {
			return true
		}
	}
	return false
}

// KeycloakValidator handles Keycloak token validation
type KeycloakValidator struct {
	realmURL     string
	realm        string
	clientID     string
	publicKey    *rsa.PublicKey
	issuer       string
	lastFetch    time.Time
	cacheTimeout time.Duration
}

// NewKeycloakValidator creates a new Keycloak validator
func NewKeycloakValidator(realmURL, realm, clientID string) *KeycloakValidator {
	return &KeycloakValidator{
		realmURL:     realmURL,
		realm:        realm,
		clientID:     clientID,
		cacheTimeout: 24 * time.Hour,
	}
}

// fetchPublicKey retrieves the public key from Keycloak
func (v *KeycloakValidator) fetchPublicKey() error {
	// Check if we need to refresh the public key
	if v.publicKey != nil && time.Since(v.lastFetch) < v.cacheTimeout {
		return nil
	}

	url := fmt.Sprintf("%s/realms/%s", v.realmURL, v.realm)
	fmt.Printf("Fetching public key from: %s\n", url)

	// Fetch the realm configuration
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch realm configuration: %v", err)
	}
	defer resp.Body.Close()

	var config struct {
		PublicKey    string `json:"public_key"`
		TokenService string `json:"token-service"`
		Issuer       string `json:"issuer"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return fmt.Errorf("failed to decode realm configuration: %v", err)
	}

	fmt.Printf("Received token service: %s\n", config.TokenService)
	fmt.Printf("Received public key: %s\n", config.PublicKey)

	// Parse the public key
	block, _ := pem.Decode([]byte("-----BEGIN PUBLIC KEY-----\n" + config.PublicKey + "\n-----END PUBLIC KEY-----"))
	if block == nil {
		return fmt.Errorf("failed to decode public key PEM")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %v", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("public key is not RSA")
	}

	v.publicKey = rsaPublicKey
	// Set the issuer to the token service URL without the protocol/openid-connect part
	v.issuer = strings.TrimSuffix(config.TokenService, "/protocol/openid-connect")
	v.lastFetch = time.Now()

	fmt.Printf("Set issuer to: %s\n", v.issuer)

	return nil
}

// ValidateToken validates a Keycloak JWT token
func (v *KeycloakValidator) ValidateToken(tokenString string) (*KeycloakClaims, error) {
	// Fetch the public key if needed
	if err := v.fetchPublicKey(); err != nil {
		return nil, err
	}

	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &KeycloakClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Print token header for debugging
		fmt.Printf("Token header: %+v\n", token.Header)

		return v.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	// Type assert the claims
	claims, ok := token.Claims.(*KeycloakClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Print claims for debugging
	fmt.Printf("Token claims: %+v\n", claims)

	// Validate the issuer
	if claims.Issuer != v.issuer {
		return nil, fmt.Errorf("invalid issuer: got %s, want %s", claims.Issuer, v.issuer)
	}

	// Validate the audience
	if !claims.VerifyAudience(v.clientID, true) {
		return nil, fmt.Errorf("invalid audience: got %v, want %s", claims.Audience, v.clientID)
	}

	// Validate the expiration
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token has expired")
	}

	return claims, nil
}

// ExtractTokenFromHeader extracts the token from the Authorization header
func ExtractTokenFromHeader(header string) (string, error) {
	if header == "" {
		return "", fmt.Errorf("authorization header is empty")
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("authorization header format must be Bearer {token}")
	}

	return parts[1], nil
}
