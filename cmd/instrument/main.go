/*
 * @Author: xhang 1263403710@qq.com
 * @Date: 2025-04-01 11:21:20
 * @LastEditors: xhang 1263403710@qq.com
 * @LastEditTime: 2025-04-01 11:22:12
 * @FilePath: /github.com/instrument_trace/cmd/instrument/main.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/bigwhite/instrument_trace/instrumenter"
)

var (
	wrote bool
	pkg   string
)

func init() {
	flag.BoolVar(&wrote, "w", false, "write result to (source) file instead of stdout")
	flag.StringVar(&pkg, "pkg", "", "package to be instrumented")
}

func usage() {
	fmt.Println("instrument [-w] [-pkg package] [path...]")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if len(flag.Args()) == 0 {
		usage()
		return
	}

	for _, path := range flag.Args() {
		path, err := filepath.Abs(path)
		if err != nil {
			fmt.Printf("get absolute path for %s error: %v\n", path, err)
			continue
		}

		err = instrumenter.Instrument(path, pkg, wrote)
		if err != nil {
			fmt.Printf("instrument %s error: %v\n", path, err)
		}
	}
}
