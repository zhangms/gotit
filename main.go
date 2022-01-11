//go:generate goversioninfo

package main

import (
	"gotit/nav"
	"os"
)

func main() {
	nav.Start(os.Args...)
}
