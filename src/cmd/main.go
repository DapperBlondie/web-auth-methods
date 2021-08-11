package main

import (
	"github.com/DapperBlondie/web-auth-methods/src/handlers"
	"github.com/DapperBlondie/web-auth-methods/src/repo"
	"github.com/DapperBlondie/web-auth-methods/src/routes"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	scsManager := scs.New()
	scsManager.Lifetime = 24 * time.Hour
	scsManager.Cookie.Persist = true
	scsManager.Cookie.SameSite = http.SameSiteLaxMode
	scsManager.Cookie.Secure = false

	dbRepo, err := repo.NewDB("appdb.db")
	if err != nil {
		return
	}
	err = dbRepo.CreateUserDataModelMethod()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	handlers.NewConfiguration(scsManager, dbRepo)

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt)
	srv := &http.Server{
		Addr:              "localhost:8080",
		Handler:           routes.AppRoutes(),
		TLSConfig:         nil,
		ReadTimeout:       time.Second * 10,
		ReadHeaderTimeout: time.Second * 5,
		WriteTimeout:      time.Second * 15,
		IdleTimeout:       time.Second * 7,
	}

	go func() {
		log.Println("Http Server is listening on localhost:8080 ...")
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err.Error())
			return
		}
	}()

	<-sigC
}
