package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	port = "80"
)

func main() {
	log.Println("Initializing Badges server...")
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/health", healthHandler)
	log.Printf("Initialziation complete. Serving on port %v.\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func healthHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.WriteHeader(http.StatusOK)
}

func rootHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		getBadge(responseWriter, request.URL.Query())
	} else if request.Method == http.MethodPost {
		postBadge(responseWriter, request.URL.Query())
	} else {
		responseWriter.WriteHeader(http.StatusBadRequest)
	}
	return
}

func postBadge(responseWriter http.ResponseWriter, parameters map[string][]string) {
	expectedParameters := []string{"project", "item", "value", "color"}
	if !parametersAreValid(parameters, expectedParameters) {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	project := parameters["project"][0]
	item := parameters["item"][0]
	value := parameters["value"][0]
	color := parameters["color"][0]
	fileName := generateFileName(project, item)
	data := formatData(value, color)
	os.Mkdir("data", 0777)
	err := ioutil.WriteFile(fileName, data, 0777)
	if err != nil {
		log.Printf("The following error occured while writing to the data file (%v).\n", fileName)
		log.Println(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
	} else {
		responseWriter.WriteHeader(http.StatusOK)
	}
}

func getBadge(responseWriter http.ResponseWriter, parameters map[string][]string) {
	expectedParameters := []string{"project", "item"}
	if !parametersAreValid(parameters, expectedParameters) {
		writeUnknownBadge(responseWriter)
		return
	}
	project := parameters["project"][0]
	item := parameters["item"][0]
	fileName := generateFileName(project, item)
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("The following error occured while finding the badge: '%v'\n", fileName)
		log.Println(err)
		writeUnknownBadge(responseWriter)
		return
	}
	writeBadge(responseWriter, item, data)
}

func parametersAreValid(parameters map[string][]string, expectedKeys []string) bool {
	for _, key := range expectedKeys {
		if _, found := parameters[key]; !found {
			return false
		}
	}
	return true
}

func writeUnknownBadge(responseWriter http.ResponseWriter) {
	newURI := generateShieldURI("error", "unknownbadge", "red")
	sendResponse(responseWriter, newURI)
}

func writeBadge(responseWriter http.ResponseWriter, item string, data []byte) {
	parsedData := parseData(data)
	value := parsedData["value"]
	color := parsedData["color"]
	newURI := generateShieldURI(item, value, color)
	sendResponse(responseWriter, newURI)
}

func sendResponse(rw http.ResponseWriter, URI string) {
	// Intercept response to control cache-related headers.
	// Reference: https://github.com/github/markup/issues/224
	r, err := http.Get(URI)
	if err != nil {
		log.Println(err)
		return
	}
	etag := fmt.Sprintf("\"%v\"", time.Now().UnixNano())
	now := time.Now().UTC().Format(time.RFC1123)
	rw.Header().Set("Cache-Control", "max-age=1")
	rw.Header().Add("Cache-Control", "no-cache")
	rw.Header().Add("Expires", now)
	rw.Header().Add("Last-Modified", now)
	rw.Header().Add("Pragma", "no-cache")
	rw.Header().Set("ETag", etag)

	// Forward original headers and body
	rw.Header().Set("Date", r.Header.Get("Date"))
	rw.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	rw.Header().Set("Server", r.Header.Get("Server"))
	rw.Header().Set("Cf-Ray", r.Header.Get("Cf-Ray"))
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.Bytes()
	rw.WriteHeader(http.StatusOK)
	rw.Write(body)
}

func generateShieldURI(item, value, color string) string {
	shieldsIoFormat := "https://img.shields.io/badge/%v-%v-%v.svg"
	return fmt.Sprintf(shieldsIoFormat, item, value, color)
}

func generateFileName(project, item string) string {
	dataDirectory := "data"
	return fmt.Sprintf("%v/%v%v", dataDirectory, project, item)
}

func formatData(value, color string) []byte {
	return []byte(fmt.Sprintf("%v|%v", value, color))
}

func parseData(data []byte) map[string]string {
	dataString := string(data)
	parsedData := make(map[string]string)
	parsedData["value"] = strings.Split(dataString, "|")[0]
	parsedData["color"] = strings.Split(dataString, "|")[1]
	return parsedData
}
