////go:generate goversioninfo

package main

import (
	"fmt"
	"gotit/dirgroup"
	_ "gotit/img/compress"
	_ "gotit/img/resolution"
	_ "gotit/img/watermark"
	"os"
	"path/filepath"
)

func main() {
	//cmd.Interact(os.Args)
	path, _ := os.Executable()
	fmt.Println(path)
	err := dirgroup.Exec(filepath.Dir(path))
	if err != nil {
		panic(err)
	}
}
