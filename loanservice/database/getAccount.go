package database

type GetAccountReqDb struct {
	CustomerXid string
}
type GetAccountResDb struct {
	AccountId   int64  `db:"id"`
	Email       string `db:"email"`
	Name        string `db:"name"`
	CustomerXid string `db:"customer_xid"`
}

func (pg *PostgresServer) GetAccount(req GetAccountReqDb) (GetAccountResDb, error) {
	res := GetAccountResDb{}
	err := pg.Db.QueryRow(`SELECT id, email, name, customer_xid FROM accounts WHERE customer_xid = $1`, req.CustomerXid).Scan(&res.AccountId, &res.Name, &res.Email, &res.CustomerXid)
	if err != nil {
		return GetAccountResDb{}, err
	}
	return res, nil
}
