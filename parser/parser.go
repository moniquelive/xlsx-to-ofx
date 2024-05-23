// Package parser is responsible for parsing XLSX files
package parser

import (
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx/v3"
)

// Record represents a line in the spreadsheet
type Record struct {
	Date    time.Time
	History string
	Doc     string
	Value   float64
	Balance float64
}

// ParseReader parses the XLSX file and returns a slice of Records
func ParseReader(reader io.ReaderAt, size int64) (records []Record, err error) {
	// open an existing file
	var wb *xlsx.File
	wb, err = xlsx.OpenReaderAt(reader, size)
	if err != nil {
		return nil, err
	}

	sh, ok := wb.Sheet["Extrato_Santander"]
	if !ok {
		return nil, err
	}

	records = make([]Record, 0, sh.MaxRow)
	err = sh.ForEachRow(func(r *xlsx.Row) error {
		record := Record{}
		col := 0
		err = r.ForEachCell(func(c *xlsx.Cell) error {
			var v string
			v, err = c.FormattedValue()
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
