// Package geosite provides a geo site library for Go.
//
//	package main
//
//	import (
//		"context"
//
//		"github.com/phuslu/geosite"
//	)
//
//	func main() {
//		dlc := &geosite.DomainListCommunity{Transport: http.DefaultTransport}
//		dlc.Load(context.Backgroud(), geosite.OnlineTarball)
//		println(dlc.Site("chat.openai.com"))
//	}
//
//	// Output: openai
package geosite

import (
	"context"
	_ "embed" // for domain-list-community.tar.gz
	"io"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
)

const Version = "v1.0.20241001"

//go:embed domain-list-community.tar.gz
var InlineTarball string

const OnlineTarball = "https://codeload.github.com/v2fly/domain-list-community/legacy.tar.gz/refs/heads/master"

type DomainListCommunity struct {
	Transport http.RoundTripper
	dlc       atomic.Value // *dlc
}

// Load loads dlc data from repo url to memory.
func (d *DomainListCommunity) Load(ctx context.Context, tarball string) error {
	var data []byte
	var err error

	switch {
	case strings.HasPrefix(tarball, "\x1f\x8b\x08"):
		data = []byte(tarball)
	case strings.HasPrefix(tarball, "https://") || strings.HasPrefix(tarball, "http://"):
		transport := d.Transport
		if transport == nil {
			transport = http.DefaultTransport
		}
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, tarball, nil)
		resp, err := transport.RoundTrip(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		data, err = io.ReadAll(resp.Body)
	default:
		data, err = os.ReadFile(tarball)
	}

	if err != nil {
		return err
	}

	v, err := parse(data)
	if err != nil {
		return err
	}

	d.dlc.Store(v)

	return nil
}

// Site return geo site of domain.
func (d *DomainListCommunity) Site(domain string) (site string) {
	v, ok := d.dlc.Load().(*dlc)
	if !ok || v == nil {
		return
	}

	if site = v.full[domain]; site != "" {
		return
	}

	for _, pair := range v.regex {
		if pair.regex.MatchString(domain) {
			site = pair.site
			return
		}
	}

	s := domain
	for {
		i := strings.IndexByte(s, '.')
		if i < 0 || i+1 > len(s) {
			break
		}
		s = s[i+1:]
		if site = v.suffix[s]; site != "" {
			return
		}
	}

	return
}

// SiteAttrs return geo site of domain and its attrs.
func (d *DomainListCommunity) SiteAttrs(domain string) (site string, attrs []string) {
	v, ok := d.dlc.Load().(*dlc)
	if !ok || v == nil {
		return
	}

	if site = v.full[domain]; site != "" {
		attrs = v.attrs[domain]
		return
	}

	for _, pair := range v.regex {
		if pair.regex.MatchString(domain) {
			site = pair.site
			attrs = pair.attrs
			return
		}
	}

	s := domain
	for {
		i := strings.IndexByte(s, '.')
		if i < 0 || i+1 > len(s) {
			break
		}
		s = s[i+1:]
		if site = v.suffix[s]; site != "" {
			attrs = v.attrs[s]
			return
		}
	}

	return
}
