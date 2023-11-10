package main

import (
	"fmt"
	"text/template"
	"time"
)

var ofx = `
OFXHEADER:100
DATA:OFXSGML
VERSION:102
SECURITY:NONE
ENCODING:USASCII
CHARSET:1252
COMPRESSION:NONE
OLDFILEUID:NONE
NEWFILEUID:NONE
<OFX>
	<SIGNONMSGSRSV1>
		<SONRS>
			<STATUS>
				<CODE>0
				<SEVERITY>INFO
			</STATUS>
			<DTSERVER>{{.Now | fmtdate}}
			<LANGUAGE>ENG
			<FI>
				<ORG>SANTANDER
				<FID>SANTANDER
			</FI>
		</SONRS>
	</SIGNONMSGSRSV1>
	<BANKMSGSRSV1>
		<STMTTRNRS>
			<TRNUID>1
			<STATUS>
				<CODE>0
				<SEVERITY>INFO
			</STATUS>
			<STMTRS>
				<CURDEF>BRL
				<BANKACCTFROM>
					<BANKID>033
					<ACCTID>{{.Agencia}}{{.Conta}}
					<ACCTTYPE>CHECKING
				</BANKACCTFROM>
				<BANKTRANLIST>
				<DTSTART>{{.Now | fmtdate}}
				<DTEND>{{.Now | fmtdate}}
{{- range .Records }}
				<STMTTRN>
					<TRNTYPE>OTHER
					<DTPOSTED>{{.Date | fmtdate}}
					<TRNAMT>{{.Value | fmtfloat}}
					<FITID>{{.Doc}}
					<CHECKNUM>{{.Doc}}
					<PAYEEID>0
					<MEMO>{{.History}}
				</STMTTRN>
{{- end }}
				<BANKTRANLIST>
				<LEDGERBAL>
					<BALAMT>{{.Balance}}
					<DTASOF>{{.Now | fmtdate}}
				</LEDGERBAL>
			</STMTRS>
		</STMTTRNRS>
	</BANKMSGSRSV1>
</OFX>`

var fm = template.FuncMap{
	"fmtdate": func(dt time.Time) string {
		return fmt.Sprintf("%04d%02d%02d%02d%02d%02d[-3:GMT]",
			dt.Year(), dt.Month(), dt.Day(),
			dt.Hour(), dt.Minute(), dt.Second())
	},
	"fmtfloat": func(f float64) string { return fmt.Sprintf("%.02f", f) },
}

var tmpl = template.Must(template.New("ofx").Funcs(fm).Parse(ofx))
