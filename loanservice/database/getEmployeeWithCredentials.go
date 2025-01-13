package database

type GetEmployeeWithCredentialsReqDb struct {
	Email    string
	Password string
}
type GetEmployeeWithCredentialsResDb struct {
	EmployeeXid string `db:"employee_xid"`
}

func (pg *PostgresServer) GetEmployeeWithCredentials(req GetEmployeeWithCredentialsReqDb) (GetEmployeeWithCredentialsResDb, error) {
	res := GetEmployeeWithCredentialsResDb{}
	err := pg.Db.QueryRow(`SELECT employee_xid FROM employees WHERE email = $1 and password = $2`, req.Email, req.Password).Scan(&res.EmployeeXid)
	if err != nil {
		return GetEmployeeWithCredentialsResDb{}, err
	}
	return res, nil
}
