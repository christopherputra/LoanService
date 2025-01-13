package database

type InsertLoanReqDb struct {
	LoanXid    string
	BorrowerId int64
	Amount     int64
	Rate       int64
	Roi        int64
	Status     string
}
type InsertLoanResDb struct {
	LoanDb
}

func (pg *PostgresServer) InsertLoan(req InsertLoanReqDb) (InsertLoanResDb, error) {
	res := InsertLoanResDb{}
	rows, err := pg.Db.Query(`INSERT INTO loans (loan_xid, borrower_id, principal_amount, rate, roi, status)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING loan_xid, borrower_id, principal_amount, rate, roi, status`, req.LoanXid, req.BorrowerId, req.Amount, req.Rate, req.Roi, req.Status)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		err := rows.Scan(&res.LoanXid, &res.BorrowerId, &res.PrincipalAmount, &res.Rate, &res.Roi, &res.Status)
		if err != nil {
			return InsertLoanResDb{}, err
		}
	}
	return res, nil
}
