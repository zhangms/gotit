package cmd

import (
	"bufio"
	"fmt"
	resources "gotit/res"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Actor struct {
	Id     string
	Index  int
	Name   string
	Action func(args []string) error
}

func (a *Actor) usage() {
	data, err := resources.ReadData(fmt.Sprintf("usage/%s.md", a.Id))
	if err == nil {
		fmt.Println(string(data))
	}
}

var commands []*Actor

func init() {
	commands = make([]*Actor, 0)
}

func RegisterActor(actor *Actor) {
	commands = append(commands, actor)
}

func Interact(args []string) {
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Index < commands[j].Index
	})
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
	err := resources.ApesTesting()
	if err != nil {
		fmt.Println(err)
		return true
	}

	i := strings.ToLower(strings.ReplaceAll(input, " ", ""))
	switch i {
	case "iloveu", "iloveyou", "loveyou", "loveu", "爱你":
		dt, _ := resources.ReadData("egg/love")
		play(string(dt))
		return true
	default:
		return false
	}
}

func play(text string) {
	arr := strings.Split(text, "\n")
	for _, str := range arr {
		time.Sleep(500 * time.Millisecond)
		fmt.Println(str)
	}
}

func exec(input string, args []string) {
	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(commands) {
		fmt.Println("请输入正确的数字")
		return
	}
	n := commands[index-1]
	fmt.Println("开始执行：", n.Name, "(姜姜出品)")
	time.Sleep(time.Second)
	err = n.Action(args)
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
	fmt.Println(n.Name)
	n.usage()
}

func scanInput() string {
	fmt.Print(" ➜ ")
	reader := bufio.NewReader(os.Stdin)
	line, _, _ := reader.ReadLine()
	return strings.TrimSpace(strings.TrimSpace(string(line)))
}

func printHeader() {
	fmt.Printf("\n")
	for i, n := range commands {
		fmt.Printf("输入%2d 执行:%s\n", i+1, n.Name)
	}
	fmt.Printf("输入 h+数字 查看数字对应功能的帮助信息，例如: h1\n")
	fmt.Printf("输入 q 退出\n")
}
