package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"

	httptreemux "github.com/dimfeld/httptreemux/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
	loginConstructor "littlerollingsushi.com/example/usecase/login/constructor"
	registrationConstructor "littlerollingsushi.com/example/usecase/registration/constructor"
)

type SqlConfig struct {
	Driver   string `envconfig:"DRIVER" default:"mysql"`
	Host     string `envconfig:"HOST" default:"127.0.0.1"`
	Port     int    `envconfig:"PORT" default:"3306"`
	Username string `envconfig:"USERNAME" default:"example"`
	Password string `envconfig:"PASSWORD" default:"example"`
	Database string `envconfig:"DATABASE" default:"example"`
}

func (c *SqlConfig) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
}

func main() {
	_ = godotenv.Load(".env")

	sqlConfig := SqlConfig{}
	envconfig.Process("sql", &sqlConfig)
	db, err := sql.Open(sqlConfig.Driver, sqlConfig.DSN())
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	defer db.Close()

	handler := httptreemux.New()
	handler.POST("/v1/register", registrationConstructor.ConstructRegisterHandler(db).Register)
	handler.POST("/v1/login", loginConstructor.ConstructLoginHandler(db).Login)

	server := &http.Server{
		Addr:           ":7070",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Error from closing listeners, or context timeout: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Error starting or closing listener: %v", err)
	}

	<-idleConnsClosed
}
