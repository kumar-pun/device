package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Device struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Device  string `json:"device"`
	Project string `json:"project"`
}

var devices = []Device{
	Device{Id: 1, Name: "Jan Doe", Device: "mac", Project: "Vivint"},
}

type Response struct {
	Method  string `json:"method"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func NewResponse(method, message string, status int) Response {

	return Response{Method: method, Message: message, Status: status}

}

func HttpInfo(r *http.Request) {

	fmt.Printf("%s\t %s\t %s%s\r\n", r.Method, r.Proto, r.Host, r.URL)

}

func main() {
	fmt.Println("Api rodando na porta 3000...")

	r := mux.NewRouter().StrictSlash(true)

	headers := handlers.AllowedHeaders([]string{"X-Request", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	r.HandleFunc("/devices", getDevices).Methods("GET")

	r.HandleFunc("/devices/{id}", getDevice).Methods("GET")

	r.HandleFunc("/devices", postDevice).Methods("POST")

	r.HandleFunc("/devices/{id}", putDevice).Methods("PUT")

	r.HandleFunc("/devices/{id}", deleteDevice).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(headers, methods, origins)(r)))

}
func setJsonHeader(w http.ResponseWriter) {

	w.Header().Set("Content-type", "application/json")

}

func getDevices(w http.ResponseWriter, r *http.Request) {

	setJsonHeader(w)

	HttpInfo(r)

	json.NewEncoder(w).Encode(devices)

}

func getDevice(w http.ResponseWriter, r *http.Request) {

	setJsonHeader(w)

	HttpInfo(r)

	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])

	for _, device := range devices {

		if device.Id == id {

			json.NewEncoder(w).Encode(device)
			return

		}

	}

	json.NewEncoder(w).Encode(NewResponse(r.Method, "failed", 400))
}

func postDevice(w http.ResponseWriter, r *http.Request) {

	setJsonHeader(w)

	HttpInfo(r)

	body, _ := ioutil.ReadAll(r.Body)

	var device Device

	err := json.Unmarshal(body, &device)

	if err != nil {

		json.NewEncoder(w).Encode(NewResponse(r.Method, "failed", 400))
		return

	}

	devices = append(devices, device)

	json.NewEncoder(w).Encode(NewResponse(r.Method, "success", 201))

}

func putDevice(w http.ResponseWriter, r *http.Request) {

	setJsonHeader(w)

	HttpInfo(r)

	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])

	body, _ := ioutil.ReadAll(r.Body)

	var device Device

	err := json.Unmarshal(body, &device)

	if err != nil {

		log.Fatal(err)

	}

	for index, _ := range devices {

		if devices[index].Id == id {

			devices[index] = device
			json.NewEncoder(w).Encode(NewResponse(r.Method, "success", 200))
			return

		}

	}

	json.NewEncoder(w).Encode(NewResponse(r.Method, "failed", 400))

}

func deleteDevice(w http.ResponseWriter, r *http.Request) {

	setJsonHeader(w)

	HttpInfo(r)

	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])

	for index, _ := range devices {

		if devices[index].Id == id {

			devices = append(devices[:index], devices[index+1:]...)
			json.NewEncoder(w).Encode(NewResponse(r.Method, "success", 200))
			return

		}

	}

	json.NewEncoder(w).Encode(NewResponse(r.Method, "failed", 400))

}
