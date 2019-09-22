# VideoWeb

Stream Video Web in golang

# System Design

![g0](https://raw.githubusercontent.com/rexllz/VideoWeb/master/img/g0.jpg)

## RESTful API

* based on http request / json data format
* Unified interface
* no status
* cacheable
* API document wit 
* crud to resource with 4 different method(get,post,delete,put)
* return code to describe the resource status

![g1](https://raw.githubusercontent.com/rexllz/VideoWeb/master/img/g1.jpg)

# Project Basic

## Httprouter
github.com/julienschmidt/httprouter

```go
func RegisterHandlers() *httprouter.Router {

	router := httprouter.New()
	router.POST("/user",CreateUser)
	return router
}

func CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params)  {

	io.WriteString(w,"create user!")
}

func main()  {

	r := RegisterHandlers()
	http.ListenAndServe(":8000",r)
}
```
## Save Login Status by Session
create sessions table in DB to save the user login status
{sessionID, loginName}
also need set TTL

because the RESTful API is stateless, we use session to save the user's status

Session saved in server , and the sessionId saved in client(cookie), used to check status

![g4](https://raw.githubusercontent.com/rexllz/VideoWeb/master/img/g4.jpg)

**Session Based Authentication**
In the session based authentication, the server will create a session for the user after the user logs in. The session id is then stored on a cookie on the user’s browser. While the user stays logged in, the cookie would be sent along with every subsequent request. The server can then compare the session id stored on the cookie against the session information stored in the memory to verify user’s identity and sends response with the corresponding state!

![g5](https://raw.githubusercontent.com/rexllz/VideoWeb/master/img/g5.png)

**Token Based Authentication**
Many web applications use JSON Web Token (JWT) instead of sessions for authentication. In the token based application, the server creates JWT with a secret and sends the JWT to the client. The client stores the JWT (usually in local storage) and includes JWT in the header with every request. The server would then validate the JWT with every request from the client and sends response.

![g6](https://raw.githubusercontent.com/rexllz/VideoWeb/master/img/g6.png)

## Main DB tables

![g2](https://raw.githubusercontent.com/rexllz/VideoWeb/master/img/g2.jpg)

create connect in init function to save resource

```go
package dbops

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	dbConn *sql.DB
	err error
)

func init()  {
	dbConn, err = sql.Open("mysql",
		"root:root@tcp(127.0.0.1:3306)/video_web?charset=utf8")
	if err!=nil {
		panic(err.Error())
	}
}
```

avoid write sql with '+'
use dbConn.prepare("sql")

```go
func AddUserCredential(loginName string, pwd string) error {

	stmtIns, err := dbConn.Prepare("INSERT INTO users (login_name, pwd) VALUES (?,?)")
	if err!=nil {
		return err
	}

	stmtIns.Exec(loginName,pwd)
	stmtIns.Close()
	return nil
}

func GetUserCredential(loginName string) (string, error) {

	stmtOut, err := dbConn.Prepare("SELECT pwd FROM users WHERE login_name = ?")
	if err!=nil {
		log.Println("%s",err)
		return "",err
	}

	var pwd string
	stmtOut.QueryRow(loginName).Scan(&pwd)
	stmtOut.Close()
	return pwd,nil
}

func DeleteUser(loginName string,pwd string) error {

	stmtDel, err := dbConn.Prepare("DELETE FROM users WHERE login_name = ? AND pwd = ?")
	if err!=nil {
		log.Println("delete user error: %s", err)
		return err
	}

	stmtDel.Exec(loginName,pwd)
	stmtDel.Close()
	return nil
}
```

## GoTest Format

* init (dblogin, truncate tables)
* test
* clear data(truncate tables)

init -> TestMain

go test function , improve efficient to develop 

can do some init things and test in order

```go
package dbops

import(
	"testing"
	)

func clearTables()  {

	dbConn.Exec("TRUNCATE users")
	dbConn.Exec("TRUNCATE video_info")
	dbConn.Exec("TRUNCATE comments")
	dbConn.Exec("TRUNCATE sessions")
}

func TestMain(m *testing.M)  {

	clearTables()
	m.Run()
	clearTables()
}

func TestUserWorkFlow(t *testing.T) {
	t.Run("add",testAddUserCredential)
	t.Run("get",testGetUserCredential)
	t.Run("delete",testDeleteUser)
	t.Run("reGet",testRegetUser)
}

func testAddUserCredential(t *testing.T) {
	err := AddUserCredential("rex","123")
	if err!=nil {
		t.Errorf("Error of Add User : %v", err)
	}
}

func testGetUserCredential(t *testing.T) {
	pwd,err := GetUserCredential("rex")
	if err!=nil {
		t.Errorf("Error of Get User : %v", err)
	}
	t.Logf("User's pwd: %v",pwd)
}

func testDeleteUser(t *testing.T) {
	err := DeleteUser("rex", "123")
	if err!=nil {
		t.Errorf("Error of Delete User : %v", err)
	}
}

func testRegetUser(t *testing.T)  {
	pwd,err := GetUserCredential("rex")
	if err!=nil {
		t.Errorf("Error of reGet User : %v", err)
	}
	if pwd!= "" {
		t.Error("Error of reGet User : reGet pwd is not null")
	}
}
```
![g3](https://raw.githubusercontent.com/rexllz/VideoWeb/master/img/g3.jpg)

## Http Header Response Handler

add a middle ware handler to check the session

```go
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
```

Eg. in create user handler , 
response the different status row by row  to make sure clinet can get the right information

```go
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
```


# Stream

## Limiter

we need to limit the traffic to avoid some problem
use bucket-token method to solve this

golang can use channel sharing instead of memory sharing

```go
import "log"

//bucket
type ConnLimiter struct {

	concurrentConn int
	bucket chan int

}

//constructor func
func NewConnLimiter(cc int) *ConnLimiter {

	return &ConnLimiter{
		concurrentConn:cc,
		bucket:make(chan int, cc),
	}
}

func (cl *ConnLimiter) GetConn() bool {

	if len(cl.bucket) >= cl.concurrentConn {
		log.Printf("Reach the limitation")
		return false
	}
	cl.bucket <- 1
	return true
}

func (cl *ConnLimiter) ReleaseConn() {

	c :=<- cl.bucket
	log.Printf("New connection coming: %d", c)
}
```
Then
implement the traffic limit in main.go

```go
package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type middleWareHandler struct {
	r *httprouter.Router
	l *ConnLimiter
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !m.l.GetConn() {
		sendErrorResponse(w, http.StatusTooManyRequests, "Too Many Request")
		return
	}
	m.r.ServeHTTP(w,r)
	defer m.l.ReleaseConn()
}

func NewMiddleWareHandler(r *httprouter.Router, cc int) http.Handler {
	m := middleWareHandler{}
	m.r = r
	m.l = NewConnLimiter(cc)
	return m
}

func RegisterHandlers() *httprouter.Router{
	router := httprouter.New()

	router.GET("/videos/:vid-id", streamHandler)
	router.POST("/upload/:vid-id",uploadHandler)

	return router
}

func main(){
	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r, 2)
	http.ListenAndServe(":9000", mh)
}
```

## Play a Video

```go
func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params)  {

	vid := p.ByName("vid-id")
	vl := VIDEO_DIR + vid

	video, err := os.Open(vl)
	if err != nil{
		sendErrorResponse(w,http.StatusInternalServerError,err.Error())
		return
	}
	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, r, "", time.Now(), video)
	defer video.Close()
}
```

## Upload Video

Use html/ template to parse a upload page(html), and upload a video 
```go
func testPageHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params){

	t, _ := template.ParseFiles("streamserver/upload.html")
	t.Execute(w, nil)
}
```
# Scheduler

Using scheduler to finish some Asynchronous tasks

**Scheduler Model**

![g7](https://raw.githubusercontent.com/rexllz/VideoWeb/master/img/g7.jpg)

We use a table to save the videos we want to delete, and wait for delete by batch

```go
func DelVideoDeletionRecord(vid string) error {

	stmtDel, err := dbConn.Prepare("DELETE FROM video_del_rec WHERE video_id = ?")
	if err != nil {
		return err
	}

	_, err = stmtDel.Exec(vid)
	if err != nil {
		log.Printf("Deleting VideoDeletionRecord Error : %v", err)
		return err
	}

	defer stmtDel.Close()
	return nil

}
```

We delete the video asynchronously (using a buffer), and by batch

```go
func VideoClearDispatcher(dc dataChan) error {

	res, err := dbops.ReadVideoDeletionRecord(3)
	if err != nil{
		log.Printf("Video Clear Dispatcher Error : %v", err)
		return err
	}

	if len(res) == 0 {
		return errors.New("all task finished")
	}

	for _,id := range res {
		dc <- id
	}

	return nil
}

func VideoClearExecutor(dc dataChan) error {

	errMap := &sync.Map{}
	var err error

	forloop:
		for {
			select {
			case vid :=<- dc :
				go func(id interface{}) {
					if err := deleteVideo(id.(string)); err != nil{
						errMap.Store(id, err)
						return
					}
					if err := dbops.DelVideoDeletionRecord(id.(string)); err != nil{
						errMap.Store(id, err)
						return
					}
				}(vid)
			default:
				break forloop
			}
		}

	errMap.Range(func(key, value interface{}) bool {
		err = value.(error)
		if err !=nil {
			return false
		}
		return true
	})

	return err
}
```

# Timer

use timer to start the work :

```go
package taskrunner

import "time"

type Worker struct {
	ticker *time.Ticker
	runner *Runner
}

func NewWorker(interval time.Duration, r *Runner) *Worker {

	return &Worker{
		ticker: time.NewTicker(interval * time.Second),
		runner: r,
	}
}

func (w *Worker) startWorker() {

	//not good , this (range) is a synchronization method
	//for c = range w.ticker.C {
	//
	//}

	for  {
		select {
		// get ticker 's channel signal
		case <- w.ticker.C:
			go w.runner.StartAll()
		}
	}
}

func Start()  {

	r := NewRunner(3, true, VideoClearDispatcher, VideoClearDispatcher)
	w := NewWorker(3, r)
	go w.startWorker()

}
```

The details of the process :

step 1 :
user -> api service -> delete video

step 2 :
api service -> scheduler -> write video deletion record 

step 3 :
create a timer 

step 4 :
timer -> runner -> read deletion record -> exec -> delete video file 

When we start a goroutine,  we can use this method to create a block situation

```go
c := make(chan int)
	//something
<- c
```
or

```go
for{
    //something
}
```

# FontEnd

template engine can transfer the elements to the final page
Golang has 2 :
html/template   and   text/template    (dynamic generate)

So, how to use this template?
the following code merge the html and dynamic elements together 
```html
<div class="topnav">
    <a class="active" href="#home">Home</a>
    <a href="#news">{{.Name}}</a>
    <a href="#about">Help</a>
```
```go
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
```

We can get some info from request cookie:

if we can not get the correct info from cookie , we need goto the login page, and if we have the login info, the user should be lead to the userhome page (user redirect function)

```go
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
```

the same things for the situation where user request the /userhome, maybe we need to lead the user to the home page

```go
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
```

