package asynchronousreport

import (
	"FlexiProxyHub/internal/database"
	utils "FlexiProxyHub/internal/utils/generic"
	"FlexiProxyHub/pkg/asynchronousReport/download"
	"FlexiProxyHub/pkg/asynchronousReport/interceptor"
	"FlexiProxyHub/pkg/asynchronousReport/jswebsocket"
	"FlexiProxyHub/pkg/asynchronousReport/websocket"
	"net/http"

	"go.uber.org/zap"
)

var db *database.Database
var log *zap.Logger
var config *utils.Configuration
var route map[string]http.Handler

type HostSwitch = interceptor.HostSwitch

func Init(database *database.Database, logger *zap.Logger, configuration *utils.Configuration) {
	db = database
	log = logger
	config = configuration
}

// Middleware to verify if host is config host

func Start(upstreamAddr string, hostnamePerPath utils.HostnameAndPortPerPath, router *http.ServeMux) map[string]http.Handler {
	log.Info("Starting Asynchronous Report")
	interceptor.DB = db
	websocket.DB = db
	download.DB = db
	interceptor.Init(log, config, upstreamAddr)
	websocket.SetLog(log)
	download.SetLog(log)
	log.Info("Add routes")
	router.Handle("/js-websocket/", http.HandlerFunc(jswebsocket.Send))
	router.Handle("/download/", http.HandlerFunc(download.Download))
	router.Handle("/websocket/", http.HandlerFunc(websocket.WebSocketHandler))
	hs := make(HostSwitch)

	for path, hostnameToServe := range hostnamePerPath {
		for h := range hostnameToServe {
			router.Handle(path, http.HandlerFunc(hs.ServeHTTP))
			log.Info("Route added", zap.String("hostname", h))
			hs[h] = router
		}
	}
	//Log config object
	// log.Info("Config", zap.Any("config", config))
	return route
}
