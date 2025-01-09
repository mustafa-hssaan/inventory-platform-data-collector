package url

import (
	"net/url"
	"strings"
)

type Builder struct {
	scheme string
	host   string
	path   []string
	query  url.Values
}

func NewBuilder() *Builder {
	return &Builder{
		scheme: "https",
		query:  make(url.Values),
	}
}
func (b *Builder) WithHost(host string) *Builder {
	b.host = strings.TrimRight(host, "/")
	return b
}
func (b *Builder) WithPath(paths ...string) *Builder {
	b.path = append(b.path, paths...)
	return b
}
func (b *Builder) WithParam(key, value string) *Builder {
	if value != "" {
		b.query.Add(key, value)
	}
	return b
}
func (b *Builder) WithParams(params map[string]string) *Builder {
	for k, v := range params {
		b.WithParam(k, v)
	}
	return b
}
func (b *Builder) Build() string {
	u := &url.URL{
		Scheme: b.scheme,
		Host:   b.host,
		Path:   "/" + strings.Join(b.path, "/"),
	}

	if len(b.query) > 0 {
		u.RawQuery = b.query.Encode()
	}

	return u.String()
}
