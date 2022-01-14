package cmd

import (
	"bufio"
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

var commands []*actor

func init() {
	commands = make([]*actor, 0)
	commands = append(commands, &actor{
		id:   "watermarker",
		name: "图片加水印",
		action: func(args []string) error {
			return watermark.Do(args)
		},
	})
	commands = append(commands, &actor{
		id:   "img_compressor",
		name: "图片压缩",
		action: func(args []string) error {
			return compress.Do(args)
		},
	})
}

func Interact(args []string) {
	printHeader()
	for {
		input := scanInput()
		if len(input) == 0 {
			continue
		}
		if "q" == input {
			os.Exit(0)
			return
		}
		if egg(input) {
			continue
		}
		if strings.HasPrefix(input, "h") {
			usage(input)
			continue
		}
		exec(input, args)
	}
}

func egg(input string) bool {
	i := strings.ToLower(strings.ReplaceAll(input, " ", ""))
	switch i {
	case "iloveu", "iloveyou", "loveyou", "loveu", "爱你":
		dt, _ := resources.ReadData("egg/love")
		fmt.Println(string(dt))
		return true
	default:
		return false
	}
}

func exec(input string, args []string) {
	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(commands) {
		fmt.Println("请输入正确的数字")
		return
	}
	n := commands[index-1]
	fmt.Println("开始执行：", n.name, "(姜姜出品)")
	time.Sleep(time.Second)
	err = resources.ApesTesting()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = n.action(args)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Second)
		n.usage()
	}
	printHeader()
}

func usage(input string) {
	input = string([]rune(input)[1:])
	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(commands) {
		fmt.Println("请输入正确的数字")
		printHeader()
		return
	}
	n := commands[index-1]
	fmt.Println(n.name)
	n.usage()
}

func scanInput() string {
	fmt.Print("➜")
	reader := bufio.NewReader(os.Stdin)
	line, _, _ := reader.ReadLine()
	return strings.TrimSpace(strings.TrimSpace(string(line)))
}

func printHeader() {
	fmt.Printf("\n")
	for i, n := range commands {
		fmt.Printf("输入%2d 执行:%s\n", i+1, n.name)
	}
	fmt.Printf("输入 h+数字 查看数字对应功能的帮助信息，例如: h1\n")
	fmt.Printf("输入 q 退出\n")
}
