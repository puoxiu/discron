package handler

import (
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"github.com/puoxiu/discron/common/pkg/etcdclient"
)

//获取node节点的信息
func WatchSystem(nodeUUID string) clientv3.WatchChan {
	return etcdclient.Watch(fmt.Sprintf(etcdclient.KeyEtcdSystemSwitch, nodeUUID), clientv3.WithPrefix())
}
