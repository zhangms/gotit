//go:generate goversioninfo

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"test/watermark"
)

func main() {
	action := flag.String("action", "watermark", "watermark")
	workspace := flag.String("workspace", ".", "workspace")
	flag.Parse()
	switch *action {
	case "watermark":
		watermark.Usage()
		watermark.DoMark(*workspace)
	}

	fmt.Println("----------------------")
	fmt.Println("按 ctrl + c 退出")
	c := make(chan os.Signal)
	signal.Notify(c)
	s := <-c
	fmt.Println("stop,signal : ", s)
}
