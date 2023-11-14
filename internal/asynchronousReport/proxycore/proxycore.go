package proxycore

import (
	utils "FlexiProxyHub/internal/utils/generic"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"

	"go.uber.org/zap"
)

func Proxy_Pass(url *url.URL, r *http.Request, w http.ResponseWriter, log *zap.Logger, config *utils.Configuration) {
	proxy := httputil.NewSingleHostReverseProxy(url)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = url.Host
		for header, values := range r.Header {
			for _, value := range values {
				req.Header.Add(header, value)
			}
		}
	}
	proxy.ServeHTTP(w, r)
}

func Proxy_Pass_With_Recorder(url *url.URL, r *http.Request, recorder *httptest.ResponseRecorder, log *zap.Logger, config *utils.Configuration) {
	log.Info("url", zap.Any("url", url))
	proxy := httputil.NewSingleHostReverseProxy(url)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = url.Host
		for header, values := range r.Header {
			for _, value := range values {
				req.Header.Add(header, value)
			}
		}
	}
	proxy.ServeHTTP(recorder, r)
}
