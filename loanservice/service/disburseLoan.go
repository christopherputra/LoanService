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

type DisburseLoanReq struct {
	LoanXid string `json:"loan_xid" binding:"required"`
}

type DisburseLoanRes struct {
	LoanXid          string `json:"loan_xid"`
	AmountDisbursed  int64  `json:"amount_disbursed"`
	DisbursedDate    string `json:"disbursed_date"`
	AgreementPdfPath string `json:"agreement_pdf_path"`
}

func (w *Service) DisburseLoan(authEmployee AuthenticatedEmployee, req DisburseLoanReq, res *DisburseLoanRes) (int, error) {
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
	if loan.Status != constants.LoanStatusInvested {
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
	pdf_agreement := "/temp_agreement/agreement.pdf"
	reqDisbursementDb := database.InsertLoanDisbursementReqDb{
		LoanId:           loan.LoanId,
		EmployeeId:       authEmployee.EmployeeId,
		AgreementPdfPath: pdf_agreement,
		DisbursedDate:    time.Now(),
	}
	inserted, err := w.PgServer.InsertLoanDisbursement(reqDisbursementDb, tx)
	if err != nil || !inserted {
		return http.StatusInternalServerError, err
	}

	reqUpdateLoanDb := database.UpdateLoanStatusReqDb{
		LoanId: loan.LoanId,
		Status: constants.LoanStatusDisbursed,
	}
	resUpdateDB, err := w.PgServer.UpdateLoanStatus(reqUpdateLoanDb, tx)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	*res = DisburseLoanRes{
		LoanXid:          resUpdateDB.LoanXid,
		AmountDisbursed:  resUpdateDB.PrincipalAmount,
		DisbursedDate:    time.Now().Format("2006-01-02"),
		AgreementPdfPath: pdf_agreement,
	}

	err = w.PgServer.CommitTransaction(tx)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
