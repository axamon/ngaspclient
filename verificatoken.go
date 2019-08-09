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

func verificatoken(ctx context.Context, vault string) (err error) {

	// err = gonfig.GetConf(configfile, &configuration)
	// if err != nil {
	// 	log.Printf("ERROR Problema con il file di configurazione conf.json: %s\n", err.Error())
	// 	return
	// }

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
	key := []byte(createHash(passphrase))
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
