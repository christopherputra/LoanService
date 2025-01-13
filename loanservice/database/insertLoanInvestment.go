package database

import "database/sql"

type InsertLoanInvestmentReqDb struct {
	LoanId           int64
	Amount           int64
	InvestedBy       int64
	AgreementPdfPath string
}

func (pg *PostgresServer) InsertLoanInvestment(req InsertLoanInvestmentReqDb, tx *sql.Tx) (bool, error) {
	res, err := pg.Db.Exec(`INSERT INTO loan_investments (loan_id, amount, invested_by, agreement_pdf_path)
	VALUES ($1, $2, $3, $4) RETURNING id, loan_id, amount, invested_by, agreement_pdf_path`, req.LoanId, req.Amount, req.InvestedBy, req.AgreementPdfPath)
	if err != nil {
		return false, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected > 0 {
		return true, nil
	}
	return false, nil
}
