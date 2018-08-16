package env_strings

import (
	"errors"
	"fmt"
	"text/template"

	"github.com/hoisie/redis"
)

type ExtFuncsRedis struct {
	client redis.Client
	prefix string
}

func NewExtFuncsRedis(options map[string]interface{}) ExtFuncs {
	storage := new(ExtFuncsRedis)

	var addr string

	if v, exist := options["address"]; !exist {
		panic("option of address not exist")
	} else if strAddr, ok := v.(string); !ok {
		panic("option of address must be string")
	} else {
		addr = strAddr
	}

	var db float64
	if v, exist := options["db"]; !exist {
		db = 0
	} else if intDb, ok := v.(float64); !ok {
		panic("option of db must be int")
	} else {
		db = intDb
	}

	var password string
	if v, exist := options["password"]; !exist {
		password = ""
	} else if strPassword, ok := v.(string); !ok {
		panic("option of password must be string")
	} else {
		password = strPassword
	}

	var poolSize float64
	if v, exist := options["pool_size"]; !exist {
		poolSize = 0
	} else if intPoolSize, ok := v.(float64); !ok {
		panic("option of poolSize must be int")
	} else {
		poolSize = intPoolSize
	}

	var prefix string

	if v, exist := options["prefix"]; !exist {
		prefix = ""
	} else if strPrefix, ok := v.(string); !ok {
		panic("option of prefix must be string")
	} else {
		prefix = strPrefix
	}

	client := redis.Client{
		Addr:        addr,
		Db:          int(db),
		Password:    password,
		MaxPoolSize: int(poolSize),
	}

	storage.client = client
	storage.prefix = prefix

	return storage
}

func (p *ExtFuncsRedis) GetFuncs() template.FuncMap {
	funcs := make(template.FuncMap)

	funcs["redis_get"] = p.Get
	funcs["redis_hget"] = p.HGet

	return funcs
}

func (p *ExtFuncsRedis) Get(args ...interface{}) (ret interface{}, err error) {
	if len(args) < 1 {
		err = errors.New("args need 1 or 2 args")
		return
	}

	key := args[0].(string)

	if key == "" {
		err = errors.New("key could not be empty")
		return
	}

	if p.prefix != "" {
		key = p.prefix + "/" + key
	}

	if v, e := p.client.Get(key); e != nil {
		if len(args) >= 2 {
			ret = args[1]
		} else {
			err = e
		}
		return
	} else {
		ret = string(v)
	}
	return
}

func (p *ExtFuncsRedis) HGet(args ...interface{}) (ret interface{}, err error) {
	if len(args) < 2 {
		err = errors.New("args need 2 or 3 args")
		return
	}

	key := args[0].(string)

	if key == "" {
		err = errors.New("key could not be empty")
		return
	}

	if p.prefix != "" {
		key = p.prefix + "/" + key
	}

	field := args[1].(string)

	if field == "" {
		err = fmt.Errorf("field could not be empty, key: %s", key)
		return
	}

	if v, e := p.client.Hget(key, field); e != nil {
		if len(args) >= 3 {
			ret = args[2]
			return
		} else {
			err = fmt.Errorf("%s, key: %s, field: %s", e.Error(), key, field)
			return
		}
	} else {
		ret = string(v)
	}
	return
}
