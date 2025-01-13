package constants

const LoanStatusProposed = "PROPOSED"
const LoanStatusApproved = "APPROVED"
const LoanStatusInvested = "INVESTED"
const LoanStatusDisbursed = "DISBURSED"

var LoanStatusPriority = map[string]int{
	LoanStatusProposed:  1,
	LoanStatusApproved:  2,
	LoanStatusInvested:  3,
	LoanStatusDisbursed: 4,
}
