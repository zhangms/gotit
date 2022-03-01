//go:generate goversioninfo

package main

import (
	"gotit/cmd"
	_ "gotit/img/compress"
	_ "gotit/img/resolution"
	_ "gotit/img/watermark"
	"os"
)

func main() {
	cmd.Interact(os.Args)
}
