package main

import (
	"encoding/json"
	"fmt"
	//"github.com/gorilla/mux"
	"github.com/jrjsb4/tumblebus/client/db"
	"net/http"
)

type TumbleBusAPI struct {
	myconnection *db.MongoConnection
}

type ClientForm struct {
	FirstName   string `json:firstname`
	LastName    string `json:lastname`
	Address     string `json:address`
	City        string `json:city`
	State       string `json:state`
	ZipCode     string `json:zipcode`
	MobilePhone string `json:mobilephone`
	HomePhone   string `json:homephone`
	Email       string `json:email`
}

type ChildForm struct {
	Id         string `json:id`
	FirstName  string `json:firstname`
	LastName   string `json:lastname`
	BirthYear  int    `json:birthyear`
	BirthMonth int    `json:birthmonth`
}

type SchoolForm struct {
	Id          string `json:id`
	Name        string `json:name`
	Address     string `json:address`
	City        string `json:city`
	State       string `json:state`
	ZipCode     string `json:zipcode`
	MainPhone   string `json:mainphone`
	ContactName string `json:contactname`
	Url         string `json:url`
}

type APIResponse struct {
	StatusMessage string `json:statusmessage`
	StatusId      string `json:statusid`
}

func NewTumbleBusAPI() *TumbleBusAPI {
	TB := &TumbleBusAPI{
		myconnection: db.NewConnection(),
	}
	return TB
}

func (Tb *TumbleBusAPI) TumbleBusRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello and welcome to the TumbleBus API \n"+
		"Do a POST request with client information to add a client to the database\n"+
		"Do a POST request with school information to add a school to the database \n"+
		"Do a GET request to list school and client information\n")
}

// AddParent is a POST request API interface to add parent contact information to a Client collection in the database.
// If the Client collection does not exists in the database, a new Client collection is created to then allow
// the parent contact information to be added.
func (Tb *TumbleBusAPI) AddParent(w http.ResponseWriter, r *http.Request) {
	reqBodyStruct := new(ClientForm)
	fmt.Println("Add Parent")
	responseEncoder := json.NewEncoder(w)
	if err := json.NewDecoder(r.Body).Decode(&reqBodyStruct); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := responseEncoder.Encode(&APIResponse{StatusMessage: err.Error(), StatusId: "Jason"}); err != nil {
			fmt.Fprintf(w, "Error occured while processing post request %v \n", err.Error())
		}
		return
	}

	exist, clientId := Tb.myconnection.ClientExist(reqBodyStruct.FirstName, reqBodyStruct.LastName)
	if exist == false {
		var err error
		if clientId, err = Tb.myconnection.CreateClient(); err != nil {
			w.WriteHeader(http.StatusConflict)
			if err := responseEncoder.Encode(&APIResponse{StatusMessage: err.Error(), StatusId: ""}); err != nil {
				fmt.Fprintf(w, "Error %s occured while trying to create the client \n", err.Error())
			}
			return
		}
	}
	if err := Tb.myconnection.AddParent(clientId,
		reqBodyStruct.FirstName, reqBodyStruct.LastName,
		reqBodyStruct.Address, reqBodyStruct.City,
		reqBodyStruct.State, reqBodyStruct.ZipCode,
		reqBodyStruct.HomePhone, reqBodyStruct.MobilePhone,
		reqBodyStruct.Email); err != nil {
		w.WriteHeader(http.StatusConflict)
		if err := responseEncoder.Encode(&APIResponse{StatusMessage: err.Error(), StatusId: ""}); err != nil {
			fmt.Fprintf(w, "Error %s occured while trying to add the parent information \n", err.Error())
		}
		return
	}
	responseEncoder.Encode(&APIResponse{StatusMessage: "Ok", StatusId: clientId})
}

// AddParent is a POST request API interface to add parent contact information to a Client collection in the database.
// If the Client collection does not exists in the database, a new Client collection is created to then allow
// the parent contact information to be added.
func (Tb *TumbleBusAPI) AddChild(w http.ResponseWriter, r *http.Request) {
	reqBodyStruct := new(ClientForm)
	responseEncoder := json.NewEncoder(w)
	if err := json.NewDecoder(r.Body).Decode(&reqBodyStruct); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := responseEncoder.Encode(&APIResponse{StatusMessage: err.Error(), StatusId: ""}); err != nil {
			fmt.Fprintf(w, "Error occured while processing post request %v \n", err.Error())
		}
		return
	}

	exist, clientId := Tb.myconnection.ClientExist(reqBodyStruct.FirstName, reqBodyStruct.LastName)
	if exist == false {
		var err error
		if clientId, err = Tb.myconnection.CreateClient(); err != nil {
			w.WriteHeader(http.StatusConflict)
			if err := responseEncoder.Encode(&APIResponse{StatusMessage: err.Error(), StatusId: ""}); err != nil {
				fmt.Fprintf(w, "Error %s occured while trying to create the client \n", err.Error())
			}
			return
		}
	}
	if err := Tb.myconnection.AddParent(clientId,
		reqBodyStruct.FirstName, reqBodyStruct.LastName,
		reqBodyStruct.Address, reqBodyStruct.City,
		reqBodyStruct.State, reqBodyStruct.ZipCode,
		reqBodyStruct.HomePhone, reqBodyStruct.MobilePhone,
		reqBodyStruct.Email); err != nil {
		w.WriteHeader(http.StatusConflict)
		if err := responseEncoder.Encode(&APIResponse{StatusMessage: err.Error(), StatusId: ""}); err != nil {
			fmt.Fprintf(w, "Error %s occured while trying to add the parent information \n", err.Error())
		}
		return
	}
	responseEncoder.Encode(&APIResponse{StatusMessage: "Ok", StatusId: clientId})
}

/*
func (Ls *TumbleBusAPI) UrlShow(w http.ResponseWriter, r *http.Request) {
	//retrieve the variable from the request
	vars := mux.Vars(r)
	sUrl := vars["shorturl"]
	if len(sUrl) > 0 {
		//find long url that corresponds to the short url
		lUrl, err := Ls.myconnection.FindlongUrl(sUrl)
		if err != nil {
			fmt.Fprintf(w, "Could not find saved long url that corresponds to the short url %s \n", sUrl)
			return
		}
		//Ensure we are dealing with an absolute path
		http.Redirect(w, r, lUrl, http.StatusFound)
	}
}
*/
