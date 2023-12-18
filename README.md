# geosite - a golang geosite library

[![godoc][godoc-img]][godoc] [![release][release-img]][release] [![goreport][goreport-img]][goreport]

### Getting Started

try on https://play.golang.org/p/Hp5zEc2XhwM
```go
package main

import (
	"context"
	"net/http"

	"github.com/phuslu/geosite"
)

func main() {
	dlc := &geosite.DomainListCommunity{Transport: http.DefaultTransport}
	dlc.Load(context.Background(), geosite.OnlineTarball)
	
	println(dlc.Site("chat.openai.com"))
}

// Output: openai
```

### Benchmarks
```
BenchmarkSite-8         63770426              24.70 ns/op          0 B/op        0 allocs/op
```

### Acknowledgment
This site or product uses dlc data available from http://github.com/v2fly/domain-list-community

[godoc-img]: http://img.shields.io/badge/godoc-reference-blue.svg
[godoc]: https://godoc.org/github.com/phuslu/geosite
[release-img]: https://img.shields.io/github/v/tag/phuslu/geosite?label=release
[release]: https://github.com/phuslu/geosite/releases
[goreport-img]: https://goreportcard.com/badge/github.com/phuslu/geosite
[goreport]: https://goreportcard.com/report/github.com/phuslu/geosite
