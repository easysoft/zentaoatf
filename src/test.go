package main

import (
	"fmt"
	"path"
)

func main() {
	a := path.Base("ss/1.txt")

	fmt.Println(a)
}
