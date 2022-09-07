package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Must specify a filepath. Exiting.\n")
		return
	}

	filepath := os.Args[1]
	if is_jack_file(filepath) {
		a := new(JackAnalyzer).Init()
		a.analyze_file(filepath)
	} else {
		files, err := ioutil.ReadDir(filepath)
		if err != nil {
			fmt.Printf("Can't read directory.\n")
			return
		}

		if string(filepath[len(filepath)-1]) != "/" {
			filepath += "/"
		}
		for _, file := range files {
			filename := filepath + file.Name()
			if is_jack_file(filename) {
				a := new(JackAnalyzer).Init()
				a.analyze_file(filename)
			}
		}
	}
}

/** helpers **/

func is_jack_file(filename string) bool {
	return filename[len(filename)-5:] == ".jack"
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
