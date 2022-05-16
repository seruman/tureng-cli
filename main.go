package main

import (
	"fmt"
	"os"
	"strings"
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
		v.Category = strings.ReplaceAll(v.Category, " ", "")
		fmt.Fprintf(w, "%s\t%s\t%s\n", v.Category, v.Term, v.TypeEN)
	}
	w.Flush()
}
