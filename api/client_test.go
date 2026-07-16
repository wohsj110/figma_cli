package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/wohsj110/figma_cli/internal/cache"
)

func TestMeSendsFigmaTokenHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/me" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Figma-Token"); got != "tok" {
			t.Fatalf("X-Figma-Token = %q", got)
		}
		_, _ = w.Write([]byte(`{"id":"u1","email":"u@example.com","handle":"User"}`))
	}))
	defer server.Close()

	client := New("tok")
	client.BaseURL = server.URL
	me, _, err := client.Me(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if me.Handle != "User" {
		t.Fatalf("handle = %q", me.Handle)
	}
}

func TestNodesBuildsIDsQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/files/abc/nodes" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("ids"); got != "1:2,3:4" {
			t.Fatalf("ids = %q", got)
		}
		_, _ = w.Write([]byte(`{"nodes":{"1:2":{"document":{"id":"1:2","name":"Frame","type":"FRAME"}}}}`))
	}))
	defer server.Close()

	client := New("tok")
	client.BaseURL = server.URL
	resp, _, err := client.Nodes(context.Background(), "abc", []string{"1:2", "3:4"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Nodes["1:2"].Document.Name != "Frame" {
		t.Fatalf("node not decoded: %#v", resp.Nodes)
	}
}

func TestAPIErrorIncludesStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"err":"forbidden"}`))
	}))
	defer server.Close()

	client := New("tok")
	client.BaseURL = server.URL
	_, _, err := client.Me(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestVariablesEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/files/abc/variables/local" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"meta":{"variables":{"v1":{"id":"v1","name":"Color/Primary","resolvedType":"COLOR","variableCollectionId":"c1"}},"variableCollections":{"c1":{"id":"c1","name":"Colors"}}}}`))
	}))
	defer server.Close()

	client := New("tok")
	client.BaseURL = server.URL
	resp, _, err := client.Variables(context.Background(), "abc")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Meta.Variables["v1"].Name != "Color/Primary" {
		t.Fatalf("variables not decoded: %#v", resp.Meta.Variables)
	}
}

func TestClientCacheAvoidsSecondRequest(t *testing.T) {
	var hits int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		_, _ = w.Write([]byte(`{"id":"u1","email":"u@example.com","handle":"User"}`))
	}))
	defer server.Close()

	client := New("tok")
	client.BaseURL = server.URL
	client.Cache = cache.Store{Dir: t.TempDir(), TTL: time.Hour}
	if _, _, err := client.Me(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, _, err := client.Me(context.Background()); err != nil {
		t.Fatal(err)
	}
	if hits != 1 {
		t.Fatalf("hits = %d, want 1", hits)
	}
}
