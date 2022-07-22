package sitereader

import (
	"net/http"
	"net/url"
	"regexp"

	"github.com/pvr1/sitemaps/src/internal/wrapper"
	"github.com/sirupsen/logrus"
)

const (
	OK         = 200
	LinkRegExp = "<a.*?href=\"(.*?)\""
)

type SiteReader interface {
	ReadHtml(url string) ([]string, error)
}

type siteReader struct {
	logger *logrus.Entry
	regExp *regexp.Regexp
	url    wrapper.UrlWrapper
	http   wrapper.HttpWrapper
	ioutil wrapper.IoutilWrapper
}

func NewSiteReader(logger *logrus.Entry, url wrapper.UrlWrapper, http wrapper.HttpWrapper, ioutil wrapper.IoutilWrapper) SiteReader {
	return &siteReader{
		logger: logger,
		url:    url,
		http:   http,
		ioutil: ioutil,
		regExp: regexp.MustCompile(LinkRegExp),
	}
}

func (s *siteReader) ReadHtml(link string) ([]string, error) {
	var parsedUrl *url.URL
	var err error
	var content string
	var links []string

	l := s.logger.WithField("URL", link)
	l.Info("reading html data from url...")

	if parsedUrl, err = s.url.Parse(link); err != nil {
		l.WithError(err).Error("there was an error parsing url")

		return []string{}, err
	}

	if content, err = s.getUrlContent(parsedUrl.String()); err != nil {
		l.WithError(err).Error("there was an error getting url content")

		return []string{}, err
	}

	if links, err = s.processLinks(parsedUrl, content); err != nil {
		l.WithError(err).Error("there was an error getting url content")

		return []string{}, err
	}

	return links, nil
}

func (s *siteReader) getUrlContent(urlToGet string) (string, error) {
	var (
		err     error
		content []byte
		resp    *http.Response
	)

	if resp, err = s.http.Get(urlToGet); err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != OK {
		return "", err
	}

	if content, err = s.ioutil.ReadAll(resp.Body); err != nil {
		return "", err
	}

	return string(content), err
}

func (s *siteReader) processLinks(link *url.URL, content string) ([]string, error) {
	var (
		err     error
		links   []string = make([]string, 0)
		matches [][]string
	)

	matches = s.regExp.FindAllStringSubmatch(content, -1)

	for _, val := range matches {
		var linkUrl *url.URL

		if linkUrl, err = url.Parse(val[1]); err != nil {
			return links, err
		}

		links = append(links, linkUrl.String())
	}

	return links, err
}
