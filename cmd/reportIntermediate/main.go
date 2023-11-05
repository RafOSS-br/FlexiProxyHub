package main

import (
	"net/http"
	"reportIntermediate/internal/utils/logs"
	"reportIntermediate/pkg/interceptor"
)

func main() {
	logs.Init()
	logs.SetDebugMode()
	http.Handle("/", logs.LogRequestMiddleware(http.HandlerFunc(interceptor.ServeHTTP)))
	http.ListenAndServe(":8080", nil)
}
