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

func TestSiteAttrs(t *testing.T) {
	cases := []struct {
		Domain string
		Site   string
		Attrs  []string
	}{
		{"phus.lu", "", nil},
		{"www.asus.com.cn", "asus", []string{"cn"}},
	}

	for _, c := range cases {
		site, attrs := d.SiteAttrs(c.Domain)
		if site != c.Site {
			t.Errorf("SiteAttrs(%#v) return site \"%s\", expect %#v", c.Domain, site, c.Site)
		}
		if len(attrs) != len(c.Attrs) || (len(attrs) > 0 && attrs[0] != c.Attrs[0]) {
			t.Errorf("SiteAttrs(%#v) return attrs \"%s\", expect %#v", c.Domain, attrs, c.Attrs)
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

func BenchmarkSiteAttrs(b *testing.B) {
	domain := "www.asus.com.cn"

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.SiteAttrs(domain)
	}
}
