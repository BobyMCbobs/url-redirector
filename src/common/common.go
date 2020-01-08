package common

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/olekukonko/tablewriter"
	"gitlab.com/bobymcbobs/url-redirector/src/types"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	logger                = Logger()
	appPort               = GetAppPort()
	appPortTLS            = GetAppPortTLS()
	appUseTLS             = GetAppUseTLS()
	appTLSpublicCert      = GetAppTLSpublicCert()
	appTLSprivateCert     = GetAppTLSprivateCert()
	appConfigYAMLlocation = GetDeploymentLocation()
	appUseLogging         = GetUseLogging()
	appLogFileLocation    = GetLogFileLocation()
)

// if an environment variable exists return it, otherwise return a default value
func GetEnvOrDefault(envName string, defaultValue string) (output string) {
	output = os.Getenv(envName)
	if output == "" {
		output = defaultValue
	}
	return output
}

// determine the port for the app to run on
func GetAppPort() (output string) {
	return GetEnvOrDefault("APP_PORT", ":8080")
}

// determine the tls port for the app to run on
func GetAppPortTLS() (output string) {
	return GetEnvOrDefault("APP_PORT_TLS", ":4433")
}

// determine if the app should host with TLS
func GetAppUseTLS() (output string) {
	return GetEnvOrDefault("APP_USE_TLS", "false")
}

// determine path to the public SSL cert
func GetAppTLSpublicCert() (output string) {
	return GetEnvOrDefault("APP_TLS_PUBLIC_CERT", "server.crt")
}

// determine path to the private SSL cert
func GetAppTLSprivateCert() (output string) {
	return GetEnvOrDefault("APP_TLS_PRIVATE_CERT", "server.key")
}

// determine the location of the config.yaml
func GetDeploymentLocation() (output string) {
	return GetEnvOrDefault("APP_CONFIG_YAML", "./config.yaml")
}

func GetUseLogging() (output string) {
	return GetEnvOrDefault("APP_USE_LOGGING", "false")
}

// determine the location of the log file
func GetLogFileLocation() (output string) {
	return GetEnvOrDefault("APP_LOG_FILE", "./redirector.log")
}

// returns the request host
func GetRequestHost(r *http.Request) string {
	return r.Host
}

// determine if there is config available for the host
func GetConfigHost(configYAML types.RouteHosts, r *http.Request) (useHost string) {
	requestHost := GetRequestHost(r)
	// if the host has no routes
	if len(configYAML[requestHost].Routes) == 0 {
		// if the wildcard host has no routes
		if len(configYAML["*"].Routes) == 0 && configYAML["*"].Root == "" && configYAML["*"].Wildcard == "" {
			return ""
		} else {
			return "*"
		}
	}
	return requestHost
}

// log all requests
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%v %v %v %v %v %v %v %v", r.Header["User-Agent"], r.Method, r.Host, r.URL, r.Proto, r.Response, r.RemoteAddr, r.Header)
		next.ServeHTTP(w, r)
	})
}

// load and parse the config.yaml file
func ReadConfigYAML() (output types.RouteHosts) {
	yamlFile, err := ioutil.ReadFile(appConfigYAMLlocation)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return
	}
	err = yaml.Unmarshal(yamlFile, &output)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
		return
	}
	return output
}

// check if the config.yaml exists
func CheckForConfigYAML() {
	if _, err := os.Stat(appConfigYAMLlocation); err != nil {
		logger.Fatalf("File %v does not exist, please create or mount it.\n", appConfigYAMLlocation)
	}
}

// handle the url variables on /{link}
func APIroutesHandler(w http.ResponseWriter, r *http.Request) {
	configYAML := ReadConfigYAML()
	vars := mux.Vars(r)
	requestHost := GetConfigHost(configYAML, r)
	redirectURL := configYAML[requestHost].Routes[vars["link"]]
	fmt.Println(redirectURL)
	if redirectURL == "" {
		if configYAML[requestHost].Wildcard == "" {
			w.WriteHeader(404)
			w.Write([]byte(`404 page not found`))
			return
		} else {
			http.Redirect(w, r, configYAML[requestHost].Wildcard, 302)
			return
		}
	}
	http.Redirect(w, r, redirectURL, 302)
}

// handle root requests
func APIrootRouteHandler(w http.ResponseWriter, r *http.Request) {
	configYAML := ReadConfigYAML()
	requestHost := GetConfigHost(configYAML, r)
	if configYAML[requestHost].Root == "" {
		w.WriteHeader(404)
		w.Write([]byte(`404 page not found`))
		return
	}
	http.Redirect(w, r, configYAML[requestHost].Root, 302)
}

// print a table of the environment variables
func PrintEnvConfig() {
	fmt.Println()
	data := [][]string{
		[]string{"APP_PORT", appPort},
		[]string{"APP_CONFIG_YAML", appConfigYAMLlocation},
		[]string{"APP_USE_LOGGING", appUseLogging},
		[]string{"APP_LOG_FILE", appLogFileLocation},
		[]string{"APP_USE_TLS", appUseTLS},
		[]string{"APP_TLS_PUBLIC_CERT", appTLSpublicCert},
		[]string{"APP_TLS_PRIVATE_CERT", appTLSprivateCert},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Value"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data)
	table.Render()
	fmt.Println()
}

// request file logger
func Logger() *log.Logger {
	logger := log.New(os.Stdout, "", 3)
	if appUseLogging == "true" {
		appLogFileLocation := GetLogFileLocation()
		loggerOutFile, _ := os.OpenFile(appLogFileLocation, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		logger = log.New(loggerOutFile, "", 3)
		multiWriter := io.MultiWriter(os.Stdout, loggerOutFile)
		logger.SetOutput(multiWriter)
	}
	return logger
}

// manage starting of webserver
func HandleWebserver() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./robots.txt")
	})
	router.HandleFunc("/{link:[a-zA-Z0-9]+}", APIroutesHandler)
	router.HandleFunc("/", APIrootRouteHandler)
	router.Use(RequestLogger)
	srv := &http.Server{
		Handler:      router,
		Addr:         appPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if appUseTLS == "true" {
		logger.Println("Listening on", appPortTLS)
		srv.Addr = appPortTLS
		logger.Fatal(srv.ListenAndServeTLS(appTLSpublicCert, appTLSprivateCert))
	} else {
		logger.Println("Listening on", appPort)
		logger.Fatal(srv.ListenAndServe())
	}
}
