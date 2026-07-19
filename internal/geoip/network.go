package geoip

import (
	"net"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

type NetworkInfo struct {
	Country, Region, City, TimeZone, ASN, Organization, ProviderClass string
}

type NetworkClassifier struct {
	city, asn *geoip2.Reader
}

func OpenNetwork(cityPath, asnPath string) (*NetworkClassifier, error) {
	classifier := &NetworkClassifier{}
	var err error
	if cityPath != "" {
		classifier.city, err = geoip2.Open(cityPath)
		if err != nil {
			return nil, err
		}
	}
	if asnPath != "" {
		classifier.asn, err = geoip2.Open(asnPath)
		if err != nil {
			classifier.Close()
			return nil, err
		}
	}
	return classifier, nil
}

func (c *NetworkClassifier) Close() error {
	if c == nil {
		return nil
	}
	if c.city != nil {
		_ = c.city.Close()
	}
	if c.asn != nil {
		_ = c.asn.Close()
	}
	return nil
}

func (c *NetworkClassifier) Classify(ip net.IP) NetworkInfo {
	var info NetworkInfo
	if c == nil || ip == nil {
		return info
	}
	if c.city != nil {
		if record, err := c.city.City(ip); err == nil {
			info.Country = record.Country.IsoCode
			if len(record.Subdivisions) > 0 {
				info.Region = localizedName(record.Subdivisions[0].Names)
			}
			info.City = localizedName(record.City.Names)
			info.TimeZone = record.Location.TimeZone
		}
	}
	if c.asn != nil {
		if record, err := c.asn.ASN(ip); err == nil {
			if record.AutonomousSystemNumber > 0 {
				info.ASN = "AS" + uintToString(record.AutonomousSystemNumber)
			}
			info.Organization = record.AutonomousSystemOrganization
		}
	}
	info.ProviderClass = classifyProvider(info.Organization)
	return info
}

func localizedName(names map[string]string) string {
	if value := names["zh-CN"]; value != "" {
		return value
	}
	if value := names["en"]; value != "" {
		return value
	}
	for _, value := range names {
		return value
	}
	return ""
}

func classifyProvider(organization string) string {
	value := strings.ToLower(organization)
	for _, marker := range []string{"amazon", "google", "microsoft", "oracle", "cloudflare", "digitalocean", "hetzner", "alibaba", "tencent", "vultr", "linode"} {
		if strings.Contains(value, marker) {
			return "cloud"
		}
	}
	for _, marker := range []string{"hosting", "datacenter", "data center", "server", "colo"} {
		if strings.Contains(value, marker) {
			return "hosting"
		}
	}
	if strings.Contains(value, "mobile") || strings.Contains(value, "wireless") {
		return "mobile"
	}
	for _, marker := range []string{"telecom", "communications", "broadband", "internet", "cable"} {
		if strings.Contains(value, marker) {
			return "isp"
		}
	}
	return "unknown"
}

func uintToString(value uint) string {
	if value == 0 {
		return "0"
	}
	var digits [20]byte
	index := len(digits)
	for value > 0 {
		index--
		digits[index] = byte('0' + value%10)
		value /= 10
	}
	return string(digits[index:])
}
