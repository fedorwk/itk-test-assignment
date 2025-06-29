package tests

import (
	"context"
	"itk-assignment/infra"
	domain "itk-assignment/wallet"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestConn(t *testing.T) {
	db, err := infra.ConnectDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateWallets(t *testing.T) {
	db, err := infra.ConnectDB()
	if err != nil {
		t.Fatal(err)
	}
	err = setupDatabase(db)
	if err != nil {
		t.Fatal(err)
	}

	service := domain.NewSQLWalletService(db)

	wallets := make([]domain.Wallet, 0, 5)

	for i := 0; i < 5; i++ {
		wallet, err := service.CreateWallet()
		if err != nil {
			t.Error(err)
		}
		wallets = append(wallets, wallet)
	}

	defer func() {
		for _, one := range wallets {
			err = deleteWallet(db, one.Id())
			if err != nil {
				t.Error(err)
			}
		}
	}()

	gotWallets := make([]domain.Wallet, 0, 5)
	for _, inwallet := range wallets {
		outwallet, err := service.Get(inwallet.Id())
		if err != nil {
			t.Error(err)
		}
		gotWallets = append(gotWallets, outwallet)
	}

	for i := range wallets {
		if wallets[i].Id() != gotWallets[i].Id() {
			t.Error("wrong id returned")
		}
	}

	for _, wallet := range wallets {
		err = service.Process(wallet.Id(), domain.NewOperation(domain.OperationDeposit, domain.Amount(10)))
		if err != nil {
			t.Error(t)
		}
	}

	for _, wallet := range wallets {
		err = service.Process(wallet.Id(), domain.NewOperation(domain.OperationWithdraw, domain.Amount(5)))
		if err != nil {
			t.Error(t)
		}
	}
	for _, wallet := range wallets {
		err = service.Process(wallet.Id(), domain.NewOperation(domain.OperationWithdraw, domain.Amount(10)))
		if err == nil {
			t.Error("withdraw more than avaliable should produce error")
		}
		if err != domain.ErrInsufficientFunds {
			t.Error("unexpected error", err)
		}
	}
}

func setupDatabase(db *sqlx.DB) error {
	createTableStmt, err := os.ReadFile("../sql/001_create_wallet.up.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(createTableStmt))
	if err != nil {
		return err
	}
	return nil
}

func deleteWallet(db *sqlx.DB, id domain.WalletId) error {
	res, err := db.ExecContext(context.TODO(), "DELETE from wallet WHERE id = $1", id.UUID())
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrIDNotFound
	}
	return nil
}
