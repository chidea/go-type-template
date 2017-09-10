package main

import (
	"io/ioutil"
	//"log"
	"fmt"
	"os"
	//"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var debug bool

func main() {
	os.Args = os.Args[1:]
	if debug = len(os.Args) > 1 && (os.Args[0] == "--debug" || os.Args[0] == "-d"); debug {
		os.Args = os.Args[1:]
		debug = true
	}
	if debug {
		//log.Println(os.Args)
		fmt.Fprintln(os.Stderr, os.Args)
	}
	var files []string
	var i int
	for j, arg := range os.Args {
		if strings.HasSuffix(arg, ".go") {
			files = append(files, arg)
		} else {
			i = j
			break
		}
	}
	if len(files) == 0 {
		g, e := filepath.Glob("./*_generate.go")
		if e != nil {
			fmt.Fprintln(os.Stderr, "no files specified to template")
			os.Exit(1)
		}
		for _, path := range g {
			files = append(files, path)
		}
	}
	os.Args = os.Args[i:]
	if len(os.Args) == 0 {
		fmt.Fprintln(os.Stderr, "no types specified to template")
		os.Exit(1)
	}
	if debug {
		fmt.Fprintln(os.Stderr, os.Args)
		//log.Println(os.Args)
		fmt.Fprintln(os.Stderr, files)
		//log.Println(files)
	}
	for _, file := range files {
		b, e := ioutil.ReadFile(file)
		if e != nil {
			//log.Panic(e)
			fmt.Fprintln(os.Stderr, e)
			os.Exit(1)
		}
		code := string(b)
		code = regexp.MustCompile("// \\+build generate(\r|\n\r|\n)").ReplaceAllString(code, "")
		code = regexp.MustCompile("//[ ]?go:generate [^\n\r]+(\r|\n\r|\n)").ReplaceAllString(code, "")
		varrgxp := regexp.MustCompile("var .+_T_ = .+_T_")
		funcblockrgxp := regexp.MustCompile("\nfunc [^(]+_T_\\(.+\\)[^}]+{\r?\n([^}].+(\r|\n\r|\n))+}(\r|\n\r|\n)")
		code = varrgxp.ReplaceAllStringFunc(code, replacefn(2))
		code = casergxp.ReplaceAllStringFunc(code, replaceCase)
		code = funcblockrgxp.ReplaceAllStringFunc(code, replacefn(1))
		code = `// Code generated with github.com/chidea/go-type-template DO NOT EDIT.
		` + code
		if strings.HasSuffix(file, "_generate.go") {
			ioutil.WriteFile(file[:len(file)-12]+".go", []byte(code), 733)
		} else {
			ioutil.WriteFile(file[:len(file)-3]+"-generated.go", []byte(code), 733)
		}
	}
}
func typeNameRule(typename string) string {
	rst := strings.ToUpper(string(t[0])) + t[1:]
	if typename == "string" {
		rst = rst[:3]
	}
	return rst
}

var trgxp = regexp.MustCompile("(]T[ ,){\r\n(]|[( ]T\\)|\\(T,| T,)")

func replacefn(typenamecnt int) func(string) string {
	return func(v string) string {
		var rst string
		if debug {
			fmt.Fprintf(os.Stderr, "found function:\n%s", v)
			//log.Println("found code:")
			//log.Println(v)
		}
		for i, t := range os.Args {
			appendum := v[0:]
			typename := typeNameRule(t)
			appendum = strings.Replace(appendum, "_T_", typename, typenamecnt)
			appendum = trgxp.ReplaceAllStringFunc(appendum, templateReplace)
			if i > 0 {
				rst += "\n"
			}
			if debug {
				fmt.Fprintf(os.Stderr, "replaced T to %s:\n%s", t, appendum)
				//log.Printf("replaced T to %s:\n%s", t, appendum)
			}
			rst += appendum
		}
		return rst
	}
}

//var casergxp = regexp.MustCompile("(\r|\n\r|\n)[ ]+switch [^.]+.\\(type\\) {(\r|\n\r|\n)[ ]+case [^:]+T:(\r|\n\r|\n).+(\r|\n\r|\n)[ ]+}(\r|\n\r|\n)")
var casergxp = regexp.MustCompile("(\r|\n\r|\n)[\t ]+case .+T:(\r|\n\r|\n).+(\r|\n\r|\n)[\t ]+(}(\r|\n\r|\n)|case )")

func replaceCase(v string) string {
	var rst string
	if debug {
		fmt.Fprintf(os.Stderr, "found case:\n%s", v)
	}
	for i, t := range os.Args {
		appendum := v[0:]
		typename := typeNameRule(t)
		appendum = strings.Replace(appendum, "_T_", typename, -1)
		appendum = trgxp.ReplaceAllStringFunc(appendum, templateReplace)
		if debug {
			fmt.Fprintf(os.Stderr, "replaced T to %s:\n%s", t, appendum)
		}
		if i > 0 {
			rst += "\n"
		}
		rst += appendum
	}
}

func templateReplace(v string) string {
	return string(v[0]) + t + string(v[len(v)-1])
}
