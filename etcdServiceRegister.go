package clux

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	// "go.etcd.io/etcd/clientv3/concurrency"
)

//创建租约注册服务
type ServiceRegister struct {
	etcdClient *clientv3.Client             //etcd client
	lease      clientv3.Lease               //租约
	leaseResp  *clientv3.LeaseGrantResponse //设置租约时间返回
	canclefunc func()                       //租约撤销
	//租约keepalieve相应chan
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string //注册的key
}

func RegisterLease(key, value string, TTL int64, opts ...EtcdOpt) {
	//注册模块
	if r, err := NewServiceRegister(TTL, opts...); err != nil {
		panic(err)
	} else if err := r.PutService(key, value); err != nil {
		panic(err)
	} else {
		r.ListenLeaseRespChan()
	}
}

// func RegisterTryLocker(pfx string, TTL int, opts ...EtcdOpt) (*concurrency.Mutex, *concurrency.Session, error) {
// 	if cli, err := GetEtcdClient(opts...); err != nil {
// 		return nil, nil, err
// 	} else if s1, err := concurrency.NewSession(cli, concurrency.WithTTL(TTL)); err != nil {
// 		return nil, nil, err
// 	} else {
// 		// defer s1.Close()
// 		// s1.Done()
// 		m1 := concurrency.NewMutex(s1, pfx)
// 		if err := m1.TryLock(context.TODO()); err != nil {
// 			return nil, nil, err
// 		}
// 		return m1, s1, nil
// 	}
// }

func NewServiceRegister(timeNum int64, opts ...EtcdOpt) (*ServiceRegister, error) {
	ser := &ServiceRegister{}
	if client, err := GetEtcdClient(opts...); err != nil {
		return nil, err
	} else {
		ser.etcdClient = client
	}
	//申请租约设置时间keepalive
	if err := ser.setLease(timeNum); err != nil {
		return nil, err
	}

	//监听续租相应chan
	go ser.ListenLeaseRespChan()
	return ser, nil
}

//设置租约
func (this *ServiceRegister) setLease(timeNum int64) error {
	//申请租约
	lease := clientv3.NewLease(this.etcdClient)

	//设置租约时间
	leaseResp, err := lease.Grant(context.TODO(), timeNum)
	if err != nil {
		return err
	}

	//设置续租 定期发送需求请求
	ctx, cancelFunc := context.WithCancel(context.TODO())
	if leaseRespChan, err := lease.KeepAlive(ctx, leaseResp.ID); err != nil {
		// cancelFunc()
		return err
	} else {
		this.lease = lease
		this.leaseResp = leaseResp
		this.canclefunc = cancelFunc
		this.keepAliveChan = leaseRespChan
	}
	return nil
}

//监听 续租情况
func (this *ServiceRegister) ListenLeaseRespChan() {
	for {
		select {
		case leaseKeepResp := <-this.keepAliveChan:
			if leaseKeepResp == nil {
				fmt.Printf("lease Closed\n")
				return
			}
		}
	}
}

//通过租约 注册服务
func (this *ServiceRegister) PutService(key, val string) error {
	//带租约的模式写入数据即注册服务
	kv := clientv3.NewKV(this.etcdClient)
	_, err := kv.Put(context.TODO(), key, val, clientv3.WithLease(this.leaseResp.ID))
	// res.
	return err
}

//撤销租约
func (this *ServiceRegister) RevokeLease() error {
	this.canclefunc()
	time.Sleep(2 * time.Second)
	_, err := this.lease.Revoke(context.TODO(), this.leaseResp.ID)
	return err
}
