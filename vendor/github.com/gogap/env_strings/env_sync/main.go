package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/gogap/env_strings"
	"github.com/hoisie/redis"
)

var (
	client redis.Client
)

type redisConfig struct {
	db       int32
	password string
	poolSize int32
	address  string
	prefix   string
}

type syncData struct {
	key   string
	field string
	value string
}

func main() {
	config, err := getRedisConfig()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	client = redis.Client{
		Addr:        config.address,
		Db:          int(config.db),
		Password:    config.password,
		MaxPoolSize: 3,
	}
	data, err := prepare()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = set(data)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return
}

func getRedisConfig() (config redisConfig, err error) {
	path := os.Getenv(env_strings.ENV_STRINGS_CONFIG_KEY)
	if path == "" {
		path = env_strings.ENV_STRINGS_CONF
	}

	f, err := os.Open(path)
	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	var storageConfig env_strings.EnvStringConfig
	err = json.Unmarshal(data, &storageConfig)
	if err != nil {
		return
	}
	for _, storage := range storageConfig.Storages {
		if storage.Engine == env_strings.STORAGE_REDIS {
			if v, exist := storage.Options["address"]; !exist {
				err = errors.New("option of address not exist")
				return
			} else if strAddr, ok := v.(string); !ok {
				err = errors.New("option of address must be string")
				return
			} else {
				config.address = strAddr
			}

			if v, exist := storage.Options["db"]; !exist {
				config.db = 0
			} else if intDb, ok := v.(float64); !ok {
				err = errors.New("option of db must be int")
				return
			} else {
				config.db = int32(intDb)
			}

			if v, exist := storage.Options["password"]; !exist {
				config.password = ""
			} else if strPassword, ok := v.(string); !ok {
				err = errors.New("option of password must be string")
				return
			} else {
				config.password = strPassword
			}

			if v, exist := storage.Options["pool_size"]; !exist {
				config.poolSize = 0
			} else if intPoolSize, ok := v.(float64); !ok {
				err = errors.New("option of poolSize must be int")
				return
			} else {
				config.poolSize = int32(intPoolSize)
			}

			if v, exist := storage.Options["prefix"]; !exist {
				config.prefix = ""
			} else if strPrefix, ok := v.(string); !ok {
				err = errors.New("option of prefix must be string")
				return
			} else {
				config.prefix = strPrefix
			}
		} else {
			err = errors.New("storage only support redis now")
		}
	}
	return
}

func prepare() (sycnDatas []syncData, err error) {
	workDir, err := os.Getwd()
	if err != nil {
		return
	}
	fnWalk := func(path string, info os.FileInfo, e error) (err error) {
		if info.Name() != "data" || info.IsDir() {
			return
		}
		datafile, _ := filepath.Rel(workDir, path)
		if datafile == "." {
			return
		}
		datafileDir := filepath.Dir(datafile)

		var data []byte
		if data, err = ioutil.ReadFile(datafile); err != nil {
			return
		}
		dataKV := map[string]interface{}{}

		if err = json.Unmarshal(data, &dataKV); err != nil {
			return
		}

		for k, v := range dataKV {
			strv, err := serializeObject(v)
			if err != nil {
				return err
			}
			if datafileDir == "." {
				//SET
				sycnDatas = append(sycnDatas, syncData{key: k, field: "", value: strv})
			} else {
				//HSET
				sycnDatas = append(sycnDatas, syncData{key: datafileDir, field: k, value: strv})
			}

		}
		return
	}
	if err = filepath.Walk(workDir, fnWalk); err != nil {
		return
	}
	return
}
func serializeObject(obj interface{}) (str string, err error) {
	switch reflect.TypeOf(obj).Kind() {
	case reflect.Map, reflect.Array, reflect.Slice:
		{
			data, err := json.Marshal(&obj)
			if err != nil {
				return str, err
			}
			str = string(data)
		}
	default:
		{
			str = fmt.Sprintf("%v", obj)
		}
	}
	return
}

func set(data []syncData) (err error) {
	for _, v := range data {
		if v.field == "" {
			origin, e := client.Get(v.key)
			if e == nil && string(origin) == v.value {
				continue
			}
			err = client.Set(v.key, []byte(v.value))
			if err != nil {
				return
			}
		} else {
			origin, e := client.Hget(v.key, v.field)
			if e == nil && string(origin) == v.value {
				continue
			}
			_, err = client.Hset(v.key, v.field, []byte(v.value))
			if err != nil {
				return
			}
		}
	}
	return
}
