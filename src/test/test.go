package main

import (
	"fmt"
	"github.com/ajg/form"
)

func main() {

	user := map[string]interface{}{
		"Name":   "joeybloggs",
		"Age":    3,
		"Gender": "Male",
		"steps":  map[string]string{"1": "false", "2": "true"},
	}

	val, _ := form.EncodeToValues(user)
	fmt.Printf("%s\n", val.Encode())
}
