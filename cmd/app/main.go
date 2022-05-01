package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"

	// Logger
	_loggerUcase "github.com/wascript3r/cryptopay/pkg/logger/usecase"

	// Request
	_requestHandler "github.com/wascript3r/anomaly/pkg/request/delivery/http"
	_requestRepo "github.com/wascript3r/anomaly/pkg/request/repository"
	_requestUcase "github.com/wascript3r/anomaly/pkg/request/usecase"
	_requestValidator "github.com/wascript3r/anomaly/pkg/request/validator"

	// Graph
	_graphHandler "github.com/wascript3r/anomaly/pkg/graph/delivery/http"
	_graphRepo "github.com/wascript3r/anomaly/pkg/graph/repository"
	_graphUcase "github.com/wascript3r/anomaly/pkg/graph/usecase"
	_graphValidator "github.com/wascript3r/anomaly/pkg/graph/validator"

	// CORS
	_corsMid "github.com/wascript3r/anomaly/pkg/cors/delivery/http/middleware"
)

const (
	// Database
	DatabaseDriver = "postgres"

	// App
	AppLoggerPrefix = "[APP]"
)

var (
	WorkDir string
	Cfg     *Config

	flagIndexPath = flag.String("i", "public/index.html", "index.html path")
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var err error

	// Get working directory
	WorkDir, err = os.Getwd()
	if err != nil {
		fatalError(err)
	}

	// Parse config file
	cfgPath, err := getConfigPath()
	if err != nil {
		fatalError(err)
	}

	Cfg, err = parseConfig(filepath.Join(WorkDir, cfgPath))
	if err != nil {
		fatalError(err)
	}
}

func fatalError(err any) {
	fmt.Println(err)
	os.Exit(1)
}

func main() {
	flag.Parse()
	if *flagIndexPath == "" {
		fatalError("Missing index.html path (-i public/index.html)")
	}

	// Logging
	logFlags := 0
	if Cfg.Log.ShowTimestamp {
		logFlags = log.Ltime
	}
	logger := _loggerUcase.New(
		AppLoggerPrefix,
		log.New(os.Stdout, "", logFlags),
	)

	// Startup message
	logger.Info("... Starting app ...")

	// Database connection
	dbConn, err := openDatabase(DatabaseDriver, Cfg.Database.Postgres.DSN)
	if err != nil {
		fatalError(err)
	}

	fuzzyUcase, err := getFuzzyUcase()
	if err != nil {
		log.Fatalln(err)
	}

	// Request
	requestRepo := _requestRepo.NewPgRepo(dbConn)
	requestValidator := _requestValidator.New(Cfg.Request.DateTimeFormat)
	requestUcase := _requestUcase.New(
		requestRepo,
		Cfg.Database.Postgres.QueryTimeout.Duration,

		fuzzyUcase,
		requestValidator,
	)

	// Graph
	graphRepo := _graphRepo.NewPgRepo(dbConn)
	graphValidator := _graphValidator.New()
	graphUcase := _graphUcase.New(
		graphRepo,
		Cfg.Database.Postgres.QueryTimeout.Duration,

		graphValidator,
	)

	// Graceful shutdown
	stopSig := make(chan os.Signal, 1)
	signal.Notify(stopSig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// HTTP server
	httpRouter := httprouter.New()
	httpRouter.MethodNotAllowed = MethodNotAllowedHnd
	httpRouter.NotFound = NotFoundHnd

	if Cfg.HTTP.EnablePprof {
		// pprof
		httpRouter.Handler(http.MethodGet, "/debug/pprof/*item", http.DefaultServeMux)
	}

	_requestHandler.NewHTTPHandler(httpRouter, requestUcase)
	_graphHandler.NewHTTPHandler(httpRouter, graphUcase)

	// index.html file
	httpRouter.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.ServeFile(w, r, *flagIndexPath)
	})

	httpServer := &http.Server{
		Addr:    ":" + Cfg.HTTP.Port,
		Handler: _corsMid.NewHTTPMiddleware().EnableCors(httpRouter),
	}

	if Cfg.HTTP.TLS != nil {
		httpServer.Addr = ":443"
	}

	logger.Info("Listening on port %s", httpServer.Addr)

	// Graceful shutdown
	gracefulShutdown := func() {
		if err := httpServer.Shutdown(context.Background()); err != nil {
			logger.Error("Cannot shutdown HTTP server: %s", err)
		}

		logger.Info("... Exited ...")
		os.Exit(0)
	}

	go func() {
		<-stopSig
		logger.Info("... Received stop signal ...")
		gracefulShutdown()
	}()

	if Cfg.HTTP.TLS != nil {
		err = httpServer.ListenAndServeTLS(Cfg.HTTP.TLS.CertFile, Cfg.HTTP.TLS.KeyFile)
	} else {
		err = httpServer.ListenAndServe()
	}

	if err != nil {
		if err != http.ErrServerClosed {
			fmt.Println(err)
			gracefulShutdown()
		}
	}

	var c chan struct{}
	<-c
}
