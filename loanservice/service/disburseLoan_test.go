package service

import (
	"database/sql"
	"loanservice/constants"
	"loanservice/database"
	mock_database "loanservice/database/mock_database"
	"testing"

	gomock "github.com/golang/mock/gomock"
	assert "github.com/stretchr/testify/assert"
)

func TestDisburseLoan(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	postgresClientMock := mock_database.NewMockPostgresClient(mockCtrl)
	authEmployee := AuthenticatedEmployee{
		EmployeeId:  1,
		EmployeeXid: "employee_xid_1",
		Email:       "email_1",
		Name:        "name_1",
	}
	// errSample := fmt.Errorf("error")

	t.Run("When postgres client is running properly", func(t *testing.T) {
		resDb := database.GetLoanResDb{
			LoanDb: database.LoanDb{
				LoanId:          2,
				LoanXid:         "loan_xid_1",
				BorrowerId:      1,
				PrincipalAmount: 100000,
				Rate:            12,
				Roi:             0,
				Status:          constants.LoanStatusInvested,
			},
		}
		postgresClientMock.EXPECT().GetLoan(gomock.Any()).Do(func(reqDb database.GetLoanReqDb) {
			assert.Equal(t, reqDb.LoanXid, "loan_xid_1", "loan_xid is not the same")
		}).Return(resDb, nil).AnyTimes()

		postgresClientMock.EXPECT().StartTransaction().Return(&sql.Tx{}, nil)
		postgresClientMock.EXPECT().CommitTransaction(gomock.Any()).Return(nil)

		postgresClientMock.EXPECT().InsertLoanDisbursement(gomock.Any(), gomock.Any()).Return(true, nil)

		resUpdateStatusDb := database.UpdateLoanStatusResDb{
			LoanDb: database.LoanDb{
				LoanId:          2,
				LoanXid:         "loan_xid_1",
				BorrowerId:      1,
				PrincipalAmount: int64(100000),
				Rate:            int64(100),
				Roi:             0,
				Status:          constants.LoanStatusDisbursed,
			},
		}
		postgresClientMock.EXPECT().UpdateLoanStatus(gomock.Any(), gomock.Any()).Return(resUpdateStatusDb, nil).Times(1)

		req := DisburseLoanReq{
			LoanXid: "loan_xid_1",
		}
		res := &DisburseLoanRes{}

		loanService := Service{PgServer: postgresClientMock}
		_, err := loanService.DisburseLoan(authEmployee, req, res)

		assert.Equal(t, err, nil, "Error should be nil")
		assert.Equal(t, res.LoanXid, "loan_xid_1", "Loan xid is not the same")
		assert.Equal(t, res.AmountDisbursed, int64(100000), "Amount is not the same")

	})
}
