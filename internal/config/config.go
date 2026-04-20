// Package config loads and validates the aggregator's YAML configuration.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Sources   []Source        `yaml:"sources"`
	Probe     ProbeConfig     `yaml:"probe"`
	Verify    VerifyConfig    `yaml:"verify"`
	Aggregate AggregateConfig `yaml:"aggregate"`
	Output    OutputConfig    `yaml:"output"`
	GeoIP     GeoIPConfig     `yaml:"geoip"`
	Readme    ReadmeConfig    `yaml:"readme"`
}

type VerifyConfig struct {
	Enabled        bool     `yaml:"enabled"`
	CandidatePool  int      `yaml:"candidate_pool"`
	BatchSize      int      `yaml:"batch_size"`
	BasePort       int      `yaml:"base_port"`
	Concurrency    int      `yaml:"concurrency"`
	TimeoutMS      int      `yaml:"timeout_ms"`
	Rounds         int      `yaml:"rounds"`
	RoundGapMS     int      `yaml:"round_gap_ms"`
	Targets        []string `yaml:"targets"`
	SingBoxBin     string   `yaml:"sing_box_bin"`
	StartupTimeoutMS int    `yaml:"startup_timeout_ms"`
}

type GeoIPConfig struct {
	Enabled       bool   `yaml:"enabled"`
	DBURL         string `yaml:"db_url"`
	DBPath        string `yaml:"db_path"`
	MinPerCountry int    `yaml:"min_per_country"`
}

// Source describes a single upstream subscription feed.
type Source struct {
	Name    string `yaml:"name"`
	URL     string `yaml:"url"`
	Format  string `yaml:"format"`  // uri-list | base64 | clash
	Enabled bool   `yaml:"enabled"`
}

type ProbeConfig struct {
	TimeoutMS         int  `yaml:"timeout_ms"`
	Concurrency       int  `yaml:"concurrency"`
	MaxNodesPerSource int  `yaml:"max_nodes_per_source"`
	TLSVerify         bool `yaml:"tls_verify"`
}

type AggregateConfig struct {
	TopN           int      `yaml:"top_n"`
	MaxRTTMS       int      `yaml:"max_rtt_ms"`
	MinPerProtocol int      `yaml:"min_per_protocol"`
	Protocols      []string `yaml:"protocols"`
}

type OutputConfig struct {
	Dir     string     `yaml:"dir"`
	Formats []string   `yaml:"formats"`
	Pages   PagesConfig `yaml:"pages"`
}

type PagesConfig struct {
	Enabled bool   `yaml:"enabled"`
	Dir     string `yaml:"dir"`
	SiteURL string `yaml:"site_url"`
}

type ReadmeConfig struct {
	Title   string `yaml:"title"`
	RepoURL string `yaml:"repo_url"`
}

// Load reads and validates a YAML config file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config read %q: %w", path, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config parse: %w", err)
	}
	applyDefaults(&cfg)
	if err := validate(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func applyDefaults(c *Config) {
	if c.Probe.TimeoutMS == 0 {
		c.Probe.TimeoutMS = 3000
	}
	if c.Probe.Concurrency == 0 {
		c.Probe.Concurrency = 50
	}
	if c.Probe.MaxNodesPerSource == 0 {
		c.Probe.MaxNodesPerSource = 1000
	}
	if c.Aggregate.TopN == 0 {
		c.Aggregate.TopN = 100
	}
	if c.Aggregate.MaxRTTMS == 0 {
		c.Aggregate.MaxRTTMS = 2000
	}
	if c.Output.Dir == "" {
		c.Output.Dir = "output"
	}
	if len(c.Output.Formats) == 0 {
		c.Output.Formats = []string{"clash", "v2ray-base64"}
	}
	if c.GeoIP.DBURL == "" {
		c.GeoIP.DBURL = "https://github.com/P3TERX/GeoLite.mmdb/releases/latest/download/GeoLite2-Country.mmdb"
	}
	if c.GeoIP.DBPath == "" {
		c.GeoIP.DBPath = "output/.cache/GeoLite2-Country.mmdb"
	}
	if c.GeoIP.MinPerCountry == 0 {
		c.GeoIP.MinPerCountry = 3
	}
	if c.Output.Pages.Dir == "" {
		c.Output.Pages.Dir = "docs"
	}
	if c.Verify.CandidatePool == 0 {
		c.Verify.CandidatePool = 600
	}
	if c.Verify.BatchSize == 0 {
		c.Verify.BatchSize = 40
	}
	if c.Verify.BasePort == 0 {
		c.Verify.BasePort = 20000
	}
	if c.Verify.Concurrency == 0 {
		c.Verify.Concurrency = 20
	}
	if c.Verify.TimeoutMS == 0 {
		c.Verify.TimeoutMS = 6000
	}
	if c.Verify.Rounds == 0 {
		c.Verify.Rounds = 2
	}
	if c.Verify.RoundGapMS == 0 {
		c.Verify.RoundGapMS = 45000
	}
	if len(c.Verify.Targets) == 0 {
		c.Verify.Targets = []string{
			"http://www.gstatic.com/generate_204",
			"https://www.cloudflare.com/cdn-cgi/trace",
		}
	}
	if c.Verify.SingBoxBin == "" {
		c.Verify.SingBoxBin = "sing-box"
	}
	if c.Verify.StartupTimeoutMS == 0 {
		c.Verify.StartupTimeoutMS = 10000
	}
}

func validate(c *Config) error {
	if len(c.Sources) == 0 {
		return fmt.Errorf("config: no sources defined")
	}
	seen := make(map[string]bool)
	for i, s := range c.Sources {
		if s.Name == "" {
			return fmt.Errorf("config: sources[%d] missing name", i)
		}
		if seen[s.Name] {
			return fmt.Errorf("config: duplicate source name %q", s.Name)
		}
		seen[s.Name] = true
		if s.URL == "" {
			return fmt.Errorf("config: sources[%d] %q missing url", i, s.Name)
		}
		switch s.Format {
		case "uri-list", "base64", "clash":
		default:
			return fmt.Errorf("config: sources[%d] %q invalid format %q (want uri-list|base64|clash)", i, s.Name, s.Format)
		}
	}
	return nil
}
