package database

import (
	"database/sql"
	"time"
)

type InsertLoanDisbursementReqDb struct {
	LoanId           int64
	EmployeeId       int64
	AgreementPdfPath string
	DisbursedDate    time.Time
}

func (pg *PostgresServer) InsertLoanDisbursement(req InsertLoanDisbursementReqDb, tx *sql.Tx) (bool, error) {
	res, err := pg.Db.Exec(`INSERT INTO loan_disbursements (loan_id, employee_id, agreement_pdf_path, disbursed_date)
	VALUES ($1, $2, $3, $4)`, req.LoanId, req.EmployeeId, req.AgreementPdfPath, req.DisbursedDate)
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
