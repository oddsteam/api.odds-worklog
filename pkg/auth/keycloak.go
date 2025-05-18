package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
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
}

// VerifyAudience implements the jwt.Claims interface
func (c *KeycloakClaims) VerifyAudience(claim string, required bool) bool {
	if !required {
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
	clientID     string
	publicKey    *rsa.PublicKey
	issuer       string
	lastFetch    time.Time
	cacheTimeout time.Duration
}

// NewKeycloakValidator creates a new Keycloak validator
func NewKeycloakValidator(realmURL, clientID string) *KeycloakValidator {
	return &KeycloakValidator{
		realmURL:     realmURL,
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

	// Fetch the realm configuration
	resp, err := http.Get(fmt.Sprintf("%s/realms/%s", v.realmURL, "master"))
	if err != nil {
		return fmt.Errorf("failed to fetch realm configuration: %v", err)
	}
	defer resp.Body.Close()

	var config struct {
		PublicKey string `json:"public_key"`
		Issuer    string `json:"issuer"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return fmt.Errorf("failed to decode realm configuration: %v", err)
	}

	// Decode the public key
	decodedKey, err := base64.StdEncoding.DecodeString(config.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to decode public key: %v", err)
	}

	// Parse the public key
	modulus := new(big.Int)
	modulus.SetBytes(decodedKey)
	v.publicKey = &rsa.PublicKey{
		N: modulus,
		E: 65537, // Standard RSA public exponent
	}
	v.issuer = config.Issuer
	v.lastFetch = time.Now()

	fmt.Printf("v.publicKey %#v\n", v.publicKey)
	fmt.Printf("v.issuer %#v\n", v.issuer)

	return nil
}

// ValidateToken validates a Keycloak JWT token
func (v *KeycloakValidator) ValidateToken(tokenString string) (*KeycloakClaims, error) {
	// Fetch the public key if needed
	if err := v.fetchPublicKey(); err != nil {
		return nil, err
	}

	// Parse and validate the token
	token, _ := jwt.ParseWithClaims(tokenString, &KeycloakClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return v.publicKey, nil
	})

	// if err != nil {
	// 	return nil, fmt.Errorf("failed to parse token: %v", err)
	// }

	// Type assert the claims
	claims, ok := token.Claims.(*KeycloakClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Validate the issuer
	// if claims.Issuer != v.issuer {
	// return nil, fmt.Errorf("invalid issuer")
	// }

	// Validate the audience
	// if !claims.VerifyAudience(v.clientID, true) {
	// 	return nil, fmt.Errorf("invalid audience")
	// }

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
