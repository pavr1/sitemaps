package wrapper

import "net/url"

type UrlWrapper interface {
	Parse(rawURL string) (*url.URL, error)
}

type urlWrapper struct {
}

func NewUrlWrapper() UrlWrapper {
	return &urlWrapper{}
}

func (u *urlWrapper) Parse(rawURL string) (*url.URL, error) {
	return url.Parse(rawURL)
}
