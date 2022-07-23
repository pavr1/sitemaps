package sitereader

import (
	"fmt"
	"io"
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
		methodReturn struct {
			url *url.URL
			err error
		}
	}
	expectedHttp *struct {
		method       string
		methodReturn struct {
			response *http.Response
			err      error
		}
	}
	expectedIoutil *struct {
		method       string
		methodReturn struct {
			response *http.Response
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
			methodReturn struct {
				url *url.URL
				err error
			}
		}{
			method: "Parse",
			methodReturn: struct {
				url *url.URL
				err error
			}{
				url: nil,
				err: fmt.Errorf("custom error"),
			},
		},
	},
}

func Test_GetPageLinks(t *testing.T) {
	log := logrus.WithFields(logrus.Fields{
		"App": "Sitemaps",
	})

	urlWrapperMock := urlWrapperMock{}
	httpWrapperMock := httpWrapperMock{}
	ioutilWrapperMock := ioutilWrapperMock{}

	reader := NewSiteReader(3, log, urlWrapperMock, httpWrapperMock, ioutilWrapperMock)

	for _, tt := range tests {
		if tt.expectedURL != nil {
			urlWrapperMock.On(tt.expectedURL.method).Return(tt.expectedURL.methodReturn.url, tt.expectedURL.methodReturn.err)
		}
		if tt.expectedURL != nil {
			httpWrapperMock.On(tt.expectedHttp.method).Return(tt.expectedHttp.methodReturn.response, tt.expectedHttp.methodReturn.err)
		}
		if tt.expectedIoutil != nil {
			ioutilWrapperMock.On(tt.expectedHttp.method, mock.AnythingOfType("io.ReadCloser")).Return(tt.expectedHttp.methodReturn.response, tt.expectedHttp.methodReturn.err)
		}

		node, err := reader.GetPageLinks(tt.link, tt.parentLink, tt.depth)

		if tt.expectedError == nil && err != nil {
			t.Errorf("expected error nil, but got %s", err)
			return
		} else if tt.expectedError != nil {
			assert.EqualError(t, tt.expectedError, err.Error())
			return
		}

		assert.Equal(t, tt.expectedResult, node)
	}
}
