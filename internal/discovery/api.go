package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v3"
)

// APIDiscoverer extracts URL strings recursively from public JSON/YAML APIs.
type APIDiscoverer struct{ Client *http.Client }

func (d APIDiscoverer) Discover(ctx context.Context, seed Seed) ([]Candidate, error) {
	body, _, err := fetch(ctx, d.Client, seed.URL, readOptionalCredential(seed.TokenFile))
	if err != nil {
		return nil, err
	}
	var value any
	if json.Unmarshal(body, &value) != nil {
		if err := yaml.Unmarshal(body, &value); err != nil {
			return nil, &Error{Code: "invalid_document", Err: err}
		}
	}
	var candidates []Candidate
	walkURLs(value, func(found string) {
		candidates = append(candidates, Candidate{URL: found, Kind: "public-api", ParentURL: seed.URL, Evidence: "structured-value", Depth: seed.Depth + 1})
	})
	return deduplicate(candidates), nil
}

func walkURLs(value any, add func(string)) {
	switch typed := value.(type) {
	case string:
		for _, found := range absoluteURLPattern.FindAllString(typed, -1) {
			add(found)
		}
	case []any:
		for _, item := range typed {
			walkURLs(item, add)
		}
	case map[string]any:
		for _, item := range typed {
			walkURLs(item, add)
		}
	case map[any]any:
		for _, item := range typed {
			walkURLs(item, add)
		}
	default:
		_ = fmt.Sprint(typed)
	}
}
