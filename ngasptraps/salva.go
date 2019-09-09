// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ngasptraps

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/tealeg/xlsx"
)

func salvaXLSX(ctx context.Context, tgu string, fruizioni map[string]Output, trapsRAW []TrapRAW) error {

	// file xlsx dove salvare i risultati
	var file *xlsx.File

	var err error

	// file nuovo xlsx dove salvare i risultati
	file = xlsx.NewFile()

	// trova gli hash chiave delle fruizioni
	chiaviFruizioni := make([]string, 0, len(fruizioni))
	for key := range fruizioni {
		chiaviFruizioni = append(chiaviFruizioni, key)
	}

	// Ordina la slice alfabeticamente per avere export consistenti e non
	// randomici come accadrebbe senza sort.
	sort.Strings(chiaviFruizioni)

	// fmt.Println(len(fruizioni))
	// fmt.Println(fruizioni)

	fruizione, err := file.AddSheet("Fruzioni_" + tgu)
	if err != nil {
		fmt.Printf(err.Error())
	}

	var nRow int

	// scrive sulla riga 0 la prima
	fruizione.Cell(nRow, 0).Value = "TGU"
	fruizione.Cell(nRow, 1).Value = tgu

	nRow++
	var nCol int
	// scrive sulla riga 1 la seconda
	fruizione.Cell(nRow, nCol).Value = "Cpeid"
	nCol++
	fruizione.Cell(nRow, nCol).Value = "Mode"
	nCol++
	fruizione.Cell(nRow, nCol).Value = "ModelName"
	nCol++
	fruizione.Cell(nRow, nCol).Value = "Start"
	nCol++
	fruizione.Cell(nRow, nCol).Value = "End"
	nCol++
	fruizione.Cell(nRow, nCol).Value = "VideoTitle"
	nCol++
	fruizione.Cell(nRow, nCol).Value = "Buffering"
	nCol++
	fruizione.Cell(nRow, nCol).Value = "Long Buffering"
	nCol++
	fruizione.Cell(nRow, nCol).Value = "PayerErrors"
	nCol++
	fruizione.Cell(nRow, nCol).Value = "FQDN"
	nCol++
	fruizione.Cell(nRow, nCol).Value = "NetworkType"

	// cicla su fruizioni per chiave hash
	for _, key := range chiaviFruizioni {

		// Non mostra le fruizioni senza titolo
		if fruizioni[key].VideoTitle == "" {
			continue
		}

		// ricomincia a scrivere dalla prima colonna 0
		nCol = 0
		// passa a una riga sotto
		nRow++

		//row := fruizione.AddRow()
		//row.AddCell()

		fruizione.Cell(nRow, nCol).Value = fruizioni[key].CpeID
		nCol++
		fruizione.Cell(nRow, nCol).Value = fruizioni[key].Mode
		nCol++
		fruizione.Cell(nRow, nCol).Value = fruizioni[key].ModelName
		nCol++
		fruizione.Cell(nRow, nCol).Value = fruizioni[key].StartTS
		nCol++
		fruizione.Cell(nRow, nCol).Value = fruizioni[key].EndTS
		nCol++
		fruizione.Cell(nRow, nCol).Value = fruizioni[key].VideoTitle
		nCol++
		fruizione.Cell(nRow, nCol).SetInt(fruizioni[key].Buffering)
		nCol++
		fruizione.Cell(nRow, nCol).SetInt(fruizioni[key].LongBuffering)
		nCol++
		fruizione.Cell(nRow, nCol).SetInt(fruizioni[key].PlayerError)
		nCol++
		fruizione.Cell(nRow, nCol).Value = fruizioni[key].FQDN
		nCol++
		fruizione.Cell(nRow, nCol).Value = fruizioni[key].NetworkType

	}

	/*
		GESTIONE DELLE TRAP

	*/

	// for k,v := range raw {
	// 	fmt.Println(k,v)
	// }

	// Creo un nuovo sheet dove mettere le trap
	traps, err := file.AddSheet("Trap_" + tgu)
	if err != nil {
		fmt.Printf(err.Error())
	}

	var ntRow int

	// scrive sulla riga 0 la prima
	traps.Cell(ntRow, 0).Value = "TGU"
	traps.Cell(ntRow, 1).Value = tgu

	ntRow++
	for n, v := range trapsRAW {
		if n == 0 {
			var ntCol int
			for _, value := range v.Indici {
				// scrive sulla riga 1 la seconda
				if strings.Contains(value, ">") {
					value = strings.Replace(value, "<traps_fields>", "", 1)
					value = strings.Replace(value, "</traps_fields>", "", 1)
				}
				traps.Cell(1, ntCol).Value = value
				ntCol++
			}
		}
		ntRow++
	}

	for _, v := range trapsRAW {
		ntRow = 2
		for _, value := range v.Traps {
			var ntCol int
			values := strings.Split(value, ",")
			if len(values) < 2 {
				continue
			}
			for i := 0; i < len(values); i++ {
				traps.Cell(ntRow, ntCol).Value = values[i]
				ntCol++
			}
			ntRow++
		}
	}

	// Salva il tutto su file
	err = file.Save("fruizioni_" + tgu + ".xlsx")
	if err != nil {
		fmt.Printf(err.Error())

	}
	return err
}
