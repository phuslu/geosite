package geosite

import (
	"context"
	"net/http"
	"testing"
)

var d = func() *DomainListCommunity {
	dlc := &DomainListCommunity{Transport: http.DefaultTransport}
	if err := dlc.Load(context.Background(), InlineTarball); err != nil {
		panic(err)
	}
	return dlc
}()

func TestSite(t *testing.T) {
	cases := []struct {
		Domain string
		Site   string
	}{
		{"phus.lu", ""},
		{"apple.com.co", "apple"},
		{"www.google.com", "google"},
		{"chat.openai.com", "openai"},
	}

	for _, c := range cases {
		site := d.Site(c.Domain)
		if site != c.Site {
			t.Errorf("Site(%#v) return \"%s\", expect %#v", c.Domain, site, c.Site)
		}
	}
}

func BenchmarkSite(b *testing.B) {
	domain := "chat.openai.com"

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Site(domain)
	}
}
