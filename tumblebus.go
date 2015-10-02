package main

import "net/http"

/*
	This is the entry point for the RESTful API.
	The program makes use of the gorilla mux library for routing as well as
	the mgo library to interface with mongo database backend.
*/

func main() {
	//Create a new API shortner API
	TumbleBus := NewTumbleBusAPI()
	//Create the needed routes for the API
	routes := CreateRoutes(TumbleBus)
	//Initiate the API routers
	router := NewTumbleBusRouter(routes)
	//This will start the web server on local port 5100
	http.ListenAndServe(":5100", router)
}
