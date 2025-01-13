package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresClient interface {
	StartTransaction() (*sql.Tx, error)
	CommitTransaction(*sql.Tx) error
	CreateOrInsertAccount(req CreateAccountReqDb) error
	CreateOrInsertEmployee(req CreateEmployeeReqDb) error
	GetAccount(req GetAccountReqDb) (GetAccountResDb, error)
	GetEmployee(req GetEmployeeReqDb) (GetEmployeeResDb, error)
	GetAccountWithCredentials(req GetAccountWithCredentialsReqDb) (GetAccountWithCredentialsResDb, error)
	GetEmployeeWithCredentials(req GetEmployeeWithCredentialsReqDb) (GetEmployeeWithCredentialsResDb, error)
	GetLoan(req GetLoanReqDb) (GetLoanResDb, error)
	InsertLoan(req InsertLoanReqDb) (InsertLoanResDb, error)
	InsertLoanApproval(req InsertLoanApprovalReqDb, tx *sql.Tx) error
	UpdateLoanStatus(req UpdateLoanStatusReqDb, tx *sql.Tx) (UpdateLoanStatusResDb, error)
	GetLoanInvestments(req GetLoanInvestmentsReqDb) (GetLoanInvestmentsResDb, error)
	InsertLoanInvestment(req InsertLoanInvestmentReqDb, tx *sql.Tx) (bool, error)
	InsertLoanDisbursement(req InsertLoanDisbursementReqDb, tx *sql.Tx) (bool, error)
}
type PostgresServer struct {
	Db *sql.DB
}

func NewPostgresClient(host string, port int, user string, password string, dbname string) PostgresClient {
	conn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	dbserver, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}
	return &PostgresServer{
		Db: dbserver,
	}
}
