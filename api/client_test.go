package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
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
