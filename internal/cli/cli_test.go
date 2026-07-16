package cli

import (
	"reflect"
	"testing"
)

func TestInterspersedMovesFlagsBeforePositionals(t *testing.T) {
	got := interspersed(
		[]string{"FILE", "--node", "1:2", "--raw"},
		map[string]bool{"node": true},
	)
	want := []string{"--node", "1:2", "--raw", "FILE"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestInterspersedPreservesEqualsFlags(t *testing.T) {
	got := interspersed(
		[]string{"FILE", "--depth=2"},
		map[string]bool{"depth": true},
	)
	want := []string{"--depth=2", "FILE"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}
