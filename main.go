//go:generate goversioninfo

package main

import (
	"gotit/cmd"
	"os"
)

func main() {

	cmd.Interact(os.Args)
}
