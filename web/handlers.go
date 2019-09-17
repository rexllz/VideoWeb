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

func homeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {

	p := &HomePage{Name:"test"}
	t,e := template.ParseFiles("./web/template/home.html")
	if e != nil {
		log.Printf("Parsing template home.html err %s", e)
		return
	}
	//merge the p(name) and template together
	t.Execute(w, p)
}