GO DVB-MABR decoders
---

This library helps you handling a FLUTE (File Delivery over Unidirectional Transport) or a common nCDN data transfer over UDP/DVB-x

This library was built reverse-engineering the udp-multicast payloads from Hotbird 10949V (PID 2011 for Infos, 2012 for Data) and 12073V (PIDs 1011 for Infos, 1012 for Data).

All rights belong to their respective owners

Example: 
```go
package main

import (
    "fmt"
    mabrc "github.com/0fabris/go-dvb-mabr/classes"
	mabr "github.com/0fabris/go-dvb-mabr"
)

func yourCallbackFunction(f *mabrc.MABRFile) error {
    // do something with your file
    // f.Content, f.HTTPHeaders, f.Location
    return nil
}

func main(){

    // Decoding FLUTE payloads
    fluteDecoder := mabr.FlutePayloadDecoder(yourCallbackFunction)

    // ... or Decoding nCDNPayloads
    nCDNDecoder := mabr.NCDNPayloadDecoder(yourCallbackFunction)

    for {
        // in a stream
        packet := []byte{/*...*/}

        /* filter by multicast group */

        udpPayload := []byte{/*...*/}

        if err := fluteDecoder(udpPayload); err != nil{
            fmt.Println(err.Error())
        }
        if err := nCDNDecoder(udpPayload); err != nil{
            fmt.Println(err.Error())
        }
    }

}

```
