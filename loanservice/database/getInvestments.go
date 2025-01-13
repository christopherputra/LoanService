package database

type GetLoanInvestmentsReqDb struct {
	LoanId int64
}
type GetLoanInvestmentsResDb struct {
	LoanInvestments []LoanInvestmentDb
}
type LoanInvestmentDb struct {
	LoanInvestmentId int64  `db:"id"`
	LoanId           int64  `db:"loan_id"`
	InvestedBy       int64  `db:"invested_by"`
	Amount           int64  `db:"amount"`
	AgreementPdfPath string `db:"agreement_pdf_path"`
}

func (pg *PostgresServer) GetLoanInvestments(req GetLoanInvestmentsReqDb) (GetLoanInvestmentsResDb, error) {
	res := GetLoanInvestmentsResDb{}
	rows, err := pg.Db.Query(`SELECT id, loan_id, amount, invested_by, agreement_pdf_path FROM loan_investments WHERE loan_id = $1`, req.LoanId)
	if err != nil {
		return GetLoanInvestmentsResDb{}, err
	}
	loanInvestments := make([]LoanInvestmentDb, 0)
	for rows.Next() {
		rowResult := LoanInvestmentDb{}
		err := rows.Scan(&rowResult.LoanInvestmentId, &rowResult.LoanId, &rowResult.Amount, &rowResult.InvestedBy, &rowResult.AgreementPdfPath)
		if err != nil {
			return GetLoanInvestmentsResDb{}, err
		}
		loanInvestments = append(loanInvestments, rowResult)
	}
	res.LoanInvestments = loanInvestments
	return res, nil
}
