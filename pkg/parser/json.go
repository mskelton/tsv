package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

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
