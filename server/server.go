package server

import (
	"context"
	"itk-assignment/infra"
	"itk-assignment/wallet"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func Start() error {
	server := setupServer()

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-osSignal

	return shutdownServer(server)
}

func setupServer() *http.Server {
	db, err := infra.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	if gin.Mode() == gin.TestMode || gin.Mode() == gin.DebugMode {
		err = setupDatabase(db)
		if err != nil {
			log.Fatal(err)
		}
	}

	service := wallet.NewSQLWalletService(db)

	router := gin.Default()
	router.POST("api/v1/wallet", newOperationHandler(service))
	router.GET("/api/v1/wallets/:walletId", newBalanceHandler(service))
	router.POST("api/v1/wallet/create", newCreateHandler(service))

	server := http.Server{
		Addr:    ":" + Config().Port,
		Handler: router,
	}
	return &server
}

func shutdownServer(server *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return server.Shutdown(ctx)
}

func setupDatabase(db *sqlx.DB) error {
	createTableStmt, err := os.ReadFile("sql/001_create_wallet.up.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(createTableStmt))
	if err != nil {
		return err
	}

	addUsrStmt, err := os.ReadFile("sql/002_test_wallet.up.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(addUsrStmt))
	if err != nil {
		return err
	}
	return nil
}
