//go:generate goversioninfo

package main

import (
	"flag"
	"fmt"
	"gotit/watermark"
	"os"
	"os/signal"
	"path/filepath"
)

func main() {
	action := flag.String("action", "watermark", "watermark")
	workspace := flag.String("workspace", filepath.Dir(os.Args[0]), "workspace")
	routine := flag.Int("routine", 16, "routine count")
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
	for true {
		sig := <-c
		if "interrupt" == sig.String() {
			os.Exit(0)
		}
	}

}
