package database

type CreateEmployeeReqDb struct {
	EmployeeXid string
	Name        string
	Email       string
	Password    string
}

func (pg *PostgresServer) CreateOrInsertEmployee(req CreateEmployeeReqDb) error {
	_, err := pg.Db.Query(`INSERT INTO employees (employee_xid, name, email, password)
	VALUES ($1, $2, $3, $4)`, req.EmployeeXid, req.Name, req.Email, req.Password)
	if err != nil {
		return err
	}
	return nil
}
