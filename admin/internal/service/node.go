package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/puoxiu/discron/admin/internal/model/request"
	"github.com/puoxiu/discron/common/models"
	"github.com/puoxiu/discron/common/pkg/config"
	"github.com/puoxiu/discron/common/pkg/dbclient"
	"github.com/puoxiu/discron/common/pkg/etcdclient"
	"github.com/puoxiu/discron/common/pkg/logger"
	"github.com/puoxiu/discron/common/pkg/notify"
	"github.com/puoxiu/discron/common/pkg/utils"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
)

type NodeWatcherService struct {
	client   *etcdclient.Client
	nodeList map[string]models.Node
	lock     sync.Mutex
}

var DefaultNodeWatcher *NodeWatcherService

func NewNodeWatcherService() *NodeWatcherService {
	return &NodeWatcherService{
		client:   etcdclient.GetEtcdClient(),
		nodeList: make(map[string]models.Node),
	}
}

// Watch 启动节点监控：先拉取现有节点，再启动监听协程
func (n *NodeWatcherService) Watch() error {
	resp, err := n.client.Get(context.Background(), etcdclient.KeyEtcdNodeProfile, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	_ = n.extractNodes(resp)

	go n.watcher()
	return nil
}

// watcher 持续监听 etcd 中节点的 PUT/DELETE 事件
func (n *NodeWatcherService) watcher() {
	rch := n.client.Watch(context.Background(), etcdclient.KeyEtcdNodeProfile, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				n.setNodeList(n.GetUUID(string(ev.Kv.Key)), string(ev.Kv.Value))
			case mvccpb.DELETE:
				uuid := n.GetUUID(string(ev.Kv.Key))
				n.delNodeList(uuid)
				logger.GetLogger().Warn(fmt.Sprintf("Cronix node[%s] DELETE event detected", uuid))

				node := &models.Node{UUID: uuid}
				err := node.FindByUUID()
				if err != nil {
					logger.GetLogger().Error(fmt.Sprintf("Cronix node[%s] find by uuid  error:%s", uuid, err.Error()))
					return
				}

				success, fail, err := n.FailOver(uuid)
				if err != nil {
					logger.GetLogger().Error(fmt.Sprintf("Cronix node[%s] fail over error:%s", uuid, err.Error()))
					return
				}
				// if the failover is all successful, delete the node in the database
				if fail.Count() == 0 {
					err = node.Delete()
					if err != nil {
						logger.GetLogger().Error(fmt.Sprintf("Cronix node[%s] delete by uuid  error:%s", uuid, err.Error()))
					}
				}
				//Node inactivation information defaults to email.
				msg := &notify.Message{
					Type:      notify.NotifyTypeMail,
					IP:        fmt.Sprintf("%s:%s", node.IP, node.PID),
					Subject:   "节点失活报警",
					Body:      fmt.Sprintf("[Cronix Warning]Cronix node[%s] in the cluster has failed,，fail over success count:%d jobID are :%s ,fail count:%d jobID are :%s ", uuid, success.Count(), success.String(), fail.Count(), fail.String()),
					To:        config.GetConfigModels().Email.To,
					OccurTime: time.Now().Format(utils.TimeFormatSecond),
				}

				go notify.Send(msg)
			}
		}
	}
}

// extractNodes 从etcd中提取节点信息 存储到内存
func (n *NodeWatcherService) extractNodes(resp *clientv3.GetResponse) []string {
	nodes := make([]string, 0)
	if resp == nil || resp.Kvs == nil {
		return nodes
	}
	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			n.setNodeList(n.GetUUID(string(resp.Kvs[i].Key)), string(resp.Kvs[i].Value))
			nodes = append(nodes, string(v))
		}
	}
	return nodes
}

