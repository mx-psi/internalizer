package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/mx-psi/internalizer/graph"
	"github.com/mx-psi/internalizer/internalizer"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [go module folder]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "missing input folder")
		os.Exit(1)
	}

	g, err := graph.FromFolder(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build graph for %q: %s", args[0], err)
		os.Exit(1)
	}

	moves, err := internalizer.Internalize(g)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to calculate internalizers for %q: %s", args[0], err)
		os.Exit(1)
	}

	pkgs := make([]string, 0, len(moves))
	for pkg := range moves {
		pkgs = append(pkgs, pkg)
	}
	sort.Strings(pkgs)

	for _, pkg := range pkgs {
		moved := moves[pkg]
		if strings.HasPrefix(moved, "github.com/DataDog/datadog-agent/internal") ||
			strings.HasPrefix(moved, "github.com/DataDog/datadog-agent/pkg/internal") ||
			strings.HasPrefix(moved, "github.com/DataDog/datadog-agent/cmd/internal") {
			continue
		}

		fmt.Printf("%s -> %s\n", pkg, moved)
	}
}
