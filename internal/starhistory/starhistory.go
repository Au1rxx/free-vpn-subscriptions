package starhistory

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	acceptHeader = "application/vnd.github.star+json"
	maxPages     = 100
)

// Fetch returns every stargazer timestamp for repo using GitHub's paginated API.
func Fetch(ctx context.Context, client *http.Client, apiBase, repo, token string) ([]time.Time, error) {
	owner, name, err := splitRepo(repo)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, fmt.Errorf("HTTP client is required")
	}
	if strings.TrimSpace(token) == "" {
		return nil, fmt.Errorf("GitHub token is required")
	}
	base, err := url.Parse(strings.TrimRight(apiBase, "/"))
	if err != nil {
		return nil, fmt.Errorf("parse API base: %w", err)
	}
	base.Path = strings.TrimRight(base.Path, "/") + "/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(name) + "/stargazers"
	query := base.Query()
	query.Set("per_page", "100")
	base.RawQuery = query.Encode()

	next := base.String()
	stars := make([]time.Time, 0, 256)
	for page := 0; next != ""; page++ {
		if page >= maxPages {
			return nil, fmt.Errorf("stargazer pagination exceeds %d pages", maxPages)
		}
		pageURL, err := url.Parse(next)
		if err != nil || pageURL.Scheme != base.Scheme || !strings.EqualFold(pageURL.Host, base.Host) {
			return nil, fmt.Errorf("stargazer pagination origin differs from API base")
		}
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, next, nil)
		if err != nil {
			return nil, fmt.Errorf("create stargazer request: %w", err)
		}
		request.Header.Set("Accept", acceptHeader)
		request.Header.Set("Authorization", "Bearer "+token)
		response, err := client.Do(request)
		if err != nil {
			return nil, fmt.Errorf("fetch stargazers: %w", err)
		}
		if response.StatusCode < 200 || response.StatusCode >= 300 {
			message, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
			response.Body.Close()
			return nil, fmt.Errorf("fetch stargazers: HTTP %d: %s", response.StatusCode, strings.TrimSpace(string(message)))
		}
		var items []struct {
			StarredAt time.Time `json:"starred_at"`
		}
		err = json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&items)
		response.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("decode stargazers: %w", err)
		}
		for _, item := range items {
			if item.StarredAt.IsZero() {
				return nil, fmt.Errorf("decode stargazers: missing starred_at")
			}
			stars = append(stars, item.StarredAt.UTC())
		}
		next = nextLink(response.Header.Get("Link"))
	}
	sort.Slice(stars, func(i, j int) bool { return stars[i].Before(stars[j]) })
	return stars, nil
}

// RenderSVG renders a deterministic cumulative chart grouped by UTC date.
func RenderSVG(repo string, stars []time.Time) ([]byte, error) {
	if _, _, err := splitRepo(repo); err != nil {
		return nil, err
	}
	if len(stars) == 0 {
		return nil, fmt.Errorf("at least one star is required")
	}
	sorted := append([]time.Time(nil), stars...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Before(sorted[j]) })

	type point struct {
		date  time.Time
		count int
	}
	points := make([]point, 0, len(sorted))
	for index, value := range sorted {
		day := value.UTC().Truncate(24 * time.Hour)
		if len(points) == 0 || !points[len(points)-1].date.Equal(day) {
			points = append(points, point{date: day, count: index + 1})
		} else {
			points[len(points)-1].count = index + 1
		}
	}

	const width, height = 960, 320
	const left, right, top, bottom = 64, 28, 36, 54
	plotWidth := width - left - right
	plotHeight := height - top - bottom
	first, last := points[0].date, points[len(points)-1].date
	spanDays := int(last.Sub(first) / (24 * time.Hour))
	if spanDays < 1 {
		spanDays = 1
	}
	coordinates := make([]string, 0, len(points))
	for _, item := range points {
		day := int(item.date.Sub(first) / (24 * time.Hour))
		x := left + plotWidth*day/spanDays
		y := top + plotHeight - plotHeight*item.count/len(sorted)
		coordinates = append(coordinates, strconv.Itoa(x)+","+strconv.Itoa(y))
	}

	repoText := html.EscapeString(repo)
	firstText := first.Format("2006-01-02")
	lastText := last.Format("2006-01-02")
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" role="img" viewBox="0 0 %d %d">
<title>Star history for %s</title>
<desc>%d stars from %s to %s</desc>
<rect width="100%%" height="100%%" rx="12" fill="#0f172a"/>
<text x="%d" y="24" fill="#f8fafc" font-family="system-ui,sans-serif" font-size="16" font-weight="600">%s · %d stars</text>
<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#475569"/>
<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#475569"/>
<polyline points="%s" fill="none" stroke="#facc15" stroke-width="3" stroke-linejoin="round" stroke-linecap="round"/>
<text x="%d" y="%d" fill="#94a3b8" font-family="system-ui,sans-serif" font-size="12">%s</text>
<text x="%d" y="%d" text-anchor="end" fill="#94a3b8" font-family="system-ui,sans-serif" font-size="12">%s</text>
<text x="%d" y="%d" text-anchor="end" fill="#94a3b8" font-family="system-ui,sans-serif" font-size="12">%d</text>
</svg>
`, width, height, repoText, len(sorted), firstText, lastText,
		left, repoText, len(sorted),
		left, top+plotHeight, left+plotWidth, top+plotHeight,
		left, top, left, top+plotHeight,
		strings.Join(coordinates, " "),
		left, height-18, firstText,
		left+plotWidth, height-18, lastText,
		left-8, top+4, len(sorted))
	return []byte(svg), nil
}

func splitRepo(repo string) (string, string, error) {
	parts := strings.Split(strings.TrimSpace(repo), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" || strings.ContainsAny(repo, "?#") {
		return "", "", fmt.Errorf("repository must be owner/repo")
	}
	return parts[0], parts[1], nil
}

func nextLink(header string) string {
	for _, value := range strings.Split(header, ",") {
		segments := strings.Split(value, ";")
		if len(segments) < 2 || !strings.Contains(strings.Join(segments[1:], ";"), `rel="next"`) {
			continue
		}
		return strings.Trim(strings.TrimSpace(segments[0]), "<>")
	}
	return ""
}
