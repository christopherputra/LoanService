package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"loanservice/database"
	middlewares "loanservice/middlewares"
	service "loanservice/service"
	h "loanservice/service/handler"
	"log"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Postgres struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Dbname   string `json:"dbname"`
	} `json:"postgres"`
	Host string `json:"host"`
	Port string `json:"port"`
}

func LoadConfig(file string) (Config, error) {
	var config Config
	configFile, err := ioutil.ReadFile(file)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func main() {

	//LOAD CONFIG
	config, err := LoadConfig("./config.json")
	if err != nil {
		fmt.Println(err.Error())
	}

	//POSTGRES
	postgresClient := database.NewPostgresClient(
		config.Postgres.Host,
		config.Postgres.Port,
		config.Postgres.User,
		config.Postgres.Password,
		config.Postgres.Dbname,
	)
	loanService := service.NewLoanService(postgresClient)

	handler := h.Handler{
		LoanService: loanService,
	}
	router := gin.New()
	router.POST("/api/v1/account/create", handler.InitAccount)
	router.POST("/api/v1/account/token", handler.TokenAccount)
	router.POST("/api/v1/employee/create", handler.InitEmployee)
	router.POST("/api/v1/employee/token", handler.TokenEmployee)

	rg1 := router.Group("/api/v1/loan", middlewares.TokenAuthenticateMiddleware(postgresClient))
	rg1.POST("/propose", handler.ProposeLoan)
	rg1.POST("/invest", handler.InvestLoan)

	rg2 := router.Group("/api/v1/loan", middlewares.TokenEmployeeAuthenticateMiddleware(postgresClient))
	rg2.POST("/approve", handler.ApproveLoan)
	rg2.POST("/disburse", handler.DisburseLoan)

	if err := router.Run(fmt.Sprintf(":%s", config.Port)); err != nil {
		log.Fatal(context.Background(), err.Error())
	}

}
