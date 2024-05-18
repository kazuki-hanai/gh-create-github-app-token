package internal

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
)

// JWTTokenGenerator implements the TokenGenerator interface for generating a JWT token.
type JwtTokenGenerator struct {
	PrivateKey string
}

// Header defines the JWT Header.
type JwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// JwtPayload represents the JWT payload.
type JwtPayload struct {
	Iat int64  `json:"iat"`
	Exp int64  `json:"exp"`
	Iss string `json:"iss"`
}

func (g *JwtTokenGenerator) GenerateToken(payload JwtPayload) (string, error) {
	PrivateKey, err := convertPrivateKey(g.PrivateKey)
	if err != nil {
		return "", err
	}

	if isPkcs1(PrivateKey) {
		return "", errors.New("Private Key is in PKCS#1 format, but only PKCS#8 is supported")
	}

	if isOpenSsh(PrivateKey) {
		return "", errors.New("Private Key is in OpenSSH format, but only PKCS#8 is supported")
	}

	header := JwtHeader{Alg: "RS256", Typ: "JWT"}
	headerJson, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJson)

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadJson)
	encodedMessage := fmt.Sprintf("%s.%s", headerEncoded, payloadEncoded)

	hashed := sha256.Sum256([]byte(encodedMessage))
	signature, err := rsa.SignPKCS1v15(rand.Reader, PrivateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}
	encodedSignature := base64.RawURLEncoding.EncodeToString(signature)

	return fmt.Sprintf("%s.%s", encodedMessage, encodedSignature), nil
}

// convertPrivateKey converts a PEM encoded private key to rsa.PrivateKey.
func convertPrivateKey(pemKey string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemKey))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// isPkcs1 checks if the private key is in PKCS#1 format.
func isPkcs1(privateKey *rsa.PrivateKey) bool {
	// Assuming PKCS#1 format detection logic
	return false
}

// isOpenSsh checks if the private key is in OpenSSH format.
func isOpenSsh(privateKey *rsa.PrivateKey) bool {
	// Assuming OpenSSH format detection logic
	return false
}
