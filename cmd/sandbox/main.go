package main

import (
	"fmt"
	"github.com/OwO-Network/gdeeplx"
)

func main() {
	result, err := gdeeplx.Translate("Hello World!", "EN", "ZH", 0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(result)
}
