package main

import (
	"flag"
	"fmt"
	"os"
)

type Logout struct {
}

func (cmd *Logout) Name() string                 { return "logout" }
func (cmd *Logout) Synopsis() string             { return "End session with MBO" }
func (cmd *Logout) DefineFlags(fs *flag.FlagSet) {}
func (cmd *Logout) Run() {
	filename := fmt.Sprintf("%s/.mindbodyonline", os.Getenv("HOME"))
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error opening %s: %s", filename, err)
	}
	defer file.Close()
}
