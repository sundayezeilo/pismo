package app

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/sundayezeilo/pismo/config"
	"github.com/sundayezeilo/pismo/db"
	"github.com/sundayezeilo/pismo/handlers"
	"github.com/sundayezeilo/pismo/repositories"
	"github.com/sundayezeilo/pismo/services"
)

type App struct {
	Store              *sql.DB
	AccountHandler     *handlers.AccountHandler
	TransactionHandler *handlers.TransactionHandler
	Router             *http.ServeMux
}

func NewApp(cfg *config.Config) *App {
	db, err := db.NewPostgresDB("postgres", cfg.PostgresURL)

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	accountRepo := repositories.NewAccountRepository(db)
	txnRepo := repositories.NewTxnRepository(db)
	opTypeRepo := repositories.NewOpTypeTypeRepository(db)

	accSrv := services.NewAccountService(accountRepo)
	opTypeSrv := services.NewOpTypeService(opTypeRepo)
	txnSrv := services.NewTransactionService(txnRepo, accSrv, opTypeSrv)

	accountHandler := handlers.NewAccountHandler(accSrv)
	txnHandler := handlers.NewTransactionHandler(txnSrv)

	app := &App{
		Store:              db,
		AccountHandler:     accountHandler,
		TransactionHandler: txnHandler,
		Router:             http.NewServeMux(),
	}

	app.setupRoutes()

	return app
}

func (a *App) setupRoutes() {
	a.Router.HandleFunc("/accounts", a.AccountHandler.CreateAccount)
	a.Router.HandleFunc("/accounts/{accountId}", a.AccountHandler.GetAccount)
	a.Router.HandleFunc("/transactions", a.TransactionHandler.CreateTransaction)
}

func (a *App) Cleanup() {
	a.Store.Close()
}