// setNodeList 新增/更新节点到内存，并为未分配任务分配节点
func (n *NodeWatcherService) setNodeList(key, val string) {
	var node models.Node
	err := json.Unmarshal([]byte(val), &node)
	if err != nil {
		logger.GetLogger().Warn(fmt.Sprintf("discover node[%s] json error:%s", key, err.Error()))
		return
	}
	n.lock.Lock()
	n.nodeList[key] = node
	n.lock.Unlock()
	logger.GetLogger().Debug(fmt.Sprintf("discover node node[%s] with pid[%s]", key, val))
	//Wait for the node to be fully started and assign the node
	time.Sleep(5 * time.Second)
	//find unassigned job
	jobs, err := DefaultJobService.GetNotAssignedJob()
	if err != nil {
		logger.GetLogger().Warn(fmt.Sprintf("discover node[%s],pid[%s] and get not assigned job err:%s", key, val, err.Error()))
		return
	}
	for _, job := range jobs {
		/*if job.Type == models.JobTypeCmd && !config.GetConfigModels().System.CmdAutoAllocation {
			logger.GetLogger().Warn(fmt.Sprintf("assign unassigned job[%d]  don't support cmd type", job.ID))
			continue
		}*/
		err = job.Unmarshal()
		if err != nil {
			logger.GetLogger().Warn(fmt.Sprintf("assign unassigned job[%d] json unmarshal error:%s", job.ID, err.Error()))
			continue
		}
		oldUUID := job.RunOn
		nodeUUID := DefaultJobService.AutoAllocateNode()
		if nodeUUID == "" {
			//If automatic allocation fails, it will be directly assigned to the new node.
			nodeUUID = key
		}
		err = n.assignJob(nodeUUID, &job)
		if err != nil {
			logger.GetLogger().Warn(fmt.Sprintf("assign unassigned job[%d]  error:%s", job.ID, err.Error()))
			continue
		}
		//debug Delete the key value if the transfer is successful
		_, err = etcdclient.Delete(fmt.Sprintf(etcdclient.KeyEtcdJob, oldUUID, job.ID))
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("node[%s] job[%d] fail over etcd delete job error:%s", nodeUUID, job.ID, err.Error()))
			continue
		}
	}
}

func (n *NodeWatcherService) delNodeList(key string) {
	n.lock.Lock()
	defer n.lock.Unlock()
	delete(n.nodeList, key)
	logger.GetLogger().Debug(fmt.Sprintf("delelte node[%s]", key))
}

func (n *NodeWatcherService) List2Array() []string {
	n.lock.Lock()
	defer n.lock.Unlock()
	nodes := make([]string, 0)

	for k, _ := range n.nodeList {
		nodes = append(nodes, k)
	}
	return nodes
}

func (n *NodeWatcherService) Close() error {
	return nil
}

func (n *NodeWatcherService) GetUUID(key string) string {
	// /Cronix/node/<node_uuid>
	index := strings.LastIndex(key, "/")
	if index == -1 {
		return ""
	}
	return key[index+1:]
}

func (n *NodeWatcherService) Search(s *request.ReqNodeSearch) ([]models.Node, int64, error) {
	db := dbclient.GetMysqlDB().Table(models.CronixNodeTableName)
	if len(s.UUID) > 0 {
		db = db.Where("uuid = ?", s.UUID)
	}
	if len(s.IP) > 0 {
		db.Where("ip = ?", s.IP)
	}
	if s.Status > 0 {
		db.Where("status = ?", s.Status)
	}
	if s.UpTime > 0 {
		db.Where("up > ?", s.UpTime)
	}
	nodes := make([]models.Node, 2)
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(s.PageSize).Offset((s.Page - 1) * s.PageSize).Order("up desc").Find(&nodes).Error
	if err != nil {
		return nil, 0, err
	}

	return nodes, total, nil
}

func (n *NodeWatcherService) GetJobCount(nodeUUID string) (int, error) {
	resps, err := etcdclient.Get(fmt.Sprintf(etcdclient.KeyEtcdJobProfile, nodeUUID), clientv3.WithPrefix(), clientv3.WithCountOnly())
	if err != nil {
		return 0, err
	}
	return int(resps.Count), nil
}

type Result []int

func (r Result) Count() (count int) {
	for _, v := range r {
		if v != 0 {
			count++
		}
	}
	return
}
func (r Result) String() (str string) {
	str = "["
	for _, v := range r {
		if v != 0 {
			str += fmt.Sprintf("%d,", v)
		}
	}
	str += "]"
	return
}

