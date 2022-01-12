package nav

import (
	"fmt"
	"gotit/img/compress"
	"gotit/img/watermark"
	resources "gotit/res"
	"os"
	"strconv"
	"strings"
	"time"
)

type actor struct {
	id     string
	name   string
	action func(args []string) error
}

func (a *actor) usage() {
	data, err := resources.ReadData(fmt.Sprintf("usage/%s.md", a.id))
	if err == nil {
		fmt.Println(string(data))
	}
}

var navigates []*actor

func init() {
	navigates = make([]*actor, 0)
	navigates = append(navigates, &actor{
		id:   "watermarker",
		name: "图片加水印",
		action: func(args []string) error {
			return watermark.Do(args)
		},
	})
	navigates = append(navigates, &actor{
		id:   "img_compressor",
		name: "图片压缩",
		action: func(args []string) error {
			return compress.Do(args)
		},
	})
}

func Start(args ...string) {
	printHeader := func() {
		fmt.Printf("\n")
		for i, n := range navigates {
			fmt.Printf("输入%2d 执行:%s\n", i+1, n.name)
		}
		fmt.Printf("输入 q 退出\n")
	}

	scanInput := func() string {
		fmt.Printf("Opt>")
		var input string
		_, _ = fmt.Scanln(&input)
		return strings.TrimSpace(input)
	}

	printHeader()
	for {
		input := scanInput()
		switch input {
		case "q":
			os.Exit(0)
		case "":
			continue
		default:
			index, err := strconv.Atoi(input)
			if err != nil || index < 1 || index > len(navigates) {
				fmt.Println("请输入正确的数字")
				continue
			}
			n := navigates[index-1]
			fmt.Println("开始执行：", n.name, "(姜姜出品)")
			time.Sleep(time.Second)
			err = n.action(args)
			if err != nil {
				fmt.Println(err)
				time.Sleep(time.Second)
				n.usage()
				printHeader()
			}
		}
	}
}
