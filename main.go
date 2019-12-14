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
	"gopkg.in/yaml.v2"
)

type DeploymentYAML map[string]string

var (
	deploymentYAMLlocation = GetDeploymentLocation()
	logFileLocation        = GetLogFileLocation()
	outfile, _             = os.OpenFile(logFileLocation, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	logger                 = log.New(outfile, "", 3)
	multiWriter            = io.MultiWriter(os.Stdout, outfile)
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
func ReadDeploymentYAML() (output DeploymentYAML) {
	yamlFile, err := ioutil.ReadFile(deploymentYAMLlocation)
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
func CheckForDeploymentYAML() {
	if _, err := os.Stat(deploymentYAMLlocation); err != nil {
		logger.Fatalf("File %v does not exist, please create or mount it.\n", deploymentYAMLlocation)
	}
	logger.Printf("Using file %v as the deployment configuration\n", deploymentYAMLlocation)
}

// handle the url variables on /{link}
func APIshortLink(w http.ResponseWriter, r *http.Request) {
	deploymentYAML := ReadDeploymentYAML()
	vars := mux.Vars(r)
	redirectURL := deploymentYAML[vars["link"]]
	if redirectURL == "" {
		return
	}
	http.Redirect(w, r, redirectURL, 302)
}

// manage starting of webserver
func HandleWebserver() {
	port := GetAppPort()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./robots.txt")
	})
	router.HandleFunc("/{link:[a-zA-Z0-9]+}", APIshortLink)
	router.Use(RequestLogger)
	srv := &http.Server{
		Handler:      router,
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Println("Listening on", port)
	logger.Fatal(srv.ListenAndServe())
}

func main() {
	logger.SetOutput(multiWriter)
	CheckForDeploymentYAML()
	HandleWebserver()
}
