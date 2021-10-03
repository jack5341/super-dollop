package main

import (
	"github.com/jack5341/super-dollop/cmd"
	"log"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	cmd.Execute()
}
