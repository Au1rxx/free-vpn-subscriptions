package readme

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/aggregate"
	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

func testInput() Input {
	return Input{
		Title:   "Free VPN Subscriptions",
		RepoURL: "https://github.com/example/repo",
		Nodes:   []*node.Node{{Protocol: node.ProtoVLESS, Server: "one.example", Port: 443, Country: "US"}},
		Summary: aggregate.Summary{
			TotalAlive:      900,
			TotalSelected:   120,
			ByCountry:       map[string]int{"US": 40},
			ByProtocol:      map[string]int{"vless": 120},
			MedianLatencyMS: 210,
			GeneratedAtUnix: time.Date(2026, 7, 25, 13, 30, 0, 0, time.UTC).Unix(),
		},
		CountryEnabled: true,
		MinPerCountry:  3,
	}
}

// Every Locale must populate every string field: a half-translated README
// mixes languages, which is exactly what the hreflang setup is meant to avoid.
func TestEveryLocaleIsFullyTranslated(t *testing.T) {
	// GuideLocaleSuffix is legitimately empty for every locale whose guides
	// fall back to the English pages.
	optional := map[string]bool{"GuideLocaleSuffix": true}
	for _, loc := range Locales() {
		value := reflect.ValueOf(loc)
		for i := 0; i < value.NumField(); i++ {
			field := value.Type().Field(i)
			if field.Type.Kind() != reflect.String || optional[field.Name] {
				continue
			}
			if strings.TrimSpace(value.Field(i).String()) == "" {
				t.Errorf("locale %q has empty %s", loc.Code, field.Name)
			}
		}
	}
}

func TestGenerateRendersStarConversionElements(t *testing.T) {
	in := testInput()
	for _, loc := range Locales() {
		body := Generate(in, loc)
		if !strings.Contains(body, "https://img.shields.io/github/stars/example/repo") {
			t.Errorf("locale %q missing live star badge", loc.Code)
		}
		if !strings.Contains(body, "https://github.com/example/repo/stargazers") {
			t.Errorf("locale %q star badge does not link to stargazers", loc.Code)
		}
		if !strings.Contains(body, loc.StarCTA) {
			t.Errorf("locale %q missing star call to action", loc.Code)
		}
		if !strings.Contains(body, "https://img.shields.io/github/license/example/repo") {
			t.Errorf("locale %q missing license badge", loc.Code)
		}
	}
}

func TestGenerateOmitsStarBadgeForNonGitHubRepo(t *testing.T) {
	in := testInput()
	in.RepoURL = "https://git.example.org/team/repo"
	body := Generate(in, Locales()[0])
	if strings.Contains(body, "img.shields.io/github/stars") {
		t.Error("star badge rendered for a non-GitHub repository URL")
	}
}

func TestGenerateUsesSelfHostedStarHistoryForEveryLocale(t *testing.T) {
	const image = "https://github.com/example/repo/raw/main/assets/star-history.svg"
	for _, loc := range Locales() {
		body := Generate(testInput(), loc)
		if !strings.Contains(body, "[![Star History Chart]("+image+")]") {
			t.Errorf("locale %q does not use the self-hosted chart", loc.Code)
		}
		if strings.Contains(body, "api.star-history.com") {
			t.Errorf("locale %q still uses the third-party chart", loc.Code)
		}
		if !strings.Contains(body, "https://github.com/example/repo/stargazers") {
			t.Errorf("locale %q chart does not link to stargazers", loc.Code)
		}
	}
}

func TestRepoSlug(t *testing.T) {
	for input, want := range map[string]string{
		"https://github.com/example/repo":      "example/repo",
		"https://github.com/example/repo/":     "example/repo",
		"https://git.example.org/example/repo": "",
		"https://github.com/example":           "",
	} {
		if got := repoSlug(input); got != want {
			t.Errorf("repoSlug(%q)=%q want %q", input, got, want)
		}
	}
}

// The reader came for a subscription URL. Keeping it above the pitch is the
// whole point of the layout, so assert it rather than trusting section order.
func TestSubscribeTableComesBeforeThePitch(t *testing.T) {
	in := testInput()
	for _, loc := range Locales() {
		body := Generate(in, loc)
		subscribe := strings.Index(body, loc.SubscribeHeading)
		why := strings.Index(body, loc.WhyHeading)
		if subscribe < 0 || why < 0 {
			t.Fatalf("locale %q is missing a section", loc.Code)
		}
		if subscribe > why {
			t.Errorf("locale %q puts the subscription table after the pitch", loc.Code)
		}
	}
}

func TestSubscribeURLLandsInTheFirstScreen(t *testing.T) {
	const budget = 40 // roughly one screen of rendered README
	for _, loc := range Locales() {
		body := Generate(testInput(), loc)
		found := -1
		for i, line := range strings.Split(body, "\n") {
			if strings.Contains(line, "output/clash.yaml") {
				found = i + 1
				break
			}
		}
		switch {
		case found < 0:
			t.Errorf("locale %q renders no subscription URL", loc.Code)
		case found > budget:
			t.Errorf("locale %q puts the first subscription URL on line %d, past the %d-line fold", loc.Code, found, budget)
		}
	}
}

// The verification section is the longest by far; collapsed it stops pushing
// everything else below the fold.
func TestVerificationDetailIsCollapsed(t *testing.T) {
	for _, loc := range Locales() {
		body := Generate(testInput(), loc)
		summary := "<summary><b>" + headingText(loc.VerificationHeading) + "</b></summary>"
		if !strings.Contains(body, summary) {
			t.Errorf("locale %q does not collapse the verification section", loc.Code)
		}
	}
}

func TestHeadingText(t *testing.T) {
	for input, want := range map[string]string{
		"## 🔬 How we verify": "🔬 How we verify",
		"# Title":            "Title",
		"Already plain":      "Already plain",
	} {
		if got := headingText(input); got != want {
			t.Errorf("headingText(%q)=%q want %q", input, got, want)
		}
	}
}
