package service

import (
	"fmt"
	"loanservice/constants"
	"loanservice/database"
	mock_database "loanservice/database/mock_database"
	"testing"

	gomock "github.com/golang/mock/gomock"
	assert "github.com/stretchr/testify/assert"
)

func TestProposeLoan(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	postgresClientMock := mock_database.NewMockPostgresClient(mockCtrl)
	authUser := AuthenticatedUser{
		AccountId:   1,
		CustomerXid: "customer_xid_1",
		Email:       "email_1",
		Name:        "name_1",
	}
	errSample := fmt.Errorf("error")

	t.Run("When postgres client is running properly", func(t *testing.T) {
		resDb := database.InsertLoanResDb{
			LoanDb: database.LoanDb{
				LoanId:          2,
				LoanXid:         "loan_xid_1",
				BorrowerId:      1,
				PrincipalAmount: 12345,
				Rate:            12,
				Roi:             0,
				Status:          constants.LoanStatusProposed,
			},
		}
		postgresClientMock.EXPECT().InsertLoan(gomock.Any()).Do(func(reqDb database.InsertLoanReqDb) {
			assert.Equal(t, reqDb.BorrowerId, int64(1), "Borrower is not the same")
		}).Return(resDb, nil)

		req := ProposeLoanReq{
			Amount: 12345,
		}
		res := &ProposeLoanRes{}

		loanService := Service{PgServer: postgresClientMock}
		err := loanService.ProposeLoan(authUser, req, res)

		assert.Equal(t, err, nil, "Error should be nil")
		assert.Equal(t, res.Amount, int64(12345), "Amount is not the same")

	})

	t.Run("When loan status is not correct", func(t *testing.T) {
		resDb := database.InsertLoanResDb{
			LoanDb: database.LoanDb{
				LoanId:          2,
				LoanXid:         "loan_xid_1",
				BorrowerId:      1,
				PrincipalAmount: 12345,
				Rate:            12,
				Roi:             0,
				Status:          constants.LoanStatusProposed,
			},
		}
		postgresClientMock.EXPECT().InsertLoan(gomock.Any()).Do(func(reqDb database.InsertLoanReqDb) {
			assert.Equal(t, reqDb.BorrowerId, int64(1), "Borrower is not the same")
		}).Return(resDb, errSample)

		req := ProposeLoanReq{
			Amount: 12345,
		}
		res := &ProposeLoanRes{}

		loanService := Service{PgServer: postgresClientMock}
		err := loanService.ProposeLoan(authUser, req, res)

		assert.NotNil(t, err, "Error should not be nil")
	})
}
