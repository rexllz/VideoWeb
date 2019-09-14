package taskrunner

import (
	"VideoWeb/scheduler/dbops"
	"errors"
	"log"
	"os"
	"sync"
)

func deleteVideo(vid string) error {

	println("delete file!")
	err := os.Remove(VIDEO_PATH + vid)
	log.Print("delete ok")
	if os.IsNotExist(err) {
		log.Printf("file not exist")
	}
	if err != nil && !os.IsNotExist(err){
		log.Printf("Deleting Video File Error : %v", err)
		return err
	}
	log.Printf("delete finish")

	return nil
}

func VideoClearDispatcher(dc dataChan) error {

	println("read delete record")
	res, err := dbops.ReadVideoDeletionRecord(3)
	if err != nil{
		log.Printf("Video Clear Dispatcher Error : %v", err)
		return err
	}

	if len(res) == 0 {
		println("no record")
		return errors.New("all task finished")

	}
	println("get delete record")
	for _,id := range res {
		dc <- id
		log.Printf("send delete signal : %v to dc",id)
	}

	return nil
}

func VideoClearExecutor(dc dataChan) error {

	println("executor")
	errMap := &sync.Map{}
	var err error

	forloop:
		for {
			select {
			case vid :=<- dc :
				println("get delete signal")
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
				println("no signal to executor")
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
	println(err)
	return err
}


