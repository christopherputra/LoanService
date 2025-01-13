package service

import (
	"loanservice/database"
)

type AuthenticatedUser struct {
	AccountId   int64  `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	CustomerXid string `json:"customer_xid"`
}

type AuthenticatedEmployee struct {
	EmployeeId  int64  `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	EmployeeXid string `json:"employee_xid"`
}

type Wallet struct {
	WalletId string `json:"id"`
	OwnerId  string `json:"owned_by"`
	Status   string `json:"status"`
	Balance  int64  `json:"balance"`
}

type Service struct {
	PgServer database.PostgresClient
}

func NewLoanService(pgserver database.PostgresClient) Service {
	return Service{
		PgServer: pgserver,
	}
}
