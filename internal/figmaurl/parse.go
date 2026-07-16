package figmaurl

import (
	"errors"
	"net/url"
	"strings"
)

// Target is the Figma resource encoded by a URL or raw file key.
type Target struct {
	FileKey string
	NodeID  string
	Raw     string
}

var ErrMissingFileKey = errors.New("missing Figma file key")

// Parse accepts raw file keys and Figma file/design/proto URLs.
func Parse(input string) (Target, error) {
	raw := strings.TrimSpace(input)
	if raw == "" {
		return Target{}, ErrMissingFileKey
	}
	if !strings.Contains(raw, "://") {
		return Target{FileKey: raw, Raw: raw}, nil
	}

	u, err := url.Parse(raw)
	if err != nil {
		return Target{}, err
	}
	parts := splitPath(u.Path)
	if len(parts) >= 2 && isDesignPath(parts[0]) && parts[1] != "" {
		return Target{
			FileKey: parts[1],
			NodeID:  normalizeNodeID(u.Query().Get("node-id")),
			Raw:     raw,
		}, nil
	}
	return Target{}, ErrMissingFileKey
}

func isDesignPath(part string) bool {
	return part == "file" || part == "design" || part == "proto"
}

func splitPath(path string) []string {
	raw := strings.Split(strings.Trim(path, "/"), "/")
	out := make([]string, 0, len(raw))
	for _, p := range raw {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func normalizeNodeID(id string) string {
	id = strings.TrimSpace(id)
	if id == "" {
		return ""
	}
	return strings.ReplaceAll(id, "-", ":")
}
