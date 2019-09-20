package main

import (
	"fmt"
	"regexp"
)

var content = `[4. steps] 
    step4...
    line1...
    line2... 
  [4. expects] 
    
  [5. steps] 
    step5... 
  [5. expects] 
    
  [6. steps] 
    step6... 
  [6. expects]
?>`

func main() {
	myExp := regexp.MustCompile(`(?U)\[.*steps\]([\S\s]+)\[(.*steps.*|$)\]`)

	arr := myExp.FindAllStringSubmatch(content, -1)
	fmt.Printf("%+v", arr)
}
