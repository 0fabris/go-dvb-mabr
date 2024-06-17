package flute

import (
	"io"
)

// from RFC 5651
type CodePoint uint8

const (
	CP_RES                            CodePoint = 0x00
	CP_NRT_FILE                       CodePoint = 0x01
	CP_NRT_ENTITY                     CodePoint = 0x02
	CP_NRT_UNSIGNED_PACKAGE           CodePoint = 0x03
	CP_NRT_SIGNED_PACKAGE             CodePoint = 0x04
	CP_NEWIS_TL_CHANGED               CodePoint = 0x05
	CP_NEWIS_TL_CONTINUED             CodePoint = 0x06
	CP_REDUNDANTIS                    CodePoint = 0x07
	CP_MS_FILE                        CodePoint = 0x08
	CP_MS_ENTITY                      CodePoint = 0x09
	CP_MS_FILE_CMAF_RAND_ACCESS_CHUNK CodePoint = 0x0A
)

type HeaderExtensionType uint8

const (
	HET_EXT_NOP  HeaderExtensionType = 0x00
	HET_EXT_AUTH HeaderExtensionType = 0x01
	HET_EXT_TIME HeaderExtensionType = 0x02
	HET_EXT_FTI  HeaderExtensionType = 0x40
	HET_EXT_FDT  HeaderExtensionType = 0xC0
	HET_EXT_CENC HeaderExtensionType = 0xC1
)

type HeaderExtension struct {
	Type    HeaderExtensionType
	Length  uint8
	Content []uint8
}

// from RFC 5651
type LCTHeader struct {
	V         uint8     // 4 bits
	CCf       uint8     // 2 bits, C field
	PSI       uint8     // 2 bits
	TSIf      uint8     // 1 bit, S field
	TOIf      uint8     // 2 bits, O field
	HWf       uint8     // 1 bit, H field
	Reserved  uint8     // 2 bits,
	CSf       uint8     // 1 bit
	COf       uint8     // 1 bit
	HeaderLen uint8     // 8 bits // Total length of the LCT header in units of 32-bit words.
	CP        CodePoint // 8 bits
	CCI       []uint8   //[4 * 4]uint8 //[4]uint32 // size is 32*(C+1)
	TSI       []uint8   //[3 * 2]uint8 //[3]uint16 // size is 32*S + 16*H
	TOI       []uint8   //[7 * 2]uint8 //[7]uint16 // size is 32*O + 16*H
	HET       HeaderExtension
}

func BuildLCTHeader(r io.Reader) (LCTHeader, error) {
	var buf []byte
	var header LCTHeader
	buf = make([]byte, 4)

	if _, err := r.Read(buf); err != nil {
		return header, err
	}

	// parsing
	header = LCTHeader{
		V:         (buf[0] & 0xF0) >> 4,
		CCf:       (buf[0] & 0x0C) >> 2,
		PSI:       (buf[0] & 0x03),
		TSIf:      (buf[1] & 0x80) >> (3 + 4),
		TOIf:      (buf[1] & 0x60) >> (1 + 4),
		HWf:       (buf[1] & 0x10) >> 4,
		Reserved:  (buf[1] & 0x0C) >> 2,
		CSf:       (buf[1] & 0x02) >> 1,
		COf:       (buf[1] & 0x01),
		HeaderLen: (buf[2]),
		CP:        CodePoint(buf[3]),
		HET:       HeaderExtension{},
	}

	// Parsing CCI
	buf = make([]byte, 4*(header.CCf+1))
	if _, err := r.Read(buf); err != nil {
		return header, err
	}
	header.CCI = buf

	// Parsing TSI
	buf = make([]byte, 4*header.TSIf+2*header.HWf)
	if _, err := r.Read(buf); err != nil {
		return header, err
	}
	header.TSI = buf

	// Parsing TOI
	buf = make([]byte, 4*header.TOIf+2*header.HWf)
	if _, err := r.Read(buf); err != nil {
		return header, err
	}

	header.TOI = buf

	return header, parseHeaderExtension(r, &header.HET)
}

// Extension Headers parsing
func parseHeaderExtension(r io.Reader, HET *HeaderExtension) error {
	buf := make([]byte, 1)
	if _, err := r.Read(buf); err != nil {
		return err
	}
	HET.Type = HeaderExtensionType(buf[0])

	if HET.Type == HET_EXT_NOP {
		// Just stopping here, as RFC says
		return nil
	}

	var toRead uint8 = 3
	if HET.Type < 0x80 {
		// HET.Length is variable
		// Read 1
		if _, err := r.Read(buf); err != nil {
			return err
		}
		HET.Length = buf[0] // [32bit-words]

		// Read max (2^8)-1 Bytes
		toRead = (4 * HET.Length) - 2 // [bytes]
	} else {
		// Read 3 (HET.Length has fixed size)
		HET.Length = toRead
	}

	// reading HET content based on HET.Type
	buf = make([]byte, toRead)
	if _, err := r.Read(buf); err != nil {
		return err
	}
	HET.Content = buf
	return nil
}
