package database

type GetAccountWithCredentialsReqDb struct {
	Email    string
	Password string
}
type GetAccountWithCredentialsResDb struct {
	CustomerXid string `db:"customer_xid"`
}

func (pg *PostgresServer) GetAccountWithCredentials(req GetAccountWithCredentialsReqDb) (GetAccountWithCredentialsResDb, error) {
	res := GetAccountWithCredentialsResDb{}
	err := pg.Db.QueryRow(`SELECT customer_xid FROM accounts WHERE email = $1 and password = $2`, req.Email, req.Password).Scan(&res.CustomerXid)
	if err != nil {
		return GetAccountWithCredentialsResDb{}, err
	}
	return res, nil
}
