package service

import (
	"database/sql"
	"fmt"
	"loanservice/constants"
	"loanservice/database"
	"log"
	"net/http"

	"time"
)

type ApproveLoanReq struct {
	LoanXid string `json:"loan_xid" binding:"required"`
}

type ApproveLoanRes struct {
	LoanXid string `json:"loan_xid"`
	Amount  int64  `json:"amount"`
	Rate    int64  `json:"rate"`
	Status  string `json:"status"`
}

func (w *Service) ApproveLoan(authEmployee AuthenticatedEmployee, req ApproveLoanReq, res *ApproveLoanRes) (int, error) {
	reqDb := database.GetLoanReqDb{
		LoanXid: req.LoanXid,
	}
	loan, err := w.PgServer.GetLoan(reqDb)
	if err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, fmt.Errorf("loan id not found")
		}
		return http.StatusInternalServerError, err
	}
	if loan.Status != constants.LoanStatusProposed {
		return http.StatusBadRequest, fmt.Errorf("invalid loan id")
	}

	tx, err := w.PgServer.StartTransaction()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer func() {
		// If an error occurs, Rollback the transaction
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("transaction rollback failed: %v", rbErr)
			}
			log.Printf("transaction rollbacked")
		}
	}()
	reqApprovalDb := database.InsertLoanApprovalReqDb{
		LoanId:                 loan.LoanId,
		ApprovedBy:             authEmployee.EmployeeId,
		ApprovedDate:           time.Now(),
		ApprovalProofImagePath: "temporary_image_path/image.jpg",
	}
	err = w.PgServer.InsertLoanApproval(reqApprovalDb, tx)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	reqUpdateLoanDb := database.UpdateLoanStatusReqDb{
		LoanId: loan.LoanId,
		Status: constants.LoanStatusApproved,
	}
	resUpdateDB, err := w.PgServer.UpdateLoanStatus(reqUpdateLoanDb, tx)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	*res = ApproveLoanRes{
		LoanXid: resUpdateDB.LoanXid,
		Amount:  resUpdateDB.PrincipalAmount,
		Rate:    resUpdateDB.Rate,
		Status:  resUpdateDB.Status,
	}

	err = w.PgServer.CommitTransaction(tx)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
