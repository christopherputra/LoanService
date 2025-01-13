package handler

import (
	"fmt"
	"loanservice/service"
	"net/http"
	"reflect"
	"strings"

	"clevergo.tech/jsend"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	LoanService service.Service
}

func parseError(err error, blueprint interface{}) map[string][]string {
	out := make(map[string][]string)
	switch typedError := any(err).(type) {
	case validator.ValidationErrors:
		for _, e := range typedError {
			fieldName := e.Field()
			field, _ := reflect.TypeOf(blueprint).Elem().FieldByName(fieldName)
			fieldJSONName, _ := field.Tag.Lookup("json")
			out[fieldJSONName] = parseFieldError(e)
		}
	}

	return out
}

func parseFieldError(e validator.FieldError) []string {
	tag := strings.Split(e.Tag(), "|")[0]
	switch tag {
	case "required":
		return []string{"Missing data for required field."}
	default:
		return []string{fmt.Sprintf("Field '%s' failed on tag '%s'", e.Field(), e.Tag())}
	}
}

func (h *Handler) InitAccount(c *gin.Context) {
	var req service.InitAccountReq
	res := &service.InitAccountRes{}
	if err := c.BindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, jsend.NewFail(parseError(err, &req)))
		return
	}

	status_code, err := h.LoanService.InitAccount(req, res)

	if err != nil {
		fmt.Println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, jsend.NewFail(map[string]string{
			"error": err.Error(),
		}))
		c.Abort()
		return
	}
	c.IndentedJSON(status_code, jsend.New(res))
	return
}

func (h *Handler) TokenAccount(c *gin.Context) {
	var req service.TokenAccountReq
	res := &service.TokenAccountRes{}
	if err := c.BindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, jsend.NewFail(parseError(err, &req)))
		return
	}

	status_code, err := h.LoanService.TokenAccount(req, res)

	if err != nil {
		fmt.Println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, jsend.NewFail(map[string]string{
			"error": err.Error(),
		}))
		c.Abort()
		return
	}
	c.IndentedJSON(status_code, jsend.New(res))
	return
}

func (h *Handler) InitEmployee(c *gin.Context) {
	var req service.InitEmployeeReq
	res := &service.InitEmployeeRes{}
	if err := c.BindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, jsend.NewFail(parseError(err, &req)))
		return
	}

	status_code, err := h.LoanService.InitEmployee(req, res)

	if err != nil {
		fmt.Println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, jsend.NewFail(map[string]string{
			"error": err.Error(),
		}))
		c.Abort()
		return
	}
	c.IndentedJSON(status_code, jsend.New(res))
	return
}

func (h *Handler) TokenEmployee(c *gin.Context) {
	var req service.TokenEmployeeReq
	res := &service.TokenEmployeeRes{}
	if err := c.BindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, jsend.NewFail(parseError(err, &req)))
		return
	}

	status_code, err := h.LoanService.TokenEmployee(req, res)

	if err != nil {
		c.IndentedJSON(status_code, jsend.NewFail(map[string]string{
			"error": err.Error(),
		}))
		c.Abort()
		return
	}
	c.IndentedJSON(status_code, jsend.New(res))
	return
}

func (h *Handler) ProposeLoan(c *gin.Context) {
	var req service.ProposeLoanReq
	res := &service.ProposeLoanRes{}
	if err := c.BindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, jsend.NewFail(parseError(err, &req)))
		return
	}

	user, _ := c.Get("user")
	u, ok := user.(service.AuthenticatedUser)
	if !ok {
		c.IndentedJSON(http.StatusInternalServerError, jsend.NewFail(map[string]string{
			"error": "internal server error",
		}))
		return
	}

	err := h.LoanService.ProposeLoan(u, req, res)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, jsend.NewFail(map[string]string{
			"error": err.Error(),
		}))
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusCreated, jsend.New(res))
	return
}

func (h *Handler) ApproveLoan(c *gin.Context) {
	var req service.ApproveLoanReq
	res := &service.ApproveLoanRes{}
	if err := c.BindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, jsend.NewFail(parseError(err, &req)))
		return
	}
	user, _ := c.Get("user")
	employee, ok := user.(service.AuthenticatedEmployee)
	if !ok {
		c.IndentedJSON(http.StatusInternalServerError, jsend.NewFail(map[string]string{
			"error": "internal server error",
		}))
		return
	}

	status_code, err := h.LoanService.ApproveLoan(employee, req, res)
	if err != nil {
		fmt.Println(err.Error())
		c.IndentedJSON(status_code, jsend.NewFail(map[string]string{
			"error": err.Error(),
		}))
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusCreated, jsend.New(res))
	return
}

func (h *Handler) InvestLoan(c *gin.Context) {
	var req service.InvestLoanReq
	res := &service.InvestLoanRes{}
	if err := c.BindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, jsend.NewFail(parseError(err, &req)))
		return
	}
	user, _ := c.Get("user")
	u, ok := user.(service.AuthenticatedUser)
	if !ok {
		c.IndentedJSON(http.StatusInternalServerError, jsend.NewFail(map[string]string{
			"error": "internal server error",
		}))
		return
	}

	status_code, err := h.LoanService.InvestLoan(u, req, res)
	if err != nil {
		fmt.Println(err.Error())
		c.IndentedJSON(status_code, jsend.NewFail(map[string]string{
			"error": err.Error(),
		}))
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusCreated, jsend.New(res))
	return
}

func (h *Handler) DisburseLoan(c *gin.Context) {
	var req service.DisburseLoanReq
	res := &service.DisburseLoanRes{}
	if err := c.BindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, jsend.NewFail(parseError(err, &req)))
		return
	}
	user, _ := c.Get("user")
	employee, ok := user.(service.AuthenticatedEmployee)
	if !ok {
		c.IndentedJSON(http.StatusInternalServerError, jsend.NewFail(map[string]string{
			"error": "internal server error",
		}))
		return
	}

	status_code, err := h.LoanService.DisburseLoan(employee, req, res)
	if err != nil {
		fmt.Println(err.Error())
		c.IndentedJSON(status_code, jsend.NewFail(map[string]string{
			"error": err.Error(),
		}))
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusCreated, jsend.New(res))
	return
}
