package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type customResponse struct {
	Status  string      `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type registration struct {
	Name     string `json:"name,omitempty"`
	Address  string `json:"address,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
}

type login struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

const (
	PORT = "8200"
)

func main() {
	err := dbConnect()
	if err != nil {
		fmt.Printf("Error in initializing dbclient and queries :[%s]\n Exiting..\n", err.Error())
		return
	}
	defer dbClient.Close()
	fmt.Println("Server Started ...")
	fmt.Printf("Listening on  PORT  %s\n", PORT)
	r := mux.NewRouter()
	r.HandleFunc("/login", loginHandler).Methods(http.MethodPost)
	r.HandleFunc("/registration", registrationHandler).Methods(http.MethodPost)
	r.HandleFunc("/forgot", forgotHandler).Methods(http.MethodPut)
	r.HandleFunc("/delete", deleteHandler).Methods(http.MethodDelete)
	http.Handle("/", r)
	err = http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		fmt.Printf("Error in starting server : %s", err.Error())
	}
}

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var response customResponse
		setupResponse(&w, r)
		regReq, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Error in reading body : %s", err.Error())
			response.Message = "Error in reading data"
			writeSuccessMessage(w, r, response)
			return
		}
		defer r.Body.Close()
		var reg registration
		err = json.Unmarshal(regReq, &reg)
		if err != nil {
			fmt.Printf("Error in unmarshalling body : %s", err.Error())
			response.Message = "Error in parsing JSON"
			writeSuccessMessage(w, r, response)
			return
		}
		rows, err := createStmt.Query(reg.Name, reg.Email, reg.Password, reg.Address)
		if err != nil {
			fmt.Printf("Error in creating user : [%s]\n", err.Error())
			response.Message = "Error creating user"
			writeSuccessMessage(w, r, response)
			return
		}
		defer rows.Close()
		response.Message = "Registration Successful"
		fmt.Printf("Registration Success for user %s\n", reg.Name)
		writeSuccessMessage(w, r, response)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var response customResponse
		setupResponse(&w, r)
		req, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Error in reading body : %s", err.Error())
			response.Message = "Error in reading data"
			writeSuccessMessage(w, r, response)
			return
		}
		defer r.Body.Close()
		var log login
		err = json.Unmarshal(req, &log)
		if err != nil {
			fmt.Printf("Error in unmarshalling body : %s", err.Error())
			response.Message = "Error in parsing JSON"
			writeSuccessMessage(w, r, response)
			return
		}
		rows := readStmt.QueryRow(log.Email, log.Password)
		var email, passowrd string
		err = rows.Scan(&email, &passowrd)
		if email == "" || passowrd == "" {
			fmt.Printf("Error in getting user : [%s]\n", err.Error())
			response.Message = "Error getting user"
			writeSuccessMessage(w, r, response)
		}
		response.Message = "User exists Successful"
		fmt.Printf("Login Success for user %s\n", log.Email)
		writeSuccessMessage(w, r, response)
	}
}

func forgotHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		var response customResponse
		setupResponse(&w, r)
		req, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Error in reading body : %s", err.Error())
			response.Message = "Error in reading data"
			writeSuccessMessage(w, r, response)
			return
		}
		defer r.Body.Close()
		var forg login
		err = json.Unmarshal(req, &forg)
		if err != nil {
			fmt.Printf("Error in unmarshalling body : %s", err.Error())
			response.Message = "Error in parsing JSON"
			writeSuccessMessage(w, r, response)
			return
		}
		var email string
		err = updateStmt.QueryRow(forg.Email, forg.Password).Scan(&email)
		if err != nil && email != "" {
			fmt.Printf("Error in creating user : [%s]\n", err.Error())
			response.Message = "Error creating user"
			writeSuccessMessage(w, r, response)
		}

		response.Message = "Registration Successful"
		fmt.Printf("Login Success for user %s\n	", forg.Email)
		writeSuccessMessage(w, r, response)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		var response customResponse
		setupResponse(&w, r)
		req, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Error in reading body : %s", err.Error())
			response.Message = "Error in reading data"
			writeSuccessMessage(w, r, response)
			return
		}
		defer r.Body.Close()
		var forg login
		err = json.Unmarshal(req, &forg)
		if err != nil {
			fmt.Printf("Error in unmarshalling body : %s", err.Error())
			response.Message = "Error in parsing JSON"
			writeSuccessMessage(w, r, response)
			return
		}
		_, err = deleteStmt.Exec(forg.Email)
		if err != nil {
			fmt.Printf("Error in deleting user : [%s]\n", err.Error())
			response.Message = "Error deleting user"
			writeSuccessMessage(w, r, response)
		}
		response.Message = "Deleting Successful"
		fmt.Printf("Login Success for user %s\n	", forg.Email)
		writeSuccessMessage(w, r, response)
	}
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("content-type", "application/json")
}
func writeSuccessMessage(w http.ResponseWriter, r *http.Request, data interface{}) {
	fmt.Printf(
		"%s %s \n",
		r.Method,
		r.RequestURI,
	)
	w.WriteHeader(http.StatusOK)
	body, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write(body)
}
