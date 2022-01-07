package watermark

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func Usage() {
	msg := make([]string, 0)

	msg = append(msg, "----------------------------------------------------")
	msg = append(msg, "图片加水印（姜姜出品）")
	msg = append(msg, "  STEP1: 创建文件夹 workspace")
	msg = append(msg, "  STEP2: 将水印图片放到 workspace/mark.png")
	msg = append(msg, "  STEP3: 创建文件夹 workspace/mark")
	msg = append(msg, "  STEP4: 将所有需要加水印的图片放到 workspace/mark 下，可以带目录")
	msg = append(msg, "  STEP5: 将该程序放到 workspace 目录下，双击执行，加水印后的图片会生成到 workspace/mark_new 文件夹下")
	msg = append(msg, "  ")
	msg = append(msg, "  最终目录结构如下：")
	msg = append(msg, "  ")
	msg = append(msg, "  workspace")
	msg = append(msg, "  |-mark.png")
	msg = append(msg, "  |-mark.exe")
	msg = append(msg, "  |-mark")
	msg = append(msg, "  | |-dir1")
	msg = append(msg, "  |   |-a.jpg")
	msg = append(msg, "  |   |-b.jpeg")
	msg = append(msg, "  | |-dir2")
	msg = append(msg, "  |   |-c.png")
	msg = append(msg, "  |   |-d.PNG")
	msg = append(msg, "----------------------------------------------------")

	for _, m := range msg {
		fmt.Println(m)
	}
}

func DoMark(workspace string) {
	jobs, err := CreateJobs(workspace)
	if err != nil {
		fmt.Println(err)
		return
	}

	total := len(jobs)
	var complete int32

	var wg sync.WaitGroup
	jobChan := make(chan *Job)
	for i := 0; i < 10; i++ {
		go func() {
			for {
				job, ok := <-jobChan
				if !ok || job == nil {
					break
				}
				atomic.AddInt32(&complete, 1)
				fmt.Printf("[INFO ] 开始处理(%d/%d) : %s\n", complete, total, job.getDest())
				err := job.Do()
				wg.Done()
				if err != nil {
					fmt.Println("[ERROR] 发生错误:", err)
				}
			}
		}()
	}
	for _, job := range jobs {
		wg.Add(1)
		jobChan <- job
	}
	close(jobChan)
	wg.Wait()
}
