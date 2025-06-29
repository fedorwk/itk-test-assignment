package wallet

import "strconv"

const (
	OperationUndefined operationType = ""
	OperationWithdraw  operationType = "WITHDRAW"
	OperationDeposit   operationType = "DEPOSIT"
)

type operationType string

func (t operationType) String() string { return string(t) }

type amount uint64

func Amount(val int64) amount {
	return amount(val)
}

func (am amount) Int64() int64   { return int64(am) }
func (am amount) String() string { return strconv.FormatInt(am.Int64(), 10) }

type Operation struct {
	t      operationType
	amount amount
}

func (op Operation) Type() operationType {
	return op.t
}
func (op Operation) Amount() amount {
	return op.amount
}

func NewOperation(op operationType, amount amount) Operation {
	return Operation{
		t:      op,
		amount: amount,
	}
}
