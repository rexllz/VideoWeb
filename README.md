# VideoWeb

Stream Video Web in golang

# 整体设计

![g0](https://raw.githubusercontent.com/rexllz/VideoWeb/master/img/g0.jpg)

## RESTful API

* 以http为通信协议，json为数据格式
* 统一接口
* 无状态
* 可缓存
* 以URL设计API
* 通过四个不同的method（get、put、post、delete）来区分对资源的crud
* 返回码返回对资源的描述

![g1](https://raw.githubusercontent.com/rexllz/VideoWeb/master/img/g1.jpg)

# 项目基础

## httprouter
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
## session保存登陆状态
create sessions table in DB to save the user login status
{sessionID, loginName}
并且设置TTL，定期过期

## 主要DB Tables

![g2](https://raw.githubusercontent.com/rexllz/VideoWeb/master/img/g2.jpg)

在init中创建全局conn，节省开销

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

避免用加号操作sql，易被攻击
使用dbConn.prepare("sql")

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

## GoTest规范

* init (dblogin, truncate tables)
* test
* clear data(truncate tables)

init -> TestMain

go test 功能，提高了开发和测试的效率。 
有时会遇到这样的场景： 
进行测试之前需要初始化操作(例如打开连接)，测试结束后，需要做清理工作(例如关闭连接)等等。这个时候就可以使用TestMain()

测试从TestMain进入，依次执行测试用例，最后从TestMain退出。
保证既定顺序

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







