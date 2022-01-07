//go:generate goversioninfo

package main

import (
	"flag"
	"fmt"
	"gotit/watermark"
	"os"
	"os/signal"
)

func main() {
	action := flag.String("action", "watermark", "watermark")
	workspace := flag.String("workspace", ".", "workspace")
	routine := flag.Int("routine", 10, "routine count")
	flag.Parse()
	switch *action {
	case "watermark":
		watermark.Usage()
		watermark.DoMark(*workspace, *routine)
	}

	fmt.Println("----------------------")
	fmt.Println("按 ctrl + c 退出")
	c := make(chan os.Signal)
	signal.Notify(c)
	s := <-c
	fmt.Println("stop,signal : ", s)
}
