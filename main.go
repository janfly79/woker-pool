package main

import (
	. "scraptch/logic"
	"time"
)

func main(){
	dispatcher := NewDispatcher(2)
	dispatcher.Run()

	go func(){
		for {
			time.Sleep( 1 * time.Second)
			NewJob()
		}
	}()

	time.Sleep( 86400 * time.Second)
}
