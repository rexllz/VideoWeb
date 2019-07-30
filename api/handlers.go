package main

import (
	"VideoWeb/api/dbops"
	"VideoWeb/api/defs"
	"VideoWeb/api/session"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params)  {
	//create user is post method which have a body can read
	res,_ := ioutil.ReadAll(r.Body)
	ubody := &defs.UserCredential{}


	//request body is incorrect
	if err := json.Unmarshal(res,ubody); err!=nil {
		sendErrorResponse(w,defs.ErrorRequestBodyParseFailed)
		return
	}

	// add this user to DB
	if err:= dbops.AddUserCredential(ubody.Username, ubody.Pwd); err!=nil{
		sendErrorResponse(w,defs.ErrorDBError)
		return
	}

	//generate new session
	id := session.GenerateNewSessionId(ubody.Username)
	//create response body
	su := &defs.SignedUp{Success:true,SessionId:id}

	if resp,err := json.Marshal(su); err!=nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
		return
	}else {
		sendNormalResponse(w,string(resp),201)
	}
}

func Login(w http.ResponseWriter, r *http.Request, p httprouter.Params)  {

	uname := p.ByName("user_name")
	io.WriteString(w,uname)
}


