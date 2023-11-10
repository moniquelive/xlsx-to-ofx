// package comment
package main

import (
	"log"
	"os"
	"time"
)

type Record struct {
	Date    time.Time
	History string
	Doc     string
	Value   float64
	Balance float64
}

type OFXData struct {
	Now     time.Time
	Agencia string
	Conta   string
	Records []Record
	Balance float64
}

func main() {
	var filename string
	if len(os.Args) > 1 {
		filename = os.Args[1]
	} else {
		filename = "./Extrato_Santander.setembro e outubro 2023.xlsx"
	}
	records, err := parseFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	data := OFXData{
		Now:     time.Now(),
		Agencia: "1053",
		Conta:   "130002469",
		Records: records[1:],
		Balance: records[len(records)-1].Balance,
	}
	// Agencia 1053
	// Conta brand 13000220-1
	// Conta 21212 s.a. 13000105-9
	// Conta 21212AFN 13000246-9
	if err = tmpl.Execute(os.Stdout, data); err != nil {
		log.Fatal(err)
	}
}
