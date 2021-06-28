// Package ports includes common concerns across multiple ports implementations.
package ports

import "github.com/dgrijalva/jwt-go"

// JWTKey is a secret used to sign the JWT Token.
type JWTKey = []byte

// JWTConfig represents the information necessary to sign/verify JWT Tokens.
// Note: Variants of HMAC signing are the only supported methods currently.
type JWTConfig struct {
	Method *jwt.SigningMethodHMAC
	Key    JWTKey
}

// DefaultJWTConfig creates a new JWTConfig with the HS256 signing method and the provide key.
func DefaultJWTConfig(key string) JWTConfig {
	return JWTConfig{
		jwt.SigningMethodHS256,
		[]byte(key),
	}
}
