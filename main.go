package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
)

type ConfigYAML map[string]string

var (
	appPort               = GetAppPort()
	appConfigYAMLlocation = GetDeploymentLocation()
	appUseLogging         = os.Getenv("APP_USE_LOGGING")
	appLogFileLocation    = GetLogFileLocation()
	loggerOutFile, _      = os.OpenFile(appLogFileLocation, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	logger                = log.New(loggerOutFile, "", 3)
	multiWriter           = io.MultiWriter(os.Stdout, loggerOutFile)
)

// determine the port for the app to run on
func GetAppPort() (output string) {
	output = os.Getenv("APP_PORT")
	if output == "" {
		output = ":8080"
	}
	return output
}

// determine the location of the config.yaml
func GetDeploymentLocation() (output string) {
	output = os.Getenv("APP_CONFIG_YAML")
	if output == "" {
		output = "./config.yaml"
	}
	return output
}

// determine the location of the config.yaml
func GetLogFileLocation() (output string) {
	output = os.Getenv("APP_LOG_FILE")
	if output == "" {
		output = "./redirector.log"
	}
	return output
}

// log all requests
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%v %v %v %v %v %v %v", r.Header["User-Agent"], r.Method, r.URL, r.Proto, r.Response, r.RemoteAddr, r.Header)
		next.ServeHTTP(w, r)
	})
}

// load and parse the config.yaml file
func ReadConfigYAML() (output ConfigYAML) {
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
	redirectURL := configYAML[vars["link"]]
	if redirectURL == "" {
		return
	}
	http.Redirect(w, r, redirectURL, 302)
}

// manage starting of webserver
func HandleWebserver() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./robots.txt")
	})
	router.HandleFunc("/{link:[a-zA-Z0-9]+}", APIshortLink)
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

func PrintEnvConfig() {
	data := [][]string{
		[]string{"APP_PORT", appPort},
		[]string{"APP_CONFIG_YAML", appConfigYAMLlocation},
		[]string{"APP_USE_LOGGING", appUseLogging},
		[]string{"APP_LOG_FILE", appLogFileLocation},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Value"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(data)
	table.Render()
	fmt.Println()
}

func main() {
	if appUseLogging != "false" {
		logger.SetOutput(multiWriter)
	}
	logger.Println("Warming up")
	CheckForConfigYAML()
	PrintEnvConfig()
	HandleWebserver()
}
