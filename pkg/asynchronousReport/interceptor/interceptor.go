package interceptor

import (
	"FlexiProxyHub/internal/asynchronousReport/proxycore"
	utils "FlexiProxyHub/internal/utils/generic"
	"FlexiProxyHub/pkg/database"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"

	"go.uber.org/zap"
)

const routes = "/teste.txt"

var config *utils.Configuration

type HostSwitch map[string]http.Handler

var DB *database.Database
var hostname string
var log *zap.Logger
var WG sync.WaitGroup

var routesList []string

func Init(logging *zap.Logger, configuration *utils.Configuration, host string) {
	log = logging
	config = configuration
	hostname = host
}

func routesSplittedByCommaToRoutesList() {
	routesList = strings.Split(routes, ",")
}
func pathIsInRoutesList(path string) bool {
	for _, route := range routesList {
		if route == path {
			return true
		}
	}
	return false
}

func (hs HostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	routesSplittedByCommaToRoutesList()
	host := r.Host
	log.Debug("Serve host", zap.String("host", host))
	if hs[host] != nil {
		log.Debug("Serve host", zap.String("host", host))
		if pathIsInRoutesList(r.URL.Path) {
			// proxy request and download to local file /tmp/$(hash) and return the hash in cookie
			remote, err := url.Parse(hostname)
			if err != nil {
				log.Error("Error parsing url", zap.Error(err))
			}
			// generate hash from ci_session cookie, timestamp and path
			filename, err := utils.CreateFileName(r)
			if err != nil {
				log.Error("Error creating filename", zap.Error(err))
			}
			// w.Header().Set("Content-Type", "text/plain")
			// w.Write([]byte(filename))
			// if f, ok := w.(http.Flusher); ok {
			// 	f.Flush()
			// }
			// return
			cookie := http.Cookie{Name: "filename", Value: filename}
			http.SetCookie(w, &cookie)
			database.DB = DB
			database.FlagHashAsDownloading(filename)
			http.Redirect(w, r, "/js-websocket/"+filename, http.StatusFound)
			go func(req *http.Request) {
				recorder := httptest.NewRecorder()
				proxycore.Proxy_Pass_With_Recorder(remote, req, recorder, log, config)
				// Check if the recorder has any data
				if recorder.Body.Len() == 0 {
					log.Error("Recorder is empty")
				} else {
					// Save the response body to file
					utils.SaveResponseBodyToFile(filename, recorder, log)
					// delete file flag and complete
					database.SaveFromHashAndComplete(filename, recorder.Header().Get("Content-Type"), recorder.Header().Get("Content-Disposition"))
				}
			}(r.Clone(context.Background()))
		}
		// else {
		// 	remote, err := url.Parse(hostname)
		// 	if err != nil {
		// 		log.Error("Error parsing url", zap.Error(err))
		// 	}
		// 	proxycore.Proxy_Pass(remote, r, w, log, config)

		// }
	}
}
