package database

import (
	"database/sql"
	"time"
)

type InsertLoanApprovalReqDb struct {
	LoanId                 int64
	ApprovedBy             int64
	ApprovedDate           time.Time
	ApprovalProofImagePath string
}

func (pg *PostgresServer) InsertLoanApproval(req InsertLoanApprovalReqDb, tx *sql.Tx) (err error) {
	query := `INSERT INTO loan_approvals (loan_id, approved_by, approved_date, approval_proof_image_path)
	VALUES ($1, $2, $3, $4)`
	if tx != nil {
		_, err = tx.Exec(query, req.LoanId, req.ApprovedBy, req.ApprovedDate, req.ApprovalProofImagePath)
	} else {
		_, err = pg.Db.Exec(query, req.LoanId, req.ApprovedBy, req.ApprovedDate, req.ApprovalProofImagePath)
	}

	if err != nil {
		return err
	}
	return nil
}
