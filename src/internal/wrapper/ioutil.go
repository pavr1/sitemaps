package wrapper

import (
	"io"
	"io/ioutil"
)

type IoutilWrapper interface {
	ReadAll(r io.Reader) ([]byte, error)
}

type ioutilWrapper struct {
}

func NewIoutilWrapper() IoutilWrapper {
	return &ioutilWrapper{}
}

func (i *ioutilWrapper) ReadAll(r io.Reader) ([]byte, error) {
	return ioutil.ReadAll(r)
}
