package traps

import (
	"context"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// XMLResponse è la trasposizione in struct della struttura XML di risposta.
type XMLResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	SOAPENC string   `xml:"SOAP-ENC,attr"`
	SOAPENV string   `xml:"SOAP-ENV,attr"`
	Xsd3    string   `xml:"xsd3,attr"`
	Xsi3    string   `xml:"xsi3,attr"`
	Body    struct {
		Text                              string `xml:",chardata"`
		Xmlns                             string `xml:"xmlns,attr"`
		GetTimVisionTrapsResponseResponse struct {
			Text   string `xml:",chardata"`
			Return struct {
				Text       string `xml:",chardata"`
				DeviceRows struct {
					Text string `xml:",chardata"`
					Item struct {
						Text           string `xml:",chardata"`
						Cpeid          string `xml:"cpeid"`
						DeviceTypeName string `xml:"device_type_name"`
						FwName         string `xml:"fw_name"`
						Traps          struct {
							Text string   `xml:",chardata"`
							Item []string `xml:"item"`
						} `xml:"traps"`
						TrapsFields string `xml:"traps_fields"`
						Vendor      string `xml:"vendor"`
					} `xml:"item"`
				} `xml:"device_rows"`
				Result struct {
					Text      string `xml:",chardata"`
					ErrorCode string `xml:"errorCode"`
					ErrorDesc string `xml:"errorDesc"`
				} `xml:"result"`
			} `xml:"return"`
		} `xml:"getTimVisionTrapsResponseResponse"`
	} `xml:"Body"`
}

type Output struct {
	CpeID         string
	Mode          string
	StartTS       string
	EndTS         string
	ChiusoDa      string
	VideoTitle    string
	QoV           string
	NetworkType   string
	IsAlice       string
	Buffering     string
	LongBuffering int
	PlayerError   string
	StreamingType string
	TGU           string
	FQDN          string
}

// Parse parsa le trap di risposta dell'API
func Parse(ctx context.Context, response []byte) {

	// file, err := os.Open("result.xml")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// data, err := ioutil.ReadAll(response)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	var elementi XMLResponse
	errxml := xml.Unmarshal(response, &elementi)
	if errxml != nil {
		log.Fatal(errxml.Error())
	}

	// fields := elementi.Body.GetTimVisionTrapsResponseResponse.Return.DeviceRows.Item.TrapsFields
	// fieldsSlice := strings.Split(fields, ",")
	// for n, field := range fieldsSlice {
	// 	fmt.Println(n, field)
	// }
	/*
	   0 creation_time
	   1 deviceId
	   2 deviceType
	   3 mode
	   4 modelName
	   5 originIPAddress
	   6 trap.orig_timestamp
	   7 trap.bitRateKbps
	   8 trap.body.averageBitrate
	   9 trap.body.avgSSKbps
	   10 trap.body.bufferingDuration
	   11 trap.body.callerClass
	   12 trap.body.callErrorCode
	   13 trap.body.callErrorMessage
	   14 trap.body.callErrorType
	   15 trap.body.callUrl
	   16 trap.body.errorDesc
	   17 trap.body.errorReason
	   18 trap.body.eventName
	   19 trap.body.fwUpgradeVersion
	   20 trap.body.levelBitrates
	   21 trap.body.lineSpeedKbps
	   22 trap.body.loUpgradeVersion
	   23 trap.body.maxSSChunkKbps
	   24 trap.body.maxSSKbps
	   25 trap.body.minSSKbps
	   26 trap.body.playingInterval
	   27 trap.body.raUpgradeVersion
	   28 trap.body.resType
	   29 trap.body.streamingType
	   30 trap.body.timeValue
	   31 trap.body.videoDuration
	   32 trap.body.videoPosition
	   33 trap.body.videoTitle
	   34 trap.body.videoType
	   35 trap.body.videoUrl
	   36 trap.eventType
	   37 trap.fwVersion
	   38 trap.loVersion
	   39 trap.networkType
	   40 trap.raVersion
	   41 trap.SNR
	   42 trap.timestamp
	   43 provider

	*/

	// cpeid := elementi.Body.GetTimVisionTrapsResponseResponse.Return.DeviceRows.Item.Cpeid
	traps := elementi.Body.GetTimVisionTrapsResponseResponse.Return.DeviceRows.Item.Traps.Item

	location, err := time.LoadLocation("Europe/Rome")
	if err != nil {
		panic(err)
	}

	var fruizioni = make(map[string]Output)

	var out Output
	for _, trap := range traps {

		trapSlice := strings.Split(trap, ",")
		fmt.Println(trapSlice) // Debug

		// Creo hash per archiaviare le fruizioni
		hash := md5.New()
		// Come valori per l'hash aggiungo cpeid ip e videotitle
		hash.Write([]byte(trapSlice[1] + trapSlice[5] + trapSlice[33]))
		nomehash := fmt.Sprintf("%x", hash.Sum(nil))

		out.CpeID = trapSlice[1]
		out.Mode = trapSlice[3]
		out.VideoTitle = trapSlice[33]
		out.NetworkType = trapSlice[39]
		starTS, err := time.Parse("2006/01/02_15:04:05", trapSlice[6])
		if err != nil {
			log.Println(err.Error())
		}
		out.StartTS = starTS.In(location).Format("2006-01-02T15:04:05-0700")
		if trapSlice[18] == "STOP" {
			endTS, err := time.Parse("2006/01/02_15:04:05", trapSlice[6])
			if err != nil {
				log.Println(err.Error())
			}
			out.EndTS = endTS.In(location).Format("2006-01-02T15:04:05-0700")
		}

		if trapSlice[18] == "buffering" {
			out.Buffering = "Sì"
		} else {
			out.Buffering = "No"
		}

		if trapSlice[10] != "" {
			v, err := strconv.Atoi(trapSlice[10])
			if err != nil {
				log.Printf("ERROR Impossibile estrarre tempo di buffering: %s", err.Error())
			}
			out.LongBuffering = out.LongBuffering + v
		}

		fruizioni[nomehash] = out

		fmt.Println(starTS, out.StartTS)
		fmt.Println(out)
	}

	fmt.Println(len(fruizioni), fruizioni)

	return
}
