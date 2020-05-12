package model

import "encoding/xml"

type TestData struct {
	XMLName    xml.Name `xml:"testdata"`
	Title string `json:"title" xml:"title"`
	Table  Table `json:"table" xml:"table"`
}

type Table struct {
	Rows  []Row `json:"row" xml:"row"`
}

type Row struct {
	Cols  []string `json:"col" xml:"col"`
}