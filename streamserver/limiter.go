package main

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


