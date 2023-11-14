package router

import (
	utils "FlexiProxyHub/internal/utils/generic"
	"FlexiProxyHub/internal/utils/logs"
	asynchronousreport "FlexiProxyHub/pkg/asynchronousReport"
	"FlexiProxyHub/pkg/database"
	"net/http"

	"go.uber.org/zap"
)

var config *utils.Configuration
var log *zap.Logger

func Router(c *utils.Configuration, l *zap.Logger) {
	config = c
	log = l
	database.Init()
	router := http.NewServeMux()
	ret := getHostnameAndPortPerPathPerProxyModeAndUpstreamAddr()
	for proxy := range config.Proxy {
		switch config.Proxy[proxy].Mode {
		case utils.Normal:
			// normalProxy()
			continue
		case utils.AsynchronousReport:
			asynchronousReport(ret[int(utils.AsynchronousReport)], router)
		}
	}
	// Run server
	log.Info("Starting server " + config.Listener.Address + ":" + config.Listener.Port)
	http.ListenAndServe(config.Listener.Address+":"+config.Listener.Port, logs.LogRequestMiddleware(router))
}

// Aggregate host and port per path and proxy type
// return
// utils.HostnameAndPortPerPathPerProxyModeAndUpstreamAddr
// or
// map[proxyMode]map[proxy_to]map[path]map[host]bool
func getHostnameAndPortPerPathPerProxyModeAndUpstreamAddr() utils.HostnameAndPortPerPathPerProxyModeAndUpstreamAddr {
	hostnameAndPortPerPathPerProxyModeAndUpstreamAddr := make(map[int]map[string]map[string]map[string]bool)
	for _, proxy := range config.Proxy {
		proxyMode := proxy.Mode
		log.Debug("Proxy Mode", zap.Int("proxyMode", int(proxyMode)))
		hostnameAndPortPerPathPerProxyModeAndUpstreamAddr[int(proxyMode)] = make(map[string]map[string]map[string]bool)
		for _, routes := range proxy.Routes {
			hostnameAndPortPerPathPerProxyModeAndUpstreamAddr[int(proxyMode)][proxy.ProxyTo] = make(map[string]map[string]bool)
			for _, path := range routes.Path {
				hostnameAndPortPerPathPerProxyModeAndUpstreamAddr[int(proxyMode)][proxy.ProxyTo][path] = make(map[string]bool)
				hostnameAndPortPerPathPerProxyModeAndUpstreamAddr[int(proxyMode)][proxy.ProxyTo][path][routes.Host] = true
			}
		}
	}
	return hostnameAndPortPerPathPerProxyModeAndUpstreamAddr
}

// Convert hostnameAndPortPerPathPerProxyModeAndUpstreamAddr to ListProxyConfig
// map[proxyType]map[proxy_to]map[path]map[counter 0++]host:port
// to
// ListProxyConfig
// func (MapHostnamePortByPathAndProxyMode *hostnameAndPortPerPathPerProxyModeAndUpstreamAddr) to() ListProxyConfig {
// 	var listProxyConfig ListProxyConfig
// 	for proxyMode, mapProxyTo := range *MapHostnamePortByPathAndProxyMode {
// 		for proxyTo, mapPath := range mapProxyTo {
// 			for path, listHostnameAndPort := range mapPath {
// 				for _, hostnameAndPort := range listHostnameAndPort {
// 					listProxyConfig = append(listProxyConfig, utils.utils.HostnameAndPortPerPathAndUpstreamAddr{
// 						Hostname: hostnameAndPort,
// 						Path:     path,
// 						Upstream: proxyTo,
// 					})
// 				}
// 			}
// 		}
// 	}
// }

// Create Asynchronous Report Handler
func asynchronousReport(list utils.HostnameAndPortPerPathAndUpstreamAddr, router *http.ServeMux) {
	// Create Asynchronous Report routes
	log.Debug("List", zap.Any("HostnameAndPortPerPathAndUpstreamAddr", list))
	for upstreamAddr, hostnameAndPortPerPath := range list {
		log.Debug("Upstream addr", zap.String("upstreamAddr", upstreamAddr))
		asynchronousreport.Init(database.DB, log, config)
		asyncReportRoutes := asynchronousreport.Start(upstreamAddr, hostnameAndPortPerPath, router)
		for route, handler := range asyncReportRoutes {
			http.Handle(route, handler)
		}

	}
}
