package service

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	// "github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"github.com/puoxiu/discron/common/pkg/etcdclient"
	"log"
	"sync"
)

// NodeWatcher 用于监听 etcd 中节点（Node）的变化（新增/删除），维护节点列表
type NodeWatcher struct {
	client     *etcdclient.Client
	serverList map[string]string		// 存储节点信息：key=etcd中的键，value=节点数据（如节点ID/IP）
	lock       sync.Mutex
}

func NewNodeWatcher() *NodeWatcher {
	return &NodeWatcher{
		client:     etcdclient.GetEtcdClient(),
		serverList: make(map[string]string),
		lock:       sync.Mutex{},
	}
}

// Watch 初始化：先获取 etcd 中已存在的节点数据，再启动监听（异步）
func (s *NodeWatcher) Watch() error {
	// 1. 先查询 etcd 中前缀为 KeyEtcdNode 的所有已有数据（初始化节点列表）
	resp, err := s.client.Get(context.Background(), etcdclient.KeyEtcdNode, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	_ = s.extractAddrs(resp)

	// 2. 启动监听 etcd 中 KeyEtcdNode 前缀的变化（新增/删除）
	go s.watcher()
	return nil
}

// watcher 异步监听 etcd 数据变化，处理 PUT/DELETE 事件
func (s *NodeWatcher) watcher() {
	rch := s.client.Watch(context.Background(), etcdclient.KeyEtcdNode, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				//todo
				// 新增/更新节点：将 etcd 中的键值对存入 serverList
				s.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE:
				fmt.Println("server delete")
				s.DelServiceList(string(ev.Kv.Key))
			}
		}
	}
}

// extractAddrs 解析 etcd Get 响应，提取节点数据并填充到 serverList
func (s *NodeWatcher) extractAddrs(resp *clientv3.GetResponse) []string {
	addrs := make([]string, 0)
	if resp == nil || resp.Kvs == nil {
		return addrs
	}
	// 遍历所有 KV 对，存入 serverList
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
	log.Println("set data key :", key, "val:", val)
}

func (s *NodeWatcher) DelServiceList(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.serverList, key)
	log.Println("del data key:", key)
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
