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
	/*if debug {
		//log.Println(os.Args)
		fmt.Fprintln(os.Stderr, os.Args)
	}*/
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
		fmt.Fprintln(os.Stderr, "template files:", files)
		fmt.Fprintln(os.Stderr, "generating types:", os.Args)
		//log.Println(os.Args)
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
		casergxp := regexp.MustCompile("(\r|\n\r|\n)[\t ]+case T:([\r\n\t ]+[^\r\n]+)*?[\r\n\t ]+(}(\r|\n\r|\n| |//|/*)|case )")
		code = casergxp.ReplaceAllStringFunc(code, replaceCase())
		trgxp := regexp.MustCompile("(]T[ ,){\r\n(]|[( ]T\\)|\\(T,| T,)")
		code = varrgxp.ReplaceAllStringFunc(code, replacefn("var", 2, trgxp))
		funcblockrgxp := regexp.MustCompile("\nfunc [^(]+_T_\\(.+\\)[^}]+{(\r|\n\r|\n)([^}].+(\r|\n\r|\n))+}(\r|\n\r|\n)")
		code = funcblockrgxp.ReplaceAllStringFunc(code, replacefn("function", 1, trgxp))
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
	rst := strings.ToUpper(string(typename[0])) + typename[1:]
	if typename == "string" {
		rst = rst[:3]
	}
	return rst
}

func replacefn(name string, typenamecnt int, trgxp *regexp.Regexp) func(string) string {
	return func(v string) string {
		var rst string
		if debug {
			fmt.Fprintf(os.Stderr, "found %s:\n%s\n", name, v)
			//log.Println("found code:")
			//log.Println(v)
		}
		for i, t := range os.Args {
			appendum := v[0:]
			typename := typeNameRule(t)
			appendum = strings.Replace(appendum, "_T_", typename, typenamecnt)
			appendum = trgxp.ReplaceAllStringFunc(appendum, templateReplace(t))
			if i > 0 {
				rst += "\n"
			}
			if debug {
				fmt.Fprintf(os.Stderr, "replaced T to %s:\n%s\n", t, appendum)
				//log.Printf("replaced T to %s:\n%s", t, appendum)
			}
			rst += appendum
		}
		return rst
	}
}

//var casergxp = regexp.MustCompile("(\r|\n\r|\n)[ ]+switch [^.]+.\\(type\\) {(\r|\n\r|\n)[ ]+case [^:]+T:(\r|\n\r|\n).+(\r|\n\r|\n)[ ]+}(\r|\n\r|\n)")

func replaceCase() func(string) string {
	fn := replacefn("case", -1, regexp.MustCompile("( T:|.T\\))"))
	endrgxp := regexp.MustCompile("(\r|\n\r|\n)[\t ]+(}(\r|\n\r|\n| |//|/*)|case )")
	return func(v string) string {
		f := endrgxp.FindAllString(v, -1)
		end := f[len(f)-1]
		v = v[:len(v)-len(end)]
		return fn(v) + end
	}
}

func templateReplace(typename string) func(string) string {
	return func(v string) string { return v[:1] + typename + v[len(v)-1:] }
}
