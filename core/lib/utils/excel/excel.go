package excel

import (
	"time"

	"github.com/tealeg/xlsx"
)

func BuildExcelHeaders(sheet *xlsx.Sheet, headers []string) {
	// header row
	header := sheet.AddRow()
	for k := range headers {
		cell := header.AddCell()
		cell.Value = headers[k]
	}
}

func AddCellString(row *xlsx.Row, value string) {
	cell := row.AddCell()
	cell.Value = value
}

func AddCellMoney(row *xlsx.Row, value float64) {
	cell := row.AddCell()
	cell.SetFloatWithFormat(value, "_(\"$\"* #,##0.00_);_(\"$\"* \\(#,##0.00\\);_(\"$\"* \"-\"??_);_(@_)")
}

func AddCellPercent(row *xlsx.Row, value float64) {
	cell := row.AddCell()
	cell.SetFloatWithFormat(value, "0.00%")
}

func AddCellInteger(row *xlsx.Row, value int64) {
	cell := row.AddCell()
	cell.SetInt64(value)
}

func AddCellDate(row *xlsx.Row, value time.Time) {
	timeLocationUTC, _ := time.LoadLocation("UTC")
	cell := row.AddCell()
	cell.SetDateWithOptions(value, xlsx.DateTimeOptions{
		Location:        timeLocationUTC,
		ExcelTimeFormat: "mm-dd-yy",
	})
}
