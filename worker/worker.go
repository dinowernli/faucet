package main

import (
	"fmt"

	"dinowernli.me/faucet/demo"
	pb_config "dinowernli.me/faucet/proto/config"
	pb_worker "dinowernli.me/faucet/proto/service/worker"
)

func main() {
	_ = &pb_config.Configuration{}
	_ = &pb_worker.StatusRequest{}
	fmt.Println(demo.Foo())
}
