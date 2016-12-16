package main

import (
	"fmt"
	"strings"
)

func main() {

	s1 := "abcd"
	byte1 := []byte(s1)
	fmt.Println(byte1)

	byte2 :=[]byte("get")
	s2 := string(byte2)
	if strings.HasPrefix(s2, "get"){
		fmt.Println("true")
	}
	fmt.Println(s2)
}