package service

import (
	"fmt"
	"loanservice/database"
	funcs "loanservice/functions"
	"net/http"

	"github.com/lib/pq"
	"github.com/rs/xid"
)

type InitAccountReq struct {
	Name     string `json:"name" binding:"required,min=3,max=36"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

type InitAccountRes struct {
	Token string `json:"token"`
}

func (w *Service) InitAccount(req InitAccountReq, res *InitAccountRes) (int, error) {
	statusCode := http.StatusInternalServerError
	var customerXid = xid.New().String()

	reqDb := database.CreateAccountReqDb{
		Name:        req.Name,
		CustomerXid: customerXid,
		Email:       req.Email,
		Password:    req.Password,
	}
	err := w.PgServer.CreateOrInsertAccount(reqDb)
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
	token, err := funcs.GenerateToken(customerXid)
	if err != nil {
		return statusCode, err
	}
	*res = InitAccountRes{
		Token: token,
	}
	return http.StatusOK, nil
}
