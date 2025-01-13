package database

type CreateAccountReqDb struct {
	CustomerXid string
	Name        string
	Email       string
	Password    string
}

func (pg *PostgresServer) CreateOrInsertAccount(req CreateAccountReqDb) error {
	_, err := pg.Db.Query(`INSERT INTO accounts (customer_xid, name, email, password)
	VALUES ($1, $2, $3, $4)`, req.CustomerXid, req.Name, req.Email, req.Password)
	if err != nil {
		return err
	}
	return nil
}
