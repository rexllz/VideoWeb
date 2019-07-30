package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type middleWareHandler struct {
	r *httprouter.Router
}


func NewMiddleWareHanler(r *httprouter.Router) http.Handler {
	m := middleWareHandler{}
	m.r = r
	return m
}

//implement the interface
func (m middleWareHandler) ServeHTTP(w http.ResponseWriter,r *http.Request) {
	//check session
	validateUserSession(r)

	m.r.ServeHTTP(w,r)
}

func RegisterHandlers() *httprouter.Router {

	router := httprouter.New()
	router.POST("/user", CreateUser)
	router.POST("/user/:user_name", Login)

	return router
}


func main()  {

	r := RegisterHandlers()
	mh := NewMiddleWareHanler(r)
	http.ListenAndServe(":8000",mh)
}