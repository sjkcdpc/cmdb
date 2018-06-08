package model

import (
	"github.com/jinzhu/gorm"
)

type AccountState int

const (
	Active AccountState = iota
	InActive
	Banned
	BanDied
)

type Hosts struct {
	gorm.Model
	GlobalId       int64        `json:"GlobalId" gorm:"index"`
	WanIp          string       `json:"WanIp"`
	LanIp          string       `json:"LanIp"`
	Conf      	   string       `json:"Conf"`
	HostName       string       `json:"HostName"`
	Os             string       `json:"Os"`
	User1          string       `json:"User1"`
	User2          string       `json:"User2"`
	SshPort        int64        `json:"SshPort"`
	Tags           string       `json:"Tags"`
	Remark         string       `json:"Remark"`
}

