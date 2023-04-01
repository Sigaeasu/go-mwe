package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"github.com/Sigaeasu/go-mwe/config"
	"github.com/Sigaeasu/go-mwe/config/database"
	"github.com/Sigaeasu/go-mwe/service"
	"github.com/Sigaeasu/go-mwe/repository"
	"github.com/Sigaeasu/go-mwe/handler"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func RunServer() {
	postgresConfig := config.Config.PostgresCfg
	configDatabase := database.ParametersConnection{
		Username: postgresConfig.Username,
		Password: postgresConfig.Password,
		Host: postgresConfig.Host,
		Port: postgresConfig.Port,
		Database: postgresConfig.Database,
		MaxConnection: postgresConfig.MaxConn,
		MinIdleConnection: postgresConfig.MinIdleConn,
		MaxRetries: postgresConfig.MaxRetries,
	}

	db := database.DatabaseConnection(configDatabase)

	err := db.Ping(context.Background())
	if err != nil {
		logrus.Fatalf("Ping DB error: %v", err)
	}

	m := mux.NewRouter()
	miniWalletDatabase := repository.MiniWalletRepository(db)
	handlerAPI := handler.MiniWalletHandler(miniWalletDatabase)

	api := m.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/init", handlerAPI.AuthMiniWallet).Methods(http.MethodPost)
	api.HandleFunc("/wallet", handlerAPI.ViewMiniWalletBalance).Methods(http.MethodGet)
	api.HandleFunc("/wallet/transactions", handlerAPI.ViewTransactions).Methods(http.MethodGet)
	api.HandleFunc("/wallet", handlerAPI.EnableMiniWallet).Methods(http.MethodPost)
	api.HandleFunc("/wallet", handlerAPI.DisableMiniWallet).Methods(http.MethodPatch)
	api.HandleFunc("/wallet/deposits", handlerAPI.DepositToMiniWallet).Methods(http.MethodPost)
	api.HandleFunc("/wallet/withdrawals", handlerAPI.WithdrawFromMiniWallet).Methods(http.MethodPost)
	m.Use(mux.CORSMethodMiddleware(m))
	m.Use(service.AuthMiddlewareService())

	srvr := &http.Server{
		Handler:      m,
		Addr:         ":5000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logrus.Info("Starting on port 5000")

	go func() {
		if err := srvr.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	srvr.Shutdown(ctx)
	os.Exit(0)
}