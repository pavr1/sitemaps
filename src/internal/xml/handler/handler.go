package handler

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"github.com/pvr1/sitemaps/src/internal/model"
	"github.com/sirupsen/logrus"
)

const (
	FilePath = "../internal/xml/output/%s.xml"
)

type XMLHandler interface {
	WriteXML(nodes *model.Node) error
}

type xmlHandler struct {
	xmlFileName string
	logger      *logrus.Entry
}

func NewXmlHandler(xmlFileName string, logger *logrus.Entry) XMLHandler {
	return &xmlHandler{
		xmlFileName: xmlFileName,
		logger:      logger,
	}
}
func (h *xmlHandler) WriteXML(nodes *model.Node) error {
	h.logger.Info("generating xml out from nodes...")
	fp := fmt.Sprintf(FilePath, h.xmlFileName)

	file, err := os.Create(fp)
	if err != nil {
		h.logger.WithField("Path", fp).WithError(err).Error("there was an error while creating xml file")
		return err
	}

	xmlWriter := io.Writer(file)
	enc := xml.NewEncoder(xmlWriter)
	enc.Indent("", "	")

	if err := enc.Encode(nodes); err != nil {
		h.logger.WithError(err).Error("there was an error while encoding nodes")

		return err
	}

	return nil
}
