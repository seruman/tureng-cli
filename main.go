package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/seruman/tureng-cli/tureng"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "" {
		fmt.Println("Usage tureng [word] ")
		os.Exit(1)
	}

	results, err := tureng.Translate(os.Args[1])
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	for _, v := range results {
		fmt.Fprintf(w, "%s\t%s(%s)\t%s(%s)\n", v.CategoryTextB, v.TermA, v.TermTypeTextA, v.TermB, v.TermTypeTextB)
	}
	w.Flush()
}
