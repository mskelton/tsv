package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func ParseJson(r io.Reader) []map[string]string {
	var rawRows []map[string]any
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(&rawRows); err != nil {
		log.Fatalf("Error reading input: %v\n", err)
		return nil
	}

	var rows []map[string]string
	for _, rawRow := range rawRows {
		row := make(map[string]string)

		for key, value := range rawRow {
			row[key] = fmt.Sprintf("%v", value) // Convert any value to a string
		}

		rows = append(rows, row)
	}

	return rows
}
