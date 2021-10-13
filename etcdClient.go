package clux

import (
	"crypto/tls"
	"errors"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var _defaultEtcdClient *clientv3.Client

type etcdOpts struct {
	clientv3.Config
	cli *clientv3.Client
}

// etcdOptFn is a function
type etcdOptFn func(opts *etcdOpts) error

func (opt etcdOptFn) config(opts *etcdOpts) error {
	return opt(opts)
}

// EtcdOpt configures options
type EtcdOpt interface {
	config(opts *etcdOpts) error
}

// EtcdConfigEndpoints sets the endpoints of etcd
func WithEtcdConfigEndpoints(endpoints []string) EtcdOpt {
	return etcdOptFn(func(opts *etcdOpts) error {
		opts.Endpoints = endpoints
		return nil
	})
}

// EtcdConfigUserPwd sets the username password of etcd
func WithEtcdConfigUserPwd(username, password string) EtcdOpt {
	return etcdOptFn(func(opts *etcdOpts) error {
		opts.Username = username
		opts.Password = password
		return nil
	})
}

// EtcdConfigTls sets the tls of etcd
func WithEtcdConfigTls(tls *tls.Config) EtcdOpt {
	return etcdOptFn(func(opts *etcdOpts) error {
		opts.TLS = tls
		return nil
	})
}

// EtcdConfigTls sets the tls of etcd
func WithEtcdConfigClient(cli *clientv3.Client) EtcdOpt {
	return etcdOptFn(func(opts *etcdOpts) error {
		opts.cli = cli
		return nil
	})
}

func GetEtcdClient(opts ...EtcdOpt) (*clientv3.Client, error) {
	cfg := &etcdOpts{}
	for i := 0; i < len(opts); i++ {
		opts[i].config(cfg)
	}
	if cfg.cli != nil {
		return cfg.cli, nil
	}
	if cfg.Endpoints == nil {
		cfg.Endpoints = []string{"localhost:2379"}
	}
	if cfg.DialTimeout == 0 {
		cfg.DialTimeout = 5 * time.Second
	}

	return clientv3.New(cfg.Config)
}

func GetEtcdDefaultClient() (*clientv3.Client, error) {
	if _defaultEtcdClient == nil {
		return _defaultEtcdClient, errors.New("not init")
	}
	return _defaultEtcdClient, nil
}
