package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	var dir string

	arg1 := os.Args[0]

	if strings.Index(arg1, "build") > -1 {
		dir, _ = os.Getwd()
	} else {
		dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}

	fmt.Println(dir)

}
