package watermark

import (
	"fmt"
	"gotit/parallel"
	"path/filepath"
)

func Do(args []string) {
	usage()
	doMark(filepath.Dir(args[0]), 16)
}

func usage() {
	msg := make([]string, 0)
	msg = append(msg, "----------------------------------------------------")
	msg = append(msg, "STEP1: 创建文件夹 workspace")
	msg = append(msg, "STEP2: 水印必须是png格式，水印图片文件名改为 mark.png 放到 workspace")
	msg = append(msg, "STEP3: 在 workspace 下创建文件夹 mark, 将所有需要加水印的图片或文件夹放入其中 ")
	msg = append(msg, "STEP4: 将该程序放到 workspace 下，双击执行，加水印后的图片会生成到 workspace 下的 mark_new 内")
	msg = append(msg, "")
	msg = append(msg, "最终目录结构如下：")
	msg = append(msg, "")
	msg = append(msg, "workspace")
	msg = append(msg, "|-mark.png")
	msg = append(msg, "|-gotit.exe")
	msg = append(msg, "|-mark")
	msg = append(msg, "  |-1.jpg")
	msg = append(msg, "  |-2.png")
	msg = append(msg, "  |-dir1")
	msg = append(msg, "    |-a.jpg")
	msg = append(msg, "    |-b.jpeg")
	msg = append(msg, "  |-dir2")
	msg = append(msg, "    |-c.png")
	msg = append(msg, "    |-d.PNG")
	msg = append(msg, "----------------------------------------------------")
	for _, m := range msg {
		fmt.Println(m)
	}
}

func doMark(workspace string, routine int) {
	jobs, err := CreateJobs(workspace)
	if err != nil {
		fmt.Println(err)
		return
	}
	parallel.Do(jobs, routine)
}
