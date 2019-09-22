package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var httpClient *http.Client

func init()  {
	httpClient = &http.Client{}
}

func request(b *ApiBody, w http.ResponseWriter, r *http.Request)  {

	var resp *http.Response
	var err error

	//the body has three methods
	switch b.Method {
	case http.MethodGet:
		//first , prepare the request
		//second , send the request (by client.do) to the api server
		req, _ := http.NewRequest("GET",b.Url,nil)
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil{
			log.Printf(err.Error())
			return
		}
		normalResponse(w,resp)
	case http.MethodPost:
		req, _ := http.NewRequest("POST",b.Url,bytes.NewBuffer([]byte(b.ReqBody)))
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil{
			log.Printf(err.Error())
			return
		}
		normalResponse(w,resp)
	case http.MethodDelete:
		req, _ := http.NewRequest("DELETE",b.Url,nil)
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil{
			log.Printf(err.Error())
			return
		}
		normalResponse(w,resp)
	default:
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w,"Bad api request")
		return
	}
}

// this function transfer the response from api server to the local system response
func normalResponse(w http.ResponseWriter, r *http.Response)  {
	//check response first
	res, err := ioutil.ReadAll(r.Body)
	if err != nil{
		re, _ := json.Marshal(ErrorInternalFaults)
		w.WriteHeader(500)
		io.WriteString(w, string(re))
	}
	w.WriteHeader(r.StatusCode)
	io.WriteString(w,string(res))
}