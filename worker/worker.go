package main

import (
	"fmt"

	"dinowernli.me/faucet/demo"
	pb_config "dinowernli.me/faucet/proto/config"
)

func main() {
	_ = &pb_config.Configuration{}
	fmt.Println(demo.Foo())
}
