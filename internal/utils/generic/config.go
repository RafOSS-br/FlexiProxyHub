package utils

import (
	"encoding/json"
	"os"
	"strings"

	"go.uber.org/zap"
)

type ProxyMode int

const (
	Normal ProxyMode = iota
	AsynchronousReport
)
// PROXY MODE, UPSTREAM ADDR, HOSTNAME, PATH = BOOL
type HostnameAndPortPerPathPerProxyModeAndUpstreamAddr map[int]map[string]map[string]map[string]bool
// UPSTREAM ADDR, HOSTNAME, PATH = BOOL
type HostnameAndPortPerPathAndUpstreamAddr map[string]map[string]map[string]bool
// HOSTNAME, PATH = BOOL
type HostnameAndPortPerPath map[string]map[string]bool

type ProxyToList string

// "routes": [{"host": "localhost", "port": "8080", "path": ["/teste.txt"]}]
type Routes struct {
	Host string   `json:"host"`
	Path []string `json:"path"`
}

type Router struct {
	Routes  []Routes  `json:"routes"`
	Mode    ProxyMode `json:"mode"`
	ProxyTo string    `json:"proxy_to"`
}

type Listener struct {
	Address string
	Port    string
}

type Configuration struct {
	LogLevel          string
	VisibleHeaders    string
	HeaderToReplicate Headers
	BodyMaxLen        int
	Proxy             []Router
	Listener          Listener
}

// Create configuration from environment variables
func CreateConfigFromEnv(log *zap.Logger) *Configuration {
	return &Configuration{
		LogLevel:          verifyLogLevel(os.Getenv("LOG_LEVEL")),
		VisibleHeaders:    verifyVisibleHeaders(os.Getenv("VISIBLE_HEADERS")),
		HeaderToReplicate: verifyHeaderToReplicate(os.Getenv("HEADER_TO_REPLICATE")),
		BodyMaxLen:        verifyBodyMaxLen(os.Getenv("LOG_BODY_MAX_SIZE"), log),
		Proxy:             ConvertJSONToRouter(os.Getenv("PROXY_CONFIGURATION"), log),
		Listener:          verifyListener(os.Getenv("LISTEN_PORT"), os.Getenv("LISTEN_HOST"), log),
	}
}

type Headers []string

// Append HEADER_TO_REPLICATE to Headers type
func (h Headers) append(header string) Headers {
	return append(h, header)
}

// Verify if HEADER_TO_REPLICATE is valid and return Headers type result from HEADER_TO_REPLICATE splitted by ","
// If not, not set.
func verifyHeaderToReplicate(headerToReplicate string) Headers {
	if headerToReplicate == "" {
		return Headers{}
	}
	var headers Headers
	//Check if each one is valid
	for _, header := range strings.Split(headerToReplicate, ",") {
		if header != "" {
			headers = headers.append(header)
		}
	}
	return headers
}

// Verify LISTEN_PORT and LISTEN_HOST and return Listener struct
// If not, set to "8080" and "localhost"
func verifyListener(listenPort string, listenHost string, log *zap.Logger) Listener {
	if listenPort == "" {
		listenPort = "8080"
	}
	if listenHost == "" {
		listenHost = "localhost"
	}
	listenPortInt, err := ConvertStringToInt(listenPort)
	if err != nil {
		log.Warn("Invalid LISTEN_PORT")
		listenPortInt = 8080
	}
	return Listener{
		Address: listenHost,
		Port:    ConvertIntToString(listenPortInt),
	}
}

// Verify if LOG_BODY_MAX_SIZE is valid
// If not, set to 255
func verifyBodyMaxLen(bodyMaxLen string, log *zap.Logger) int {
	if bodyMaxLen == "" {
		log.Warn("Invalid LOG_BODY_MAX_SIZE")
		return 255
	}
	if bodyManLenInt, err := ConvertStringToInt(bodyMaxLen); err != nil {
		log.Warn("Invalid LOG_BODY_MAX_SIZE")
		return 255
	} else {
		return bodyManLenInt
	}
}

// Verify if VisibleHeaders is valid
// If not, set to "host,x-request-id,x-real-ip,content-length,user-agent,accept-encoding,content-type,custom-app-headers"
func verifyVisibleHeaders(visibleHeaders string) string {
	if visibleHeaders == "" {
		return "host,x-request-id,x-real-ip,content-length,user-agent,accept-encoding,content-type,custom-app-headers"
	}
	return visibleHeaders
}

// Verify if LOG_LEVEL is valid
// If not, set to "info"
func verifyLogLevel(logLevel string) string {
	switch strings.ToLower(logLevel) {
	case "debug":
		return "debug"
	case "info":
		return "info"
	case "warn":
		return "warn"
	case "error":
		return "error"
	default:
		return "info"
	}
}

// Convert JSON to list of router struct 
// [{"host": "localhost", "routes": ["/teste.txt"], "mode": 0}]
// [{"host": "localhost", "routes": ["/teste.txt"], "mode": 1}]
func ConvertJSONToRouter(jsonData string, log *zap.Logger) []Router {
	var routers []Router
	err := json.Unmarshal([]byte(jsonData), &routers)
	if err != nil {
		log.Fatal("Error parsing JSON", zap.Error(err))
	}
	return routers
}
