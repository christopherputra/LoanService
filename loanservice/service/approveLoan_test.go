package service

import (
	"database/sql"
	"fmt"
	"loanservice/constants"
	"loanservice/database"
	mock_database "loanservice/database/mock_database"
	"testing"

	gomock "github.com/golang/mock/gomock"
	assert "github.com/stretchr/testify/assert"
)

func TestApproveLoan(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	postgresClientMock := mock_database.NewMockPostgresClient(mockCtrl)
	authEmployee := AuthenticatedEmployee{
		EmployeeId:  1,
		EmployeeXid: "employee_xid_1",
		Email:       "email_1",
		Name:        "name_1",
	}

	t.Run("When postgres client is running properly", func(t *testing.T) {
		resDb := database.GetLoanResDb{
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
		postgresClientMock.EXPECT().GetLoan(gomock.Any()).Do(func(reqDb database.GetLoanReqDb) {
			assert.Equal(t, reqDb.LoanXid, "loan_xid_1", "loan_Xid is not the same")
		}).Return(resDb, nil)

		postgresClientMock.EXPECT().StartTransaction().Return(&sql.Tx{}, nil)

		postgresClientMock.EXPECT().InsertLoanApproval(gomock.Any(), gomock.Any()).Return(nil)

		resUpdateStatusDb := database.UpdateLoanStatusResDb{
			LoanDb: database.LoanDb{
				LoanId:          2,
				LoanXid:         "loan_xid_1",
				BorrowerId:      1,
				PrincipalAmount: int64(12345),
				Rate:            int64(12),
				Roi:             0,
				Status:          constants.LoanStatusApproved,
			},
		}
		postgresClientMock.EXPECT().UpdateLoanStatus(gomock.Any(), gomock.Any()).Return(resUpdateStatusDb, nil)

		postgresClientMock.EXPECT().CommitTransaction(gomock.Any()).Return(nil)

		req := ApproveLoanReq{
			LoanXid: "loan_xid_1",
		}
		res := &ApproveLoanRes{}

		loanService := Service{PgServer: postgresClientMock}
		_, err := loanService.ApproveLoan(authEmployee, req, res)

		assert.Equal(t, err, nil, "Error should be nil")
		assert.Equal(t, res.LoanXid, "loan_xid_1", "loan_Xid is not the same")
		assert.Equal(t, res.Amount, int64(12345), "Amount is not the same")
		assert.Equal(t, res.Rate, int64(12), "Rate is not the same")
		assert.Equal(t, res.Status, constants.LoanStatusApproved, "Status is not the same")

	})

	t.Run("When loan status is not correct", func(t *testing.T) {
		resDb := database.GetLoanResDb{
			LoanDb: database.LoanDb{
				LoanId:          2,
				LoanXid:         "loan_xid_1",
				BorrowerId:      1,
				PrincipalAmount: 12345,
				Rate:            12,
				Roi:             0,
				Status:          constants.LoanStatusDisbursed, //wrong status
			},
		}
		postgresClientMock.EXPECT().GetLoan(gomock.Any()).Do(func(reqDb database.GetLoanReqDb) {
			assert.Equal(t, reqDb.LoanXid, "loan_xid_1", "loan_Xid is not the same")
		}).Return(resDb, nil)

		req := ApproveLoanReq{
			LoanXid: "loan_xid_1",
		}
		res := &ApproveLoanRes{}

		loanService := Service{PgServer: postgresClientMock}
		_, err := loanService.ApproveLoan(authEmployee, req, res)

		assert.Equal(t, err, fmt.Errorf("invalid loan id"), "Error should be nil")
	})
}
