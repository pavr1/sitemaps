package model

type Node struct {
	URL   string  `xml:"url"`
	Nodes []*Node `xml:"nodes"`
}
