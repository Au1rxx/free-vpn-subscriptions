// Package config loads and validates the aggregator's YAML configuration.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Sources   []Source        `yaml:"sources"`
	Database  DatabaseConfig  `yaml:"database"`
	Probe     ProbeConfig     `yaml:"probe"`
	Verify    VerifyConfig    `yaml:"verify"`
	Aggregate AggregateConfig `yaml:"aggregate"`
	Output    OutputConfig    `yaml:"output"`
	GeoIP     GeoIPConfig     `yaml:"geoip"`
	Readme    ReadmeConfig    `yaml:"readme"`
}

// DatabaseConfig contains non-secret MySQL connection settings. The password
// is always read from PasswordFile at runtime and must never be stored here.
type DatabaseConfig struct {
	Enabled      bool   `yaml:"enabled"`
	Address      string `yaml:"address"`
	Name         string `yaml:"name"`
	User         string `yaml:"user"`
	PasswordFile string `yaml:"password_file"`
	TLSMode      string `yaml:"tls_mode"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
}

type VerifyConfig struct {
	Enabled          bool     `yaml:"enabled"`
	CandidatePool    int      `yaml:"candidate_pool"`
	BatchSize        int      `yaml:"batch_size"`
	BasePort         int      `yaml:"base_port"`
	Concurrency      int      `yaml:"concurrency"`
	TimeoutMS        int      `yaml:"timeout_ms"`
	Rounds           int      `yaml:"rounds"`
	RoundGapMS       int      `yaml:"round_gap_ms"`
	Targets          []string `yaml:"targets"`
	SingBoxBin       string   `yaml:"sing_box_bin"`
	StartupTimeoutMS int      `yaml:"startup_timeout_ms"`
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
	Format  string `yaml:"format"` // uri-list | base64 | clash
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
	Dir     string      `yaml:"dir"`
	Formats []string    `yaml:"formats"`
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
	if c.Database.Address == "" {
		c.Database.Address = "127.0.0.1:13306"
	}
	if c.Database.Name == "" {
		c.Database.Name = "vpn_nodes"
	}
	if c.Database.User == "" {
		c.Database.User = "adminai"
	}
	if c.Database.TLSMode == "" {
		c.Database.TLSMode = "required"
	}
	if c.Database.MaxOpenConns == 0 {
		c.Database.MaxOpenConns = 20
	}
	if c.Database.MaxIdleConns == 0 {
		c.Database.MaxIdleConns = 10
	}
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
	switch c.Database.TLSMode {
	case "required", "verify-ca", "verify-identity":
	default:
		return fmt.Errorf("config: database tls_mode %q is invalid", c.Database.TLSMode)
	}
	if c.Database.MaxOpenConns < 1 {
		return fmt.Errorf("config: database max_open_conns must be positive")
	}
	if c.Database.MaxIdleConns < 0 || c.Database.MaxIdleConns > c.Database.MaxOpenConns {
		return fmt.Errorf("config: database max_idle_conns must be between 0 and max_open_conns")
	}
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
