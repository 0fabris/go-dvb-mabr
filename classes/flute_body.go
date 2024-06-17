package classes

import (
	"bytes"

	route "github.com/0fabris/go-dvb-route"
)

func (fh *FluteHeader) Parse(b []byte) error {
	r := bytes.NewReader(b)

	// Parsing LCT Header
	header, err := route.BuildLCTHeader(r)
	if err != nil {
		return err
	}
	fh.LCT = header

	// Parsing Header Extension Types
	switch header.HET.Type {
	case route.HET_EXT_FDT:
		/*
			V := (header.HET.Content[0] & 0xF0) >> 4
			header.HET.Content[0] &= 0x0F
			var instanceId uint32 = (uint32(header.HET.Content[0])<<16 + uint32(header.HET.Content[1])<<8 + uint32(header.HET.Content[2]))
		*/
	case route.HET_EXT_CENC:
		/*
			CENC := header.HET.Content[0]
			Reserved := header.HET.Content[1:]
		*/
	case route.HET_EXT_FTI:
		if err := route.ParseFECOTIfromHeaderExtension(header.HET, &fh.FECOTI); err != nil {
			return err
		}
	}

	// Getting length in bytes (uint32 word number to uint8 byte number)
	fh.Length = uint16(fh.LCT.HeaderLen) * 4 //uint16(len(b) - r.Len()) //- fh.FECOTI.EncodingSymbolLen

	if err := route.ParseFECPayloadID(r, &fh.FECPayload, fh.FECOTI.M); err != nil {
		return err
	}
	fh.Length += 4 // due to FECPayloadID

	fh.TOI = convertToUInt64(fh.LCT.TOI)

	fh.Data = b[:fh.Length]
	return nil
}

func (fp *FlutePacket) Parse(b []byte) error {

	// Parsing Packet Header
	fp.Header = FluteHeader{}
	if err := fp.Header.Parse(b); err != nil {
		return err
	}

	// Striping header from data
	fp.Data = b[fp.Header.Length:]

	return nil
}
