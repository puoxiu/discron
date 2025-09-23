package etcdclient

import (
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"github.com/puoxiu/discron/common/models"
	"github.com/puoxiu/discron/common/pkg/logger"
	"time"
)

var _defalutEtcd *clientv3.Client

func Init(e models.Etcd) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   e.Endpoints,
		DialTimeout: time.Duration(e.DialTimeout)*time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	_defalutEtcd=cli
}

func GetEtcdClient() *clientv3.Client {
	if _defalutEtcd==nil{
		logger.Errorf("etcd client is not initialized")
		return nil
	}
	return _defalutEtcd
}