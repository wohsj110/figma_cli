package credential

import (
	"errors"
	"fmt"
	"os"
	"strings"

	rawkeyring "github.com/byteness/keyring"
)

const (
	serviceName = "figma-cli"
	tokenKey    = "figma_token"
)

var ErrTokenNotFound = errors.New("no Figma token found; run `figma-cli init --token-stdin` or set FIGMA_TOKEN")

// Resolve returns FIGMA_TOKEN when set, otherwise the OS keyring token.
func Resolve() (string, string, error) {
	if token := strings.TrimSpace(os.Getenv("FIGMA_TOKEN")); token != "" {
		return token, "environment", nil
	}
	kr, err := open()
	if err != nil {
		return "", "", err
	}
	item, err := kr.Get(tokenKey)
	if errors.Is(err, rawkeyring.ErrKeyNotFound) {
		return "", "", ErrTokenNotFound
	}
	if err != nil {
		return "", "", fmt.Errorf("read token from keyring: %w", err)
	}
	token := strings.TrimSpace(string(item.Data))
	if token == "" {
		return "", "", ErrTokenNotFound
	}
	return token, "keyring", nil
}

// Store writes token to the OS keyring.
func Store(token string) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return errors.New("refusing to store an empty Figma token")
	}
	kr, err := open()
	if err != nil {
		return err
	}
	return kr.Set(rawkeyring.Item{
		Key:         tokenKey,
		Data:        []byte(token),
		Label:       "Figma API token",
		Description: "Token used by figma-cli to call the Figma REST API",
	})
}

func open() (rawkeyring.Keyring, error) {
	kr, err := rawkeyring.Open(rawkeyring.Config{ServiceName: serviceName})
	if err != nil {
		return nil, fmt.Errorf("open OS keyring: %w", err)
	}
	return kr, nil
}
