package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gocraft/web"
	"github.com/jasonlvhit/gocron"

	"ucoi4tzlito52agmquc6oopn3zpmr6djz5vvfabtgrhyc6hufpzjtnad.onion/Tochka/tochka-free-market/modules/settings"
	"ucoi4tzlito52agmquc6oopn3zpmr6djz5vvfabtgrhyc6hufpzjtnad.onion/Tochka/tochka-free-market/modules/webapp"
)

var APP_SETTINGS = settings.GetSettings()

func runCrons() {
	if !APP_SETTINGS.Debug {
		webapp.StartTransactionsCron()
		webapp.StartWalletsCron()
		webapp.StartStatsCron()
		webapp.StartSERPCron()
		webapp.StartMessageboardCron()
	}

	webapp.StartCurrencyCron()

	<-gocron.Start()
}

func runWebserver() {
	// Root router
	rootRouter := web.New(webapp.Context{})

	rootRouter.OptionsHandler((*webapp.Context).OptionsHandler)

	rootRouter.Middleware(web.StaticMiddleware("public"))
	// webapp router
	webapp.ConfigureRouter(rootRouter.Subrouter(webapp.Context{}, "/"))
	// Start HTTP server

	startHttpServer := func() {
		address := fmt.Sprintf("%s:%s", APP_SETTINGS.Host, APP_SETTINGS.Port)
		println(fmt.Sprintf("Running a HTTP server or, %s", address))

		srv := &http.Server{
			ReadTimeout:       60 * time.Second,
			WriteTimeout:      60 * time.Second,
			IdleTimeout:       120 * time.Second,
			ReadHeaderTimeout: 60 * time.Second,
			Handler:           rootRouter,
			Addr:              address,
		}

		if err := srv.ListenAndServe(); err != nil {
			fmt.Printf("Server failed: %s\n", err)
		}

		// http.ListenAndServe(address, rootRouter)
	}
	// Start HTTPs server
	startHttpServer()
}

func runServer() {
	go runCrons()
	runWebserver()
}
