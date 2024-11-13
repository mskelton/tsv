package tsv

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func Parse(r io.Reader) [][]string {
	rows := [][]string{}
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		columns := strings.Split(line, "\t")
		rows = append(rows, columns)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}

	return rows
}
