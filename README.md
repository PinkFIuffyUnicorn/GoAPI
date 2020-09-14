# GoAPI

## Usage
- Usage of the API is documented with go-swagger
- You can access to the Swagger UI by going to this URL (after you've sucessfully ran the API on Docker): `http://localhost:8080/docs`

## Used components with versions
- Docker 2.3.0.4(46911)
- GO Language go1.15.1 windows/amd64

## Prerequisites
- Have the latest version of Docker Desktop installed

## Running the API
1. Download the source code from GitHub and unzip it in a desired location
2. In CMD (Command Line Terminal) move to the downloaded repository on your local machine (to the GoAPI folder)
3. Run the following command `docker-compose up --build` to build and set up the API in your Docker

## Improvements
- Add Encryption/Decryption to Passwords
- Move some functions in other classes for better readability