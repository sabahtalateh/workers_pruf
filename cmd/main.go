package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	WorkerStatusFree    = 0
	WorkerStatusWorking = 1
)

type WorkersList struct {
	mu   sync.Mutex
	list map[int64]int64
	work func()
}

func initWorkersList(work func()) WorkersList {
	list := make(map[int64]int64)
	for i := 0; i < 2; i++ {
		list[int64(i)] = WorkerStatusFree
	}

	return WorkersList{
		list: list,
		work: work,
	}
}

func (wl *WorkersList) checkWorkers() {
	wl.mu.Lock()
	defer wl.mu.Unlock()

	for id, status := range wl.list {
		if status == WorkerStatusFree {
			fmt.Printf("-> worker [%d] is free\n", id)
			go work(id, wl)
			wl.list[id] = WorkerStatusWorking
			fmt.Printf("-> worker [%d] got work\n", id)
		}
	}
}

func (wl *WorkersList) free(workerId int64) {
	wl.mu.Lock()
	defer wl.mu.Unlock()

	wl.list[workerId] = WorkerStatusFree
	fmt.Printf("-> worker [%d] free and can work again\n", workerId)
}

func work(workerId int64, wl *WorkersList) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("(!) worker [%d] panicked with message [%s]\n", workerId, r)
			wl.free(workerId)
		}
	}()

	count := int64(0)
	for {
		if count-1 == workerId {
			panic(fmt.Sprintf("(x) AAA!! Mr. [%d] is dead", workerId))
		}

		fmt.Printf("(.) worker [%d] working its work\n", workerId)
		wl.work()

		count++
	}
}

func main() {
	wl := initWorkersList(func() { fmt.Printf("SOME WORK!!!"); time.Sleep(20 * time.Second) })

	for {
		println("-> check workers")
		wl.checkWorkers()
		time.Sleep(15 * time.Second)
	}
}
