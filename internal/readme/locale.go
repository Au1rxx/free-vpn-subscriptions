package readme

// Locale holds every user-facing string rendered into a README for a single
// language. A new Locale must populate every field; partial translations
// would mix languages which hurts SEO (Google sees the page as non-cohesive).
type Locale struct {
	Code         string // "en", "zh", "ja", ...
	DisplayName  string // "English", "简体中文", "日本語"
	FileName     string // "README.md", "README_CN.md"
	LangAttr     string // HTML lang attribute: "en", "zh-Hans", "ja"

	// Badges
	BadgeNodes   string // "nodes"
	BadgeAlive   string // "alive"
	BadgeMedian  string // "median-rtt"
	BadgeUpdated string // "updated"

	// Hook sentences (rendered as blockquotes)
	Hook1 string
	Hook2 string
	// Extra SEO-keyword line in native language.
	KeywordLine string

	// Why This Project?
	WhyHeading string
	WhyBody    string

	// Verification — explains what we actually probe and what we can't.
	// VerificationBody may contain Markdown (headings, fenced blocks, lists).
	VerificationHeading string
	VerificationBody    string

	// One-click subscribe
	SubscribeHeading     string
	SubscribeIntro       string
	SubscribeColClient   string
	SubscribeColFormat   string
	SubscribeColURL      string

	// Clients
	ClientsHeading  string
	ClientsWindows  string
	ClientsMacOS    string
	ClientsIOS      string
	ClientsAndroid  string
	ClientsLinux    string

	// Stats
	StatsHeading     string
	StatsNodes       string
	StatsAlive       string
	StatsFastest     string
	StatsMedian      string
	StatsUpdated     string
	ProtocolMixLabel string
	SourcesLabel     string

	// By country
	ByCountryHeading string
	ByCountryIntro   string
	ByCountryColCC   string
	ByCountryColN    string

	// Guides (client tutorials)
	GuidesHeading string
	GuidesIntro   string
	// GuideLocaleSuffix is appended to guide filenames for non-English
	// locales: "" for English, ".zh" for Chinese. Keeps English at the
	// canonical URL while Chinese lives at clash-verge.zh.html.
	GuideLocaleSuffix string

	// FAQ
	FAQHeading string
	FAQ1Q      string
	FAQ1A      string
	FAQ2Q      string
	FAQ2A      string
	FAQ3Q      string
	FAQ3A      string
	FAQ4Q      string
	FAQ4A      string

	// Contributing
	ContributingHeading string
	ContributingBody    string

	// Disclaimer
	DisclaimerHeading string
	DisclaimerBody    string

	// Star History
	StarHistoryHeading string
	FinalCTA           string
}

// Locales returns every supported locale in fixed order so switcher bars are
// consistent across files.
func Locales() []Locale {
	return []Locale{EN, ZH, JA, KO, ES, PT, RU}
}

// Default is the locale written to README.md.
func Default() Locale { return EN }
