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

# 搭建项目

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



