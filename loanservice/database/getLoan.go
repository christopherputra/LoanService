package database

type GetLoanReqDb struct {
	LoanXid string
}
type GetLoanResDb struct {
	LoanDb
}
type LoanDb struct {
	LoanId          int64  `db:"id"`
	LoanXid         string `db:"loan_xid"`
	BorrowerId      int64  `db:"borrower_id"`
	PrincipalAmount int64  `db:"principal_amount"`
	Rate            int64  `db:"rate"`
	Roi             int64  `db:"roi"`
	Status          string `db:"status"`
}

func (pg *PostgresServer) GetLoan(req GetLoanReqDb) (GetLoanResDb, error) {
	res := GetLoanResDb{}
	err := pg.Db.QueryRow(`SELECT id, loan_xid, borrower_id, principal_amount, rate, roi, status FROM loans WHERE loan_xid = $1`, req.LoanXid).Scan(&res.LoanId, &res.LoanXid, &res.BorrowerId, &res.PrincipalAmount, &res.Rate, &res.Roi, &res.Status)
	if err != nil {
		return GetLoanResDb{}, err
	}
	return res, nil
}
