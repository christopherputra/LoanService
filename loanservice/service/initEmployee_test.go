package service

import (
	"errors"
	"loanservice/database"
	mock_database "loanservice/database/mock_database"
	"testing"

	gomock "github.com/golang/mock/gomock"
	assert "github.com/stretchr/testify/assert"
)

func TestInitEmployee(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	postgresClientMock := mock_database.NewMockPostgresClient(mockCtrl)
	errSample := errors.New("Error")

	t.Run("When postgres client is running properly", func(t *testing.T) {
		postgresClientMock.EXPECT().CreateOrInsertEmployee(gomock.Any()).Do(func(reqDb database.CreateEmployeeReqDb) {
			assert.Equal(t, reqDb.Email, "email-test-1", "Email is not the same")
			assert.Equal(t, reqDb.Password, "password-test-1", "Password is not the same")
		}).Return(nil)

		req := InitEmployeeReq{
			Name:     "name-test-1",
			Email:    "email-test-1",
			Password: "password-test-1",
		}
		res := &InitEmployeeRes{}

		loanService := Service{PgServer: postgresClientMock}
		_, err := loanService.InitEmployee(req, res)

		assert.Equal(t, err, nil, "Error should be nil")
		assert.NotNil(t, res.Token, "Token should not be nil")

	})

	t.Run("When postgres client is not running properly", func(t *testing.T) {
		postgresClientMock.EXPECT().CreateOrInsertEmployee(gomock.Any()).Do(func(reqDb database.CreateEmployeeReqDb) {
			assert.Equal(t, reqDb.Email, "email-test-1", "Email is not the same")
			assert.Equal(t, reqDb.Password, "password-test-1", "Password is not the same")
		}).Return(errSample)

		req := InitEmployeeReq{
			Name:     "name-test-1",
			Email:    "email-test-1",
			Password: "password-test-1",
		}
		res := &InitEmployeeRes{}

		loanService := Service{PgServer: postgresClientMock}
		_, err := loanService.InitEmployee(req, res)

		assert.NotNil(t, err, "Error should not be nil")

	})

}
