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











