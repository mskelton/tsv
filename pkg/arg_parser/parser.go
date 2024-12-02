package arg_parser

import (
	"log"
	"strconv"
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
	Key      string
	Name     string
	Type     ColumnType
	Format   ColumnFormat
	Color    CellColor
	Truncate int
}

type TableConfig struct {
	Columns []ColumnConfig
}

type InputFormat string

const (
	InputFormatTSV  InputFormat = "tsv"
	InputFormatJson InputFormat = "json"
)

type Args struct {
	InputHasHeader bool
	InputFormat    InputFormat
	TableConfig    TableConfig
}

func Parse(args []string) Args {
	return parseFlags(parseArgs(args))
}

type flag struct {
	Name string
	Args []string
}

func parseArgs(args []string) []flag {
	var flags []flag
	var currentFlag *flag

	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			// Add the current flag to the list
			if currentFlag != nil {
				flags = append(flags, *currentFlag)
			}

			// Create a new flag
			currentFlag = &flag{
				Name: strings.TrimPrefix(arg, "--"),
				Args: []string{},
			}
		} else {
			if currentFlag == nil {
				log.Fatalf("unexpected argument: %s", arg)
			}

			currentFlag.Args = append(currentFlag.Args, arg)
		}
	}

	// Add the last flag to the list
	if currentFlag != nil {
		flags = append(flags, *currentFlag)
	}

	// Process single "=" format
	for i := range flags {
		parts := strings.SplitN(flags[i].Name, "=", 2)

		if len(parts) == 2 {
			flags[i].Name = parts[0]
			flags[i].Args = append([]string{parts[1]}, flags[i].Args...)
		}
	}

	return flags
}

func parseFlags(flags []flag) Args {
	args := Args{InputFormat: InputFormatTSV}

	for _, flag := range flags {
		switch flag.Name {
		case "column":
			column := &ColumnConfig{
				Key:      strconv.Itoa(len(args.TableConfig.Columns)),
				Name:     flag.Args[0],
				Type:     ColumnTypeText,
				Color:    CellColorDefault,
				Format:   ColumnFormatDefault,
				Truncate: 0,
			}

			// Parse key=value pairs
			for _, arg := range flag.Args[1:] {
				parts := strings.SplitN(arg, "=", 2)
				key, value := parts[0], parts[1]

				switch key {
				case "key":
					column.Key = value

					// If any column contains a key, set the format to JSON automatically
					args.InputFormat = InputFormatJson
				case "color":
					column.Color = CellColor(value)
				case "type":
					column.Type = ColumnType(value)
				case "format":
					column.Format = ColumnFormat(value)
				case "trunc":
					trunc, err := strconv.Atoi(value)
					if err != nil {
						log.Fatalf("invalid trunc value: %s", value)
					}

					column.Truncate = trunc
				}
			}

			args.TableConfig.Columns = append(args.TableConfig.Columns, *column)
		case "format":
			args.InputFormat = InputFormat(flag.Args[0])
		case "header":
			args.InputHasHeader = true
		}
	}

	return args
}
