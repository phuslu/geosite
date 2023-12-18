package geosite

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"path"
	"regexp"
	"strings"
)

type dlc struct {
	full   map[string]string
	suffix map[string]string
	regex  []struct {
		regex *regexp.Regexp
		site  string
	}
}

func parse(data []byte) (*dlc, error) {
	v := &dlc{}
	v.full = make(map[string]string)
	v.suffix = make(map[string]string)

	var err error
	err = walktarball(bytes.NewReader(data), 1, func(header *tar.Header, r io.Reader) bool {
		if !strings.HasPrefix(header.Name, "data/") || header.Typeflag == tar.TypeDir {
			return true
		}

		site := path.Base(header.Name)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "#") {
				continue
			}
			if strings.HasPrefix(line, "include:") {
				continue
			}
			if !strings.Contains(line, ".") {
				continue
			}
			line = strings.Split(line, " ")[0]
			switch {
			case strings.HasPrefix(line, "regexp:"):
				var re *regexp.Regexp
				re, err = regexp.Compile(line[len("regexp:"):])
				if err != nil {
					re, err = regexp.CompilePOSIX(line[len("regexp:"):])
				}
				if err != nil {
					return true
				}
				v.regex = append(v.regex, struct {
					regex *regexp.Regexp
					site  string
				}{
					regex: re,
					site:  site,
				})
			case strings.HasPrefix(line, "full:"):
				v.full[line[len("full:"):]] = site
			case strings.HasPrefix(line, "domain:"):
				line = line[len("domain:"):]
				fallthrough
			default:
				v.suffix[line] = site
			}
		}
		if err = scanner.Err(); err != nil {
			return false
		}

		return true
	})
	if err != nil {
		return nil, err
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
