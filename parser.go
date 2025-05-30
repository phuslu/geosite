package geosite

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"path"
	"regexp"
	"strings"
	"unsafe"
)

type dlc struct {
	full   map[string]string
	suffix map[string]string
	attrs  map[string][]string
	regex  []struct {
		regex *regexp.Regexp
		site  string
		attrs []string
	}
}

func parse(reader io.Reader) (*dlc, error) {
	v := &dlc{}
	v.full = make(map[string]string)
	v.suffix = make(map[string]string)
	v.attrs = make(map[string][]string)

	reattr := regexp.MustCompile(`@\S+`)

	var err error
	err = walktarball(reader, 1, func(header *tar.Header, r io.Reader) bool {
		if !strings.HasPrefix(header.Name, "data/") || header.Typeflag == tar.TypeDir {
			return true
		}

		site := path.Base(header.Name)
		data, err := io.ReadAll(r)
		if err != nil {
			return false
		}
		content := *(*string)(unsafe.Pointer(&data))

		for {
			i := strings.IndexByte(content, '\n')
			if i < 0 {
				break
			}

			line := strings.TrimSpace(content[:i])
			content = content[i+1:]

			if strings.HasPrefix(line, "#") {
				continue
			}
			if strings.HasPrefix(line, "include:") {
				continue
			}
			if !strings.Contains(line, ".") {
				continue
			}
			var attrs []string
			for _, m := range reattr.FindAllStringSubmatch(line, -1) {
				attrs = append(attrs, m[0][1:])
			}
			if i := strings.IndexByte(line, ' '); i > 0 {
				line = line[:i]
			}
			switch {
			case strings.HasPrefix(line, "regexp:"):
				var re *regexp.Regexp
				re, err = regexp.Compile(line[len("regexp:"):])
				if err != nil {
					return true
				}
				v.regex = append(v.regex, struct {
					regex *regexp.Regexp
					site  string
					attrs []string
				}{
					regex: re,
					site:  site,
					attrs: attrs,
				})
			case strings.HasPrefix(line, "full:"):
				line = line[len("full:"):]
				v.full[line] = site
				v.attrs[line] = attrs
			case strings.HasPrefix(line, "domain:"):
				line = line[len("domain:"):]
				v.suffix[line] = site
				v.attrs[line] = attrs
			default:
				v.suffix[line] = site
				v.attrs[line] = attrs
			}
		}

		return true
	})
	if err != nil {
		return nil, err
	}

	for key, value := range v.suffix {
		if _, ok := v.full[key]; !ok {
			v.full[key] = value
		}
	}

	return v, nil
}

func walktarball(r io.Reader, strip int, f func(*tar.Header, io.Reader) bool) error {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)

	for {
		header, err := tr.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}
		if strip > 0 {
			if i := strings.IndexByte(header.Name, '/'); i > 0 {
				header.Name = header.Name[i+1:]
			}
		}
		if !f(header, tr) {
			return nil
		}
	}
}
