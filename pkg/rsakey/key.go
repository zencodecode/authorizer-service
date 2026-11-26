package rsakey

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"regexp"
)

var (
	pemBeginRegex = regexp.MustCompile(`-----BEGIN [^-]+-----`)
	pemEndRegex   = regexp.MustCompile(`-----END [^-]+-----`)
)

func cleanKeyBytes(keyBytes []byte) []byte {
	// Strip surrounding whitespace and quotes (", ', `)
	keyBytes = bytes.Trim(keyBytes, "\"'` \t\r\n")

	// Replace escaped newlines (both \\n and \n) and carriage returns
	keyBytes = bytes.ReplaceAll(keyBytes, []byte(`\\n`), []byte("\n"))
	keyBytes = bytes.ReplaceAll(keyBytes, []byte(`\n`), []byte("\n"))
	keyBytes = bytes.ReplaceAll(keyBytes, []byte("\r\n"), []byte("\n"))
	keyBytes = bytes.ReplaceAll(keyBytes, []byte("\r"), []byte("\n"))

	return keyBytes
}

func decodePEM(keyBytes []byte) (*pem.Block, error) {
	keyBytes = cleanKeyBytes(keyBytes)
	if len(keyBytes) == 0 {
		return nil, fmt.Errorf("empty key data")
	}

	// 1. Try direct pem.Decode
	block, _ := pem.Decode(keyBytes)
	if block != nil {
		return block, nil
	}

	// 2. Check if keyBytes is a file path that exists on disk
	if len(keyBytes) < 1024 && !bytes.Contains(keyBytes, []byte("\n")) {
		pathCandidate := string(keyBytes)
		if fileData, err := os.ReadFile(pathCandidate); err == nil {
			return decodePEM(fileData)
		}
	}

	// 3. Check if PEM was formatted on a single line with spaces instead of newlines
	beginMatch := pemBeginRegex.FindIndex(keyBytes)
	endMatch := pemEndRegex.FindIndex(keyBytes)
	if beginMatch != nil && endMatch != nil && beginMatch[1] <= endMatch[0] {
		header := keyBytes[beginMatch[0]:beginMatch[1]]
		footer := keyBytes[endMatch[0]:endMatch[1]]
		body := bytes.TrimSpace(keyBytes[beginMatch[1]:endMatch[0]])

		reconstructed := bytes.Buffer{}
		reconstructed.Write(header)
		reconstructed.WriteByte('\n')
		reconstructed.Write(body)
		reconstructed.WriteByte('\n')
		reconstructed.Write(footer)
		reconstructed.WriteByte('\n')

		block, _ = pem.Decode(reconstructed.Bytes())
		if block != nil {
			return block, nil
		}
	}

	// 4. Try base64 decoding (in case entire PEM was base64 encoded)
	compactB64 := bytes.Map(func(r rune) rune {
		if r == ' ' || r == '\n' || r == '\r' || r == '\t' {
			return -1
		}
		return r
	}, keyBytes)

	if decoded, err := base64.StdEncoding.DecodeString(string(compactB64)); err == nil {
		cleanDecoded := cleanKeyBytes(decoded)
		block, _ = pem.Decode(cleanDecoded)
		if block != nil {
			return block, nil
		}
	}

	return nil, fmt.Errorf("failed to decode PEM block")
}

func parseRSAPrivateKeyFromDER(der []byte) (*rsa.PrivateKey, error) {
	privateKey, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		privateKey, err = x509.ParsePKCS1PrivateKey(der)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
	}
	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not RSA private key")
	}
	return rsaKey, nil
}

func parseRSAPublicKeyFromDER(der []byte) (*rsa.PublicKey, error) {
	pub, err := x509.ParsePKIXPublicKey(der)
	if err == nil {
		rsaPub, ok := pub.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("key is not RSA public key")
		}
		return rsaPub, nil
	}
	rsaPub, err := x509.ParsePKCS1PublicKey(der)
	if err == nil {
		return rsaPub, nil
	}
	return nil, fmt.Errorf("failed to parse public key: %w", err)
}

func ParsePrivateKey(keyBytes []byte) (*rsa.PrivateKey, error) {
	block, err := decodePEM(keyBytes)
	if err != nil {
		// Fallback: check if keyBytes is raw Base64-encoded DER without PEM headers
		keyBytes = cleanKeyBytes(keyBytes)
		compactB64 := bytes.Map(func(r rune) rune {
			if r == ' ' || r == '\n' || r == '\r' || r == '\t' {
				return -1
			}
			return r
		}, keyBytes)
		if der, b64Err := base64.StdEncoding.DecodeString(string(compactB64)); b64Err == nil {
			if k, parseErr := parseRSAPrivateKeyFromDER(der); parseErr == nil {
				return k, nil
			}
		}
		return nil, err
	}
	return parseRSAPrivateKeyFromDER(block.Bytes)
}

func ParsePublicKey(keyBytes []byte) (*rsa.PublicKey, error) {
	block, err := decodePEM(keyBytes)
	if err != nil {
		// Fallback: check if keyBytes is raw Base64-encoded DER without PEM headers
		keyBytes = cleanKeyBytes(keyBytes)
		compactB64 := bytes.Map(func(r rune) rune {
			if r == ' ' || r == '\n' || r == '\r' || r == '\t' {
				return -1
			}
			return r
		}, keyBytes)
		if der, b64Err := base64.StdEncoding.DecodeString(string(compactB64)); b64Err == nil {
			if k, parseErr := parseRSAPublicKeyFromDER(der); parseErr == nil {
				return k, nil
			}
		}
		return nil, err
	}
	return parseRSAPublicKeyFromDER(block.Bytes)
}

func EncodePublicKeyToPEM(pub *rsa.PublicKey) (string, error) {
	b, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", fmt.Errorf("failed to marshal public key: %w", err)
	}
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: b,
	}
	return string(pem.EncodeToMemory(block)), nil
}

func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}
	return ParsePrivateKey(keyBytes)
}

func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}
	return ParsePublicKey(keyBytes)
}
