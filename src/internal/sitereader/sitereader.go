package sitereader

import (
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

type SiteReader interface {
	ReadHtml(url string) (string, error)
}

type siteReader struct {
	logger *logrus.Entry
}

func NewSiteReader(logger *logrus.Entry) SiteReader {
	return &siteReader{
		logger: logger,
	}
}

func (s *siteReader) ReadHtml(url string) (string, error) {
	l := s.logger.WithField("URL", url)
	l.Info("reading html data from url...")

	resp, err := http.Get(url)
	if err != nil {
		l.WithError(err).Error("there was an error calling endpoint")

		return "", err
	}

	defer resp.Body.Close()
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l.WithError(err).Error("there was an error calling response body")

		return "", err
	}

	return html, nil
}
