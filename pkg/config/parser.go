package config

import (
	"strings"
)

type CellColor string

const (
	CellColorDefault CellColor = ""
	CellColorRed     CellColor = "red"
	CellColorGreen   CellColor = "green"
	CellColorYellow  CellColor = "yellow"
	CellColorBlue    CellColor = "blue"
	CellColorMagenta CellColor = "magenta"
	CellColorCyan    CellColor = "cyan"
	CellColorGray    CellColor = "gray"
	CellColorDim     CellColor = "dim"
)

type ColumnType string

const (
	ColumnTypeText   ColumnType = "text"
	ColumnTypeNumber ColumnType = "number"
	ColumnTypeDate   ColumnType = "date"
)

type ColumnFormat string

const (
	ColumnFormatDefault    ColumnFormat = ""
	ColumnFormatDecimal    ColumnFormat = "decimal"
	ColumnFormatPercent    ColumnFormat = "percent"
	ColumnFormatScientific ColumnFormat = "scientific"
	ColumnFormatAccounting ColumnFormat = "accounting"
	ColumnFormatFinancial  ColumnFormat = "financial"
	ColumnFormatCurrency   ColumnFormat = "currency"
	ColumnFormatDate       ColumnFormat = "date"
	ColumnFormatTime       ColumnFormat = "time"
	ColumnFormatDateTime   ColumnFormat = "datetime"
	ColumnFormatRelative   ColumnFormat = "relative"
)

type ColumnConfig struct {
	Key    string
	Name   string
	Type   ColumnType
	Format ColumnFormat
	Color  CellColor
}

type TableConfig struct {
	Columns []ColumnConfig
}

func Parse(args []string) TableConfig {
	table := TableConfig{}
	var column *ColumnConfig

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "--column" {
			if i+1 < len(args) {
				columnName := args[i+1]
				column = &ColumnConfig{
					Key:    columnName,
					Name:   columnName,
					Type:   ColumnTypeText,
					Color:  CellColorDefault,
					Format: ColumnFormatDefault,
				}

				table.Columns = append(table.Columns, *column)

				// Move to the next argument after the column name
				i++
			}
		} else if column != nil && strings.Contains(arg, "=") {
			// Parse key=value pairs
			parts := strings.SplitN(arg, "=", 2)
			key, value := parts[0], parts[1]

			switch key {
			case "color":
				column.Color = CellColor(value)
			case "type":
				column.Type = ColumnType(value)
			case "format":
				column.Format = ColumnFormat(value)
			}
		}
	}

	return table
}
