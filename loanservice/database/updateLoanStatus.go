package database

import (
	"database/sql"
)

type UpdateLoanStatusReqDb struct {
	LoanId int64  `db:"id"`
	Status string `db:"status"`
}
type UpdateLoanStatusResDb struct {
	LoanDb
}

func (pg *PostgresServer) UpdateLoanStatus(req UpdateLoanStatusReqDb, tx *sql.Tx) (res UpdateLoanStatusResDb, err error) {
	res = UpdateLoanStatusResDb{}
	query := `UPDATE loans SET status = $1 WHERE id = $2 RETURNING loan_xid, borrower_id, principal_amount, rate, roi, status`
	if tx != nil {
		rows, err := tx.Query(query, req.Status, req.LoanId)
		if err != nil {
			return res, err
		}
		for rows.Next() {
			err := rows.Scan(&res.LoanXid, &res.BorrowerId, &res.PrincipalAmount, &res.Rate, &res.Roi, &res.Status)
			if err != nil {
				return UpdateLoanStatusResDb{}, err
			}
		}
	} else {
		rows, err := pg.Db.Query(query, req.Status, req.LoanId)
		if err != nil {
			return res, err
		}
		for rows.Next() {
			err := rows.Scan(&res.LoanXid, &res.BorrowerId, &res.PrincipalAmount, &res.Rate, &res.Roi, &res.Status)
			if err != nil {
				return UpdateLoanStatusResDb{}, err
			}
		}
	}

	return res, nil
}
