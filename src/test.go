package main

import (
	"fmt"
	"path"
)

func main() {
	a := path.Ext("ss/1.txt")

	fmt.Println(a)
}
