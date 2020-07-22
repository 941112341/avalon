package main

import (
	"flag"
	"fmt"
	"github.com/941112341/avalon/tool"
	"path/filepath"
)

func main() {

	var err error
	i := flag.String("i", "", "input thrift file ex:message.thrift")
	o := flag.String("o", "", "out put dir")

	flag.Parse()
	if i == nil || *i == "" {
		fmt.Println("args 'file' is nil")
		return
	}
	if o == nil || *o == "" {
		fmt.Println("args 'output' is nil")
		return
	}

	input, output := *i, *o
	if !filepath.IsAbs(input) {
		input, err = filepath.Abs(input)
		if err != nil {
			fmt.Println("abs input err " + err.Error())
			return
		}
	}
	if !filepath.IsAbs(output) {
		output, err = filepath.Abs(output)
		if err != nil {
			fmt.Println("abs output err " + err.Error())
			return
		}
	}

	err = tool.CreateFile(input, output)
	if err != nil {
		fmt.Println(err.Error())
	}
}
