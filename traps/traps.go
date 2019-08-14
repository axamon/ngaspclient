// Copyright (c) 2019 Alberto Bregliano
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package traps

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"log"
	"net/url"
	"time"

	"context"
	"encoding/base64"
	"encoding/xml"
	"fmt"

	"strings"
)

const rome = "VFppZjIAAAAAAAAAAAAAAAAAAAAAAAAHAAAABwAAAAAAAACsAAAABwAAAA2AAAAAmzj4cJvVzOCcxcvwnbcAYJ6J/nCfoBzgoGCl8KF+rWCiXDdwo0waYMhsNfDM50sQzakXkM6CdODOokMQz5I0EM/jxuDQbl6Q0XIWENJM0vDTPjGQ1EnSENUd93DWKZfw1uuAkNgJlhD5M7Xw+dnE4Psc0nD7ubTw/Py0cP2ZlvD+5dDw/4KzcADFsvABYpVwApxacANCd3AEhXbwBSuT8AZuk3AHC3XwCEU68AjrV/AKLldwCss58AwOOXAMqxvwDeTg8A6K/fAPzf1wEHQacBGt33ASU/xwEs6X8BNNRBAUM/qQFSPrkBYT3JAXA82QF/O+kBjjr5AZ06CQGsORkBu8vRAcrK4QHZyfEB6MkBAffIEQIGxyECFcYxAiTFQQIzxFECQsNhAlHCcQJgwYECcFQ5An9TSQKOUlkCnVFpAqxQeQK7T4kCyk6ZAtlNqQLoTLkC90vJAwZK2QMV3ZEDJytBAzPbsQNFKWEDUdnRA2MngQNv1/EDgblJA43WEQOft2kDq9QxA721iQPKZfkD27OpA+hkGQP5sckEBmI5BBhDkQQkYFkENkGxBEJeeQRUP9EEYFyZBHI98QR+7mEEkDwRBJzsgQSuOjEEuuqhBMzL+QTY6MEE6soZBPbm4QUIyDkFFXipBSbGWQUzdskFRMR5BVF06QViwpkFb3MJBYFUYQWNcSkFn1KBBatvSQW9UKEFygERBdtOwQXn/zEF+UzhBgX9UQYX3qkGI/txBjXcyQZB+ZEGU9rpBmCLWQZx2QkGfol5Bo/XKQach5kGrdVJBrqFuQbMZxEG2IPZBuplMQb2gfkHCGNRBxUTwQcmYXEHMxHhB0RfkQdREAEHYvFZB28OIQeA73kHjQxBB57tmQerCmEHvOu5B8mcKQfa6dkH55pJB/jn+QAgECAQIBAgECAQIBAwQBAwQBAwECBAMEAwQDBAIEAwQDBAMEAwQDBAMEAwQDBAMEAwQDBAMEAwIFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgAAC7QAAAAAHCABBAAADhAACQAADhAACQAAHCABBAAAHCABBAAADhAACUxNVABDRVNUAENFVAAAAAABAQEBAAAAAAABAQpDRVQtMUNFU1QsTTMuNS4wLE0xMC41LjAvMwo="

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

// Output è il risultato mostrato a video
type Output struct {
	CpeID         string
	Mode          string
	ModelName     string
	StartTS       string
	EndTS         string
	ChiusoDa      string
	VideoTitle    string
	QoV           string
	NetworkType   string
	IsAlice       string
	Buffering     int
	LongBuffering int
	PlayerError   int
	StreamingType string
	TGU           string
	FQDN          string
}

