package nav

import (
	"fmt"
	"gotit/img/compress"
	"gotit/img/watermark"
	"os"
	"strconv"
)

type actor struct {
	index  int
	id     string
	name   string
	action func(args []string)
}

var navigates []*actor

func init() {
	navigates = make([]*actor, 0)
	navigates = append(navigates, &actor{
		index: 1,
		id:    "watermarker",
		name:  "图片加水印",
		action: func(args []string) {
			watermark.Do(args)
		},
	})
	navigates = append(navigates, &actor{
		index: 2,
		id:    "img_compressor",
		name:  "图片压缩",
		action: func(args []string) {
			compress.Do(args)
		},
	})
}

func Start(args ...string) {
	printHeader := func() {
		fmt.Printf("\n\n-------------\n")
		for _, n := range navigates {
			fmt.Printf("输入%2d 执行:%s\n", n.index, n.name)
		}
		fmt.Printf("输入 q 退出\n")
	}

	printHeader()
	for {
		fmt.Printf("Opt>")
		var input string
		_, _ = fmt.Scanln(&input)
		if "q" == input {
			fmt.Println("exit")
			os.Exit(0)
		}
		if len(input) == 0 {
			continue
		}
		index, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("请输入数字编号")
		}
		for _, n := range navigates {
			if index == n.index {
				fmt.Println("开始执行：", n.name, "(姜姜出品)")
				n.action(args)
				printHeader()
				break
			}
		}
	}
}
