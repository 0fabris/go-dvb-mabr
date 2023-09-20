package decoders

import (
	"encoding/xml"
	"fmt"

	"github.com/0fabris/go-dvb-mabr/classes"
	"github.com/0fabris/go-dvb-mabr/headers"
	"github.com/phf/go-queue/queue"
)

// This is the FLUTE payload decoder, the parameter is the callback function, called when a file is ready to use
func FlutePayloadDecoder(fHandler func(*classes.MABRFile) error) func([]byte) error {

	var dH = fluteDataPacketHandler(fHandler)
	var cH = fluteConfigurationPacketHandler(fHandler)
	var tH = fluteTimeInfoPacketHandler(fHandler)

	return func(payload []byte) error {
		// creating flute packet
		packet := classes.FlutePacket{}

		// parsing payload
		if err := packet.Parse(payload); err != nil {
			return err
		}

		switch packet.Header.PacketType {
		case classes.FLUTE_DATA_PACKET:
			return dH(&packet)
		case classes.FLUTE_CONF_PACKET:
			return cH(&packet)
		case classes.FLUTE_TIME_PACKET:
			return tH(&packet)
		default:
		}
		return fmt.Errorf("Unknown Flute Packet Type!")
	}
}

func fluteDataPacketHandler(callback func(*classes.MABRFile) error) func(*classes.FlutePacket) error {
	var infoQ = queue.New()
	var dataQ = queue.New()
	var previousInfoSession uint32 = 0
	var previousDataSession uint32 = 0
	var latestFDT *classes.FluteFDTInstance = nil
	var packetMap = map[uint32][]byte{}

	return func(packet *classes.FlutePacket) error {
		switch packet.Header.Data[2] {
		case headers.FLUTE_DATA_HEAD_PACKET:
			if packet.Header.TOI != previousInfoSession || packet.Header.LatestPosition == 0 {
				previousInfoSession = packet.Header.TOI
				// extracting data from queue
				if infoQ.Len() > 0 {
					rawData := []byte{}
					for infoQ.Len() > 0 {
						if a := infoQ.PopFront(); a != nil {
							rawData = append(rawData, a.([]byte)...)
						}
					}
					infoQ.Init()
					latestFDT = &classes.FluteFDTInstance{}
					xml.Unmarshal(rawData, latestFDT)
				}
			}
			infoQ.PushBack(packet.Data)
		case headers.FLUTE_DATA_BODY_PACKET:
			if packet.Header.TOI != previousDataSession || packet.Header.LatestPosition == 0 {
				previousDataSession = packet.Header.TOI

				if dataQ.Len() > 0 {
					// extracting data from queue
					rawData := []byte{}
					for dataQ.Len() > 0 {
						if a := dataQ.PopFront(); a != nil {
							//fmt.Printf("a: %v\n", len(a.([]byte)))
							rawData = append(rawData, a.([]byte)...)
						}
					}
					dataQ.Init()
					// if consistent packet
					if len(rawData) > 0 {
						packetMap[packet.Header.TOI-1] = rawData
						if latestFDT != nil {
							for _, f := range latestFDT.File {
								if val, ok := packetMap[f.TOI]; ok && len(val) > 0 {
									// using val
									if err := callback(&classes.MABRFile{
										Location:    f.ContentLocation,
										Content:     val,
										HTTPHeaders: nil,
									}); err != nil {
										fmt.Printf("Error during callback: %+v\n", err)
									}
									// removing from map
									packetMap[f.TOI] = nil
								}
							}
						}
					}
				}
			}
			dataQ.PushBack(packet.Data)
		}
		return nil
	}
}

/* Work In Progress */
func fluteConfigurationPacketHandler(callback func(*classes.MABRFile) error) func(*classes.FlutePacket) error {
	return fluteDataPacketHandler(callback)
}

func fluteTimeInfoPacketHandler(callback func(*classes.MABRFile) error) func(*classes.FlutePacket) error {
	return fluteDataPacketHandler(callback)
}
