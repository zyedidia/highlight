package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/zyedidia/highlight"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("No input file")
		return
	}

	gopath := os.Getenv("GOPATH")

	var syn_dir string
	if gopath == "" {
		syn_dir = os.Getenv("HOME") + "/.config/zaje/syntax_files"
	} else {
		syn_dir = gopath + "/src/github.com/jessp01/gohighlight/syntax_files"
	}

	var defs []*highlight.Def
	err, warnings := highlight.ParseSyntaxFiles(syn_dir, &defs)
	if err != nil {
		log.Fatal(err)
	}

	// it's up to you what to do with the warnings. You can print them or ignore them
	fmt.Println(warnings)

	highlight.ResolveIncludes(defs)

	fileSrc, _ := ioutil.ReadFile(os.Args[1])
	def := highlight.DetectFiletype(defs, os.Args[1], bytes.Split(fileSrc, []byte("\n"))[0])

	if def == nil {
		fmt.Println(string(fileSrc))
		return
	}

	h := highlight.NewHighlighter(def)

	matches := h.HighlightString(string(fileSrc))

	lines := strings.Split(string(fileSrc), "\n")
	for lineN, l := range lines {
		colN := 0
		for _, c := range l {
			if group, ok := matches[lineN][colN]; ok {
				switch group {
				case highlight.Groups["statement"]:
					fallthrough
				case highlight.Groups["green"]:
					color.Set(color.FgGreen)

				case highlight.Groups["identifier"]:
					fallthrough
				case highlight.Groups["blue"]:
					color.Set(color.FgHiBlue)

				case highlight.Groups["preproc"]:
					//fallthrough
					//case highlight.Groups["high.red"]:
					color.Set(color.FgHiRed)

				case highlight.Groups["special"]:
					fallthrough
				case highlight.Groups["red"]:
					color.Set(color.FgRed)

				case highlight.Groups["constant.string"]:
					fallthrough
				case highlight.Groups["constant"]:
					fallthrough
				case highlight.Groups["constant.number"]:
					fallthrough
				case highlight.Groups["cyan"]:
					color.Set(color.FgCyan)

				case highlight.Groups["constant.specialChar"]:
					fallthrough
				case highlight.Groups["magenta"]:
					color.Set(color.FgHiMagenta)

				case highlight.Groups["type"]:
					fallthrough
				case highlight.Groups["yellow"]:
					color.Set(color.FgYellow)

				case highlight.Groups["comment"]:
					fallthrough
				case highlight.Groups["high.green"]:
					color.Set(color.FgHiGreen)
				default:
					color.Unset()
				}
			}
			fmt.Print(string(c))
			colN++
		}
		if group, ok := matches[lineN][colN]; ok {
			if group == highlight.Groups["default"] || group == highlight.Groups[""] {
				color.Unset()
			}
		}

		color.Unset()
		fmt.Print("\n")
	}
}