func (n *NodeWatcherService) assignJob(nodeUUID string, job *models.Job) (err error) {
	if nodeUUID == "" {
		return fmt.Errorf("node uuid can't be null")
	}
	node, ok := n.nodeList[nodeUUID]
	if !ok {
		return fmt.Errorf("assign unassigned job[%d] but  node[%s] not exist ", job.ID, nodeUUID)
	}
	job.InitNodeInfo(models.JobStatusAssigned, node.UUID, node.Hostname, node.IP)

	b, err := json.Marshal(job)
	if err != nil {
		return
	}
	_, err = etcdclient.Put(fmt.Sprintf(etcdclient.KeyEtcdJob, nodeUUID, job.ID), string(b))
	if err != nil {
		return
	}
	err = job.Update()
	if err != nil {
		return
	}
	return
}

func (n *NodeWatcherService) FailOver(nodeUUID string) (success Result, fail Result, err error) {
	jobs, err := n.GetJobs(nodeUUID)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("node[%s] fail over get jobs error:%s", nodeUUID, err.Error()))
		return
	}
	if len(jobs) == 0 {
		return
	}
	for _, job := range jobs {
		//Determine whether shell command failover is supported
		/*if job.Type == models.JobTypeCmd && !config.GetConfigModels().System.CmdAutoAllocation {
			logger.GetLogger().Warn(fmt.Sprintf("node[%s] job[%d] fail over don't support cmd type", nodeUUID, job.ID))
			fail = append(fail, job.ID)
			continue
		}*/
		oldUUID := job.RunOn
		autoUUID := DefaultJobService.AutoAllocateNode()
		if autoUUID == "" {
			logger.GetLogger().Warn(fmt.Sprintf("node[%s] job[%d] fail over auto allocate node error", nodeUUID, job.ID))
			fail = append(fail, job.ID)
			continue
		}
		err = n.assignJob(autoUUID, &job)
		if err != nil {
			logger.GetLogger().Warn(fmt.Sprintf("node[%s] job[%d] fail over assign job error:%s", nodeUUID, job.ID, err.Error()))
			fail = append(fail, job.ID)
			continue
		}
		//Delete the key value if the transfer is successful
		_, err = etcdclient.Delete(fmt.Sprintf(etcdclient.KeyEtcdJob, oldUUID, job.ID))
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("node[%s] job[%d] fail over etcd delete job error:%s", nodeUUID, job.ID, err.Error()))
			fail = append(fail, job.ID)
			continue
		}
		success = append(success, job.ID)
	}
	return
}

//get all the job under a node
func (n *NodeWatcherService) GetJobs(nodeUUID string) (jobs []models.Job, err error) {
	resps, err := etcdclient.Get(fmt.Sprintf(etcdclient.KeyEtcdJobProfile, nodeUUID), clientv3.WithPrefix())
	if err != nil {
		return
	}
	count := len(resps.Kvs)
	if count == 0 {
		return
	}
	for _, j := range resps.Kvs {
		var job models.Job
		if err := json.Unmarshal(j.Value, &job); err != nil {
			logger.GetLogger().Warn(fmt.Sprintf("job[%s] umarshal err: %s", string(j.Key), err.Error()))
			continue
		}
		jobs = append(jobs, job)
	}
	return
}

func (n *NodeWatcherService) GetNodeCount(status int) (int64, error) {
	db := dbclient.GetMysqlDB().Table(models.CronixNodeTableName)
	if status > 0 {
		db = db.Where("status = ?", status)
	}
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

func GetNodeSystemInfo(uuid string) (s *utils.Server, err error) {
	defer func() {
		_, err = etcdclient.Delete(fmt.Sprintf(etcdclient.KeyEtcdSystemSwitch, uuid))
	}()
	s = new(utils.Server)
	res, err := etcdclient.Get(fmt.Sprintf(etcdclient.KeyEtcdSystemGet, uuid), clientv3.WithPrefix())
	if err != nil || len(res.Kvs) == 0 {
		return
	}
	err = json.Unmarshal(res.Kvs[0].Value, s)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("json error:%v", err))
	}
	return
}
