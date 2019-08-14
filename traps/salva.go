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

	fruizione.Cell(0, 0).Value = "TGU"
	fruizione.Cell(0, 1).Value = tgu
	fruizione.Cell(1, 0).Value = "Cpeid"
	fruizione.Cell(1, 1).Value = "Mode"
	fruizione.Cell(1, 2).Value = "Start"
	fruizione.Cell(1, 3).Value = "End"
	fruizione.Cell(1, 4).Value = "VideoTitle"
	fruizione.Cell(1, 5).Value = "Buffering"
	fruizione.Cell(1, 6).Value = "Long Buffering"
	fruizione.Cell(1, 7).Value = "PayerErrors"
	fruizione.Cell(1, 8).Value = "FQDN"

	var nRow = 1
	var nCol int
	// cicla su fruizioni per chiave hash
	for _, key := range chiaviFruizioni {

		// Non mostra le fruizioni senza titolo
		if fruizioni[key].VideoTitle == "" {
			continue
		}

		nCol = 0
		nRow++

		//row := fruizione.AddRow()
		//row.AddCell()

		fruizione.Cell(nRow, nCol).Value = fruizioni[key].CpeID
		nCol++
		fruizione.Cell(nRow, nCol).Value = fruizioni[key].Mode
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

	}

	err = file.Save("fruizioni_" + tgu + ".xlsx")
	if err != nil {
		fmt.Printf(err.Error())

	}
	return err
}
