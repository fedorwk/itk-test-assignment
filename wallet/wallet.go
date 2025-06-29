package wallet

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type WalletId uuid.UUID

func (id WalletId) UUID() uuid.UUID { return uuid.UUID(id) }

type Wallet interface {
	Balance() amount
	Id() WalletId
}

type SimpleWallet struct {
	id      WalletId
	balance amount
}

func (w SimpleWallet) Balance() amount {
	return w.balance
}

func (w SimpleWallet) Id() WalletId {
	return w.id
}

func (w SimpleWallet) UUID() uuid.UUID {
	return uuid.UUID(w.id)
}

type WalletService interface {
	CreateWallet() (Wallet, error)
	Get(WalletId) (Wallet, error)
	Process(WalletId, Operation) error
}

type SQLWalletService struct {
	db *sqlx.DB
}

func NewSQLWalletService(db *sqlx.DB) *SQLWalletService {
	return &SQLWalletService{
		db,
	}
}

func (m *SQLWalletService) CreateWallet() (Wallet, error) {
	id := uuid.New()
	wallet := SimpleWallet{
		id: WalletId(id),
	}

	err := m.registerWallet(wallet)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (m *SQLWalletService) Get(id WalletId) (Wallet, error) {
	row := m.db.QueryRowxContext(context.TODO(), "SELECT id, balance FROM wallet WHERE id = $1", id.UUID())
	var scanId uuid.UUID
	var scanBalance int64
	err := row.Scan(&scanId, &scanBalance)
	if err != nil {
		return nil, err
	}
	return SimpleWallet{
		id:      WalletId(id),
		balance: amount(scanBalance),
	}, nil
}

var ErrUnknownOperation = errors.New("unknown operation type")
var ErrInsufficientFunds = errors.New("insufficient funds to withdraw")

func (m *SQLWalletService) Process(id WalletId, op Operation) error {
	switch op.Type() {
	case OperationWithdraw:
		err := m.tryWithdraw(id, op.Amount())
		if err != nil {
			return err
		}
	case OperationDeposit:
		err := m.deposit(id, op.Amount())
		if err != nil {
			return err
		}
	default:
		return ErrUnknownOperation
	}
	return nil
}

func (m *SQLWalletService) tryWithdraw(id WalletId, am amount) error {
	tx, err := m.db.Beginx()
	if err != nil {
		return err
	}
	row := tx.QueryRowxContext(context.TODO(), "SELECT balance FROM wallet WHERE id = $1", id.UUID())
	var rawBalance amount
	err = row.Scan(&rawBalance)
	if err != nil {
		return err
	}
	if am > rawBalance {
		return ErrInsufficientFunds
	}
	_, err = tx.ExecContext(context.TODO(), "UPDATE wallet SET balance = $1 WHERE id = $2", rawBalance-am, id.UUID())
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

var ErrIDNotFound = errors.New("no wallet with such id")

func (m *SQLWalletService) deposit(id WalletId, am amount) error {
	res, err := m.db.ExecContext(context.TODO(), "UPDATE wallet SET balance = balance + $1 WHERE id = $2", am, id.UUID())
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrIDNotFound
	}
	return nil
}

func (m *SQLWalletService) registerWallet(w Wallet) error {
	_, err := m.db.ExecContext(
		context.TODO(),
		"INSERT INTO wallet (id, balance) VALUES ($1, $2)",
		uuid.UUID(w.Id()), w.Balance(),
	)
	if err != nil {
		return err
	}
	return nil
}
