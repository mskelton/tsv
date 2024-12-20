package printer

import (
	"fmt"
	"log"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"github.com/mskelton/tsv/pkg/arg_parser"
)

func getColor(c arg_parser.ColumnColor) *color.Color {
	switch c {
	case arg_parser.ColumnColorRed:
		return color.New(color.FgRed)
	case arg_parser.ColumnColorGreen:
		return color.New(color.FgGreen)
	case arg_parser.ColumnColorYellow:
		return color.New(color.FgYellow)
	case arg_parser.ColumnColorBlue:
		return color.New(color.FgBlue)
	case arg_parser.ColumnColorMagenta:
		return color.New(color.FgMagenta)
	case arg_parser.ColumnColorCyan:
		return color.New(color.FgCyan)
	case arg_parser.ColumnColorGray:
		return color.RGB(99, 101, 123)
	case arg_parser.ColumnColorDim:
		return color.New(color.FgWhite)
	}

	return color.New(color.FgHiWhite)
}

const separator = "  "

type Table struct {
	Config arg_parser.TableConfig
	rows   []map[string]string
}

// Special implementation of string padding/truncate to account for unicode
// string width
func autosize(str string, w int, align arg_parser.ColumnAlign) string {
	sw := runewidth.StringWidth(str)

	if sw > w {
		if align == arg_parser.ColumnAlignRight {
			return runewidth.TruncateLeft(str, sw-w+1, "…")
		} else {
			return runewidth.Truncate(str, w, "…")
		}
	}

	if align == arg_parser.ColumnAlignRight {
		return runewidth.FillLeft(str, w)
	} else {
		return runewidth.FillRight(str, w)
	}
}

func (table *Table) Load(rows []map[string]string) {
	for _, row := range rows {
		newRow := make(map[string]string)

		for _, column := range table.Config.Columns {
			value, err := format(row[column.Key], column)
			if err != nil {
				log.Fatalf("error while formatting value: %v", err)
			}

			newRow[column.Key] = value
		}

		table.rows = append(table.rows, newRow)
	}
}

func (table *Table) Print() {
	widths := make([]int, len(table.Config.Columns))
	headerColor := getColor(arg_parser.ColumnColorGray).Add(color.Underline).SprintFunc()

	// Find the maximum width of each column
	for _, row := range table.rows {
		for i, column := range table.Config.Columns {
			length := runewidth.StringWidth(row[column.Key])
			widths[i] = max(widths[i], length)
		}
	}

	// Calculate the width of each column header, ignoring empty columns
	for i, column := range table.Config.Columns {
		if widths[i] > 0 {
			// Column headers never have Unicode, so `len()` is safe to use
			widths[i] = max(widths[i], len(column.Name))
		}
	}

	// Truncate columns if necessary
	for i, column := range table.Config.Columns {
		if column.Truncate > 0 {
			widths[i] = min(widths[i], column.Truncate)
		}
	}

	// Create the header row, skipping empty columns
	var header []string
	for i, column := range table.Config.Columns {
		if widths[i] > 0 {
			header = append(header, headerColor(autosize(column.Name, widths[i], column.Align)))
		}
	}

	fmt.Println(strings.Join(header, separator))

	// Print an ASCII underline if colorization is disabled
	if color.NoColor {
		var underline []string

		for _, width := range widths {
			if width > 0 {
				underline = append(underline, strings.Repeat("-", width))
			}
		}

		fmt.Println(strings.Join(underline, separator))
	}

	for _, row := range table.rows {
		var cells []string

		for i, column := range table.Config.Columns {
			if widths[i] > 0 {
				color := getColor(column.Color).SprintFunc()
				cells = append(cells, color(autosize(row[column.Key], widths[i], column.Align)))
			}
		}

		fmt.Println(strings.Join(cells, separator))
	}
}
