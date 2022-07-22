package wrapper

import "net/http"

type HttpWrapper interface {
	Get(url string) (resp *http.Response, err error)
}

type httpWrapper struct {
}

func NewHttpWrapper() HttpWrapper {
	return &httpWrapper{}
}

func (h *httpWrapper) Get(url string) (resp *http.Response, err error) {
	return http.Get(url)
}
