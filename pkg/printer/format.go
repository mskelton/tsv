package printer

import (
	"fmt"
	"strconv"
	"time"

	"github.com/mergestat/timediff"
	"github.com/mskelton/tsv/pkg/arg_parser"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func format(value string, column arg_parser.ColumnConfig) (string, error) {
	switch column.Type {
	case arg_parser.ColumnTypeNumber:
		return formatNumber(value, column.Format)
	case arg_parser.ColumnTypeDate:
		return formatDate(value, column.Format)
	}

	return value, nil
}

func formatNumber(value string, format arg_parser.ColumnFormat) (string, error) {
	if value == "" {
		return value, nil
	}

	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return "", err
	}

	switch format {
	case arg_parser.ColumnFormatDecimal:
		return formatDecimal(num), nil
	case arg_parser.ColumnFormatPercent:
		return formatPercent(num), nil
	case arg_parser.ColumnFormatScientific:
		return formatScientific(num), nil
	case arg_parser.ColumnFormatAccounting:
		return formatAccounting(num), nil
	case arg_parser.ColumnFormatFinancial:
		return formatFinancial(num), nil
	case arg_parser.ColumnFormatCurrency:
		return formatCurrency(num), nil
	}

	return value, nil
}

func formatDecimal(num float64) string {
	return message.NewPrinter(language.English).Sprintf("%.2f", num)
}

func formatPercent(num float64) string {
	return formatDecimal(num*100) + "%"
}

func formatScientific(num float64) string {
	return message.NewPrinter(language.English).Sprintf("%.2e", num)
}

func formatAccounting(num float64) string {
	return "$" + formatFinancial(num)
}

func formatFinancial(num float64) string {
	if num < 0 {
		return fmt.Sprintf("(%s)", formatDecimal(-num))
	}

	return formatDecimal(num)
}

func formatCurrency(num float64) string {
	if num < 0 {
		return "-$" + formatDecimal(-num)
	}

	return "$" + formatDecimal(num)
}

func formatDate(value string, format arg_parser.ColumnFormat) (string, error) {
	date, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return "", err
	}

	switch format {
	case arg_parser.ColumnFormatDate:
		return date.Format("Jan 2, 2006"), nil
	case arg_parser.ColumnFormatTime:
		return date.Format("3:04 PM"), nil
	case arg_parser.ColumnFormatDateTime:
		return date.Format("Jan 2, 2006 3:04 PM"), nil
	case arg_parser.ColumnFormatRelative:
		return timediff.TimeDiff(date), nil
	}

	return value, nil
}
