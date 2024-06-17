package flute

import (
	"fmt"
	"io"
)

type SourceFEC uint32

type FECObjectTransmissionInformation struct {
	//EncodingID uint8 //FEC-OTI-FEC-Encoding-ID
	TransferLength        uint64 //FEC-OTI-Transfer-Length (L)
	M                     uint8
	G                     uint8
	EncodingSymbolLen     uint16 //FEC-OTI-Encoding-Symbol-Length (E)
	MaxSourceBlockLen     uint16 //FEC-OTI-Maximum-Source-Block-Length (B)
	MaxNumEncodingSymbols uint16 //FEC-OTI-Max-Number-of-Encoding-Symbols (max_n)
	//FEC-OTI-Scheme-Specific-Info
}

func ParseFECOTIfromHeaderExtension(he HeaderExtension, fec *FECObjectTransmissionInformation) error {
	if he.Type != HET_EXT_FTI {
		return fmt.Errorf("Header Extension is not FEC OTI")
	}

	fec.TransferLength = 0
	for k, v := range he.Content[:6] {
		fec.TransferLength += uint64(v) << (8 * (5 - k))
	}

	if he.Length == 3 {
		fec.M = 0
		fec.G = 0
		fec.EncodingSymbolLen = (uint16(he.Content[6]) << 8) + uint16(he.Content[7])
		fec.MaxSourceBlockLen = uint16(he.Content[8])
		fec.MaxNumEncodingSymbols = uint16(he.Content[9])
	} else if he.Length == 4 {
		fec.M = he.Content[6]
		fec.G = he.Content[7]
		fec.EncodingSymbolLen = (uint16(he.Content[8]) << 8) + uint16(he.Content[9])
		fec.MaxSourceBlockLen = (uint16(he.Content[10]) << 8) + uint16(he.Content[11])
		fec.MaxNumEncodingSymbols = (uint16(he.Content[12]) << 8) + uint16(he.Content[13])
	} else {
		return fmt.Errorf("Expectiong length 3 or 4")
	}
	return nil
}

type FECPayloadID struct {
	SourceBlockNumber uint32
	EncodingSymbolID  uint32
}

func ParseFECPayloadID(r io.Reader, p *FECPayloadID, m uint8) error {
	buf := make([]byte, 4)
	if _, err := r.Read(buf); err != nil {
		return err
	}

	var common uint32 = uint32(buf[0])<<24 + uint32(buf[1])<<16 + uint32(buf[2])<<8 + uint32(buf[3])
	var mask = build32bitLSBMask(int(m))

	p.SourceBlockNumber = (common & ^mask) >> m
	p.EncodingSymbolID = common & mask

	return nil
}
