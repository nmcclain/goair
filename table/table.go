package table

import (
	"os"
	"reflect"

	"github.com/olekukonko/tablewriter"
)

//Table is a special type consisting of header and rows
type Table struct {
	Header  []string
	Columns []string
	RowData reflect.Value
}

//Printtable from an array of headers, columns and a reflected value of an array of structs
func (table Table) Printtable() {
	tableWriter := tablewriter.NewWriter(os.Stdout)
	tableWriter.SetHeader(table.Header)

	for i := 0; i < table.RowData.Len(); i++ {
		row := make([]string, len(table.Header))
		for j, arg := range table.Columns {
			row[j] = table.RowData.Index(i).FieldByName(arg).Interface().(string)
		}
		tableWriter.Append(row)
	}

	tableWriter.Render()

}
