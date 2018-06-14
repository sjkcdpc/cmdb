package model

import (
	"github.com/jinzhu/gorm"
)

type HostsItem int

const (
	Active HostsItem = iota
	InActive
	Banned
	BanDied
)

// 主机表
type Hosts struct {
	gorm.Model
	Wanip    string `json:"WanIp" grom:"not null;unique"`
	Lanip    string `json:"LanIp"`
	Conf     string `json:"Conf"`
	Hostname string `json:"HostName"`
	Os       string `json:"Os" grom:"unique"`
	Contact  string `json:"Contact"` // 业务线负责人
	Manager  string `json:"Manager"` // 系统负责人
	Sshport  int64  `json:"SshPort"`
	Tags     string `json:"Tags"`
	Remark   string `json:"Remark"`
}

// 用户表
type Users struct {
	gorm.Model
	Username   string    `json:"Username"`
	Password   string    `json:"Password"`
	Birthday   string `json:"Birthday"`
	Email      string    `json:"Email"`
	Phone      int64     `json:"Phone"`
	Department string    `json:"Department"`
	Enable     bool      `json:"Enable"`
	Tags       string    `json:"Tags"`
	Remark     string    `json:"Remark"`
}

// 日志表
type Logs struct {
	gorm.Model
	Content string `json:"Content"`
}

// DNS信息表
type Dnsinfo struct {
	gorm.Model
	Name string `json:"Name"`
}

// 负载均衡表
type Lbinfo struct {
	gorm.Model
	Name string `json:"Name"`
}

// 业务线信息表
type Businessinfo struct {
	gorm.Model
	Name    string `json:"Name"`
	Manager string `json:"Manager"` //业务线负责人
}

// 默认表名是Users， 修改默认表名方法如下
func (Users) TableName() string {
	return "users"
}

func (Hosts) TableName() string {
	return "hosts"
}
