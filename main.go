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

// var configfile = flag.String("c", "conf.json", "File di configurazione")
var vault = flag.String("s", "10.38.105.251:9999", "Server tokenizzatore e porta") 


func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Crea il flag per la versione app.
	version := flag.Bool("v", false, "Versione dell'APP")
	debug := flag.Bool("d", false, "Salva output xml su file")

	// Parsa i flags
	flag.Parse()

	// Parsa i non flags e recupera la tgu come primo argomento num: 0
	tguarg := flag.Arg(0)

	// fmt.Println(tguarg)

	// Se il flag version è settato mostra la versione ed esce.
	if *version {
		fmt.Printf("Ver%s\n", BuildVersion)
		os.Exit(0)
	}

	// verifica lunghezza tgu inserita
	lenghthTGU := len(tguarg)

	switch {
	case lenghthTGU == 0:
		fmt.Println("Non hai passato alcuna tgu. Rilancia con -h per sintassi.")
		os.Exit(1)
	case lenghthTGU < 12:
		for i := lenghthTGU; i < 12; i++ {
			tguarg = "0" + tguarg
		}
	case lenghthTGU > 12:
		fmt.Println("TGU troppo lunga massimo 12 numeri. Rilancia con -h per sintassi.")
		os.Exit(1)
	default:
	}

	// fmt.Println(tguarg) // debug

	// Verifica che esista un token valido
	err := verificatoken(ctx, *vault)
	if err != nil {
		log.Printf("ERROR Token non verificato: %s\n", err.Error())
		os.Exit(1)
	}

	// busta2 è il payload da passare all'endpoint con la tgu all'interno.
	busta2 := `<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">\n    
<Body>\n        
<getTimVisionTrapsRequest xmlns="urn:AxessInterface">\n
<tgu>` + tguarg + `</tgu>\n
</getTimVisionTrapsRequest>\n    
</Body>\n
</Envelope>`

	// configura il trasporto http
	transCfg := &http.Transport{
		// Ignora certificati SSL scaduti.
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Crea un client http
	client := &http.Client{Transport: transCfg}

	// Endpoint da contattare
	url := "https://ngasp-ag.tim.it/live/nbi_interfaces/soap/document_literal"

	// trasforma il payload in bytes.
	body := []byte(busta2)

	// soap action è un header per comunicare quale API interrogare.
	soapAction := `"urn:getTimVisionTraps"`

	// Crea la richista http
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("ERROR Impossibile creare richiesta: %s\n", err.Error())
	}

	// Imposta header necessari
	req.Header.Set("Content-Type", `text/xml; charset="UTF-8"`)
	// header per comunicare in linguaggio soap l'api da interrogare.
	req.Header.Set("SOAPAction", soapAction)

	// Lancia la richiesta http e recupera la response
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR Impossibile inviare richiesta http: %s\n", err.Error())
	}
	// chiude il resp.Body come da specifica
	defer resp.Body.Close()

	// legge il responsbody in bytes
	responsBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR Impossibile leggere body reqest: %s\n", err.Error())
	}

	if *debug {
		// salva il tutto su un file di appoggio
		fmt.Println("Debug attivo")
		err = ioutil.WriteFile("debug_"+tguarg+".xml", responsBody, 0666)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(responsBody)) // Debug
	}

	// Avvia elavorazione traps
	traps.Parse(ctx, responsBody, tguarg)

}
