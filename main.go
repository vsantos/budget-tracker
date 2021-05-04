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
//     Version: 0.0.2
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
	"budget-tracker-api/routes"
	"net/http"

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
	c := observability.ProvidersConfig{
		ServiceName: "budget-tracker-api",
		JaegerURL:   "http://jaeger:14268/api/traces",
		ZipkinURL:   "http://jaeger:9411/api/v2/spans",
	}

	p, err := c.InitTracerProviders()
	if err != nil {
		log.Errorln(err)
	}

	// Change provider exporter if needed. Ex: `p.Stdout`
	observability.InitGlobalTrace(p.Jaeger)
	observability.InitMetrics()

	router := mux.NewRouter()
	routes.InitRoutes(service, router)

	err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalln(err)
	}

	log.Infoln("Started Application at port", port)
}
