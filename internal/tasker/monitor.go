package tasks

import (
	"fmt"
	"runtime"
	"time"
)

func (t *Tasker) Monitor() {
	for {
		time.Sleep(time.Second * 10)
		fmt.Println("Task is working! Goroutine num:", runtime.NumGoroutine())
	}
}
