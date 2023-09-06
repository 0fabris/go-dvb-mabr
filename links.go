package mabr

import (
	"github.com/0fabris/go-dvb-mabr/classes"
	"github.com/0fabris/go-dvb-mabr/decoders"
)

func NCDNPayloadDecoder(callback func(*classes.MABRFile) error) func([]byte) error {
	return decoders.NCDNPayloadDecoder(callback)
}

func FlutePayloadDecoder(callback func(*classes.MABRFile) error) func([]byte) error {
	return decoders.FlutePayloadDecoder(callback)
}
