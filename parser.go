package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx/v3"
)

func parseFile(filename string) ([]Record, error) {
	// open an existing file
	wb, err := xlsx.OpenFile(filename)
	if err != nil {
		return nil, err
	}

	sh, ok := wb.Sheet["Extrato_Santander"]
	if !ok {
		return nil, err
	}

	records := make([]Record, 0, sh.MaxRow)
	err = sh.ForEachRow(func(r *xlsx.Row) error {
		record := Record{}
		col := 0
		err = r.ForEachCell(func(c *xlsx.Cell) error {
			v, err := c.FormattedValue()
			if err == nil {
				switch col {
				case 0:
					record.Date, _ = time.Parse("02/01/2006", v)
				case 1:
					record.History = v
				case 2:
					record.Doc = v
				case 3:
					v = strings.ReplaceAll(v, ".", "")
					v = strings.ReplaceAll(v, ",", ".")
					record.Value, _ = strconv.ParseFloat(v, 64)
				case 4:
					v = strings.ReplaceAll(v, ".", "")
					v = strings.ReplaceAll(v, ",", ".")
					record.Balance, _ = strconv.ParseFloat(v, 64)
				}
			}
			col++
			return err
		})
		records = append(records, record)
		return err
	})
	return records, err
}
