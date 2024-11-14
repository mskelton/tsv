package main

import (
	"os"

	"github.com/mskelton/tsv/pkg/arg_parser"
	"github.com/mskelton/tsv/pkg/parser"
	"github.com/mskelton/tsv/pkg/printer"
)

func main() {
	args := arg_parser.Parse(os.Args[1:])
	var rows []map[string]string

	switch args.InputFormat {
	case arg_parser.InputFormatTSV:
		rows = parser.ParseTSV(os.Stdin, args.InputHasHeader)
	case arg_parser.InputFormatJson:
		rows = parser.ParseJson(os.Stdin)
	}

	table := printer.Table{
		Config: args.TableConfig,
		Rows:   rows,
	}

	table.Print()
}
