package main

import (
	"strings"

	"github.com/pvr1/sitemaps/src/internal/sitereader"
	"github.com/pvr1/sitemaps/src/internal/wrapper"
	"github.com/pvr1/sitemaps/src/internal/xml/handler"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.WithFields(logrus.Fields{
		"App": "Sitemaps",
	})
	log.Info("starting application...")

	//reader := bufio.NewReader(os.Stdin)
	// fmt.Print("Please provide the site URL: ")
	// siteURL, err := reader.ReadString('\n')
	// if err != nil {
	// 	panic(err)
	// }

	siteURL := "https://www.sitemaps.org/"
	siteURL = strings.TrimSuffix(siteURL, "\r\n")

	maxDepth := 3

	url := wrapper.NewUrlWrapper()
	http := wrapper.NewHttpWrapper()
	ioutil := wrapper.NewIoutilWrapper()

	siteReader := sitereader.NewSiteReader(maxDepth, log, url, http, ioutil)
	xmlHandler := handler.NewXmlHandler(log)

	nodes, err := siteReader.GetPageLinks(siteURL, siteURL, 0)
	if err != nil {
		panic(err)
	}

	err = xmlHandler.WriteXML(nodes)
	if err != nil {
		panic(err)
	}
}
