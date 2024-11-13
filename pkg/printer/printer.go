package printer

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"github.com/mskelton/tsv/pkg/config"
)

func getColor(c config.CellColor) *color.Color {
	switch c {
	case config.CellColorRed:
		return color.New(color.FgRed)
	case config.CellColorGreen:
		return color.New(color.FgGreen)
	case config.CellColorYellow:
		return color.New(color.FgYellow)
	case config.CellColorBlue:
		return color.New(color.FgBlue)
	case config.CellColorMagenta:
		return color.New(color.FgMagenta)
	case config.CellColorCyan:
		return color.New(color.FgCyan)
	case config.CellColorGray:
		return color.RGB(99, 101, 123)
	case config.CellColorDim:
		return color.New(color.FgWhite)
	}

	return color.New(color.FgHiWhite)
}

const separator = "  "

type Table struct {
	Config config.TableConfig
	Rows   [][]string
}

// Special implementation of string padding to account for unicode string width
func pad(str string, w int) string {
	return str + strings.Repeat(" ", w-runewidth.StringWidth(str))
}

func (table *Table) Print() {
	widths := make([]int, len(table.Config.Columns))
	headerColor := getColor(config.CellColorGray).Add(color.Underline).SprintFunc()

	// Find the maximum width of each column
	for _, row := range table.Rows {
		for i := range table.Config.Columns {
			length := runewidth.StringWidth(row[i])
			widths[i] = max(widths[i], length)
		}
	}

	// Calculate the width of each column header, ignoring empty columns
	for i, col := range table.Config.Columns {
		if widths[i] > 0 {
			// Column headers never have Unicode, so `len()` is safe to use
			widths[i] = max(widths[i], len(col.Name))
		}
	}

	// Create the header row, skipping empty columns
	var header []string
	for i, col := range table.Config.Columns {
		if widths[i] > 0 {
			header = append(header, headerColor(pad(col.Name, widths[i])))
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

	for _, row := range table.Rows {
		var cells []string

		for i, column := range table.Config.Columns {
			if widths[i] > 0 {
				color := getColor(column.Color).SprintFunc()
				cells = append(cells, color(pad(row[i], widths[i])))
			}
		}

		fmt.Println(strings.Join(cells, separator))
	}
}
