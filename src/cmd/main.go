package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please provide the site URL: ")
	siteURL, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	siteURL = strings.TrimSuffix(siteURL, "\r\n")

	fmt.Print("Please provide max-depth: ")
	maxDept, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	maxDept = strings.TrimSuffix(maxDept, "\r\n")

	depth, err := strconv.Atoi(maxDept)
	if err != nil {
		log.WithError(err).Error("there was an error converting max-depth value, please verify it's a number")

		panic(err)
	}

	fmt.Print("Please provide the xml file name: ")
	xmlFileName, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	xmlFileName = strings.TrimSuffix(xmlFileName, "\r\n")

	url := wrapper.NewUrlWrapper()
	http := wrapper.NewHttpWrapper()
	ioutil := wrapper.NewIoutilWrapper()

	siteReader := sitereader.NewSiteReader(depth, log, url, http, ioutil)
	xmlHandler := handler.NewXmlHandler(xmlFileName, log)

	nodes, err := siteReader.GetPageLinks(siteURL, siteURL, 0)
	if err != nil {
		panic(err)
	}

	err = xmlHandler.WriteXML(nodes)
	if err != nil {
		panic(err)
	}

	log.Info("sitemap generation successfully processed")
}
