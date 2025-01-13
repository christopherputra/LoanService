package service

import (
	"database/sql"
	"fmt"
	"loanservice/constants"
	"loanservice/database"
	"log"
	"net/http"
)

type InvestLoanReq struct {
	LoanXid string `json:"loan_xid" binding:"required"`
	Amount  int64  `json:"amount" binding:"required,gt=0"`
}

type InvestLoanRes struct {
	LoanXid          string `json:"loan_xid"`
	Amount           int64  `json:"amount"`
	Status           string `json:"status"`
	InvestedAmount   int64  `json:"invested_amount"`
	AgreementPdfPath string `json:"agreement_pdf_path"`
}

func (w *Service) InvestLoan(authUser AuthenticatedUser, req InvestLoanReq, res *InvestLoanRes) (int, error) {
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
	status := loan.Status
	if status != constants.LoanStatusApproved {
		return http.StatusBadRequest, fmt.Errorf("invalid loan id")
	}

	reqInvestmentsDb := database.GetLoanInvestmentsReqDb{
		LoanId: loan.LoanId,
	}
	resInvestmentsDb, err := w.PgServer.GetLoanInvestments(reqInvestmentsDb)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	total_invested := int64(0)
	for _, investment := range resInvestmentsDb.LoanInvestments {
		total_invested = total_invested + investment.Amount
	}

	if total_invested+req.Amount > loan.PrincipalAmount {
		return http.StatusBadRequest, fmt.Errorf("amount of investment is above the remaining amount")
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
	pdf_url := "temporary_pdf_path/agreement.pdf"
	reqInvestmentDb := database.InsertLoanInvestmentReqDb{
		LoanId:           loan.LoanId,
		Amount:           req.Amount,
		InvestedBy:       authUser.AccountId,
		AgreementPdfPath: pdf_url,
	}
	inserted, err := w.PgServer.InsertLoanInvestment(reqInvestmentDb, tx)
	if err != nil || !inserted {
		return http.StatusInternalServerError, err
	}

	if total_invested+req.Amount == loan.PrincipalAmount {
		// invest done
		status = constants.LoanStatusInvested
		reqUpdateLoanDb := database.UpdateLoanStatusReqDb{
			LoanId: loan.LoanId,
			Status: status,
		}
		_, err := w.PgServer.UpdateLoanStatus(reqUpdateLoanDb, tx)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	*res = InvestLoanRes{
		LoanXid:          loan.LoanXid,
		Amount:           loan.PrincipalAmount,
		Status:           status,
		AgreementPdfPath: pdf_url,
		InvestedAmount:   total_invested + req.Amount,
	}

	err = w.PgServer.CommitTransaction(tx)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
