package compress

import (
	"fmt"
	"gotit/parallel"
	"path/filepath"
	"time"
)

func Do(args []string) {
	err := doCompress(filepath.Dir(args[0]), 16)
	if err != nil {
		fmt.Println(err)
		time.Sleep(2 * time.Second)
		usage()
	}
}

func usage() {
	msg := make([]string, 0)
	msg = append(msg, "STEP1: 创建文件夹 workspace")
	msg = append(msg, "STEP2: 打开网址：https://tinify.cn/developers 获取API key ")
	msg = append(msg, "STEP3: 在workspace下新建文本文档 apikey.txt，将 STEP2 获取到的API key 按行写入 apikey.txt")
	msg = append(msg, "STEP4: 在workspace下新建文件夹 compress, 将需要压缩的图片或文件夹放入其中")
	msg = append(msg, "STEP5: 将本程序放入 workspace，双击运行，压缩后的文件夹将放入 compress_new 内")
	msg = append(msg, "NOTE:  一个 tinypng 账号默认只能一个月免费压缩500张图片，超过500张可以付费或者换个邮箱获取apikey")
	msg = append(msg, "")
	msg = append(msg, "最终目录结构如下：")
	msg = append(msg, "")
	msg = append(msg, "workspace")
	msg = append(msg, "|-apikey.txt")
	msg = append(msg, "|-gotit.exe")
	msg = append(msg, "|-compress")
	msg = append(msg, "  |-1.jpg")
	msg = append(msg, "  |-2.png")
	msg = append(msg, "  |-dir1")
	msg = append(msg, "    |-a.jpg")
	msg = append(msg, "    |-b.jpeg")
	msg = append(msg, "  |-dir2")
	msg = append(msg, "    |-c.png")
	msg = append(msg, "    |-d.PNG")
	fmt.Println("----------------------------------------------------")
	for _, m := range msg {
		fmt.Println(m)
	}
	fmt.Println("----------------------------------------------------")
}

func doCompress(workspace string, routine int) error {
	jobs, err := CreateJobs(workspace)
	if err != nil {
		return err
	}
	parallel.Do(jobs, routine)
	return nil
}
