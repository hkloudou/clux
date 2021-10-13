package etcd

import (
	"testing"
	"time"

	"github.com/hkloudou/clux/sync"
)

func Test_etcd(t *testing.T) {
	t.Log(7)
	obj := NewSync()
	if err := obj.Lock("123", sync.LockTTL(10*time.Second), sync.LockWait(3*time.Second)); err != nil {
		t.Fatal(err)
	}
	// t.Log(11)
	// defer obj.Unlock("123")
	// time.Sleep(10 * time.Second)
	// t.Log(22)
}
