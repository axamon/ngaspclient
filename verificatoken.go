// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	passphrase = "vvkidtbcjujhtglivdjtlkgtetbtdejlivgukincfhdt"
)

// Dinamicsimmetricpass password dinamica negoziata col tokenizzatore
var Dinamicsimmetricpass string

func verificatoken(ctx context.Context, vault string) (err error) {

	// err = gonfig.GetConf(configfile, &configuration)
	// if err != nil {
	// 	log.Printf("ERROR Problema con il file di configurazione conf.json: %s\n", err.Error())
	// 	return
	// }

	dinamicsimmetricpassURL := "http://" + vault + "/dinamicsimmetricpass"

	resp1, err := http.Get(dinamicsimmetricpassURL)
	if err != nil {
		log.Fatalf("ERROR Impossibile contattare il Vault: %s\n", err.Error())
	}

	bodyBytes1, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		log.Fatalf("ERROR Impossibile leggere risposta del Vault: %s\n", err.Error())
	}

	Dinamicsimmetricpass = string(bodyBytes1)

	vaultURL := "http://" + vault + "/token"

	// fmt.Println(vaultURL)

	resp, err := http.Get(vaultURL)
	if err != nil {
		log.Fatalf("ERROR Impossibile contattare il Vault: %s\n", err.Error())
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ERROR Impossibile leggere risposta del Vault: %s\n", err.Error())
	}

	Token := string(bodyBytes)

	credBlob, _ := hex.DecodeString(Token)
	userEpass := string(decrypt(credBlob))
	credenziali := strings.Split(userEpass, " ")

	scadenza, err := strconv.Atoi(credenziali[0])
	if err != nil {
		log.Fatalf("ERROR Impossibile parsare scadenza del token: %s\n", err.Error())
	}
	// username = credenziali[1]
	// password = credenziali[2]

	oggi := time.Now().Unix()

	if oggi > int64(scadenza) {
		log.Println("Token scaduto. Impossibile proseguire.")
		os.Exit(1)
	}

	// restaunasettimana := time.Now().Add(7 * time.Hour * 24).Unix()

	// if restaunasettimana > int64(scadenza) {
	// 	log.Println("Token in scadeza. Si consiglia di rinnovarlo.")
	// }

	return err
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func decrypt(data []byte) []byte {
	// key := []byte(createHash(passphrase)) Dinamicsimmetricpass
	key := []byte(createHash(Dinamicsimmetricpass))
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Println(err.Error())
	}
	return plaintext
}
