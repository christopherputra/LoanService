package service

import (
	"fmt"
	"loanservice/database"
	funcs "loanservice/functions"
	"net/http"

	"github.com/lib/pq"
	"github.com/rs/xid"
)

type InitEmployeeReq struct {
	Name     string `json:"name" binding:"required,min=3,max=36"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

type InitEmployeeRes struct {
	Token string `json:"token"`
}

func (w *Service) InitEmployee(req InitEmployeeReq, res *InitEmployeeRes) (int, error) {
	statusCode := http.StatusInternalServerError
	var employeeXid = xid.New().String()

	reqDb := database.CreateEmployeeReqDb{
		Name:        req.Name,
		EmployeeXid: employeeXid,
		Email:       req.Email,
		Password:    req.Password,
	}
	err := w.PgServer.CreateOrInsertEmployee(reqDb)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			// PostgreSQL error code for unique violation is 23505
			if pqErr.Code == "23505" {
				statusCode = http.StatusBadRequest
				return statusCode, fmt.Errorf("email is already taken")
			}
		}
		return statusCode, err
	}
	token, err := funcs.GenerateToken(employeeXid)
	if err != nil {
		return statusCode, err
	}
	*res = InitEmployeeRes{
		Token: token,
	}
	return http.StatusOK, nil
}
