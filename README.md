# GoAPI

## Usage
- Usage of the API is documented with go-swagger
- You can access the Swagger UI by going to this URL (after sucessfully running the API on Docker): `http://localhost:8080/docs`

## Used components with versions
- Docker 2.3.0.4(46911)
- GO Language go1.15.1 windows/amd64

## Prerequisites
- Have the latest version of Docker Desktop installed

## Running the API
1. Download the source code from GitHub and unzip it in a desired location
2. In CMD (Command Line Terminal) move to the downloaded repository on your local machine (to the GoAPI folder)
3. Run the following command `docker-compose up --build` to build and set up the API in your Docker
4. After sucessfully building and deploying the API on Docker, you can close the Terminal and run the API directly from Docker

## Test cases
- All the test cases are included in `endpoints_test.go`
- You can run the tests by simply typing `go test -v` from your Terminal in the folder that you've downloaded and unzipped the source code (GoAPI)
- There are 8 test cases in total (for each function 1)

## Improvements
- Add Encryption/Decryption to Passwords
- Move some functions in other classes for better readability
- More in depth test cases