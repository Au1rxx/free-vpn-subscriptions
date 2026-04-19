// Package geoip enriches nodes with an ISO country code by resolving each
// server hostname to an IP and looking it up in a MaxMind GeoLite2-Country mmdb.
// The database is downloaded on demand from a public mirror and cached on disk.
package geoip

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/oschwald/geoip2-golang"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/node"
)

// Resolver wraps an mmdb reader plus a DNS cache.
type Resolver struct {
	db       *geoip2.Reader
	dnsCache sync.Map // host -> net.IP (nil if NXDOMAIN)
}

// EnsureDB downloads the mmdb to dbPath if it is missing.
func EnsureDB(url, dbPath string) error {
	if _, err := os.Stat(dbPath); err == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return err
	}
	client := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("geoip download %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("geoip download %s: %s", url, resp.Status)
	}
	tmp := dbPath + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, dbPath)
}

// Open loads an mmdb file into a Resolver.
func Open(dbPath string) (*Resolver, error) {
	db, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("geoip open %s: %w", dbPath, err)
	}
	return &Resolver{db: db}, nil
}

// Close releases the underlying mmdb.
func (r *Resolver) Close() error { return r.db.Close() }

// resolve returns the first IP for a hostname (IPv4 preferred), or the host
// itself if it already parses as an IP.
func (r *Resolver) resolve(host string) net.IP {
	if ip := net.ParseIP(host); ip != nil {
		return ip
	}
	if v, ok := r.dnsCache.Load(host); ok {
		if v == nil {
			return nil
		}
		return v.(net.IP)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	addrs, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil || len(addrs) == 0 {
		r.dnsCache.Store(host, (net.IP)(nil))
		return nil
	}
	// Prefer IPv4.
	var pick net.IP
	for _, a := range addrs {
		if a.IP.To4() != nil {
			pick = a.IP
			break
		}
	}
	if pick == nil {
		pick = addrs[0].IP
	}
	r.dnsCache.Store(host, pick)
	return pick
}

// Country returns the ISO code for a host/IP, or "" if unknown.
func (r *Resolver) Country(hostOrIP string) string {
	ip := r.resolve(hostOrIP)
	if ip == nil {
		return ""
	}
	rec, err := r.db.Country(ip)
	if err != nil {
		return ""
	}
	return rec.Country.IsoCode
}

// Enrich fills n.Country for every node in parallel.
func (r *Resolver) Enrich(nodes []*node.Node, concurrency int) {
	if concurrency <= 0 {
		concurrency = 50
	}
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	for _, n := range nodes {
		wg.Add(1)
		sem <- struct{}{}
		go func(n *node.Node) {
			defer wg.Done()
			defer func() { <-sem }()
			n.Country = r.Country(n.Server)
		}(n)
	}
	wg.Wait()
}
