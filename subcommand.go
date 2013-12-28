// A simple sub command parser based on the flag package
package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
)

type SubCommand interface {
	Name() string
	Synopsis() string
	DefineFlags(*flag.FlagSet)
	Run()
}

type subCommandParser struct {
	cmd SubCommand
	fs  *flag.FlagSet
}

func Parse(commands ...SubCommand) {
	scp := make(map[string]*subCommandParser, len(commands))
	for _, cmd := range commands {
		name := cmd.Name()
		scp[name] = &subCommandParser{cmd, flag.NewFlagSet(name, flag.ExitOnError)}
		cmd.DefineFlags(scp[name].fs)
	}

	oldUsage := flag.Usage
	flag.Usage = func() {
		oldUsage()
		fmt.Fprintf(os.Stderr, "Commands:\n")
		w := new(tabwriter.Writer)
		w.Init(os.Stderr, 10, 0, 2, ' ', 0)
		for name, sc := range scp {
			fmt.Fprintf(w, "   %s\t%s\n", name, sc.cmd.Synopsis())
		}
		w.Flush()
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	cmdname := flag.Arg(0)
	if sc, ok := scp[cmdname]; ok {
		sc.fs.Parse(flag.Args()[1:])
		sc.cmd.Run()
	} else {
		fmt.Fprintf(os.Stderr, "error: %s is not a valid command", cmdname)
		flag.Usage()
		os.Exit(1)
	}
}
