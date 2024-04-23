package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/smiletrl/gateway/pkg/core"
	errors "github.com/smiletrl/gateway/pkg/error"
	payment "github.com/smiletrl/gateway/service.payment/internal/payment"
)

func main() {
	// build core providers to inject them to follow services
	p := core.BuildProvider()

	// echo instance
	e := echo.New()

	// Middlewares
	// log http request
	e.Use(p.Access.Middleware())
	// log request errors
	e.Use(errors.Middleware(p.Logger))

	g := e.Group("")

	payRepo := payment.NewRepository()
	paySvc := payment.NewService(payRepo)
	payment.RegisterHandlers(g, paySvc)

	// Start rest server
	go func() {
		// @todo move port ":1323" to config to environment variables.
		err := e.Start(":1323")
		if err != nil {
			log.Print("echo server stop", "echo", err.Error())
		}
	}()

	// gracefully shutdown application
	shutdown(e)
}

func shutdown(e *echo.Echo) {
	// Handle SIGTERM
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ch:
		e.Shutdown(context.Background())
	}
}
