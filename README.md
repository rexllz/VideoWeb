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







