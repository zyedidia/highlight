package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/zyedidia/highlight"
)

func main() {
	var defs []*highlight.Def
	gopath := os.Getenv("GOPATH")
	files, _ := ioutil.ReadDir(gopath + "/src/github.com/zyedidia/highlight/syntax_files")

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".yaml") {
			input, _ := ioutil.ReadFile(gopath + "/src/github.com/zyedidia/highlight/syntax_files/" + f.Name())
			d, err := highlight.ParseDef(input)
			if err != nil {
				fmt.Println(err)
				continue
			}
			defs = append(defs, d)
		}
	}

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
				if group == "statement" {
					color.Set(color.FgGreen)
				} else if group == "identifier" {
					color.Set(color.FgBlue)
				} else if group == "preproc" {
					color.Set(color.FgHiRed)
				} else if group == "special" {
					color.Set(color.FgRed)
				} else if group == "constant.string" {
					color.Set(color.FgCyan)
				} else if group == "constant" {
					color.Set(color.FgCyan)
				} else if group == "constant.specialChar" {
					color.Set(color.FgHiMagenta)
				} else if group == "type" {
					color.Set(color.FgYellow)
				} else if group == "constant.number" {
					color.Set(color.FgCyan)
				} else if group == "comment" {
					color.Set(color.FgHiGreen)
				} else {
					color.Unset()
				}
			}
			fmt.Print(string(c))
			colN++
		}
		if group, ok := matches[lineN][colN]; ok {
			if group == "default" || group == "" {
				color.Unset()
			}
		}

		fmt.Print("\n")
	}
}
