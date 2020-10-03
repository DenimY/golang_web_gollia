package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type data struct {
	code        string
	title       string
	description string
}

var testData []data

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/test/{code}", GetData).Methods("GET")
	router.HandleFunc("/test/{category}/{id:[0-9]}",
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			w.WriteHeader(http.StatusOK)
			fmt.Println(w, "Category : %v", vars["category"])

		}).Methods("POST")

	router.HandleFunc("/", HomeHandler)
	http.Handle("/", router)
	http.ListenAndServe(":8080", router)

	testData = append(testData, data{code: "1", title: "first title", description: "test des"})
	testData = append(testData, data{code: "2", title: "second title", description: "test des2"})

}

func HomeHandler(writer http.ResponseWriter, request *http.Request) {

}

func GetData(writer http.ResponseWriter, request *http.Request) {
	p := mux.Vars(request)
	writer.WriteHeader(http.StatusOK)
	for _, i := range testData {
		if i.code == p["code"] {
			json.NewEncoder(writer).Encode(i)
			return
		}
	}
	//json.NewEncoder(writer).Encode(&evnet{})

}
