package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/pvr1/sitemaps/src/internal/sitereader"
	"github.com/pvr1/sitemaps/src/internal/wrapper"
	"github.com/pvr1/sitemaps/src/internal/xml/handler"
	"github.com/sirupsen/logrus"
)

type ExecutionMode int32

const (
	//1 = LocalConsole / 2 = HttpServer
	LocalConsole ExecutionMode = 1
	HttpServer   ExecutionMode = 2
)

var (
	executionMode ExecutionMode = 2
	log                         = logrus.WithFields(logrus.Fields{
		"App":           "Sitemaps",
		"ExecutionMode": executionMode,
	})
)

func main() {

	if executionMode == 1 {
		ExecLocalConsole()
	} else if executionMode == 2 {
		ExecHttpServer()
	}
}

func ExecLocalConsole() {
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

	err = execute(depth, xmlFileName, siteURL)
	if err != nil {
		panic(err)
	}

	log.Info("sitemap generation successfully processed")
}

func execute(depth int, xmlFileName, siteURL string) error {
	url := wrapper.NewUrlWrapper()
	http := wrapper.NewHttpWrapper()
	ioutil := wrapper.NewIoutilWrapper()

	siteReader := sitereader.NewSiteReader(depth, log, url, http, ioutil)
	xmlHandler := handler.NewXmlHandler(xmlFileName, log)
	nodes, err := siteReader.GetPageLinks(siteURL, siteURL, 0)
	if err != nil {
		log.WithError(err).Error("there was an error at GetPageLinks")
		return err
	}

	err = xmlHandler.WriteXML(nodes)
	if err != nil {
		log.WithError(err).Error("there was an error at WriteXML")
		return err
	}

	return nil
}

func ExecHttpServer() {
	log.Info("starting http server...")
	log.Info("listening to port %d", 8080)

	http.HandleFunc("/", getHttpSitemap)
	http.ListenAndServe(":8080", nil)
}

func getHttpSitemap(w http.ResponseWriter, r *http.Request) {
	url, ok := r.URL.Query()["url"]
	if !ok || len(url[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("url param not provided"))
		return
	}
	maxDepth, ok := r.URL.Query()["maxDepth"]
	if !ok || len(maxDepth[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("maxDepth param not provided"))
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
		w.Write([]byte("xmlFileName param not provided"))
		return
	}

	xmlfileName := "http_" + fileName[0]

	err = execute(depth, xmlfileName, url[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("file '%s' stored successfully", xmlfileName)))
}
