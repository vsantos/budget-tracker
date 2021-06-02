// Package classification Budget-tracker API.
//
// the purpose of this application is to provide an application
// that is using plain go code to define an API
//
// This should demonstrate all the possible comment annotations
// that are available to turn go code into a fully compliant swagger 2.0 spec
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http, https
//     Host: budget-tracker:5000
//     BasePath:
//     Version: 0.0.4
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Victor Santos<vsantos.py@gmail.com> https://github.com/vsantos
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	"budget-tracker-api/observability"
	"budget-tracker-api/server"
	"budget-tracker-api/services"
	"crypto/tls"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	port    = ":5000"
	service = "budget-tracker-api"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	o := observability.ProvidersConfig{
		ServiceName: service,
		JaegerURL:   "http://jaeger:14268/api/traces",
		ZipkinURL:   "http://jaeger:9411/api/v2/spans",
	}

	hc := server.HTTPConfig{
		Port: port,
		TLSConfig: &tls.Config{
			// In the absence of `NextProtos`, HTTP/1.1 protocol will be enabled
			NextProtos: []string{"h2"},
		},
		CertFile: "config/tls/server.crt",
		KeyFile:  "config/tls/server.key",
	}

	db, err := services.InitMongoDB()
	if err != nil {
		log.Fatalln(err)
	}
	services.NoSQLClient = db

	sv := server.Server{
		HTTPConfig: hc,
		Observability: observability.Config{
			TracerProviders: o,
		},
		Router: mux.NewRouter(),
		NoSQL: services.Storage{
			NoSQLClient: services.NoSQLClient,
		},
	}

	// In case of 'h2' (HTTP/2) the serverTLS must be set as `true`
	err = sv.Start(false)
	if err != nil {
		log.Fatalln(err)
	}
}
