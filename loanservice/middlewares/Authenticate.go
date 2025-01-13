package middlewares

import (
	"fmt"
	"loanservice/database"
	funcs "loanservice/functions"
	"loanservice/service"
	"net/http"

	"clevergo.tech/jsend"
	"github.com/gin-gonic/gin"
)

func TokenAuthenticateMiddleware(pgserver database.PostgresClient) gin.HandlerFunc {
	return func(c *gin.Context) {

		auth := c.GetHeader("Authorization")
		token := auth[len("Token "):]
		valid, customer_xid, err := funcs.VerifyToken(token)
		if auth == "" {
			c.IndentedJSON(http.StatusUnauthorized, jsend.NewFail(map[string]string{
				"error": "Unauthorized.",
			}))
			c.Abort()
			return
		}
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, jsend.NewFail(map[string]string{
				"error": "Unauthorized.",
			}))
			c.Abort()
			return
		}
		if !valid {
			c.IndentedJSON(http.StatusUnauthorized, jsend.NewFail(map[string]string{
				"error": "Unauthorized.",
			}))
			c.Abort()
			return
		}

		reqDb := database.GetAccountReqDb{
			CustomerXid: customer_xid,
		}
		resDb, err := pgserver.GetAccount(reqDb)
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, jsend.NewFail(map[string]string{
				"error": "Unauthorized.",
			}))
			c.Abort()
			return
		}

		authUser := service.AuthenticatedUser{
			AccountId:   resDb.AccountId,
			Name:        resDb.Name,
			Email:       resDb.Email,
			CustomerXid: resDb.CustomerXid,
		}

		c.Set("user", authUser)
		c.Next()
	}
}

func TokenEmployeeAuthenticateMiddleware(pgserver database.PostgresClient) gin.HandlerFunc {
	return func(c *gin.Context) {

		auth := c.GetHeader("Authorization")
		token := auth[len("Token "):]
		valid, customer_xid, err := funcs.VerifyToken(token)
		fmt.Print("lol")
		fmt.Print(customer_xid)
		if auth == "" {
			c.IndentedJSON(http.StatusUnauthorized, jsend.NewFail(map[string]string{
				"error": "Unauthorized.",
			}))
			c.Abort()
			return
		}
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, jsend.NewFail(map[string]string{
				"error": "Unauthorized.",
			}))
			c.Abort()
			return
		}
		if !valid {
			c.IndentedJSON(http.StatusUnauthorized, jsend.NewFail(map[string]string{
				"error": "Unauthorized.",
			}))
			c.Abort()
			return
		}

		reqDb := database.GetEmployeeReqDb{
			EmployeeXid: customer_xid,
		}
		resDb, err := pgserver.GetEmployee(reqDb)
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, jsend.NewFail(map[string]string{
				"error": "Unauthorized.",
			}))
			c.Abort()
			return
		}

		authUser := service.AuthenticatedEmployee{
			EmployeeId:  resDb.EmployeeId,
			Name:        resDb.Name,
			Email:       resDb.Email,
			EmployeeXid: resDb.EmployeeXid,
		}

		c.Set("user", authUser)
		c.Next()
	}
}
