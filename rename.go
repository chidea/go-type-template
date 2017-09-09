// +build ignore
package main

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var debug bool

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-d" {
		os.Args = os.Args[1:]
	}
	if debug {
		log.Println(os.Args)
	}
	/*var oldnew [][2]string
	for _, arg := range os.Args[2:] {
		splt := strings.Split(arg, "=")
		if len(splt) > 2 {
			log.Panicf("%s has multiple equals", arg)
		}
		strings.Replace(code, splt[0], splt[1], -1)
		oldnew = append(oldnew, [2]string{splt[0], splt[1]})
	}*/

	b, e := ioutil.ReadFile(os.Args[1])
	if e != nil {
		log.Panic(e)
	}
	code := string(b)
	if debug {
		log.Println("before:")
		log.Println(code)
	}
	funcblockrgxp := regexp.MustCompile("^func .+_T_(.+).+\n}")
	funcblockrgxp.ReplaceAllStringFunc(code, func(v string) string {
		var rst string
		if debug {
			log.Println("found code:")
			log.Println(v)
		}
		for i, t := range os.Args[2:] {
			appendum := v
			appendum = strings.Replace(appendum, "_T_", t, 1)
			appendum = strings.Replace(appendum, "T", t, -1)
			if i > 0 {
				rst += "\n"
			}
			rst += appendum
		}
		return rst
	})
	if debug {
		log.Println("after:")
		log.Println(code)
	}
	/*i := 0
	for {
		i = strings.Index(code[i:], "\nfunc ")
		if i < 0 {
			break
		}
			strings.Index(code[i:], "}")
	}*/
}
