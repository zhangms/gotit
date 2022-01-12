package parallel

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type Task interface {
	Total() int
	Do(index int) error
	Info(index int) string
	Summary() string
}

func Do(task Task, routine int) {
	begin := time.Now()
	routine = int(math.Max(float64(routine), 1))
	var wait sync.WaitGroup
	c := make(chan int)
	total := task.Total()
	var completeCount, successCount, errCount int32
	for i := 0; i < routine; i++ {
		go func() {
			for {
				job, ok := <-c
				if !ok {
					break
				}
				complete := atomic.AddInt32(&completeCount, 1)
				fmt.Printf("[INFO ] 开始处理(%d/%d) : %s\n", complete, total, task.Info(job))
				err := task.Do(job)
				if err != nil {
					atomic.AddInt32(&errCount, 1)
					fmt.Println("[ERROR] ", err)
				} else {
					atomic.AddInt32(&successCount, 1)
				}
				wait.Done()
			}
		}()
	}
	for i := 0; i < total; i++ {
		wait.Add(1)
		c <- i
	}
	close(c)
	wait.Wait()
	elapsed := time.Since(begin) / time.Second
	fmt.Printf("处理完成，耗时 %d 秒，共 %d 个任务，成功 %d 个，失败 %d 个, %s\n", elapsed, total, successCount, errCount, task.Summary())
}
