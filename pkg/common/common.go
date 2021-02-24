// This program is free software: you can redistribute it and/or modify
// it under the terms of the Affero GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the Affero GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Package common ...
// generally used functions
package common

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/olekukonko/tablewriter"
	"gitlab.com/bobymcbobs/url-redirector/pkg/types"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// AppVars ...
// defaults
var (
	AppBuildVersion       = "0.0.0"
	AppBuildHash          = "???"
	AppBuildDate          = "???"
	AppBuildMode          = "development"
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

// GetAppBuildVersion ...
// return the version of the current FlatTrack instance
func GetAppBuildVersion() string {
	return AppBuildVersion
}

// GetAppBuildHash ...
// return the commit which the current FlatTrack binary was built from
func GetAppBuildHash() string {
	return AppBuildHash
}

// GetAppBuildDate ...
// return the build date of FlatTrack
func GetAppBuildDate() string {
	return AppBuildDate
}

// GetAppBuildMode ...
// return the mode that the app is built in
func GetAppBuildMode() string {
	return AppBuildMode
}

// GetEnvOrDefault ...
// if an environment variable exists return it, otherwise return a default value
func GetEnvOrDefault(envName string, defaultValue string) (output string) {
	output = os.Getenv(envName)
	if output == "" {
		output = defaultValue
	}
	return output
}

// GetAppPort ...
// determine the port for the app to run on
func GetAppPort() (output string) {
	return GetEnvOrDefault("APP_PORT", ":8080")
}

// GetAppPortTLS ...
// determine the tls port for the app to run on
func GetAppPortTLS() (output string) {
	return GetEnvOrDefault("APP_PORT_TLS", ":4433")
}

// GetAppUseTLS ...
// determine if the app should host with TLS
func GetAppUseTLS() (output string) {
	return GetEnvOrDefault("APP_USE_TLS", "false")
}

// GetAppTLSpublicCert ...
// determine path to the public SSL cert
func GetAppTLSpublicCert() (output string) {
	return GetEnvOrDefault("APP_TLS_PUBLIC_CERT", "server.crt")
}

// GetAppTLSprivateCert ...
// determine path to the private SSL cert
func GetAppTLSprivateCert() (output string) {
	return GetEnvOrDefault("APP_TLS_PRIVATE_CERT", "server.key")
}

// GetDeploymentLocation ...
// determine the location of the config.yaml
func GetDeploymentLocation() (output string) {
	return GetEnvOrDefault("APP_CONFIG_YAML", "./config.yaml")
}

// GetUseLogging ...
func GetUseLogging() (output string) {
	return GetEnvOrDefault("APP_USE_LOGGING", "false")
}

// GetLogFileLocation ...
// determine the location of the log file
func GetLogFileLocation() (output string) {
	return GetEnvOrDefault("APP_LOG_FILE", "./redirector.log")
}

// GetRequestHost
// returns the request host
func GetRequestHost(r *http.Request) string {
	return r.Host
}

// GetConfigHost
// determine if there is config available for the host
func GetConfigHost(configYAML types.RouteHosts, r *http.Request) (useHost string) {
	requestHost := GetRequestHost(r)
	// if the host has no routes
	if len(configYAML[requestHost].Routes) == 0 {
		// if the wildcard host has no routes
		if len(configYAML["*"].Routes) == 0 && configYAML["*"].Root == "" && configYAML["*"].Wildcard == "" {
			return ""
		}
		return "*"
	}
	return requestHost
}

// RequestLogger ...
// log all requests
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%v %v %v %v %v %v %v %v", r.Header["User-Agent"], r.Method, r.Host, r.URL, r.Proto, r.Response, r.RemoteAddr, r.Header)
		next.ServeHTTP(w, r)
	})
}

// ReadConfigYAML ...
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

// CheckForConfigYAML ...
// check if the config.yaml exists
func CheckForConfigYAML() {
	if _, err := os.Stat(appConfigYAMLlocation); err != nil {
		logger.Fatalf("File %v does not exist, please create or mount it.\n", appConfigYAMLlocation)
	}
}

// APIroutesHandler ...
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
		}
		http.Redirect(w, r, configYAML[requestHost].Wildcard, 302)
		return
	}
	http.Redirect(w, r, redirectURL, 302)
}

// APIrootRouteHandler ...
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

// PrintEnvConfig ...
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

// Logger ...
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

// HandleWebserver ...
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
