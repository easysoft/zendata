package model

import "encoding/xml"

type XmlData struct {
	XMLName xml.Name `xml:"testdata"`
	Title   string   `xml:"title"`
	Table   XmlTable `xml:"table"`
}

type XmlTable struct {
	Rows []XmlRow `xml:"row"`
}

type XmlRow struct {
	XMLName xml.Name          `xml:"row"`
	Cols    map[string]string `xml:"col"`
}
