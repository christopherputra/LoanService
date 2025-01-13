package database

type GetEmployeeReqDb struct {
	EmployeeXid string
}
type GetEmployeeResDb struct {
	EmployeeId  int64  `db:"id"`
	Email       string `db:"email"`
	Name        string `db:"name"`
	EmployeeXid string `db:"employee_xid"`
}

func (pg *PostgresServer) GetEmployee(req GetEmployeeReqDb) (GetEmployeeResDb, error) {
	res := GetEmployeeResDb{}
	err := pg.Db.QueryRow(`SELECT id, email, name, employee_xid FROM employees WHERE employee_xid = $1`, req.EmployeeXid).Scan(&res.EmployeeId, &res.Name, &res.Email, &res.EmployeeXid)
	if err != nil {
		return GetEmployeeResDb{}, err
	}
	return res, nil
}
