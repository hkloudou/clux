package main

import (
	"log"
	"time"

	"github.com/hkloudou/clux/plugins/sync/etcd"
	"github.com/hkloudou/clux/sync"
)

func main() {
	obj := etcd.NewSync()
	if err := obj.Lock(
		"123",
		sync.LockTTL(10*time.Second),
		sync.LockWait(3*time.Second),
	); err != nil {
		log.Fatal(err)
	}
	// clux.RegisterLease("")
	// c, _ := clux.GetEtcdClient()
	// for s
}
