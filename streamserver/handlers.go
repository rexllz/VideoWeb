package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func testPageHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params){

	t, _ := template.ParseFiles("streamserver/upload.html")
	t.Execute(w, nil)
}

func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params)  {

	vid := p.ByName("vid-id")
	vl := VIDEO_DIR + vid

	video, err := os.Open(vl)
	if err != nil{
		log.Printf("Error When Try to Open Video : %v", err)
		sendErrorResponse(w,http.StatusInternalServerError,err.Error())
		return
	}
	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, r, "", time.Now(), video)
	defer video.Close()
}

func uploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params)  {

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOADSIZE)
	if err := r.ParseMultipartForm(MAX_UPLOADSIZE); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "File too big")
		return
	}
	file, _ , err := r.FormFile("file")
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Intern Error")
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Read File Error : %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Intern Error")
	}

	fn := p.ByName("vid-id")
	err = ioutil.WriteFile(VIDEO_DIR + fn, data, 0666)
	if err != nil {
		log.Printf("Write File Error : %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Intern Error")
		return
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Upload Successfully")
}
