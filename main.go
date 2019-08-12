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

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/axamon/ngasp/traps"

	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
)

// BuildVersion è la versione attuale del tool 
// valorizzato tramite:
// go  build -ldflags "-X main.BuildVersion=2"
var BuildVersion string

// Configuration contiene gli elemnti per configurare il tool.
type Configuration struct {
	Token string `json:"token"`
}

var tgu = flag.String("t", "", "TGU da controllare")

// var configfile = flag.String("c", "conf.json", "File di configurazione")
var vault = flag.String("s", "127.0.0.1:9999", "Server tokenizzatore e porta")

// var configuration Configuration

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	version := flag.Bool("v", false, "Versione dell'APP")
	flag.Parse()

	if *version {
		fmt.Printf("Ver%s\n", BuildVersion)
		os.Exit(0)
	}

	lenghthTGU := len(*tgu)

	switch {
	case lenghthTGU == 0:
		fmt.Println("Non hai passato alcuna tgu. Rilancia con -h per sintassi.")
		os.Exit(1)
	case lenghthTGU < 12:
		for i := lenghthTGU; i < 12; i++ {
			*tgu = "0" + *tgu
		}
	case lenghthTGU > 12:
		fmt.Println("TGU troppo lunga massimo 12 numeri. Rilancia con -h per sintassi.")
		os.Exit(1)
	default:
	}

	fmt.Println(*tgu)

	// Verifica che esista un token valido
	err := verificatoken(ctx, *vault)
	if err != nil {
		log.Printf("ERROR Token non verificato: %s\n", err.Error())
		os.Exit(1)
	}

	busta2 := `<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">\n    
<Body>\n        
<getTimVisionTrapsRequest xmlns="urn:AxessInterface">\n
<tgu>` + *tgu + `</tgu>\n
</getTimVisionTrapsRequest>\n    
</Body>\n
</Envelope>`

	transCfg := &http.Transport{
		// Ignora certificati SSL scaduti.
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: transCfg}
	url := "https://ngasp-ag.tim.it/live/nbi_interfaces/soap/document_literal"

	fmt.Println(url)

	body := []byte(busta2)

	// soap action
	soapAction := `"urn:getTimVisionTraps"`

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("ERROR Impossibile creare richiesta: %s\n", err.Error())
	}
	req.Header.Set("Content-Type", `text/xml; charset="UTF-8"`)
	req.Header.Set("SOAPAction", soapAction)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR Impossibile inviare richiesta http: %s\n", err.Error())
	}
	defer resp.Body.Close()

	responsBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR Impossibile leggere body reqest: %s\n", err.Error())
	}

	// fmt.Println(string(responsBody)) // Debug

	traps.Parse(ctx, responsBody)

}
