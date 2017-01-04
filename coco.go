/*
coco.go

Code Counter

Similar to the classic wc command, but counts comments and empty lines as
well. Lines that contain both code and comments are counted only as code.

Outputs "<total>: <lines> <comments> <empty lines>"
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	help    = flag.Bool("help", false, "Show usage")
	verbose = flag.Bool("verbose", false, "Output title for each count")
	single  = flag.String("single", "//", "Single-line comment")

	lines    uint32
	comments uint32
	empties  uint32
)

func init() {
	flag.BoolVar(help, "h", false, "")
	flag.BoolVar(verbose, "v", false, "")
	flag.StringVar(single, "s", "//", "Alias for -single")
}

func main() {
	flag.Parse()

	if *help {
		printHelp()
		os.Exit(1)
	}

	if flag.NArg() == 0 {
		read("<stdin>", os.Stdin)
	} else {
		for _, file := range flag.Args() {
			f, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error %s: %s\n", file, err)
			} else {
				read(file, f)
			}
			f.Close()
		}
	}

	total := lines + comments + empties
	if total > 0 {
		if *verbose {
			fmt.Println("Total:", total)
			fmt.Println("Lines:", lines)
			fmt.Println("Comments:", comments)
			fmt.Println("Empty:", empties)
		} else {
			fmt.Printf("%d: %d %d %d\n", total, lines, comments, empties)
		}
	}
}

func printHelp() {
	fmt.Println(`Usage: coco [options] [files...]
Prints a summary of total, non-commented, commented and empty lines.

The default comment style is C. To change it, use the -s and -m options.

If no files are given, standard input is used.

Example:
    coco -s "#" -m "" main.py
    cat lib.c | coco
`)
	flag.PrintDefaults()
}

func read(file string, f *os.File) {
	buf := bufio.NewReader(f)
	n := 0
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if line != "" {
					count(line, n)
				}
				return
			}
			fmt.Fprintf(os.Stderr, "error %s: %s\n", file, err)
			return
		}
		n++
		count(line, n)
	}
}

func count(line string, n int) {
	line = strings.Trim(line, " \t")

	if line == "\n" {
		empties++
	} else if strings.HasPrefix(line, *single) {
		comments++
	} else {
		lines++
	}
}
