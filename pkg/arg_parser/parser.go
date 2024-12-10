package arg_parser

import (
	"log"
	"strconv"
	"strings"
)

type ColumnColor string

const (
	ColumnColorDefault ColumnColor = ""
	ColumnColorRed     ColumnColor = "red"
	ColumnColorGreen   ColumnColor = "green"
	ColumnColorYellow  ColumnColor = "yellow"
	ColumnColorBlue    ColumnColor = "blue"
	ColumnColorMagenta ColumnColor = "magenta"
	ColumnColorCyan    ColumnColor = "cyan"
	ColumnColorGray    ColumnColor = "gray"
	ColumnColorDim     ColumnColor = "dim"
)

type ColumnType string

const (
	ColumnTypeText   ColumnType = "text"
	ColumnTypeNumber ColumnType = "number"
	ColumnTypeDate   ColumnType = "date"
)

type ColumnFormat string

const (
	ColumnFormatDefault ColumnFormat = ""

	ColumnFormatDecimal    ColumnFormat = "decimal"
	ColumnFormatPercent    ColumnFormat = "percent"
	ColumnFormatScientific ColumnFormat = "scientific"
	ColumnFormatAccounting ColumnFormat = "accounting"
	ColumnFormatFinancial  ColumnFormat = "financial"
	ColumnFormatCurrency   ColumnFormat = "currency"

	ColumnFormatDate     ColumnFormat = "date"
	ColumnFormatTime     ColumnFormat = "time"
	ColumnFormatDateTime ColumnFormat = "datetime"
	ColumnFormatRelative ColumnFormat = "relative"
)

type ColumnAlign string

const (
	ColumnAlignLeft  ColumnAlign = "left"
	ColumnAlignRight ColumnAlign = "right"
)

type ColumnConfig struct {
	Key      string
	Name     string
	Align    ColumnAlign
	Type     ColumnType
	Format   ColumnFormat
	Color    ColumnColor
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
		if strings.HasPrefix(arg, "-") {
			// Add the current flag to the list
			if currentFlag != nil {
				flags = append(flags, *currentFlag)
			}

			// Create a new flag
			currentFlag = &flag{
				Name: strings.TrimLeft(arg, "-"),
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
		case "c", "column":
			column := &ColumnConfig{
				Key:      strconv.Itoa(len(args.TableConfig.Columns)),
				Name:     flag.Args[0],
				Type:     ColumnTypeText,
				Color:    ColumnColorDefault,
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
					column.Color = ColumnColor(value)
				case "type":
					column.Type = ColumnType(value)
				case "format":
					column.Format = ColumnFormat(value)
				case "align":
					column.Align = ColumnAlign(value)
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
