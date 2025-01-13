package service

import (
	"loanservice/constants"
	"loanservice/database"
	"math"

	"github.com/rs/xid"
)

type ProposeLoanReq struct {
	Amount int64 `json:"amount" binding:"required,gte=100000"`
}

type ProposeLoanRes struct {
	LoanXid string `json:"loan_xid"`
	Amount  int64  `json:"amount"`
	Status  string `json:"status"`
}

func (w *Service) ProposeLoan(authUser AuthenticatedUser, req ProposeLoanReq, res *ProposeLoanRes) error {
	rate := int64(math.Ceil(0.1 * float64(req.Amount)))
	roi := int64(0)
	var loanXid = xid.New().String()

	reqDb := database.InsertLoanReqDb{
		LoanXid:    loanXid,
		Amount:     req.Amount,
		Rate:       rate,
		Roi:        roi,
		BorrowerId: authUser.AccountId,
		Status:     constants.LoanStatusProposed,
	}
	resDb, err := w.PgServer.InsertLoan(reqDb)
	if err != nil {
		return err
	}
	*res = ProposeLoanRes{
		LoanXid: resDb.LoanXid,
		Amount:  resDb.PrincipalAmount,
		Status:  resDb.Status,
	}

	return nil
}
