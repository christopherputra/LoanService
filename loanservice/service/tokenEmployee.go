package service

import (
	"fmt"
	"loanservice/database"
	funcs "loanservice/functions"
	"net/http"
)

type TokenEmployeeReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

type TokenEmployeeRes struct {
	Token string `json:"token"`
}

func (w *Service) TokenEmployee(req TokenEmployeeReq, res *TokenEmployeeRes) (int, error) {
	statusCode := http.StatusInternalServerError

	reqDb := database.GetEmployeeWithCredentialsReqDb{
		Email:    req.Email,
		Password: req.Password,
	}
	resDb, err := w.PgServer.GetEmployeeWithCredentials(reqDb)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("email or password is invalid")
	}
	token, err := funcs.GenerateToken(resDb.EmployeeXid)
	if err != nil {
		return statusCode, err
	}
	*res = TokenEmployeeRes{
		Token: token,
	}
	return http.StatusOK, nil
}
