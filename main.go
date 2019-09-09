// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/axamon/ngaspclient/ngasptraps"

	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
)

const endpointNGASPAPI = "https://ngasp-ag.tim.it/live/nbi_interfaces/soap/document_literal"

// BuildVersion è la versione attuale del tool
// valorizzato tramite:
// go  build -ldflags "-X main.BuildVersion=2"
var BuildVersion string

// Configuration contiene gli elemnti per configurare il tool.
// type Configuration struct {
// 	Token string `json:"token"`
// }

// var configfile = flag.String("c", "conf.json", "File di configurazione")
var vault = flag.String("s", "10.38.105.251:9999", "Server tokenizzatore e porta")

var isAllNums = regexp.MustCompile(`(?m)^\d+$`)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Crea il flag per la versione app.
	version := flag.Bool("v", false, "Versione dell'APP")
	author := flag.Bool("a", false, "autore e contatti per segnalazioni")
	debug := flag.Bool("d", false, "Salva output xml su file")

	// Parsa i flags
	flag.Parse()

	// Se il flag author è settato mostra l'autore ed esce.
	if *author {
		fmt.Println("Autore: Alberto Bregliano")
		fmt.Println("invia segnalazioni a: alberto.bregliano@telecomitalia.it")

		os.Exit(0)
	}

	// Se il flag version è settato mostra la versione ed esce.
	if *version {
		fmt.Printf("Ver%s\n", BuildVersion)
		os.Exit(0)
	}

	// Parsa i NON flags e recupera la tgu come primo argomento num: 0
	tguarg := flag.Arg(0)

	// Ripulisce la tgu da eventuali spazi alla fine ed inizio.
	tguarg = strings.TrimSpace(tguarg)

	// Se la tgu passata non contiene solo numeri esce con errore.
	if !isAllNums.Match([]byte(tguarg)) {
		log.Println("La tgu contiene caratteri non numerici")
		os.Exit(1)
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
	url := endpointNGASPAPI

	// trasforma il payload in bytes.
	body := []byte(busta2)

	// soap action è un header per comunicare quale API interrogare.
	soapAction := `"urn:getTimVisionTraps"`

	// Crea la richiesta http
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

	// Se il flag debug è settato salva output dell' API su file xml.
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
	ngasptraps.Parse(ctx, responsBody, tguarg)

}
