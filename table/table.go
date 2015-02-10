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

//PrintTable from an array of headers, columns and a reflected value of an array of structs
func (table Table) PrintTable() {
	tableWriter := tablewriter.NewWriter(os.Stdout)
	tableWriter.SetHeader(table.Header)

	for i := 0; i < table.RowData.Len(); i++ {
		row := make([]string, len(table.Header))
		for j, arg := range table.Columns {
			field := table.RowData.Index(i).FieldByName(arg)
			if field.Kind() == reflect.String {
				row[j] = field.Interface().(string)
			}
		}
		tableWriter.Append(row)
	}

	tableWriter.Render()

}

//PrintColumn from an array of headers, columns and a reflected value of a struct with an array of strings
func (table Table) PrintColumn() {
	tableWriter := tablewriter.NewWriter(os.Stdout)
	tableWriter.SetHeader(table.Header)

	for i := 0; i < table.RowData.Len(); i++ {
		row := make([]string, len(table.Header))
		row[0] = table.RowData.Index(i).Interface().(string)
		tableWriter.Append(row)
	}

	tableWriter.Render()

}
