package main

import (
	"fmt"
)

//Message 消息
type Message string

func main() {
	fmt.Print("test:", Message([]byte{49, 50, 51}))
}
