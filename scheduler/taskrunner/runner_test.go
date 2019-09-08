package taskrunner

import (
	"log"
	"testing"
	"time"
)

func TestRunner(t *testing.T)  {

	dispatcher := func(dc dataChan) error{
		for i := 0 ; i < 30 ; i++ {
			dc <- i
			log.Printf("Dispatcher send %d", i)
		}
		return nil
	}

	executor := func(dc dataChan) error{
		forloop:
			for{
				select {
				case d :=<- dc:
					log.Printf("Executor received %d", d)
				default:
					break forloop
				}
			}
		return nil
	}

	runner := NewRunner(30, false, dispatcher, executor)
	go runner.StartAll()
	time.Sleep(3 * time.Second)
}
