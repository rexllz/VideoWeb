package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
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