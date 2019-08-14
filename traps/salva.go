package traps

import (
	"context"
	"fmt"
	"sort"

	"github.com/tealeg/xlsx"
)

func salvaXLSX(ctx context.Context, tgu string, fruizioni map[string]Output) error {

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

	err = file.Save("fruizioni_" + tgu + ".xlsx")
	if err != nil {
		fmt.Printf(err.Error())

	}
	return err
}
