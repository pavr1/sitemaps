package handler

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"github.com/pvr1/sitemaps/src/internal/model"
	"github.com/sirupsen/logrus"
)

type XMLHandler interface {
	WriteXML(nodes *model.Node) (string, error)
	ReturnXML(nodes *model.Node) (string, error)
}

type xmlHandler struct {
	xmlFilePath string
	xmlFileName string
	logger      *logrus.Entry
}

func NewXmlHandler(xmlFileName, xmlFilePath string, logger *logrus.Entry) XMLHandler {
	return &xmlHandler{
		xmlFileName: xmlFileName,
		xmlFilePath: xmlFilePath,
		logger:      logger,
	}
}
func (h *xmlHandler) WriteXML(nodes *model.Node) (string, error) {
	h.logger.Info("generating xml out from nodes...")
	fp := fmt.Sprintf(h.xmlFilePath, h.xmlFileName)

	file, err := os.Create(fp)
	if err != nil {
		h.logger.WithField("Path", fp).WithError(err).Error("there was an error while creating xml file")
		return "", err
	}

	xmlWriter := io.Writer(file)
	enc := xml.NewEncoder(xmlWriter)
	enc.Indent("", "	")

	if err := enc.Encode(nodes); err != nil {
		h.logger.WithError(err).Error("there was an error while encoding nodes")

		return "", err
	}

	return fp, nil
}

func (h *xmlHandler) ReturnXML(nodes *model.Node) (string, error) {
	file, err := xml.MarshalIndent(nodes, "", " ")
	if err != nil {
		h.logger.WithError(err).Error("there was an error marshalling nodes")

		return "", err
	}

	return string(file), nil
}
