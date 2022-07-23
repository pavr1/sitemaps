package model

import "encoding/xml"

type Node struct {
	XMLName xml.Name `xml:"node"`
	URL     string   `xml:"url"`
	Nodes   []*Node  `xml:"nodes>''"`
}
