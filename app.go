package main

import (
	"fmt"
	"translator/cfg"
	"translator/tst/tt_log"
)

func main() {

	if err := cfg.GetInstance().Load(""); err != nil {
		panic(err)
	}

	tt_log.GetInstance()
	fmt.Println("hello world")
}
