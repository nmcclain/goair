package table

import (
	"fmt"
	"os"
	"reflect"

	"github.com/olekukonko/tablewriter"
)

//Table is a special type consisting of header and rows and reflected Value with array of structs
type Table struct {
	Header  []string
	Columns []string
	RowData reflect.Value
}

//KeyValue is stores a Key and a value as interface
type KeyValue struct {
	Key   string
	Value string
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
				row[j] = fmt.Sprintf("%v", field.Interface())
			}
		}
		tableWriter.Append(row)
	}

	tableWriter.SetAlignment(tablewriter.ALIGN_LEFT)
	tableWriter.Render()

}

//PrintColumn from an array of headers, and a reflected value of a struct with an array of strings
func (table Table) PrintColumn() {
	tableWriter := tablewriter.NewWriter(os.Stdout)
	tableWriter.SetHeader(table.Header)

	for i := 0; i < table.RowData.Len(); i++ {
		row := make([]string, len(table.Header))
		row[0] = fmt.Sprintf("%v", table.RowData.Index(i).Interface())
		tableWriter.Append(row)
	}
	tableWriter.SetAlignment(tablewriter.ALIGN_LEFT)
	tableWriter.Render()
}

//PrintKeyValueTable from an array of structs
func (table Table) PrintKeyValueTable() {
	keyValues := make([]KeyValue, table.RowData.NumField())
	for i := 0; i < table.RowData.NumField(); i++ {
		keyValues[i] = KeyValue{
			Key:   table.RowData.Type().Field(i).Name,
			Value: fmt.Sprintf("%v", table.RowData.Field(i).Interface()),
		}
	}
	table = Table{
		Header:  []string{"Key", "Value"},
		Columns: []string{"Key", "Value"},
		RowData: reflect.ValueOf(&keyValues).Elem(),
	}
	table.PrintTable()

}
