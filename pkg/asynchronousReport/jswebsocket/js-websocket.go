package jswebsocket

import (
	"net/http"
)

// Send js-websocket.html to client
func Send(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pkg/asynchronousReport/jswebsocket/js-websocket.html")
}
