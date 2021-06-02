package server

import (
	"budget-tracker-api/observability"
	"budget-tracker-api/routes"
	"budget-tracker-api/services"
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

// Server will initiate a server
type Server struct {
	HTTPConfig    HTTPConfig
	Observability observability.Config
	NoSQL         services.Storage
	Router        *mux.Router
}

// HTTPConfig will define server configuration for both HTTP/1 and HTTP/2 protocols
type HTTPConfig struct {
	Port      string
	TLSConfig *tls.Config
	CertFile  string
	KeyFile   string
}

// waitGracefulShutdown is blocking code to wait a SIGNAL to graceful shutdown the server
func waitGracefulShutdown(srv *http.Server, timeoutSeconds time.Duration) (err error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	os := <-quit

	log.Printf("Waiting server to shutdown due to signal '%+v'", os)
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSeconds*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func (sv *Server) setupRoutes() {
	routes.InitRoutes(sv.Observability.TracerProviders.ServiceName, sv.Router)
}

func (sv *Server) setupTracing() (err error) {

	p, err := observability.InitTracerProviders(sv.Observability.TracerProviders)
	if err != nil {
		log.Errorln(err)
	}

	observability.InitGlobalTrace(p.Jaeger)
	observability.InitMetrics()

	return nil
}

// Start will serve a HTTP/1 (h2) server with optional TLS enforcement
func (sv *Server) Start(serveTLS bool) (err error) {
	sv.setupRoutes()
	if err := sv.setupTracing(); err != nil {
		return err
	}

	srv := &http.Server{
		Addr:      sv.HTTPConfig.Port,
		Handler:   sv.Router,
		TLSConfig: sv.HTTPConfig.TLSConfig,
	}

	if serveTLS {
		log.Infoln(fmt.Sprintf("Started %s Application at port %s with TLS enabled", srv.TLSConfig.NextProtos[0], sv.HTTPConfig.Port))
		go func() {
			if err := srv.ListenAndServeTLS(sv.HTTPConfig.CertFile, sv.HTTPConfig.KeyFile); err != nil && err != http.ErrServerClosed {
				log.Panicln("Server error: ", err)
			}
		}()
	}

	if !serveTLS {
		log.Infoln(fmt.Sprintf("Started %s Application at port %s with TLS disabled", srv.TLSConfig.NextProtos[0], sv.HTTPConfig.Port))
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Panicln("Server error: ", err)
			}
		}()
	}

	if err := waitGracefulShutdown(srv, 15); err != nil {
		return err
	}

	log.Infoln("Server finished")
	return nil
}
