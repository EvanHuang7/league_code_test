package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
)

// ParseCSVFile reads csv file and returns file records.
func ParseCSVFile(r *http.Request) ([][]string, error) {
	// Get csv file
	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}
	defer file.Close()

	// Read all records from csv file
	reader := csv.NewReader(file)
	// Disable automatic field count checking to allow different row length and header cases
	// since we check them in latter ValidateSquareMatrix() to return more specific error.
    reader.FieldsPerRecord = -1
    records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse csv: %v", err)
	}

	// If csv file is empty
	if len(records) == 0 {
		return nil, fmt.Errorf("failed to parse csv: file is empty")
	}

	return records, nil
}

// ValidateSquareMatrix checks if the input matrix is square, contains only integers, and has no header row.
func ValidateSquareMatrix(records [][]string) error {
	// Empty matrix case
	if len(records) == 0 {
		return fmt.Errorf("empty matrix")
	}

	n := len(records[0])
	// Check first row for header (all values should be integers)
	for j, val := range records[0] {
		if _, err := strconv.Atoi(val); err != nil {
			return fmt.Errorf("matrix has a header row or non-integer value at row 1, column %d: %v", j+1, err)
		}
	}

	for i, row := range records {
		// Check if each row length is equal to the first row length
		if len(row) != n {
			return fmt.Errorf("matrix is not square: row %d has %d columns, expected %d", i+1, len(row), n)
		}

		for j, val := range row {
			// Check for empty value
			if val == "" {
				return fmt.Errorf("matrix has empty value at row %d, column %d", i+1, j+1)
			}
			// Check for integer
			if _, err := strconv.Atoi(val); err != nil {
				return fmt.Errorf("matrix value at row %d, column %d is not an integer: %v", i+1, j+1, err)
			}
		}
	}

	// Check if number of rows equals number of columns
	if len(records) != n {
		return fmt.Errorf("matrix is not square: %d rows and %d columns", len(records), n)
	}

	return nil
}