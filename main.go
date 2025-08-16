package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"code.selman.me/tureng-cli/tureng"
)

func main() {
	fs := flag.NewFlagSet("tureng", flag.ExitOnError)
	flagDebug := fs.Bool("debug", false, "enable debug mode to dump request/response")

	err := fs.Parse(os.Args[1:])
	if err != nil {
		fs.PrintDefaults()
		os.Exit(1)
	}

	args := fs.Args()

	if len(args) < 1 || args[0] == "" {
		fs.PrintDefaults()
		os.Exit(1)
	}

	t := tureng.NewClient(tureng.WithDebug(*flagDebug))
	results, err := t.Translate(args[0])
	if err != nil {
		os.Exit(1)
	}

	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	for _, v := range results {
		fmt.Fprintf(w, "%s\t%s(%s)\t%s(%s)\n", v.CategoryTextB, v.TermA, v.TermTypeTextA, v.TermB, v.TermTypeTextB)
	}
	w.Flush()
}
