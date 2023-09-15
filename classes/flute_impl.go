package classes

import (
	"fmt"

	"github.com/0fabris/go-dvb-mabr/headers"
)

func (fh *FluteHeader) Parse(b []byte) error {
	fh.PacketType = FlutePacketType(b[0])
	var len = headers.FLUTE_BASE_HEADER_LEN
	switch fh.PacketType {
	case FLUTE_TIME_PACKET:
		len += headers.FLUTE_TIME_HEADER_SUP_LEN
		break
	case FLUTE_DATA_PACKET:
		len += headers.FLUTE_DATA_HEADER_SUP_LEN
		if b[2] == headers.FLUTE_DATA_XML_DESCRIPTOR {
			len += headers.FLUTE_DATA_XML_DESCRIPTOR_HEADER_SUP_LEN
			fh.TOI = uint32(b[len-5])
		} else if b[2] == 0x03 {
			// common data
			fh.TOI = uint32(b[10])<<8 ^ uint32(b[11])
		}
		pos := b[len-4 : len]
		fh.LatestPosition = uint64(pos[0])<<(8*3) ^ uint64(pos[1])<<(8*2) ^ uint64(pos[2])<<(8) ^ uint64(pos[3])
		break
	case FLUTE_CONF_PACKET:
		len += headers.FLUTE_CONF_HEADER_SUP_LEN
		break
	default:
	}
	fh.Data = b[:len]
	fh.Length = len
	return nil
}

func (fp *FlutePacket) Parse(b []byte) error {

	fp.Header = FluteHeader{}
	if err := fp.Header.Parse(b); err != nil {
		return err
	}
	if fp.Header.Length > len(b) {
		return fmt.Errorf("Short len of flute packet")
	}
	fp.Data = b[fp.Header.Length:]
	if len(fp.Data) < 3 {
		return fmt.Errorf("Short len of flute packet")
	}
	fp.Data = fp.Data[:len(fp.Data)] //-3]
	return nil
}
