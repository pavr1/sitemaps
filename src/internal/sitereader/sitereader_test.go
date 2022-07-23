package sitereader

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/pvr1/sitemaps/src/internal/model"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type urlWrapperMock struct {
	mock.Mock
}

func (u urlWrapperMock) Parse(rawURL string) (*url.URL, error) {
	args := u.Called(rawURL)

	return args.Get(0).(*url.URL), args.Error(1)
}

type httpWrapperMock struct {
	mock.Mock
}

func (h httpWrapperMock) Get(url string) (resp *http.Response, err error) {
	args := h.Called(url)

	return args.Get(0).(*http.Response), args.Error(1)
}

type ioutilWrapperMock struct {
	mock.Mock
}

func (i ioutilWrapperMock) ReadAll(r io.Reader) ([]byte, error) {
	args := i.Called(r)

	return args.Get(0).([]byte), args.Error(1)
}

var tests = []struct {
	TestName    string
	link        string
	parentLink  string
	depth       int
	expectedURL *struct {
		method       string
		args         interface{}
		methodReturn struct {
			url *url.URL
			err error
		}
	}
	expectedHttp *struct {
		method       string
		args         interface{}
		methodReturn struct {
			response *http.Response
			err      error
		}
	}
	expectedIoutil *struct {
		method       string
		args         interface{}
		methodReturn struct {
			response []byte
			err      error
		}
	}
	expectedError  error
	expectedResult *model.Node
}{
	{
		TestName:   "GetPageLinksUrlParseFailed",
		link:       "test",
		parentLink: "",
		depth:      1,
		expectedURL: &struct {
			method       string
			args         interface{}
			methodReturn struct {
				url *url.URL
				err error
			}
		}{
			method: "Parse",
			args:   "test",
			methodReturn: struct {
				url *url.URL
				err error
			}{
				url: nil,
				err: fmt.Errorf("custom parse error"),
			},
		},
		expectedError: fmt.Errorf("custom parse error"),
	},
	{
		TestName:   "getUrlContentHttpGetFailed",
		link:       "url",
		parentLink: "",
		depth:      1,
		expectedURL: &struct {
			method       string
			args         interface{}
			methodReturn struct {
				url *url.URL
				err error
			}
		}{
			method: "Parse",
			args:   "url",
			methodReturn: struct {
				url *url.URL
				err error
			}{
				url: &url.URL{
					Path: "url",
				},
				err: nil,
			},
		},
		expectedHttp: &struct {
			method       string
			args         interface{}
			methodReturn struct {
				response *http.Response
				err      error
			}
		}{
			method: "Get",
			args:   "url",
			methodReturn: struct {
				response *http.Response
				err      error
			}{
				response: nil,
				err:      fmt.Errorf("custom http error"),
			},
		},
		expectedError: fmt.Errorf("custom http error"),
	},
	{
		TestName:   "getUrlContentHttpGetStatus404Failed",
		link:       "url",
		parentLink: "",
		depth:      1,
		expectedURL: &struct {
			method       string
			args         interface{}
			methodReturn struct {
				url *url.URL
				err error
			}
		}{
			method: "Parse",
			args:   "url",
			methodReturn: struct {
				url *url.URL
				err error
			}{
				url: &url.URL{
					Path: "url",
				},
				err: nil,
			},
		},
		expectedHttp: &struct {
			method       string
			args         interface{}
			methodReturn struct {
				response *http.Response
				err      error
			}
		}{
			method: "Get",
			args:   "url",
			methodReturn: struct {
				response *http.Response
				err      error
			}{
				response: &http.Response{
					StatusCode: 404,
					Body:       ioutil.NopCloser(bytes.NewBufferString("Hello World")),
				},
				err: nil,
			},
		},
		expectedError: fmt.Errorf("status code 404"),
	},
	{
		TestName:   "getUrlContentIoutilReadAllFailed",
		link:       "url",
		parentLink: "",
		depth:      1,
		expectedURL: &struct {
			method       string
			args         interface{}
			methodReturn struct {
				url *url.URL
				err error
			}
		}{
			method: "Parse",
			args:   "url",
			methodReturn: struct {
				url *url.URL
				err error
			}{
				url: &url.URL{
					Path: "url",
				},
				err: nil,
			},
		},
		expectedHttp: &struct {
			method       string
			args         interface{}
			methodReturn struct {
				response *http.Response
				err      error
			}
		}{
			method: "Get",
			args:   "url",
			methodReturn: struct {
				response *http.Response
				err      error
			}{
				response: &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewBufferString("Hello World")),
				},
				err: nil,
			},
		},
		expectedIoutil: &struct {
			method       string
			args         interface{}
			methodReturn struct {
				response []byte
				err      error
			}
		}{
			method: "ReadAll",
			args:   ioutil.NopCloser(bytes.NewBufferString("Hello World")),
			methodReturn: struct {
				response []byte
				err      error
			}{
				response: nil,
				err:      fmt.Errorf("custom ReadAll error"),
			},
		},
		expectedError: fmt.Errorf("custom ReadAll error"),
	},
	//...left processLinks test cases aside due to time
}

func Test_GetPageLinks(t *testing.T) {
	log := logrus.WithFields(logrus.Fields{
		"App": "Sitemaps",
	})

	urlWrapperMock := urlWrapperMock{}
	httpWrapperMock := httpWrapperMock{}
	ioutilWrapperMock := ioutilWrapperMock{}

	for _, tt := range tests {
		if tt.expectedURL != nil {
			urlWrapperMock.On(tt.expectedURL.method, tt.expectedURL.args).Return(tt.expectedURL.methodReturn.url, tt.expectedURL.methodReturn.err)
		}
		if tt.expectedHttp != nil {
			httpWrapperMock.On(tt.expectedHttp.method, tt.expectedHttp.args).Return(tt.expectedHttp.methodReturn.response, tt.expectedHttp.methodReturn.err)
		}
		if tt.expectedIoutil != nil {
			ioutilWrapperMock.On(tt.expectedIoutil.method, tt.expectedIoutil.args).Return(tt.expectedIoutil.methodReturn.response, tt.expectedIoutil.methodReturn.err)
		}

		reader := NewSiteReader(tt.depth, log, urlWrapperMock, httpWrapperMock, ioutilWrapperMock)
		node, err := reader.GetPageLinks(tt.link, tt.parentLink, tt.depth)

		if tt.expectedError == nil && err != nil {
			t.Errorf("expected error nil, but got '%s'", err.Error())
			return
		} else if tt.expectedError != nil {
			assert.EqualError(t, tt.expectedError, err.Error())
			return
		}

		assert.Equal(t, tt.expectedResult, node)
	}
}
