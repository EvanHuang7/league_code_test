package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type EndpointTestSuite struct {
	suite.Suite
}

// Helper function to create request with CSV file
func (s *EndpointTestSuite) createCSVRequest(endpoint string, filePath string) *http.Request {
	fileBytes, err := os.ReadFile(filePath)
	s.Require().NoError(err)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "matrix.csv")
	s.Require().NoError(err)

	_, err = part.Write(fileBytes)
	s.Require().NoError(err)

	writer.Close()

	req := httptest.NewRequest("POST", endpoint, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

// Test for echo endpoint
func (s *EndpointTestSuite) TestEchoEndpoint() {
	tests := []struct {
		name                   string
		filePath               string
		expectedStatusCode     int
		expectedResponseSubstr string
	}{
		{
			name:                   "valid 3*3 matrix",
			filePath:               "testdata/valid_3_to_3.csv",
			expectedStatusCode:     200,
			expectedResponseSubstr: "1,2,3\n4,5,6\n7,8,9\n",
		},
		{
			name:                   "valid 2*2 matrix",
			filePath:               "testdata/valid_2_to_2.csv",
			expectedStatusCode:     200,
			expectedResponseSubstr: "0,1\n2,3\n",
		},
		{
			name:                   "valid 4*4 matrix",
			filePath:               "testdata/valid_4_to_4.csv",
			expectedStatusCode:     200,
			expectedResponseSubstr: "1,2,3,4\n2,2,-1,-10\n3,3,5,-2\n4,3,2,1\n",
		},
		{
			name:                   "empty matrix",
			filePath:               "testdata/empty.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "failed to parse csv: file is empty",
		},
		{
			name:                   "matrix with header",
			filePath:               "testdata/header.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix has a header row or non-integer value at row 1, column 1",
		},
		{
			name:                   "matrix has different row length",
			filePath:               "testdata/different_row_length.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix is not square: row 2 has 2 columns, expected 3",
		},
		{
			name:                   "matrix with empty value",
			filePath:               "testdata/empty_value.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix has empty value at row 3",
		},
		{
			name:                   "matrix with non-integer value",
			filePath:               "testdata/non_integer_value.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix value at row 3, column 2 is not an integer",
		},
		{
			name:                   "matrix has more rows than columns",
			filePath:               "testdata/more_rows_than_cols.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix is not square: 3 rows and 2 columns",
		},
		{
			name:                   "matrix has more columns than rows",
			filePath:               "testdata/more_cols_than_rows.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix is not square: 2 rows and 3 columns",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			req := s.createCSVRequest("/echo", tc.filePath)
			w := httptest.NewRecorder()
			EchoHandler(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			s.Equal(tc.expectedStatusCode, resp.StatusCode)
			s.Contains(string(body), tc.expectedResponseSubstr)
		})
	}
}

// Test for invert endpoint
func (s *EndpointTestSuite) TestInvertEndpoint() {
	tests := []struct {
		name                   string
		filePath               string
		expectedStatusCode     int
		expectedResponseSubstr string
	}{
		{
			name:                   "valid 3*3 matrix",
			filePath:               "testdata/valid_3_to_3.csv",
			expectedStatusCode:     200,
			expectedResponseSubstr: "1,4,7\n2,5,8\n3,6,9\n",
		},
		{
			name:                   "valid 2*2 matrix",
			filePath:               "testdata/valid_2_to_2.csv",
			expectedStatusCode:     200,
			expectedResponseSubstr: "0,2\n1,3\n",
		},
		{
			name:                   "valid 4*4 matrix",
			filePath:               "testdata/valid_4_to_4.csv",
			expectedStatusCode:     200,
			expectedResponseSubstr: "1,2,3,4\n2,2,3,3\n3,-1,5,2\n4,-10,-2,1\n",
		},
		{
			name:                   "empty matrix",
			filePath:               "testdata/empty.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "failed to parse csv: file is empty",
		},
		{
			name:                   "matrix with header",
			filePath:               "testdata/header.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix has a header row or non-integer value at row 1, column 1",
		},
		{
			name:                   "matrix has different row length",
			filePath:               "testdata/different_row_length.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix is not square: row 2 has 2 columns, expected 3",
		},
		{
			name:                   "matrix with empty value",
			filePath:               "testdata/empty_value.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix has empty value at row 3",
		},
		{
			name:                   "matrix with non-integer value",
			filePath:               "testdata/non_integer_value.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix value at row 3, column 2 is not an integer",
		},
		{
			name:                   "matrix has more rows than columns",
			filePath:               "testdata/more_rows_than_cols.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix is not square: 3 rows and 2 columns",
		},
		{
			name:                   "matrix has more columns than rows",
			filePath:               "testdata/more_cols_than_rows.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix is not square: 2 rows and 3 columns",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			req := s.createCSVRequest("/invert", tc.filePath)
			w := httptest.NewRecorder()
			InvertHandler(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			s.Equal(tc.expectedStatusCode, resp.StatusCode)
			s.Contains(string(body), tc.expectedResponseSubstr)
		})
	}
}

// Test for flatten endpoint
func (s *EndpointTestSuite) TestFlattenEndpoint() {
	tests := []struct {
		name                   string
		filePath               string
		expectedStatusCode     int
		expectedResponseSubstr string
	}{
		{
			name:                   "valid 3*3 matrix",
			filePath:               "testdata/valid_3_to_3.csv",
			expectedStatusCode:     200,
			expectedResponseSubstr: "1,2,3,4,5,6,7,8,9",
		},
		{
			name:                   "valid 2*2 matrix",
			filePath:               "testdata/valid_2_to_2.csv",
			expectedStatusCode:     200,
			expectedResponseSubstr: "0,1,2,3",
		},
		{
			name:                   "valid 4*4 matrix",
			filePath:               "testdata/valid_4_to_4.csv",
			expectedStatusCode:     200,
			expectedResponseSubstr: "1,2,3,4,2,2,-1,-10,3,3,5,-2,4,3,2,1",
		},
		{
			name:                   "empty matrix",
			filePath:               "testdata/empty.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "failed to parse csv: file is empty",
		},
		{
			name:                   "matrix with header",
			filePath:               "testdata/header.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix has a header row or non-integer value at row 1, column 1",
		},
		{
			name:                   "matrix has different row length",
			filePath:               "testdata/different_row_length.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix is not square: row 2 has 2 columns, expected 3",
		},
		{
			name:                   "matrix with empty value",
			filePath:               "testdata/empty_value.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix has empty value at row 3",
		},
		{
			name:                   "matrix with non-integer value",
			filePath:               "testdata/non_integer_value.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix value at row 3, column 2 is not an integer",
		},
		{
			name:                   "matrix has more rows than columns",
			filePath:               "testdata/more_rows_than_cols.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix is not square: 3 rows and 2 columns",
		},
		{
			name:                   "matrix has more columns than rows",
			filePath:               "testdata/more_cols_than_rows.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix is not square: 2 rows and 3 columns",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			req := s.createCSVRequest("/flatten", tc.filePath)
			w := httptest.NewRecorder()
			FlattenHandler(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			s.Equal(tc.expectedStatusCode, resp.StatusCode)
			s.Contains(string(body), tc.expectedResponseSubstr)
		})
	}
}


// Test for sum endpoint
func (s *EndpointTestSuite) TestSumEndpoint() {
	tests := []struct {
		name                   string
		filePath               string
		expectedStatusCode     int
		expectedResponseSubstr string
	}{
		{
			name:                   "valid 3*3 matrix",
			filePath:               "testdata/valid_3_to_3.csv",
			expectedStatusCode:     200,
			expectedResponseSubstr: "45",
		},
		{
			name:                   "valid 2*2 matrix",
			filePath:               "testdata/valid_2_to_2.csv",
			expectedStatusCode:     200,
			expectedResponseSubstr: "6",
		},
		{
			name:                   "valid 4*4 matrix",
			filePath:               "testdata/valid_4_to_4.csv",
			expectedStatusCode:     200,
			expectedResponseSubstr: "22",
		},
		{
			name:                   "empty matrix",
			filePath:               "testdata/empty.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "failed to parse csv: file is empty",
		},
		{
			name:                   "matrix with header",
			filePath:               "testdata/header.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix has a header row or non-integer value at row 1, column 1",
		},
		{
			name:                   "matrix has different row length",
			filePath:               "testdata/different_row_length.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix is not square: row 2 has 2 columns, expected 3",
		},
		{
			name:                   "matrix with empty value",
			filePath:               "testdata/empty_value.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix has empty value at row 3",
		},
		{
			name:                   "matrix with non-integer value",
			filePath:               "testdata/non_integer_value.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix value at row 3, column 2 is not an integer",
		},
		{
			name:                   "matrix has more rows than columns",
			filePath:               "testdata/more_rows_than_cols.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix is not square: 3 rows and 2 columns",
		},
		{
			name:                   "matrix has more columns than rows",
			filePath:               "testdata/more_cols_than_rows.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix is not square: 2 rows and 3 columns",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			req := s.createCSVRequest("/sum", tc.filePath)
			w := httptest.NewRecorder()
			SumHandler(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			s.Equal(tc.expectedStatusCode, resp.StatusCode)
			s.Contains(string(body), tc.expectedResponseSubstr)
		})
	}
}

// Test for multiply endpoint
func (s *EndpointTestSuite) TestMultiplyEndpoint() {
	tests := []struct {
		name                   string
		filePath               string
		expectedStatusCode     int
		expectedResponseSubstr string
	}{
		{
			name:                   "valid 3*3 matrix",
			filePath:               "testdata/valid_3_to_3.csv",
			expectedStatusCode:     200,
			expectedResponseSubstr: "362880",
		},
		{
			name:                   "valid 2*2 matrix",
			filePath:               "testdata/valid_2_to_2.csv",
			expectedStatusCode:     200,
			expectedResponseSubstr: "0",
		},
		{
			name:                   "valid 4*4 matrix",
			filePath:               "testdata/valid_4_to_4.csv",
			expectedStatusCode:     200,
			expectedResponseSubstr: "-2073600",
		},
		{
			name:                   "empty matrix",
			filePath:               "testdata/empty.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "failed to parse csv: file is empty",
		},
		{
			name:                   "matrix with header",
			filePath:               "testdata/header.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix has a header row or non-integer value at row 1, column 1",
		},
		{
			name:                   "matrix has different row length",
			filePath:               "testdata/different_row_length.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix is not square: row 2 has 2 columns, expected 3",
		},
		{
			name:                   "matrix with empty value",
			filePath:               "testdata/empty_value.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix has empty value at row 3",
		},
		{
			name:                   "matrix with non-integer value",
			filePath:               "testdata/non_integer_value.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix value at row 3, column 2 is not an integer",
		},
		{
			name:                   "matrix has more rows than columns",
			filePath:               "testdata/more_rows_than_cols.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix is not square: 3 rows and 2 columns",
		},
		{
			name:                   "matrix has more columns than rows",
			filePath:               "testdata/more_cols_than_rows.csv",
			expectedStatusCode:     400,
			expectedResponseSubstr: "matrix is not square: 2 rows and 3 columns",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			req := s.createCSVRequest("/multiply", tc.filePath)
			w := httptest.NewRecorder()
			MultiplyHandler(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			s.Equal(tc.expectedStatusCode, resp.StatusCode)
			s.Contains(string(body), tc.expectedResponseSubstr)
		})
	}
}

// Run all tests
func TestEndpointTestSuite(t *testing.T) {
	suite.Run(t, new(EndpointTestSuite))
}