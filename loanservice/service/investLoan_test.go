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

func TestInvestLoan(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	postgresClientMock := mock_database.NewMockPostgresClient(mockCtrl)
	authUser := AuthenticatedUser{
		AccountId:   1,
		CustomerXid: "customer_xid_1",
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
				Status:          constants.LoanStatusApproved,
			},
		}
		postgresClientMock.EXPECT().GetLoan(gomock.Any()).Do(func(reqDb database.GetLoanReqDb) {
			assert.Equal(t, reqDb.LoanXid, "loan_xid_1", "loan_xid is not the same")
		}).Return(resDb, nil).AnyTimes()

		t.Run("When loan amount is not bigger than needed", func(t *testing.T) {

			resInvestmentDb := database.GetLoanInvestmentsResDb{}
			postgresClientMock.EXPECT().GetLoanInvestments(gomock.Any()).Do(func(reqDb database.GetLoanInvestmentsReqDb) {
				assert.Equal(t, reqDb.LoanId, int64(2), "loan_id is not the same")
			}).Return(resInvestmentDb, nil)

			postgresClientMock.EXPECT().StartTransaction().Return(&sql.Tx{}, nil)

			postgresClientMock.EXPECT().InsertLoanInvestment(gomock.Any(), gomock.Any()).Return(true, nil)

			postgresClientMock.EXPECT().CommitTransaction(gomock.Any()).Return(nil)

			t.Run("When loan amount is equal than needed invested", func(t *testing.T) {
				resUpdateStatusDb := database.UpdateLoanStatusResDb{
					LoanDb: database.LoanDb{
						LoanId:          2,
						LoanXid:         "loan_xid_1",
						BorrowerId:      1,
						PrincipalAmount: int64(100000),
						Rate:            int64(100),
						Roi:             0,
						Status:          constants.LoanStatusInvested,
					},
				}
				postgresClientMock.EXPECT().UpdateLoanStatus(gomock.Any(), gomock.Any()).Return(resUpdateStatusDb, nil).Times(1)

				req := InvestLoanReq{
					LoanXid: "loan_xid_1",
					Amount:  int64(100000),
				}
				res := &InvestLoanRes{}

				loanService := Service{PgServer: postgresClientMock}
				_, err := loanService.InvestLoan(authUser, req, res)

				assert.Equal(t, err, nil, "Error should be nil")
				assert.Equal(t, res.LoanXid, "loan_xid_1", "Loan xid is not the same")
				assert.Equal(t, res.Amount, int64(100000), "Amount is not the same")
				assert.Equal(t, res.Status, constants.LoanStatusInvested, "Status is not the same")
			})
		})

		t.Run("When loan amount is bigger than needed invested", func(t *testing.T) {
			resInvestmentDb := database.GetLoanInvestmentsResDb{
				LoanInvestments: []database.LoanInvestmentDb{
					database.LoanInvestmentDb{
						Amount: 20000,
					},
				},
			}
			postgresClientMock.EXPECT().GetLoanInvestments(gomock.Any()).Do(func(reqDb database.GetLoanInvestmentsReqDb) {
				assert.Equal(t, reqDb.LoanId, int64(2), "loan_id is not the same")
			}).Return(resInvestmentDb, nil)

			req := InvestLoanReq{
				LoanXid: "loan_xid_1",
				Amount:  int64(100000),
			}
			res := &InvestLoanRes{}

			loanService := Service{PgServer: postgresClientMock}
			_, err := loanService.InvestLoan(authUser, req, res)

			assert.Equal(t, err, fmt.Errorf("amount of investment is above the remaining amount"), "Error is not the same")
		})
	})
}
