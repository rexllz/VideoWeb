package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type HomePage struct{
	Name string
}
type UserPage struct{
	Name string
}

func homeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {

	cname, err1 := r.Cookie("username")
	sid, err2 := r.Cookie("session")
	if err1 != nil || err2 != nil{
		p := &HomePage{Name:"test"}
		t,e := template.ParseFiles("./web/template/home.html")
		if e != nil {
			log.Printf("Parsing template home.html err %s", e)
			return
		}
		//merge the p(name) and template together
		t.Execute(w, p)
		return
	}

	if len(cname.Value) != 0 && len(sid.Value) != 0 {
		http.Redirect(w,r,"/userhome",http.StatusFound)
		return
	}
}

func userHomeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {

	cname, err1 := r.Cookie("username")
	_,err2 := r.Cookie("session")

	if err1 != nil || err2 != nil{
		http.Redirect(w,r,"/",http.StatusFound)
		return
	}
	fname := r.FormValue("username")
	
	//get the user info and refresh the page
	//if the cookie has no info, go to the form value find it
	var p *UserPage
	if len(cname.Value) != 0 {
		p = &UserPage{Name:cname.Value}
	}else if len(fname) != 0 {
		p = &UserPage{Name:fname}
	}

	t,e := template.ParseFiles("./web/template/userhome.html")
	if e != nil {
		log.Printf("Parsing userhome.html error %s",e)
		return
	}

	t.Execute(w,p)
}

// we need this handler transfer the request to another server to avoid CORS problem
func apiHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {

	//before we handle the request , we must check it first
	if r.Method != http.MethodPost{
		re,_ := json.Marshal(ErrorRequestNotRecognize)
		io.WriteString(w,string(re))
		return
	}

	res,_ := ioutil.ReadAll(r.Body)
	apibody := &ApiBody{}
	if err := json.Unmarshal(res,apibody); err != nil{
		re, _ := json.Marshal(ErrorBodyParseFailed)
		io.WriteString(w,string(re))
	}

	request(apibody,w,r)
	defer r.Body.Close()
}

func proxyHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {
	// transfer the request to the following address
	u,_ = url.Parse("http://127.0.0.1:9000/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w,r)
}

























