package traps

import (
	"context"
	"fmt"

	"github.com/tealeg/xlsx"
)

func salvaXLSX(ctx context.Context, fruizioni map[string]Output) error {

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

	// fmt.Println(len(fruizioni))
	// fmt.Println(fruizioni)

	tgu, err := file.AddSheet("Fruizioni")
	if err != nil {
		fmt.Printf(err.Error())
	}

	tgu.Cell(1, 0).Value = "Cpeid"
	tgu.Cell(1, 1).Value = "Mode"
	tgu.Cell(1, 2).Value = "Start"
	tgu.Cell(1, 3).Value = "End"
	tgu.Cell(1, 4).Value = "VideoTitle"
	tgu.Cell(1, 5).Value = "Buffering"
	tgu.Cell(1, 6).Value = "Long Buffering"
	tgu.Cell(1, 7).Value = "PayerErrors"
	tgu.Cell(1, 8).Value = "FQDN"

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

		//row := tgu.AddRow()
		//row.AddCell()

		tgu.Cell(nRow, nCol).Value = fruizioni[key].CpeID
		nCol++
		tgu.Cell(nRow, nCol).Value = fruizioni[key].Mode
		nCol++
		tgu.Cell(nRow, nCol).Value = fruizioni[key].StartTS
		nCol++
		tgu.Cell(nRow, nCol).Value = fruizioni[key].EndTS
		nCol++
		tgu.Cell(nRow, nCol).Value = fruizioni[key].VideoTitle
		nCol++
		tgu.Cell(nRow, nCol).SetInt(fruizioni[key].Buffering)
		nCol++
		tgu.Cell(nRow, nCol).SetInt(fruizioni[key].LongBuffering)
		nCol++
		tgu.Cell(nRow, nCol).SetInt(fruizioni[key].PlayerError)
		nCol++
		tgu.Cell(nRow, nCol).Value = fruizioni[key].FQDN

	}
	err = file.Save("test.xlsx")
	if err != nil {
		fmt.Printf(err.Error())

	}
	return err
}

func poppo() {
	var file *xlsx.File
	// var sheet *xlsx.Sheet
	// var row *xlsx.Row
	// var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	tgu, err := file.AddSheet("TGU")
	if err != nil {
		fmt.Printf(err.Error())
	}
	cpeid, err := file.AddSheet("cpeid")
	if err != nil {
		fmt.Printf(err.Error())
	}
	//row = tgu.AddRow()
	//cell = row.AddCell()
	tgu.Cell(1, 2).Value = "daje"
	tgu.Cell(2, 2).Value = "dajedaje"
	cpeid.Cell(2, 2).Value = "dajedaje"
	//cell.Value = "I am a cell!"
	err = file.Save("pippo.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}
