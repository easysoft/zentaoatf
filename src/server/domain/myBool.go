package domain

import "fmt"

type MyBool bool

func (b *MyBool) Scan(src interface{}) error {
	str, ok := src.(int64)
	if !ok {
		return fmt.Errorf("Unexpected type for MyBool: %T", src)
	}
	switch str {
	case 0:
		v := false
		*b = MyBool(v)
	case 1:
		v := true
		*b = MyBool(v)
	}
	return nil
}
