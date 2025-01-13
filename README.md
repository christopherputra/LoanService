# Loan Service Golang
### About
Loan Service is a service that handles loans starting from init, propose, approved, invested, and disbursed. There are several features inside this service including several added features based on my own assumption:
1. Init Account & Init Employee (Assumption: needed for login to handle authorization tokens)
2. Token Account & Token Employee (Assumption: an API for token generation, this token is based on JWT implementation, put this token inside all of headers of the APIs in Authorization header)
3. Propose Loan (Account's API)
4. Approve Loan (Employee's API)
5. Invest Loan (Accounts)

### Assumption and Disclaimer
1. **JWT TOKEN SIGNING** Since this is an API design, so I decided to add authorization for the APIs using JWT token implementation. This token is used in all of the APIs headers except the Init and Token API.
2. **Sign Up & Login** Since I want to implement JWT, therefore I add Signup feature which are the IniAccount API & InitEmployee API, and the Login feature which are the TokenAccount API and TokenEmployee API.
3. **Password Encryption** I know the password should be encrypted at least using hashing functions such as sha256 etc before stored in the database. However, I haven't implemented it in this encryption.
4. **Any third party use is not implemneted** third party such as emailer and everything else is provided as diagrams only since It wont be able to be run on all local system also. SMTP gmail is not free anymore.

### Brief Architecture - Service Based API
[![ARCHITECTURE](https://imgur.com/ugDpOpr,jpg "ARCHITECTURE")](https://imgur.com/ugDpOpr.jpg "ARCHITECTURE")
### Stacks
- Golang
- JWT Token Signing https://pkg.go.dev/github.com/golang-jwt/jwt/v5
- Go Gin Routing https://github.com/gin-gonic/gin
- Postgresql

| [![go](https://i.imgur.com/miVUk6U.png "go")](https://i.imgur.com/miVUk6U.png "go")  | [![go gin](https://i.imgur.com/8OTwAo4.png "go gin")](https://i.imgur.com/8OTwAo4.png "go gin")  |  [![jwt](https://i.imgur.com/2GujZmD.png "jwt")](https://i.imgur.com/2GujZmD.png "jwt") |  [![postgres](https://i.imgur.com/dLxfiGU.png "postgres")](https://i.imgur.com/dLxfiGU.png "postgres") |
| ------------ | ------------ | ------------ | ------------ |
|   |   |   |   |    |
### Code Structure
```
loanservice
│
└───database
│   │  accessorName.go
│   │  accessorName2.go
│   │   ...
│   │
└───functions
│   │   function1.go
│   │   function2.go
│   │   ...
│   
└───middlewares
│   │   middleware1.go
│   │   middleware2.go
│   │  ...
└───service
│   │   service.go
│   │   endpoints1.go
│   │   endpoints2.go
│   │   ...
│   │
│   └───handler
│       │   handler.go
│   
│   config.json
│   go.mod
│   go.sum
│   main.go
```
## API Documentation
#### Postman Collection
**Please import this postman collection /LoanService.postman_collection.json to see the APIs including the body and the headers.**

## Getting Started
#### Golang Install
Download Go https://go.dev/doc/install

#### PostgreSQL Install
Please install postgresql on your local following these steps https://www.postgresql.org/download/macosx/

## First Setup
#### PostgreSQL Schema
make `loanservice` database
run `schema.sql` on your postgres server to create all tables
```sql
DROP TABLE IF EXISTS loan_investments;
DROP TABLE IF EXISTS loan_approvals;
DROP TABLE IF EXISTS loan_disbursements;
DROP TABLE IF EXISTS loans;
DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS accounts;

CREATE TABLE accounts (
	id BIGSERIAL PRIMARY KEY,
	customer_xid VARCHAR ( 36 ) UNIQUE NOT NULL,
	name VARCHAR ( 36 ) NOT NULL,
	email VARCHAR ( 36 ) UNIQUE NOT NULL,
	password VARCHAR ( 36 ) NOT NULL
);

CREATE TABLE employees (
	id BIGSERIAL PRIMARY KEY,
	employee_xid VARCHAR ( 36 ) UNIQUE NOT NULL,
	name VARCHAR ( 36 ) NOT NULL,
	email VARCHAR ( 36 ) UNIQUE NOT NULL,
	password VARCHAR ( 36 ) NOT NULL
);

CREATE TABLE loans (
	id BIGSERIAL PRIMARY KEY,
	loan_xid VARCHAR ( 36 ) UNIQUE NOT NULL,
	borrower_id BIGINT NOT NULL,
	principal_amount INT NOT NULL,
	rate INT NOT NULL,
	roi INT NOT NULL,
	status VARCHAR ( 10 ) NOT NULL,
	FOREIGN KEY (borrower_id) REFERENCES accounts(id)
);

CREATE TABLE loan_approvals (
	id BIGSERIAL PRIMARY KEY,
	loan_id BIGINT NOT NULL,
	approved_by BIGINT NOT NULL,
	approved_date DATE NOT NULL,
	approval_proof_image_path TEXT NOT NULL,
	FOREIGN KEY (approved_by) REFERENCES employees(id)
);

CREATE TABLE loan_investments (
	id BIGSERIAL PRIMARY KEY,
	loan_id BIGINT NOT NULL,
	amount INT NOT NULL,
	invested_by BIGINT NOT NULL,
	agreement_pdf_path TEXT NOT NULL,
	FOREIGN KEY (invested_by) REFERENCES accounts(id)
);

CREATE TABLE loan_disbursements (
	id BIGSERIAL PRIMARY KEY,
	loan_id BIGINT NOT NULL,
	employee_id BIGINT NOT NULL,
	agreement_pdf_path TEXT NOT NULL,
	disbursed_date DATE NOT NULL,
	FOREIGN KEY (employee_id) REFERENCES employees(id)
);
```

#### Setup config.json File
Change env fields based on your local env including postgres host, port, username, and password.
```json
{
    "postgres": {
        "host": "localhost",
        "port": 5432,
        "user": "<username>",
        "password": "<password>",
        "dbname": "loanservice"
    },
    "host": "localhost",
    "port": "8080"
}
```
#### Go Get Packages
Please `go get <packages not yet installed>`

## You're Good To GO - Run Service
on your console inside `/loanservice` dir
run go server using this command
`go run main.go`

## Unit Tests
All of the endpoints are already tested in unit granularity.
if you wish to run the test, go to the directive folder with the tests files and run `go test`


