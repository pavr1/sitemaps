package handler

import (
	"encoding/xml"
	"io/ioutil"

	"github.com/pvr1/sitemaps/src/internal/model"
	"github.com/sirupsen/logrus"
)

type XMLHandler interface {
	WriteXML(nodes *model.Node) error
}

type xmlHandler struct {
	logger *logrus.Entry
}

func NewXmlHandler(logger *logrus.Entry) XMLHandler {
	return &xmlHandler{
		logger: logger,
	}
}
func (h *xmlHandler) WriteXML(nodes *model.Node) error {
	h.logger.Info("generating xml out from nodes...")

	file, err := xml.MarshalIndent(&nodes, "", " ")
	if err != nil {
		h.logger.WithError(err).Error("there was an error generating xml file")
		return err
	}
	err = ioutil.WriteFile("./output/sitemap.xml", file, 0644)
	if err != nil {
		h.logger.WithError(err).Error("there was an error writing xml file")
		return err
	}

	return nil
}
