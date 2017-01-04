/*
coco.go

Code Counter

Similar to the classic wc command, but counts comments and empty lines as
well. Lines that contain both code and comments are counted only as code.

Outputs "<total> <lines> <comments> <empty lines>"
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
	multi   = flag.String("multi", "/* */", "Multi-line comment with start and end separated by a space.")

	lines    uint32
	comments uint32
	empties  uint32

	multiDelim [2]string
	inComment  bool = false
)

func init() {
	flag.BoolVar(help, "h", false, "")
	flag.BoolVar(verbose, "v", false, "")
	flag.StringVar(single, "s", "//", "Alias for -single")
	flag.StringVar(multi, "m", "/* */", "Alias for -multi")
}

func main() {
	flag.Parse()

	if *help {
		printHelp()
		os.Exit(1)
	}

	if *multi != "" {
		if split := strings.Split(*multi, " "); len(split) != 2 {
			fmt.Fprint(os.Stderr, "option error: bad format in -multi\n")
			os.Exit(1)
		} else {
			multiDelim[0] = split[0]
			multiDelim[1] = split[1]
		}
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
			fmt.Printf("%d %d %d %d\n", total, lines, comments, empties)
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

/*
* TODO: fix multiline
* */
func count(line string, n int) {
	line = strings.Trim(line, " \t")

	if inComment {
		comments++
		if strings.Contains(line, multiDelim[1]) {
			inComment = false
		}
		return
	}

	switch {
	case line == "\n":
		empties++
	case strings.HasPrefix(line, *single):
		comments++
	default:
		lines++
	}
}
