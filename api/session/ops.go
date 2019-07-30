package session

import (
	"VideoWeb/api/dbops"
	"VideoWeb/api/defs"
	"VideoWeb/api/utils"
	"sync"
	"time"
)

// the map support Concurrent
var sessionMap *sync.Map

func nowInMilli() int64 {
	return time.Now().UnixNano()/1000000
}

func deleteExpiredSession(sid string)  {
	sessionMap.Delete(sid)
	dbops.DeleteSession(sid)
}

func init()  {
	sessionMap = &sync.Map{}
	
}

func LoadSessionsFromDB()  {

	r,err := dbops.RetrieveAllSessions()
	if err!=nil {
		return
	}

	//transfer return map to cache map
	r.Range(func(key, value interface{}) bool {
		sessionMap.Store(key,value)
		return true
	})
}

func GenerateNewSessionId(un string) string {

	id,_ := utils.NewUUID()
	//nano s to ms
	ct := nowInMilli()

	//valid time ms
	ttl := ct + 30*60*1000

	ss := &defs.SimpleSession{
		Username:un,
		TTL:ttl,
	}

	sessionMap.Store(id,ss)
	dbops.InserSession(id,ttl,un)

	return id
}

func IsSessionExpired(sid string) (string,bool) {

	ss,ok := sessionMap.Load(sid)
	if ok {
		ct := nowInMilli()
		if ss.(*defs.SimpleSession).TTL <ct {
			deleteExpiredSession(sid)
			return "Session Is Expired" ,true
		}
		return ss.(*defs.SimpleSession).Username,false
	}
	return "Load Session wrong",true
}


