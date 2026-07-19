package geoip

import (
	"net"
	"testing"
)

func TestNetworkClassifyWithoutOptionalDatabases(t *testing.T) {
	classifier, err := OpenNetwork("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer classifier.Close()
	if info := classifier.Classify(net.ParseIP("2001:db8::1")); info.ProviderClass != "unknown" || info.Country != "" {
		t.Fatalf("unexpected unknown info: %+v", info)
	}
}

func TestProviderClassification(t *testing.T) {
	tests := map[string]string{"Oracle Cloud": "cloud", "Example Hosting LLC": "hosting", "Example Mobile": "mobile", "Example Telecom": "isp", "": "unknown"}
	for organization, want := range tests {
		if got := classifyProvider(organization); got != want {
			t.Fatalf("organization=%q got=%s want=%s", organization, got, want)
		}
	}
}
