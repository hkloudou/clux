package main

import (
	"log"
	"time"

	"github.com/hkloudou/clux"
)

func main() {
	if tmp, err := clux.NewServiceDiscovery(); err != nil {
		panic(err)
	} else if err := tmp.Watch("/system/proxy/nodes/", func() {
		log.Println("change")
		for _, v := range tmp.List() {
			log.Println(v)
		}
		time.Sleep(1 * time.Second)
	}); err != nil {
		panic(err)
	}
	<-make(chan bool)
}
