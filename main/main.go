package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
)

type userData struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	GrantType string `json:"grant_type"`
}

// MsgDataForm data form
type MsgDataForm struct {
	// {"clientMsgId": ""}
	ClientMsgID string `json:"clientMsgId"`
}

// User struct
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type data struct {
	code        string
	title       string
	description string
}

var testData []data

// defaultTime := time.Now()

func main() {
	router := mux.NewRouter()
	// router.HandleFunc("/test/", multiPartData).Methods("POST")
	router.HandleFunc("/oauth/token", tokenData).Methods("POST")
	router.HandleFunc("/v1/messages", getMessage).Methods("POST")
	router.HandleFunc("/v1/users/me/files", multiPartData).Methods("POST")

	// push
	router.HandleFunc("/push", multiPartData).Methods("POST")

	router.HandleFunc("/test/{code}", GetData).Methods("GET")
	router.HandleFunc("/test/{category}/{id:[0-9]}",
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			w.WriteHeader(http.StatusOK)
			fmt.Println(w, "Category : %v", vars["category"])

		}).Methods("POST")

	router.HandleFunc("/", HomeHandler)
	http.Handle("/", router)

	http.ListenAndServe(":8082", router)

	testData = append(testData, data{code: "1", title: "first title", description: "test des"})
	testData = append(testData, data{code: "2", title: "second title", description: "test des2"})

}

func getMessage(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("getMessage")
	// params := mux.Vars(r)

	// len := r.ContentLength
	// body := make([]byte, len)
	// r.Body.Read(body)
	// fmt.Println("BODY : ", string(body))
	// fmt.Println(w, string(body))

	var msgData MsgDataForm

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&msgData); err != nil {
		fmt.Println("err====================")
		panic(err)
		// respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		// return
	}

	defer r.Body.Close()

	fmt.Println("clientMsgId : ", msgData.ClientMsgID)

	// c := make(chan string)

	// c <- msgData.ClientMsgID

	// jsonBody := simplejson.New()
	jsonData := simplejson.New()

	jsonData.Set("code", "2000")
	jsonData.Set("message", "Success")

	// jsonBody.Set("message", "success")
	// jsonBody.Set("data", jsonData)

	// jsonBody.Set("push_uri", "push_uri_value")
	// jsonBody.Set("jti", "jti_value")

	payload, err := jsonData.MarshalJSON()
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)

	ch := make(chan string, 1)
	ch <- msgData.ClientMsgID

	go requestProc(ch)

}

// HomeHandler is Test Controller
func HomeHandler(writer http.ResponseWriter, request *http.Request) {

}

func requestProc(c <-chan string) {

	fmt.Println("requestProc")

	// now := time.Now()

	clientMsgID := <-c
	fmt.Println("clientMsgID", clientMsgID)

	jsonBody := simplejson.New()
	// jsonAry := simplejson.New()
	// jsonBodyArray := simplejson.New()[]

	jsonBody.Set("result", "2000")
	jsonBody.Set("message", "2000")
	jsonBody.Set("clientMsgId", clientMsgID)
	jsonBody.Set("telecom", "kt")
	jsonBody.Set("date", "2020 08 13 17:53:21.455")

	// fmt.Println("jsonBody : ", jsonBody)

	jsonBody2 := simplejson.New()
	// jsonBodyArray := simplejson.New()[]

	jsonBody2.Set("result", "2000")
	jsonBody2.Set("message", "2000")
	jsonBody2.Set("clientMsgId", clientMsgID)
	jsonBody2.Set("telecom", "kt")
	jsonBody2.Set("date", "2020 08 13 17:53:21.455")

	// fmt.Println("jsonBody2 : ", jsonBody)

	// var strArray []string

	var jsonArray []*simplejson.Json

	jsonArray = append(jsonArray, jsonBody)
	// jsonArray = append(jsonArray, jsonBody2)

	// fmt.Println(strArray)

	jsonResult := simplejson.New()

	jsonResult.Set("results", jsonArray)

	fmt.Println(jsonResult)

	// payload, err := jsonResult.MarshalJSON()
	// if err != nil {
	// panic(err)
	// }
	// fmt.Println(payload)

	postPush(jsonResult)
}

// {
// 	"
// 	"result":"
// 	"message":"
// 	"clientMsgId":"
// 	"telecom":"kt|skt|
// 	"date":"2020 08 13 17:53:21.455"
// 	}

func postPush(json *simplejson.Json) {

	// json := simplejson.New()
	// fmt.Println(j)
	payload, err := json.MarshalJSON()

	if err != nil {
		panic(err)
	}
	buff := bytes.NewBuffer(payload)

	resp, err := http.Post("http://localhost:3003/push", "application/json", buff)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	} else {
		str := string(respBody)
		println(str)
	}

}

// GetData is get data
func GetData(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("get!!")
	p := mux.Vars(request)
	writer.WriteHeader(http.StatusOK)
	for _, i := range testData {
		if i.code == p["code"] {
			fmt.Println(i.code)
			json.NewEncoder(writer).Encode(i)
			return
		}
	}
	//json.NewEncoder(writer).Encode(&evnet{})
}

func multiPartData(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("in Multipart api")

	file, handler, err := request.FormFile("file")
	fileName := request.FormValue("fileName")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	fmt.Println("fileName: ", fileName)
	// _, _ = io.WriteString(writer, "File"+fileName+"Upload successfully")
	// _, _ = io.Copy(f, file)

	json := simplejson.New()
	jsonData := simplejson.New()

	jsonData.Set("fileId", "FileId001")
	jsonData.Set("expiryDate", "FileId001")
	jsonData.Set("expiryDate", "2020 09 01 00:00:00:0000")

	json.Set("message", "success")
	json.Set("data", jsonData)

	payload, err := json.MarshalJSON()
	if err != nil {
		fmt.Println(err)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(payload)

}

func tokenData(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("token Data")

	p := mux.Vars(request)

	var u userData

	fmt.Println("1====================================================")
	userName := request.FormValue("username")
	fmt.Println("username: ", userName)
	fmt.Println("request: ", request)
	fmt.Println("body: ", request.Body)

	fmt.Println("2====================================================")
	request.ParseForm()
	fmt.Println("request: ", request)
	fmt.Println("body: ", request.Body)

	// username := p["username"]
	// passowrd := p["passowrd"]
	// granType := p["grant_type"]
	granType := p["Authorization"]

	fmt.Println("auth: ", granType)
	// fmt.Println("auth: ", json.NewDecoder(request.Body))

	// getBody := request.GetBody

	// request.ParseForm()
	// fmt.Println(request)
	// err := json.NewDecoder(request.Body).Decode(&u)
	// if err != nil {
	// fmt.Println("err: ", err)
	// panic(err)
	// }

	// fmt.Println("body : ", getBody)
	fmt.Println("body : ", request.Body)

	authKey := request.Header.Get("Authorization")

	fmt.Println("Header : ", authKey, ", username : ", u.Username, ", password : ", u.Password, ", grant_type : ", u.GrantType)

	decoder := json.NewDecoder(request.Body)
	// err := decoder.Decode(&Data)
	fmt.Println(decoder)

	json := simplejson.New()
	json.Set("access_token", "testToken11111")
	json.Set("token_type", "bearer")
	json.Set("refresh_token", "refresh_token_testToken11111")
	json.Set("expires_in", "expires_in_value")
	json.Set("scope", "scope_value")
	json.Set("ip", "ip_value")
	json.Set("push_uri", "push_uri_value")
	json.Set("jti", "jti_value")

	payload, err := json.MarshalJSON()
	if err != nil {
		fmt.Println(err)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(payload)

}
