package sitereader

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sync"

	"github.com/pvr1/sitemaps/src/internal/model"
	"github.com/pvr1/sitemaps/src/internal/wrapper"
	"github.com/sirupsen/logrus"
)

const (
	OK = 200
)

type SiteReader interface {
	GetPageLinks(url string, parentLink string, depth int) (*model.Node, error)
}

type siteReader struct {
	maxDepth int
	logger   *logrus.Entry
	regExp   *regexp.Regexp
	url      wrapper.UrlWrapper
	http     wrapper.HttpWrapper
	ioutil   wrapper.IoutilWrapper
	mu       sync.Mutex
}

func NewSiteReader(maxDepth int, linkRegExpress string, logger *logrus.Entry, url wrapper.UrlWrapper, http wrapper.HttpWrapper, ioutil wrapper.IoutilWrapper) SiteReader {
	return &siteReader{
		maxDepth: maxDepth,
		regExp:   regexp.MustCompile(linkRegExpress),
		logger:   logger,
		url:      url,
		http:     http,
		ioutil:   ioutil,
	}
}

func (s *siteReader) GetPageLinks(link string, parentLink string, depth int) (*model.Node, error) {
	depth++

	var parsedUrl *url.URL
	var err error
	var content string
	var links []string

	l := s.logger.WithFields(logrus.Fields{
		"URL":           link,
		"Current_Depth": depth,
		"Site":          parentLink,
	})
	l.Info("reading html data from url...")

	if parsedUrl, err = s.url.Parse(link); err != nil {
		l.WithError(err).Error("there was an error parsing url")

		return nil, err
	}

	if content, err = s.getUrlContent(parsedUrl.String()); err != nil {
		l.WithError(err).Error("there was an error getting url content")

		return nil, err
	}

	if links, err = s.processLinks(parsedUrl, content); err != nil {
		l.WithError(err).Error("there was an error getting url content")

		return nil, err
	}

	if len(links) == 0 {
		return nil, nil
	}

	node := &model.Node{
		URL:   link,
		Nodes: []*model.Node{},
	}

	if depth == s.maxDepth {
		return node, nil
	}

	var wg sync.WaitGroup
	wg.Add(len(links))

	for _, currentLink := range links {
		go func(currentLink string, node *model.Node) {
			defer wg.Done()

			n, err := s.GetPageLinks(currentLink, link, depth)
			if err != nil {
				l.WithError(err).Error("there was an error getting page links")
				return
			}
			if n != nil {
				s.mu.Lock()
				node.Nodes = append(node.Nodes, n)
				s.mu.Unlock()
			}
		}(currentLink, node)
	}

	wg.Wait()

	return node, nil
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
		return "", fmt.Errorf("status code %d", resp.StatusCode)
	}

	if content, err = s.ioutil.ReadAll(resp.Body); err != nil {
		return "", err
	}

	return string(content), err
}

func (s *siteReader) processLinks(link *url.URL, content string) ([]string, error) {
	var (
		err     error
		links   map[string]interface{}
		matches [][]string
	)

	matches = s.regExp.FindAllStringSubmatch(content, -1)
	links = make(map[string]interface{})

	for _, val := range matches {
		var linkUrl *url.URL

		if linkUrl, err = url.Parse(val[1]); err != nil {
			return []string{}, err
		}

		if linkUrl.Path == "" {
			continue
		}

		var currentLink string
		if linkUrl.IsAbs() {
			currentLink = linkUrl.String()
		} else {
			currentLink = link.Scheme + "://" + link.Host + "/" + linkUrl.String()
		}

		if _, found := links[currentLink]; !found {
			links[currentLink] = nil
		}
	}

	result := []string{}
	for key := range links {
		result = append(result, key)
	}

	return result, err
}
