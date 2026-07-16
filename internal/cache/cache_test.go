package cache

import (
	"testing"
	"time"
)

func TestStoreRoundTrip(t *testing.T) {
	store := Store{Dir: t.TempDir(), TTL: time.Hour}
	if err := store.Put("https://example.test/a", []byte(`{"ok":true}`)); err != nil {
		t.Fatal(err)
	}
	body, ok, err := store.Get("https://example.test/a")
	if err != nil {
		t.Fatal(err)
	}
	if !ok || string(body) != `{"ok":true}` {
		t.Fatalf("cache miss/body mismatch: ok=%v body=%q", ok, body)
	}
}

func TestStoreExpired(t *testing.T) {
	store := Store{Dir: t.TempDir(), TTL: time.Nanosecond}
	if err := store.Put("https://example.test/a", []byte(`x`)); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond)
	_, ok, err := store.Get("https://example.test/a")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected expired entry")
	}
}
