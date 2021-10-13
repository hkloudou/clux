package clux

import (
	"context"
	"errors"
	"sort"
	"sync"

	clientv3 "go.etcd.io/etcd/client/v3"
	// "go.etcd.io/etcd/mvcc/mvccpb"
	// "go.etcd.io/etcd/mvcc/mvccpb/v3"
)

type ServiceDiscovery struct {
	watched    bool
	client     *clientv3.Client
	serverList map[string]string
	lock       sync.Mutex
	onChange   func()
}

func NewServiceDiscovery(opts ...EtcdOpt) (*ServiceDiscovery, error) {
	bk := &ServiceDiscovery{
		serverList: make(map[string]string),
	}
	if client, err := GetEtcdClient(opts...); err != nil {
		return nil, err
	} else {
		bk.client = client
		return bk, nil
	}
}

func (this *ServiceDiscovery) Watch(prefix string, onChange func()) error {
	if this.watched {
		return errors.New("watched")
	}
	this.watched = true
	//使用key前桌获取所有的etcd上所有的server
	this.onChange = onChange
	resp, err := this.client.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	//解析出所有的server放入本地
	this.extractAddrs(resp)
	if this.onChange != nil {
		this.onChange()
	}
	//warch server前缀 将变更写入本地
	go this.watcher(prefix)
	return nil
}

// 监听key前缀
func (this *ServiceDiscovery) watcher(prefix string) {
	//监听 返回监听事件chan
	rch := this.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case 0: //修改或者新增
				this.setServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case 1: //删除
				this.delServiceList(string(ev.Kv.Key))
			}
		}
		if len(wresp.Events) > 0 && this.onChange != nil {
			this.onChange()
		}
	}
}

func (this *ServiceDiscovery) extractAddrs(resp *clientv3.GetResponse) []string {
	addrs := make([]string, 0)
	if resp == nil || resp.Kvs == nil {
		return addrs
	}
	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			this.setServiceList(string(resp.Kvs[i].Key), string(resp.Kvs[i].Value))
			addrs = append(addrs, string(v))
		}
	}
	return addrs
}

func (this *ServiceDiscovery) setServiceList(key, val string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.serverList[key] = string(val)
	// log.Println("set data key :", key, "val:", val)
}

func (this *ServiceDiscovery) delServiceList(key string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.serverList, key)
	// log.Println("del data key:", key)
}

func (this *ServiceDiscovery) List() []string {
	this.lock.Lock()
	defer this.lock.Unlock()
	addrs := make([]string, 0)
	// sort.Sort()
	for _, v := range this.serverList {
		addrs = append(addrs, v)
	}
	tmp := stringList(addrs)
	sort.Sort(tmp)
	return tmp
}

type stringList []string

func (I stringList) Len() int {
	return len(I)
}
func (I stringList) Less(i, j int) bool {
	return I[i] < I[j]
}
func (I stringList) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}
