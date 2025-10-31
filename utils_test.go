package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UtilsTestSuite struct {
	suite.Suite
}

// Helper function to create a mock request with or without a CSV file
func (s *UtilsTestSuite) createMockRequest(csvContent string, includeFile bool) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if includeFile {
		part, err := writer.CreateFormFile("file", "matrix.csv")
		s.Require().NoError(err)
		_, err = part.Write([]byte(csvContent))
		s.Require().NoError(err)
	}

	writer.Close()
	req := httptest.NewRequest("POST", "/echo", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}


// Test for ParseCSVFile
func (s *UtilsTestSuite) TestParseCSVFile() {
	tests := []struct {
		name        string
		csvContent  string
		includeFile bool
		expectErr   bool
		errorSubstr string
		expectRows  int
		expectCols  int
	}{
		{
			name:       "valid csv file",
			csvContent: "1,2,3\n4,5,6\n7,8,9",
			includeFile: true,
			expectErr:  false,
			expectRows: 3,
			expectCols: 3,
		},
		{
			name:        "missing file field",
			csvContent:  "",
			includeFile: false,
			expectErr:   true,
			errorSubstr: "failed to read file",
		},
		{
			name:        "invalid csv format",
			csvContent:  "\"1,2,3\n4,5,6",
			includeFile: true,
			expectErr:   true,
			errorSubstr: "failed to parse csv",
		},
		{
			name:        "empty file",
			csvContent:  "",
			includeFile: true,
			expectErr:   true,
			errorSubstr: "failed to parse csv: file is empty",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			req := s.createMockRequest(tc.csvContent, tc.includeFile)
			records, err := ParseCSVFile(req)
			
			// Check if error contain errorSubstr
			if tc.expectErr {
				s.Error(err)
				s.Contains(err.Error(), tc.errorSubstr)
			// Check row and column numbers if no error
			} else {
				s.NoError(err)
				s.Len(records, tc.expectRows)
				s.Len(records[0], tc.expectCols)
			}
		})
	}
}

// Test for ValidateSquareMatrix
func (s *UtilsTestSuite) TestValidateSquareMatrix() {
	tests := []struct {
		name        string
		matrix      [][]string
		expectErr   bool
		errorSubstr string
	}{
		{
			name: "valid square 2*2 matrix",
			matrix: [][]string{
				{"1", "2"},
				{"3", "4"},
			},
			expectErr: false,
		},
		{
			name: "valid square 3*3 matrix",
			matrix: [][]string{
				{"1", "2", "3"},
				{"4", "5", "6"},
				{"7", "8", "9"},
			},
			expectErr: false,
		},
		{
			name: "valid square 5*5 matrix",
			matrix: [][]string{
				{"1", "2", "3", "4", "5"},
				{"6", "7", "8", "9", "10"},
				{"11", "12", "13", "14", "15"},
				{"16", "17", "18", "19", "20"},
				{"21", "22", "23", "24", "25"},
			},
			expectErr: false,
		},
		{
			name: "empty matrix",
			matrix: [][]string{},
			expectErr:   true,
			errorSubstr: "empty matrix",
		},
		{
			name: "matrix has header row",
			matrix: [][]string{
				{"This is header", "row content"},
				{"1", "2"},
			},
			expectErr:   true,
			errorSubstr: "header row",
		},
		{
			name: "matrix has different row length",
			matrix: [][]string{
				{"1", "2", "3"},
				{"4", "5"},
			},
			expectErr:   true,
			errorSubstr: "matrix is not square",
		},
		{
			name: "matrix has empty value",
			matrix: [][]string{
				{"1", "2", "3"},
				{"4", "5", "6"},
				{"7", "8", ""},
			},
			expectErr:   true,
			errorSubstr: "matrix has empty value",
		},
		{
			name: "matrix has non-integer value",
			matrix: [][]string{
				{"1", "2", "3"},
				{"4", "A", "6"},
				{"7", "8", "9"},
			},
			expectErr:   true,
			errorSubstr: "is not an integer",
		},
		{
			name: "matrix has more rows than columns",
			matrix: [][]string{
				{"1", "2"},
				{"3", "4"},
				{"5", "6"},
			},
			expectErr:   true,
			errorSubstr: "matrix is not square",
		},
		{
			name: "matrix has more columns than rows",
			matrix: [][]string{
				{"1", "2", "3"},
				{"4", "5", "6"},
			},
			expectErr:   true,
			errorSubstr: "matrix is not square",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			err := ValidateSquareMatrix(tc.matrix)

			// Check if error contain errorSubstr
			if tc.expectErr {
				s.Error(err)
				if tc.errorSubstr != "" {
					s.Contains(strings.ToLower(err.Error()), strings.ToLower(tc.errorSubstr))
				}
			} else {
				s.NoError(err)
			}
		})
	}
}

// Run all tests
func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}
