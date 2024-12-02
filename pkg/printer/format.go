package printer

import (
	"strconv"
	"time"

	"github.com/mergestat/timediff"
	"github.com/mskelton/tsv/pkg/arg_parser"
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
		return strconv.FormatFloat(num, 'f', 2, 64), nil
	case arg_parser.ColumnFormatPercent:
		return strconv.FormatFloat(num*100, 'f', 2, 64) + "%", nil
	case arg_parser.ColumnFormatScientific:
		return strconv.FormatFloat(num, 'e', 2, 64), nil
	case arg_parser.ColumnFormatAccounting:
		return "$" + strconv.FormatFloat(num, 'f', 2, 64), nil
	case arg_parser.ColumnFormatFinancial:
		return strconv.FormatFloat(num, 'f', 2, 64), nil
	case arg_parser.ColumnFormatCurrency:
		return "$" + strconv.FormatFloat(num, 'f', 2, 64), nil
	}

	return value, nil
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
