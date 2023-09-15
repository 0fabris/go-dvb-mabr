package decoders

import (
	"encoding/binary"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/0fabris/go-dvb-mabr/classes"
	"github.com/0fabris/go-dvb-mabr/headers"
	"github.com/seancfoley/ipaddress-go/ipaddr"
)

const addrLen = 4

func calcAddress(bts []byte) uint32 {
	return binary.BigEndian.Uint32(bts[len(bts)-addrLen:]) // 4 bytes of addr.
}

// This is the NCDN payload decoder, the parameter is the callback function, called when a file is ready to use
func NCDNPayloadDecoder(callback func(*classes.MABRFile) error) func([]byte) error {
	sharedDataMap := &map[string][]byte{}
	sharedHeaderMap := &map[string][]byte{}
	var latestFilename *string = nil
	var latestBlock []byte = []byte{0x00, 0x00}
	return func(payload []byte) error {
		if len(payload) <= 16 {
			return nil
		}

		if slices.Equal(payload[:len(headers.NCDN_HEADER_HEAD_PACKET)], headers.NCDN_HEADER_HEAD_PACKET) {
			// received a header packet
			infoPacket, err := nCDNInfoPacketBuilder(payload, headers.NCDN_HEADER_LEN_HPACK)
			if err != nil {
				return err
			}
			if infoPacket == nil {
				return err
			}

			if !infoPacket.IsRefreshingFilename && infoPacket.Filename != nil {
				latestFilename = infoPacket.Filename
				/*
					// removing get parameters
					if strings.Contains(*latestFilename, "?") {
						fsegm := strings.Split(*infoPacket.Filename, "?")
						latestFilename = &fsegm[0]
					}
				*/
			}

			if sharedDataMap != nil && !(slices.Equal(infoPacket.GetBlock(), latestBlock)) {
				// files are ready
				for k, v := range *sharedDataMap {
					var headers []byte = nil
					if hds, ok := (*sharedHeaderMap)[k]; ok {
						headers = hds
					}
					if err := callback(&classes.MABRFile{
						Location:    k,
						Content:     v,
						HTTPHeaders: headers,
					}); err != nil {
						return err
					}
				}

				// reinitializing maps
				sharedDataMap = &map[string][]byte{}
				sharedHeaderMap = &map[string][]byte{}

				latestBlock = infoPacket.GetBlock()
			}
		} else if slices.Equal(payload[:len(headers.NCDN_HEADER_DATA_PACKET)], headers.NCDN_HEADER_DATA_PACKET) {
			// received a data packet
			dataPacket, err := nCDNDataPacketBuilder(payload, headers.NCDN_HEADER_LEN_DPACK)
			if err != nil {
				return err
			}
			if dataPacket == nil {
				return nil
			}

			if latestFilename == nil {
				return nil
			}
			if dataPacket.IsHTTPHeaders {
				(*sharedHeaderMap)[*latestFilename] = dataPacket.GetData()
			}

			if dataPacket.IsNewFile {
				(*sharedDataMap)[*latestFilename] = dataPacket.GetData()
			} else {
				(*sharedDataMap)[*latestFilename] = append((*sharedDataMap)[*latestFilename], dataPacket.GetData()...)
			}
		}
		return nil
	}
}

func nCDNDataPacketBuilder(packet []byte, shift int) (*classes.NCDNDataPacket, error) {
	var msgHeader = packet[:shift]
	var msgData = packet[shift:]
	if packet[13] == 0x05 {
		return &classes.NCDNDataPacket{
			RawHeader:         msgHeader,
			RawData:           msgData,
			StartAddress:      0,
			IsNewFile:         false,
			IsManifestRelated: packet[15] == 0x02,
			IsHTTPHeaders:     calcAddress(packet[shift-addrLen:shift]) > 0xFFFF && packet[shift-2]+packet[shift-1] == 0,
		}, nil
	} else if packet[13] == 0x03 {
		// Receiving file chunks
		// parsing start offset/address of content
		startAddr := calcAddress(packet[shift-addrLen : shift])
		return &classes.NCDNDataPacket{
			RawHeader:    msgHeader,
			RawData:      msgData,
			StartAddress: startAddr,
			IsNewFile:    startAddr == 0,
		}, nil

	}
	return nil, nil
}

func nCDNInfoPacketBuilder(packet []byte, shift int) (*classes.NCDNInfoPacket, error) {
	var msgHeader = packet[:shift]
	var msgData = packet[shift:]
	if packet[13] == 0x01 {
		var destName *string = nil
		if packet[15] != 0x05 {
			// Received first chunk of file
			dn := strings.Trim(string(packet[shift:shift+headers.NCDN_FILENAME_MAX_LEN]), "\x00")
			destName = &dn
		}
		return &classes.NCDNInfoPacket{
			RawHeader:            msgHeader,
			RawData:              msgData,
			Filename:             destName,
			IsRefreshingFilename: false,
		}, nil
	} else if packet[13] == 0x02 {
		//Refreshing segment name == following parts of segment
		destName := strings.Trim(string(packet[shift:shift+headers.NCDN_FILENAME_MAX_LEN]), "\x00")
		return &classes.NCDNInfoPacket{
			RawHeader:            msgHeader,
			RawData:              msgData,
			Filename:             &destName,
			IsRefreshingFilename: true,
		}, nil
	}
	return nil, nil
}

