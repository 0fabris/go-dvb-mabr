package decoders

import (
	"encoding/xml"
	"fmt"

	"github.com/0fabris/go-dvb-mabr/classes"
	"github.com/0fabris/go-dvb-mabr/flute"
)

// This is the FLUTE payload decoder, the parameter is the callback function, called when a file is ready to use
func FlutePayloadDecoder(fHandler func(*classes.MABRFile) error) func([]byte) error {

	var dH = fluteGenericPacketHandler(fHandler)
	return func(payload []byte) error {
		// creating flute packet
		packet := classes.FlutePacket{}
		// parsing payload
		if err := packet.Parse(payload); err != nil {
			return err
		}
		return dH(&packet)
	}
}

func fluteGenericPacketHandler(callback func(*classes.MABRFile) error) func(*classes.FlutePacket) error {
	var tmpFDTdata []byte
	var latestFDT *classes.FluteFDTInstance = nil
	var previousFDT *classes.FluteFDTInstance = nil
	var FDTend = "</FDT-Instance>"

	var processNOPPackets = processFluteDataPacketContent(callback)
	var processFTIPackets = processFluteDataPacketContent(callback)

	return func(packet *classes.FlutePacket) error {
		switch packet.Header.LCT.HET.Type {
		case flute.HET_EXT_FDT:
			{
				tmpFDTdata = append(tmpFDTdata, packet.Data...)
				// if end is FDT-Instance close tag, start parsing
				if string(packet.Data[len(packet.Data)-len(FDTend):]) == FDTend {
					previousFDT = latestFDT
					latestFDT = &classes.FluteFDTInstance{}
					xml.Unmarshal(tmpFDTdata, latestFDT)
					tmpFDTdata = []byte{}
				}
			}
		case flute.HET_EXT_NOP:
			{
				processNOPPackets(packet, previousFDT) // for other streams
			}
		case flute.HET_EXT_FTI:
			{
				processFTIPackets(packet, latestFDT) // for Inverto DVB-I streams
			}
		default:
			fmt.Printf("Unknown packet LCT Type: %02x\n", packet.Header.LCT.HET.Type)
		}
		return nil
	}
}

func processFluteDataPacketContent(callback func(*classes.MABRFile) error) func(*classes.FlutePacket, *classes.FluteFDTInstance) {
	var packetMap = map[uint64][]byte{}
	var previousTOI uint64

	return func(packet *classes.FlutePacket, fdt *classes.FluteFDTInstance) {

		if _, ok := packetMap[packet.Header.TOI]; !ok {
			packetMap[packet.Header.TOI] = []byte{}
		}

		packetMap[packet.Header.TOI] = append(packetMap[packet.Header.TOI], packet.Data...)
		if fdt != nil && previousTOI != packet.Header.TOI {
			for _, f := range fdt.File {
				if f.TOI != previousTOI {
					continue
				}

				// if existing TOI in packetMap table
				if val, ok := packetMap[f.TOI]; ok && len(val) > 0 {
					if err := callback(&classes.MABRFile{
						Location:    f.ContentLocation,
						Content:     val,
						HTTPHeaders: nil,
					}); err != nil {
						fmt.Printf("Error during callback: %+v\n", err)
					}
					delete(packetMap, f.TOI)
				}
			}
			previousTOI = packet.Header.TOI
		}
	}
}
