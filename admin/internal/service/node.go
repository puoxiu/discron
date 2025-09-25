package service

import (
	"context"
	"fmt"

	"go.etcd.io/etcd/client/v3"

	// "github.com/coreos/etcd/mvcc/mvccpb"
	"sync"

	"github.com/puoxiu/discron/common/pkg/etcdclient"
	"github.com/puoxiu/discron/common/pkg/logger"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

// NodeWatcher 用于监听 etcd 中节点（Node）的变化（新增/删除），维护节点列表
type NodeWatcher struct {
	client     *etcdclient.Client
	serverList map[string]string
	lock       sync.Mutex
}

func NewNodeWatcher() *NodeWatcher {
	return &NodeWatcher{
		client:     etcdclient.GetEtcdClient(),
		serverList: make(map[string]string),
	}
}

func (s *NodeWatcher) Watch() error {
	resp, err := s.client.Get(context.Background(), etcdclient.KeyEtcdNodeProfile, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	_ = s.extractAddrs(resp)

	go s.watcher()
	return nil
}

func (s *NodeWatcher) watcher() {
	rch := s.client.Watch(context.Background(), etcdclient.KeyEtcdNode, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				//todo
				s.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE:
				fmt.Println("server delete")
				s.DelServiceList(string(ev.Kv.Key))
			}
		}
	}
}

func (s *NodeWatcher) extractAddrs(resp *clientv3.GetResponse) []string {
	addrs := make([]string, 0)
	if resp == nil || resp.Kvs == nil {
		return addrs
	}
	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			s.SetServiceList(string(resp.Kvs[i].Key), string(resp.Kvs[i].Value))
			addrs = append(addrs, string(v))
		}
	}
	return addrs
}

func (s *NodeWatcher) SetServiceList(key, val string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.serverList[key] = val
	logger.GetLogger().Debug(fmt.Sprintf("set data key : %s val:%s", key, val))
}

func (s *NodeWatcher) DelServiceList(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.serverList, key)
	logger.GetLogger().Debug(fmt.Sprintf("del data key: %s", key))
}

func (s *NodeWatcher) SerList2Array() []string {
	s.lock.Lock()
	defer s.lock.Unlock()
	addrs := make([]string, 0)

	for _, v := range s.serverList {
		addrs = append(addrs, v)
	}
	return addrs
}

func (s *NodeWatcher) Close() error {
	return nil
}
