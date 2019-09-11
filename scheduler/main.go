package main

import (
	"VideoWeb/scheduler/taskrunner"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func RegisterHandlers() *httprouter.Router {

	router := httprouter.New()

	router.GET("/video-delete-record/:vid-id", vidDelRecHandler)

	return router
}

func main()  {

	go taskrunner.Start()
	r := RegisterHandlers()
	http.ListenAndServe(":9001", r)

	// we can use this method to create a block situation
	//c := make(chan int)
	//something
	//<- c
}