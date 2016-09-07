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

	jj, err := tureng.PrepareReq(os.Args[1])
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	results, err := tureng.Translate(jj)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 10, 0, '\t', 0)
	for _, v := range results.MobileResult.Results {
		v.Category = strings.Replace(v.Category, " ", "", -1)
		fmt.Printf("%s \t%s\n", v.Category, v.Term)
	}
	w.Flush()

}
