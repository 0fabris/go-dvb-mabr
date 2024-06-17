package classes

import (
	"encoding/xml"

	route "github.com/0fabris/go-dvb-route"
)

// Type Definition
type FlutePacket struct {
	Data   []byte // (payload)
	Header FluteHeader
}

type FluteHeader struct {
	Data       []byte
	LCT        route.LCTHeader
	FECOTI     route.FECObjectTransmissionInformation
	FECPayload route.FECPayloadID
	TOI        uint64
	Length     uint16
}

type FluteFDTInstance struct {
	XMLName       xml.Name    `xml:"FDT-Instance"`
	Text          string      `xml:",chardata"`
	Expires       int         `xml:"Expires,attr"`
	Xmlns         string      `xml:"xmlns,attr"`
	Sv            string      `xml:"sv,attr"`
	File          []FluteFile `xml:"File"`
	SchemaVersion string      `xml:"schemaVersion"`
}

type FluteFile struct {
	Text                           string `xml:",chardata"`
	ContentLength                  int    `xml:"Content-Length,attr"`
	ContentLocation                string `xml:"Content-Location,attr"`
	ContentMD5                     string `xml:"Content-MD5,attr"`
	ContentType                    string `xml:"Content-Type,attr"`
	FECOTIEncodingSymbolLength     int    `xml:"FEC-OTI-Encoding-Symbol-Length,attr"`
	FECOTIFECEncodingID            string `xml:"FEC-OTI-FEC-Encoding-ID,attr"`
	FECOTIMaximumSourceBlockLength int    `xml:"FEC-OTI-Maximum-Source-Block-Length,attr"`
	TOI                            uint64 `xml:"TOI,attr"`
	TransferLength                 string `xml:"Transfer-Length,attr"`
	Delimiter                      string `xml:"delimiter"`
}
