package models


import (
	"fmt"
	"github.com/puoxiu/discron/common/pkg/dbclient"
)

const (
	CronyGroupTypeUser = 1
	CronyGroupTypeNode = 2
)

type Group struct {
	ID   int    `json:"id" gorm:"id"`
	Name string `json:"name" gorm:"name" binding:"required"`
	//分组类型
	Type    int   `json:"type" gorm:"type" binding:"required"`
	Created int64 `json:"created" gorm:"created"`
	Updated int64 `json:"updated" gorm:"updated"`

	NodeIDs []string `json:"nids" gorm:"-"`
}

func (g *Group) Insert() (insertId int, err error) {
	err = dbclient.GetMysqlDB().Table(CronixGroupTableName).Create(g).Error
	if err == nil {
		insertId = g.ID
	}
	return
}
func (g *Group) Update() error {
	return dbclient.GetMysqlDB().Table(CronixGroupTableName).Updates(g).Error
}
func (g *Group) Delete() error {
	return dbclient.GetMysqlDB().Exec(fmt.Sprintf("delete from %s where id = ?", CronixGroupTableName), g.ID).Error
}
func (g *Group) FindById() error {
	return dbclient.GetMysqlDB().Table(CronixGroupTableName).Where("id = ? ", g.ID).First(g).Error
}

type NodeGroup struct {
	ID       int    `json:"id" gorm:"id"`
	NodeUUID string `json:"node_uuid" gorm:"node_uuid" binding:"required"`
	GroupId  int    `json:"group_id" gorm:"group_id" binding:"required"`
}

func (g *NodeGroup) Insert() (insertId int, err error) {
	err = dbclient.GetMysqlDB().Table(CronixNodeGroupTableName).Create(g).Error
	if err == nil {
		insertId = g.ID
	}
	return
}
func (g *NodeGroup) Update() error {
	return dbclient.GetMysqlDB().Table(CronixNodeGroupTableName).Updates(g).Error
}
func (g *NodeGroup) Delete() error {
	return dbclient.GetMysqlDB().Exec(fmt.Sprintf("delete from %s where node_uuid = ? and group_id = ?", CronixNodeGroupTableName), g.NodeUUID, g.GroupId).Error
}
func (g *NodeGroup) FindById() error {
	return dbclient.GetMysqlDB().Table(CronixNodeGroupTableName).Where("id = ? ", g.ID).First(g).Error
}

type UserGroup struct {
	ID      int `json:"id" gorm:"id"`
	UserId  int `json:"user_id" gorm:"user_id" binding:"required"`
	GroupId int `json:"group_id" gorm:"group_id" binding:"required" `
}

func (g *UserGroup) Insert() (insertId int, err error) {
	err = dbclient.GetMysqlDB().Table(CronixUserGroupTableName).Create(g).Error
	if err == nil {
		insertId = g.ID
	}
	return
}
func (g *UserGroup) Update() error {
	return dbclient.GetMysqlDB().Table(CronixUserGroupTableName).Updates(g).Error
}
func (g *UserGroup) Delete() error {
	return dbclient.GetMysqlDB().Exec(fmt.Sprintf("delete from %s where group_id = ? and  user_id =?", CronixUserGroupTableName), g.GroupId, g.UserId).Error
}
func (g *UserGroup) FindById() error {
	return dbclient.GetMysqlDB().Table(CronixUserGroupTableName).Where("id = ? ", g.ID).First(g).Error
}
