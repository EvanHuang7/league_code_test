package main

import (
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

// Run with
//		go run .
// Send request with:
//		curl -F 'file=@matrix.csv' "localhost:8080/echo"

// Return the matrix as a string in matrix format
func EchoHandler(w http.ResponseWriter, r *http.Request) {
	// 1st Step: get matrix from csv file
	records, err := ParseCSVFile(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// 2nd Step: validate input matrix
	err = ValidateSquareMatrix(records)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 3rd Step: build response
	var response string
	// Each row is a list of strings, such as ["1", "2", "3"]
	for _, row := range records {
		response = fmt.Sprintf("%s%s\n", response, strings.Join(row, ","))
	}

	fmt.Fprint(w, response)
}

// Return the matrix as a string in matrix format where the columns and rows are inverted
func InvertHandler(w http.ResponseWriter, r *http.Request) {
	// 1st Step: get matrix from csv file
	records, err := ParseCSVFile(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 2nd Step: validate input matrix
	err = ValidateSquareMatrix(records)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 3rd Step: build the inverted matrix
	invertedRecords := make([][]string, len(records))
	// Each row is a list of strings, such as ["1", "2", "3"]
	for _, row := range records {
		for idx, num := range row {
			invertedRecords[idx] = append(invertedRecords[idx], num)
		}
	}

	// 4th Step: build response
	var response string
	for _, row := range invertedRecords {
		response = fmt.Sprintf("%s%s\n", response, strings.Join(row, ","))
	}

	fmt.Fprint(w, response)
}

// Return the matrix as a 1 line string, with values separated by commas
func FlattenHandler(w http.ResponseWriter, r *http.Request) {
	// 1st Step: get matrix from csv file
	records, err := ParseCSVFile(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// 2nd Step: validate input matrix
	err = ValidateSquareMatrix(records)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 3rd Step: flatten the matrix to a single list
	var flattened []string
	// Each row is a list of strings, such as ["1", "2", "3"]
	for _, row := range records {
		flattened = append(flattened, row...)
	}

	// 4th Step: build response
	response := strings.Join(flattened, ",") + "\n"
	fmt.Fprint(w, response)
}

// Return the sum of the integers in the matrix
func SumHandler(w http.ResponseWriter, r *http.Request) {
	// 1st Step: get matrix from csv file
	records, err := ParseCSVFile(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// 2nd Step: validate input matrix
	err = ValidateSquareMatrix(records)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 3rd Step: calculate and return the sum
	sum := 0
	for _, row := range records {
		for _, val := range row {
			// no need for error check as we did it already in ValidateSquareMatrix()
			num, _ := strconv.Atoi(val)
			sum += num
		}
	}

	fmt.Fprintln(w, sum)
}

// Return the product of the integers in the matrix
func MultiplyHandler(w http.ResponseWriter, r *http.Request) {
	// 1st Step: get matrix from csv file
	records, err := ParseCSVFile(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// 2nd Step: validate input matrix
	err = ValidateSquareMatrix(records)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 3rd Step: calculate product using big.Int and return the product
	// Use "math/big" package and "int64" to handle large product integer case
	product := big.NewInt(1)
	for _, row := range records {
		for _, val := range row {
			// no need for error check as we did it already in ValidateSquareMatrix()
			num, _ := strconv.Atoi(val)
			product.Mul(product, big.NewInt(int64(num)))
		}
	}

	fmt.Fprintln(w, product.String())
}

func main() {
	http.HandleFunc("/echo", EchoHandler)
	http.HandleFunc("/invert", InvertHandler)
	http.HandleFunc("/flatten", FlattenHandler)
	http.HandleFunc("/sum", SumHandler)
	http.HandleFunc("/multiply", MultiplyHandler)

	port := ":8080"
	fmt.Println("Server started on port", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Println("Server failed:", err)
	}
}
