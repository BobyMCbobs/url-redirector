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

// log all requests
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%v %v %v %v %v %v %v", r.Header["User-Agent"], r.Method, r.URL, r.Proto, r.Response, r.RemoteAddr, r.Header)
		next.ServeHTTP(w, r)
	})
}

// load and parse the config.yaml file
func ReadConfigYAML() (output types.ConfigYAML) {
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
func APIshortLink(w http.ResponseWriter, r *http.Request) {
	configYAML := ReadConfigYAML()
	vars := mux.Vars(r)
	redirectURL := configYAML.Routes[vars["link"]]
	if redirectURL == "" {
		if configYAML.Wildcard == "" {
			w.WriteHeader(404)
			w.Write([]byte(`404 page not found`))
			return
		} else {
			http.Redirect(w, r, configYAML.Wildcard, 302)
			return
		}
	}
	http.Redirect(w, r, redirectURL, 302)
}

func APIroot(w http.ResponseWriter, r *http.Request) {
	configYAML := ReadConfigYAML()
	if configYAML.Root == "" {
		w.WriteHeader(404)
		w.Write([]byte(`404 page not found`))
		return
	}
	http.Redirect(w, r, configYAML.Root, 302)
}

// print a table of the environment variables
func PrintEnvConfig() {
	fmt.Println()
	data := [][]string{
		[]string{"APP_PORT", appPort},
		[]string{"APP_CONFIG_YAML", appConfigYAMLlocation},
		[]string{"APP_USE_LOGGING", appUseLogging},
		[]string{"APP_LOG_FILE", appLogFileLocation},
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
	router.HandleFunc("/{link:[a-zA-Z0-9]+}", APIshortLink)
	router.HandleFunc("/", APIroot)
	router.Use(RequestLogger)
	srv := &http.Server{
		Handler:      router,
		Addr:         appPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Println("Listening on", appPort)
	logger.Fatal(srv.ListenAndServe())
}
