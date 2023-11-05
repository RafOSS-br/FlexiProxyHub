package logs

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

var Log *zap.Logger

const visibleHeaders = "host,x-request-id,x-real-ip,content-length,user-agent,accept-encoding,content-type,custom-app-headers"

var bodyMaxLen = 255

var debug bool = false

func SetDebugMode() {
	debug = true
}

func DisableDebugMode() {
	debug = false
}

func SetBodyMaxLen(n int) {
	bodyMaxLen = n
}

func getRequestBodyFromRequest(r *http.Request) string {
	if r.Body == nil {
		return ""
	}
	if r.ContentLength > int64(bodyMaxLen) {
		return fmt.Sprintf("Body too long (%d bytes)", r.ContentLength)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	return buf.String()
}

func getResponseBodyFromResponse(r *http.Response) string {
	if r.Body == nil {
		return ""
	}
	if r.ContentLength > int64(bodyMaxLen) {
		return fmt.Sprintf("Body too long (%d bytes)", r.ContentLength)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	return buf.String()
}

type Severity struct {
	string string
	int    int
}

type SeverityList map[string]int

func (s Severity) String() string {
	return s.string
}
func (s Severity) Int() int {
	return s.int
}

func (s SeverityList) Get(key string) Severity {
	return Severity{key, s[key]}
}

func NewSeverityList() SeverityList {
	return SeverityList{
		"trace": 10,
		"debug": 20,
		"info":  30,
		"error": 40,
		"fatal": 50,
	}
}

// Truncate string if it's longer than bodyMaxLen. If truncated, add "..." to the end.
func Truncate(s string) string {
	if len(s) > bodyMaxLen {
		return s[:bodyMaxLen-3] + "..."
	}
	return s
}

// Time now in "Fri, 19 May 2023 13:26:13 GMT"
type Time string

func (t Time) String() string {
	return string(t)
}
func TimeNow() Time {
	return Time(time.Now().UTC().Format(http.TimeFormat))
}

//	"request": {
//		"ip": "127.0.0.6",
//		"url": "/proc",
//		"method": "POST",
//		"headers": {
//			"host": "dev-krakend.bolt.com.br",
//			"x-request-id": "71feefa0002fe17d469a0feb82fabb28",
//			"x-real-ip": "127.0.0.6",
//			"content-length": "217",
//			"user-agent": "curl/7.86.0",
//			"accept-encoding": "gzip",
//			"content-type": "application/json",
//			"custom-app-headers": "Be careful with credentials"
//		}
//	}
type Request struct {
	IP      string
	URL     string
	Method  string
	Headers *HeaderList
	Body    string
}
type HeaderList map[string]string

func (h HeaderList) GetVisible() HeaderList {
	headers := HeaderList{}
	for _, header := range strings.Split(visibleHeaders, ",") {
		headers[header] = h[header]
	}
	return headers
}

// "response": {
// 	"status": 500,
//  "body": "{ \"bodyContent\" }",
// 	"message": ""
// }

type Response struct {
	Status  int
	Body    string
	Message string
	Headers *HeaderList
}

type Logging struct {
	Level    *Severity
	Time     *Time
	PID      int
	Hostname string
	Request  *Request
	Message  string
	Response *Response
}

func Init() {
	var err error
	Log, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

func RequestToLog(r *http.Request) {
	var visibleBody string
	if debug {
		visibleBody = getRequestBodyFromRequest(r)
	} else {
		visibleBody = ""
	}
	request := &Request{
		IP:      r.RemoteAddr,
		URL:     r.URL.String(),
		Method:  r.Method,
		Headers: GetHeaders(r.Header),
		Body:    visibleBody,
	}
	Log.Info("request", zap.Any("request", request))
}

func ResponseToLog(r *http.Response) {
	var visibleBody string
	if debug {
		visibleBody = getResponseBodyFromResponse(r)
	} else {
		visibleBody = ""
	}
	response := &Response{
		Status:  r.StatusCode,
		Body:    visibleBody,
		Message: r.Status,
		Headers: GetHeaders(r.Header),
	}
	Log.Info("response", zap.Any("response", response))
}

// LogMiddleware is a middleware that logs requests.
func LogRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		RequestToLog(r)
		next.ServeHTTP(w, r)
	})
}


func GetHeaders(h http.Header) *HeaderList {
	headers := make(HeaderList)
	for k, v := range h {
		if k == "Authorization" {
			headers[k] = "REDACTED"
		} else {
			headers[k] = strings.Join(v, ",")
		}
	}
	return &headers
}
