package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/pvr1/sitemaps/src/internal/config"
	"github.com/pvr1/sitemaps/src/internal/sitereader"
	"github.com/pvr1/sitemaps/src/internal/wrapper"
	"github.com/pvr1/sitemaps/src/internal/xml/handler"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Entry
	cfg *config.Config
)

func main() {
	var err error

	cfg, err = config.LoadConfig()
	if err != nil {
		log.WithError(err).Error("there was an error loading configuration data")
		panic(err)
	}

	log = logrus.WithFields(logrus.Fields{
		"App":           cfg.AppName,
		"ExecutionMode": cfg.ExecutionMode,
	})

	if cfg.ExecutionMode == 1 {
		execLocalConsole()
	} else if cfg.ExecutionMode == 2 {
		execHttpServer()
	}
}

func execLocalConsole() {
	log.Info("starting LocalConsole...")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please provide the site URL: ")
	siteURL, err := reader.ReadString('\n')
	if err != nil {
		log.WithField("Input", "URL").WithError(err).Error("there was an error reading input")
		panic(err)
	}
	siteURL = strings.TrimSuffix(siteURL, "\r\n")

	fmt.Print("Please provide max-depth: ")
	maxDept, err := reader.ReadString('\n')
	if err != nil {
		log.WithField("Input", "max-depth").WithError(err).Error("there was an error reading input")
		panic(err)
	}

	maxDept = strings.TrimSuffix(maxDept, "\r\n")

	depth, err := strconv.Atoi(maxDept)
	if err != nil {
		log.WithField("max-depth", maxDept).WithError(err).Error("there was an error converting max-depth value, please verify it's a number")

		panic(err)
	}

	fmt.Print("Please provide the xml file name: ")
	xmlFileName, err := reader.ReadString('\n')
	if err != nil {
		log.WithField("Input", "XMLFIleName").WithError(err).Error("there was an error reading input")
		panic(err)
	}

	xmlFileName = strings.TrimSuffix(xmlFileName, "\r\n")

	filePath, err := execute(depth, xmlFileName, siteURL, false)
	if err != nil {
		panic(err)
	}

	log.Infof("sitemap file generated successfully. Path: %s", filePath)
}

func execute(depth int, xmlFileName, siteURL string, isHttpRequest bool) (string, error) {
	url := wrapper.NewUrlWrapper()
	http := wrapper.NewHttpWrapper()
	ioutil := wrapper.NewIoutilWrapper()

	siteReader := sitereader.NewSiteReader(depth, cfg.UrlRegExpr, log, url, http, ioutil)
	xmlHandler := handler.NewXmlHandler(xmlFileName, cfg.OutputFile, log)
	nodes, err := siteReader.GetPageLinks(siteURL, siteURL, 0)
	if err != nil {
		log.WithError(err).Error("there was an error at GetPageLinks")
		return "", err
	}

	if isHttpRequest {
		xmlData, err := xmlHandler.ReturnXML(nodes)
		if err != nil {
			log.WithError(err).Error("there was an error at ReturnXML")
			return "", err
		}

		return xmlData, nil
	} else {
		filePath, err := xmlHandler.WriteXML(nodes)
		if err != nil {
			log.WithError(err).Error("there was an error at WriteXML")
			return "", err
		}

		return filePath, nil
	}
}

func execHttpServer() {
	log.Info("starting http server...")
	log.Infof("listening to port %d", cfg.HttpPort)

	http.HandleFunc("/", getHttpSitemap)
	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.HttpPort), nil)
	if err != nil {
		log.WithError(err).Error("there was an error listening to http server")
		panic(err)
	}
	//not the best approach but putting it to avoid exiting
	select {}
}

func getHttpSitemap(w http.ResponseWriter, r *http.Request) {
	url, ok := r.URL.Query()["url"]
	if !ok || len(url[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("url param not provided"))
		if err != nil {
			log.WithError(err).Error("there was an error writting the response")
		}
		return
	}
	maxDepth, ok := r.URL.Query()["maxDepth"]
	if !ok || len(maxDepth[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("maxDepth param not provided"))
		if err != nil {
			log.WithError(err).Error("there was an error writting the response")
		}
		return
	}

	depth, err := strconv.Atoi(maxDepth[0])
	if err != nil {
		log.WithField("max-depth", maxDepth).WithError(err).Error("there was an error converting max-depth value, please verify it's a number")

		panic(err)
	}

	fileName, ok := r.URL.Query()["xmlFileName"]
	if !ok || len(fileName[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("xmlFileName param not provided"))
		if err != nil {
			log.WithError(err).Error("there was an error writting the response")
		}
		return
	}

	xmlfileName := "http_" + fileName[0]

	xmlData, err := execute(depth, xmlfileName, url[0], true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.WithError(err).Error("there was an error writting the response")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(xmlData))
	if err != nil {
		log.WithError(err).Error("there was an error writting the response")
	}
}
