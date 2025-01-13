package service

import (
	"fmt"
	"loanservice/database"
	funcs "loanservice/functions"
	"net/http"
)

type TokenAccountReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

type TokenAccountRes struct {
	Token string `json:"token"`
}

func (w *Service) TokenAccount(req TokenAccountReq, res *TokenAccountRes) (int, error) {
	statusCode := http.StatusInternalServerError

	reqDb := database.GetAccountWithCredentialsReqDb{
		Email:    req.Email,
		Password: req.Password,
	}
	resDb, err := w.PgServer.GetAccountWithCredentials(reqDb)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("email or password is invalid")
	}
	token, err := funcs.GenerateToken(resDb.CustomerXid)
	if err != nil {
		return statusCode, err
	}
	*res = TokenAccountRes{
		Token: token,
	}
	return http.StatusOK, nil
}
