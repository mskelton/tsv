package parser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func ParseTSV(r io.Reader, hasHeader bool) []map[string]string {
	rows := []map[string]string{}
	scanner := bufio.NewScanner(r)

	// Skip the header row if the `--header` flag is set
	if hasHeader {
		scanner.Scan()
	}

	for scanner.Scan() {
		line := scanner.Text()
		cells := map[string]string{}

		for i, cell := range strings.Split(line, "\t") {
			cells[strconv.Itoa(i)] = cell
		}

		rows = append(rows, cells)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}

	return rows
}