// This function gives you the stream details from the catalog csv
func ParseNCDNInfoRow(row string) classes.NCDNStream {
	ret := classes.NCDNStream{}

	urlSplitPos := strings.Index(row, headers.NCDN_INFO_DATA_SEP)
	ret.URL = row[:urlSplitPos]
	ret.BroadpeakData = row[urlSplitPos+1:]
	ret.BroadpeakData = strings.ReplaceAll(ret.BroadpeakData, headers.NCDN_INFO_DATA_SEP, ",")
	ret.BroadpeakData = strings.ReplaceAll(ret.BroadpeakData, "+", "%2B")
	bpkData, err := url.ParseQuery(ret.BroadpeakData)
	if err != nil {
		return ret
	}
	ret.ServiceType = bpkData.Get("st")
	ret.ServiceID = bpkData.Get("sri")
	ret.RFUn = bpkData.Get("rfun")
	ret.DataSpeed = bpkData.Get("dspd")

	// Parsing VideoStreams
	videoBaseIP := ipaddr.NewIPAddressString(bpkData.Get("mi"))
	ret.VideoPort, _ = strconv.ParseUint(bpkData.Get("mp"), 0, 32)
	ret.VideoStreams = map[string]string{}

	if videos := bpkData.Get("lsv"); videos != "" {
		for i, v := range strings.Split(videos, ",") {
			// building ip
			addr := videoBaseIP.GetAddress().Increment(int64(i))
			ret.VideoStreams[fmt.Sprintf("%s:%d", addr.ToCanonicalString(), ret.VideoPort)] = v
		}
	} else if videos := bpkData.Get("plav"); videos != "" {
		addr := videoBaseIP.GetAddress()
		for _, v := range strings.Split(videos, ",") {
			// building ip
			index, _ := strconv.ParseUint(v, 0, 32)
			ret.VideoStreams[fmt.Sprintf("%s:%d", addr.ToCanonicalString(), ret.VideoPort+index-1)] = v
		}
	}

	// TODO parse plv, lmv

	// Parsing AudioStreams
	audioBaseIP := ipaddr.NewIPAddressString(bpkData.Get("mia"))
	ret.AudioPort, _ = strconv.ParseUint(bpkData.Get("mpa"), 0, 32)
	ret.AudioStreams = map[string]string{}

	if audios := bpkData.Get("lsa"); audios != "" {
		for i, v := range strings.Split(audios, ",") {
			// building ip
			addr := audioBaseIP.GetAddress().Increment(int64(i))
			ret.AudioStreams[fmt.Sprintf("%s:%d", addr.ToCanonicalString(), ret.AudioPort)] = v
		}
	} else if audios := bpkData.Get("plaa"); audios != "" {
		addr := audioBaseIP.GetAddress()
		for _, v := range strings.Split(audios, ",") {
			// building ip
			index, _ := strconv.ParseUint(v, 0, 32)
			ret.AudioStreams[fmt.Sprintf("%s:%d", addr.ToCanonicalString(), ret.AudioPort+index-1)] = v
		}
	}

	// TODO handle pla, lma

	// Parsing DataStreams
	dataBaseIP := ipaddr.NewIPAddressString(bpkData.Get("mid"))
	if bpkData.Get("mid") == "" {
		dataBaseIP = ipaddr.NewIPAddressString(bpkData.Get("mi"))
	}
	ret.DataPort, _ = strconv.ParseUint(bpkData.Get("mpd"), 0, 32)
	ret.DataStreams = map[string]string{}

	if ret.DataPort > 0 {

		if datas := bpkData.Get("lsd"); datas != "" {
			for i, v := range strings.Split(datas, ",") {
				// building ip
				addr := dataBaseIP.GetAddress().Increment(int64(i))
				ret.DataStreams[fmt.Sprintf("%s:%d", addr.ToCanonicalString(), ret.DataPort)] = v
			}
		} else if datas := bpkData.Get("plad"); datas != "" {
			addr := dataBaseIP.GetAddress()
			for _, v := range strings.Split(datas, ",") {
				// building ip
				index, _ := strconv.ParseUint(v, 0, 32)
				ret.DataStreams[fmt.Sprintf("%s:%d", addr.ToCanonicalString(), ret.DataPort+index-1)] = v
			}
		}
	}
	return ret
}
