# ⚙️ Instructions - How to install and run project

Follow these steps to run the project locally on your machine.

### <a name="prerequisites">⭐ Prerequisites</a>

Make sure you have the following installed on your machine:

- Git
- Go

### <a name="clone-repo">⭐ Cloning the Repository</a>

```bash
git clone https://github.com/EvanHuang7/league_code_test
```

### <a name="install-packages">⭐ Packages Installation</a>

Install the project dependencies using npm:

```bash
cd league_code_test
go mod tidy
```

### <a name="running-project">⭐ Running the Project</a>

Open **a terminal window** and run the following commands to start the project:

**1st Terminal** – Start the Project (Go Server):

```bash
cd league_code_test
go run .
```

Open **a new terminal window** and run the following commands to send request to server:

**📌 Note**: You could use different test input files from `testdata` folder by updating the file path in CLI, such as changing `@matrix.csv` to `@testdata/valid_2_to_2.csv`.

**2nd Terminal** – Send request to Go Server:

```bash
cd league_code_test

curl -F 'file=@matrix.csv' "localhost:8080/echo"
curl -F 'file=@matrix.csv' "localhost:8080/invert"
curl -F 'file=@matrix.csv' "localhost:8080/flatten"
curl -F 'file=@matrix.csv' "localhost:8080/sum"
curl -F 'file=@matrix.csv' "localhost:8080/multiply"
```

### <a name="running-test">⭐ Run All Tests</a>

Open **a terminal window** and run the following commands to run all tests:

```bash
cd league_code_test
go test -v
```

# League Backend Challenge

In main.go you will find a basic web server written in GoLang. It accepts a single request _/echo_. Extend the webservice with the ability to perform the following operations

Given an uploaded csv file
```
1,2,3
4,5,6
7,8,9
```

1. Echo (given)
    - Return the matrix as a string in matrix format.
    
    ```
    // Expected output
    1,2,3
    4,5,6
    7,8,9
    ``` 
2. Invert
    - Return the matrix as a string in matrix format where the columns and rows are inverted
    ```
    // Expected output
    1,4,7
    2,5,8
    3,6,9
    ``` 
3. Flatten
    - Return the matrix as a 1 line string, with values separated by commas.
    ```
    // Expected output
    1,2,3,4,5,6,7,8,9
    ``` 
4. Sum
    - Return the sum of the integers in the matrix
    ```
    // Expected output
    45
    ``` 
5. Multiply
    - Return the product of the integers in the matrix
    ```
    // Expected output
    362880
    ``` 

The input file to these functions is a matrix, of any dimension where the number of rows are equal to the number of columns (square). Each value is an integer, and there is no header row. matrix.csv is example valid input.  

Run web server
```
go run .
```

Send request
```
curl -F 'file=@/path/matrix.csv' "localhost:8080/echo"
```

## What we're looking for

- The solution runs
- The solution performs all cases correctly
- The code is easy to read
- The code is reasonably documented
- The code is tested
- The code is robust and handles invalid input and provides helpful error messages
