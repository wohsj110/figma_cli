package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Entry struct {
	URL      string    `json:"url"`
	CachedAt time.Time `json:"cached_at"`
	Body     []byte    `json:"body"`
}

type Store struct {
	Dir string
	TTL time.Duration
}

func DefaultDir() string {
	dir, err := os.UserCacheDir()
	if err != nil || dir == "" {
		return filepath.Join(".", ".figma-cli-cache")
	}
	return filepath.Join(dir, "figma-cli", "responses")
}

func (s Store) Get(rawURL string) ([]byte, bool, error) {
	if s.Dir == "" || s.TTL <= 0 {
		return nil, false, nil
	}
	data, err := os.ReadFile(s.path(rawURL))
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	var e Entry
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, false, err
	}
	if e.URL != rawURL || time.Since(e.CachedAt) > s.TTL {
		return nil, false, nil
	}
	return e.Body, true, nil
}

func (s Store) Put(rawURL string, body []byte) error {
	if s.Dir == "" || s.TTL <= 0 {
		return nil
	}
	if err := os.MkdirAll(s.Dir, 0o700); err != nil {
		return err
	}
	e := Entry{URL: rawURL, CachedAt: time.Now().UTC(), Body: body}
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	tmp := s.path(rawURL) + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	if err := os.Rename(tmp, s.path(rawURL)); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}

func (s Store) path(rawURL string) string {
	sum := sha256.Sum256([]byte(rawURL))
	return filepath.Join(s.Dir, fmt.Sprintf("%s.json", hex.EncodeToString(sum[:])))
}