// Parse parsa le trap di risposta dell'API
func Parse(ctx context.Context, response []byte, tgu string) {

	// file, err := os.Open("result.xml")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// data, err := ioutil.ReadAll(response)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// var elementi XMLResponse
	// errxml := xml.Unmarshal(response, &elementi)
	// if errxml != nil {
	// 	log.Fatal(errxml.Error())
	// }

	// fields := elementi.Body.GetTimVisionTrapsResponseResponse.Return.DeviceRows.Item.TrapsFields
	// fieldsSlice := strings.Split(fields, ",")
	// for n, field := range fieldsSlice {
	// 	fmt.Println(n, field)
	// }
	// time.Sleep(5*time.Second)
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

	var fruizioni = make(map[string]Output)

	var archivio []string
	var indiceVideoTitle, indiceEventName, indiceNetworkType, indiceEventType, indiceProvider, indiceVideoURL, indiceStreamingType int
	var fieldsSlice []string

	scanner := bufio.NewScanner(bytes.NewReader(response))
	for scanner.Scan() {
		line := scanner.Text()
		// fmt.Println(line) // debug

		// Se la linea è di tipo item
		// la mette in archivio
		if strings.HasPrefix(line, "<item>") {

			// elimina il tag iniziale
			trap := strings.ReplaceAll(line, "<item>", "")
			// fmt.Println(trap)

			// Metti trap in archivio
			archivio = append(archivio, trap)
			// for n, v := range archivio {
			// 	fmt.Println(n, v)
			// }

		}

		// Se la linea è di tipo traps_fields
		// bisogna elaborare i nomi campi e le traps raccolte
		// finora in archivio salvando gli elaborati in fruzioni.
		if strings.HasPrefix(line, "<traps_fields>") {
			fieldsSlice = strings.Split(line, ",")

			// Trovo numero campo associato a quelli di interesse
			indiceVideoTitle = trovaIndice(fieldsSlice, "trap.body.videoTitle")
			indiceEventName = trovaIndice(fieldsSlice, "trap.body.eventName")
			indiceNetworkType = trovaIndice(fieldsSlice, "trap.networkType")
			indiceEventType = trovaIndice(fieldsSlice, "trap.eventType")
			indiceProvider = trovaIndice(fieldsSlice, "provider")
			indiceVideoURL = trovaIndice(fieldsSlice, "trap.body.videoUrl")
			indiceStreamingType = trovaIndice(fieldsSlice, "trap.body.streamingType")

			// aggiungere altri se servono

			// Tratto le traps
			for _, trap := range archivio {

				// Crea un Out vuoto nuovo per ogni trap
				var out Output

				// Splitta la trap nei suoi elementi creando una slice
				trapSlice := strings.Split(trap, ",")
				// for n, v := range trapSlice {
				// 	fmt.Println(n, v)
				// }

				//Salta le trap vuote
				if len(trapSlice) < 2 {
					continue
				}

				// Creo hash per archiaviare le fruizioni
				hash := md5.New()
				// Come valori per l'hash aggiungo cpeid e videotitle
				hash.Write([]byte(trapSlice[1] + trapSlice[2] + trapSlice[indiceVideoTitle]))
				// nomehash è l'has in versione esadecimale
				nomehash := fmt.Sprintf("%x", hash.Sum(nil))

				// Se l'hash esiste già recupera il valore val dalla mappa
				// e lo assegna ad out.
				if val, ok := fruizioni[nomehash]; ok {
					out = val
				}

				out.CpeID = trapSlice[1]
				out.Mode = trapSlice[3]
				out.ModelName = trapSlice[4]
				out.VideoTitle = trapSlice[indiceVideoTitle]
				out.NetworkType = trapSlice[indiceNetworkType]
				out.IsAlice = trapSlice[indiceProvider]
				out.TGU = tgu
				out.StreamingType = trapSlice[indiceStreamingType]

				// elabora il campo url
				u, err := url.Parse(trapSlice[indiceVideoURL])
				if err != nil {
					log.Println(err)
				}

				out.FQDN = u.Hostname()

				// Trova la location per trattare correttamente il passaggio
				// del fuso orario da UTC e per gestire l'ora solare/legale.
				romebytes, err := base64.StdEncoding.DecodeString(rome)
				if err != nil {
					log.Println("impossibile decodificare roma")
				}
				// location, err := tizzy.LoadLocationValue("Europe/Rome")
				location, err := time.LoadLocationFromTZData("Europe/Rome", romebytes)
				if err != nil {
					log.Println("errore con la location per i fusiorari")
				}

				// Elabora Inzio e fine fruizione
				switch trapSlice[indiceEventName] {
				case "PLAY":
					starTS, err := time.Parse("2006-01-02T15:04:05Z", trapSlice[0])
					if err != nil {
						log.Println(err.Error())
					}
					out.StartTS = starTS.In(location).Format("2006-01-02T15:04:05-0700")
				case "STOP":
					endTS, err := time.Parse("2006-01-02T15:04:05Z", trapSlice[0])
					if err != nil {
						log.Println(err.Error())
					}
					out.EndTS = endTS.In(location).Format("2006-01-02T15:04:05-0700")
				}

				// Conto buffering ed errori
				switch trapSlice[indiceEventType] {
				case "buffering":
					out.Buffering++

				case "playerError":
					out.PlayerError++
				}

				// Salva out modificato nella mappa fruizioni
				fruizioni[nomehash] = out

			}

			// Finita l'elaborazione delle trap in archivio
			// si può ripulire l'archivio e i nomi campi
			// per prepararli a un nuovo ciclo.
			// Le trap elaborate sono archiviate in fruizioni.
			archivio = []string{}
			fieldsSlice = []string{}
		}
		// Riprende il ciclo iniziale alla ricerca di altre
		// linee con tag iniziale <item>
	}

	// Costruisce una lista delle chiavi hash della mappa fruizioni.
	chiaviFruizioni := make([]string, 0, len(fruizioni))
	for key := range fruizioni {
		chiaviFruizioni = append(chiaviFruizioni, key)
	}

	// fmt.Println(len(fruizioni))
	// fmt.Println(fruizioni)

	for _, key := range chiaviFruizioni {

		// Non mostra le fruizioni senza titolo
		if fruizioni[key].VideoTitle == "" {
			continue
		}
		fmt.Println(fruizioni[key])

	}

	// for _, e := range archivio {
	// 	fmt.Println(e)
	// }

	// Avvia il salvataggio su file delle fruizioni.
	salvaXLSX(ctx, tgu, fruizioni)
	return
}

func trovaIndice(slice []string, item string) int {
	// fmt.Println(slice, item)
	// time.Sleep(3 * time.Second)
	var indice int
	for n, s := range slice {
		// fmt.Println(n, s, item)
		// time.Sleep(3 * time.Second)
		if s == item {
			indice = n
		}
	}

	return indice
}

// 0633429477
// 0633652113
// 0633429599
