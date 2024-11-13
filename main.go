package main

import (
	"os"

	"github.com/mskelton/tsv/pkg/config"
	"github.com/mskelton/tsv/pkg/printer"
	"github.com/mskelton/tsv/pkg/tsv"
)

func main() {
	tableConfig := config.Parse(os.Args[1:])
	rows := tsv.Parse(os.Stdin)
	table := printer.Table{
		Config: tableConfig,
		Rows:   rows,
	}

	table.Print()
}
