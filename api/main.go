package main

import (
	"VideoWeb/api/session"
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

	router.GET("/user/:username", GetUserInfo)

	router.POST("/user/:username/videos", AddNewVideo)

	router.GET("/user/:username/videos", ListAllVideos)

	router.DELETE("/user/:username/videos/:vid-id", DeleteVideo)

	router.POST("/videos/:vid-id/comments", PostComment)

	router.GET("/videos/:vid-id/comments", ShowComments)

	return router
}

//load the sessions info from db to the session map
//this is the first thing before we can check the user status
func Prepare() {
	session.LoadSessionsFromDB()
}

func main()  {

	r := RegisterHandlers()
	mh := NewMiddleWareHanler(r)
	http.ListenAndServe(":8000",mh)
}