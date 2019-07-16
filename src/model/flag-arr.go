package model

import "fmt"

type FlagSlice []string

func (str *FlagSlice) String() string {
	return fmt.Sprintf("%s", *str)
}

func (str *FlagSlice) Set(value string) error {
	if value != "" {
		*str = append(*str, value)
	}
	return nil
}
