package env_strings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	ENV_STRINGS_KEY = "ENV_STRINGS"
	ENV_STRINGS_EXT = ".env"

	ENV_STRINGS_CONF = "/etc/env_strings.conf"

	ENV_STRINGS_CONFIG_KEY = "ENV_STRINGS_CONF"
)

const (
	STORAGE_REDIS = "redis"
)

type EnvStringConfig struct {
	Storages []StorageConfig `json:"storages"`
}

type StorageConfig struct {
	Engine  string                 `json:"engine"`
	Options map[string]interface{} `json:"options"`
}

type option func(envStrings *EnvStrings)

type EnvStrings struct {
	envName   string
	envExt    string
	tmplFuncs *TemplateFuncs

	configFile string

	envConfig EnvStringConfig
}

func FuncMap(name string, function interface{}) option {
	return func(e *EnvStrings) {
		e.RegisterFunc(name, function)
	}
}

func EnvStringsConfig(fileName string) option {
	return func(e *EnvStrings) {
		e.configFile = fileName
	}
}

func NewEnvStrings(envName string, envExt string, opts ...option) *EnvStrings {
	if envName == "" {
		panic("env_strings: env name could not be empty")
	}

	envStrings := &EnvStrings{
		envName:    envName,
		envExt:     envExt,
		configFile: ENV_STRINGS_CONF,
		tmplFuncs:  NewTemplateFuncs(),
	}

	if opts != nil && len(opts) > 0 {
		for _, opt := range opts {
			opt(envStrings)
		}
	}

	envStringsConf := os.Getenv(ENV_STRINGS_CONFIG_KEY)
	if envStringsConf != "" {
		envStrings.configFile = envStringsConf
	}

	if envStrings.configFile != "" {
		if err := envStrings.loadConfig(envStrings.configFile); err != nil {
			if !os.IsNotExist(err) {
				panic(err)
			} else {
				return envStrings
			}
		}

		if envStrings.envConfig.Storages != nil {
			for _, storageConf := range envStrings.envConfig.Storages {
				switch storageConf.Engine {
				case STORAGE_REDIS:
					{
						extFucnRedis := NewExtFuncsRedis(storageConf.Options)
						redisFuncs := extFucnRedis.GetFuncs()

						if redisFuncs == nil {
							panic("ext funcs of redis is nil")
						}

						for funcName, fn := range redisFuncs {
							if err := envStrings.RegisterFunc(funcName, fn); err != nil {
								panic(err)
							}
						}
					}
				default:
					{
						panic("unknown storage type")
					}
				}
			}
		}
	}
	return envStrings
}

func (p *EnvStrings) Execute(str string) (ret string, err error) {
	return p.ExecuteWith(str, nil)
}

func (p *EnvStrings) ExecuteWith(str string, envValues map[string]interface{}) (ret string, err error) {
	strEnvFiles := os.Getenv(p.envName)

	envFiles := strings.Split(strEnvFiles, ";")

	if envValues == nil {
		envValues = make(map[string]interface{})
	}

	debug := false
	if os.Getenv("ENV_STRINGS_DEBUG") == "true" {
		debug = true
	}

	for _, envFile := range envFiles {

		if len(envFile) == 0 {
			continue
		}

		prefix := filepath.Base(envFile)

		prefix = filepath.Dir(prefix)

		err = p.loadEnv(debug, prefix, []string{envFile}, envValues)

		if err != nil {
			return
		}
	}

	if debug {
		debugData, _ := json.MarshalIndent(envValues, "", "    ")
		fmt.Printf("[ENV_STRINGS] final envs:\n%s\n", string(debugData))
	}

	var tpl *template.Template

	if tpl, err = template.New("tmpl:" + p.envName).Funcs(p.tmplFuncs.GetFuncMaps(p.envName)).Option("missingkey=error").Parse(str); err != nil {
		return
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, envValues); err != nil {
		return
	}

	ret = buf.String()

	if debug {
		fmt.Printf("[ENV_STRINGS] final rendered:\n%s\n", ret)
	}

	return
}

func Execute(str string) (ret string, err error) {
	envStrings := NewEnvStrings(ENV_STRINGS_KEY, ENV_STRINGS_EXT)
	return envStrings.Execute(str)
}

func ExecuteWith(str string, envValues map[string]interface{}) (ret string, err error) {
	envStrings := NewEnvStrings(ENV_STRINGS_KEY, ENV_STRINGS_EXT)
	return envStrings.ExecuteWith(str, envValues)
}

func (p *EnvStrings) RegisterFunc(name string, function interface{}) (err error) {
	return p.tmplFuncs.Register(name, function)
}

func (p *EnvStrings) FuncUsageStatic() map[string][]FuncStaticItem {
	return funcStatics
}

func (p *EnvStrings) loadEnv(debug bool, prefix string, files []string, envs map[string]interface{}) (err error) {
	for _, path := range files {

		var fi os.FileInfo
		fi, err = os.Stat(path)
		if err != nil {
			return
		}

		if strings.HasPrefix(fi.Name(), ".") {
			continue
		}

		if fi.IsDir() {

			baseName := filepath.Base(path)

			if strings.HasPrefix(baseName, ".") {
				continue
			}

			var fis []os.FileInfo
			fis, err = ioutil.ReadDir(path)

			if err != nil {
				return
			}

			var nextfiles []string

			for _, f := range fis {
				nextfiles = append(nextfiles, filepath.Join(path, f.Name()))
			}

			var nextENVs map[string]interface{}

			if envs == nil {
				envs = make(map[string]interface{})
			}

			preEnvs, exist := envs[baseName]
			if !exist {
				nextENVs = make(map[string]interface{})
				envs[baseName] = nextENVs
			} else {
				nextENVs = preEnvs.(map[string]interface{})
			}

			err = p.loadEnv(debug, prefix, nextfiles, nextENVs)
			if err != nil {
				return
			}

			continue
		}

		baseName := strings.TrimSuffix(fi.Name(), p.envExt)

		var fileEnvs map[string]interface{}
		fileEnvs, err = p.loadEnvFile(path)

		if err != nil {
			return err
		}

		if debug {
			debugData, _ := json.MarshalIndent(fileEnvs, "", "    ")
			fmt.Printf("[ENV_STRINGS] env file: %s\n%s\n", path, string(debugData))
		}

		if envs == nil {
			envs = make(map[string]interface{})
		}

		existEnvs, exist := envs[baseName]

		if exist {
			err = mergeMaps(existEnvs, fileEnvs)
			if err != nil {
				err = fmt.Errorf("merge same env key's values failure, file: %s, error: %s", path, err.Error())
				return
			}

		} else {
			envs[baseName] = fileEnvs
		}
	}

	return
}

func mergeMaps(vA, vB interface{}) (err error) {
	vAMap, okA := vA.(map[string]interface{})
	vBMap, okB := vB.(map[string]interface{})

	if okA && okB {
		for k, valB := range vBMap {
			valA, exist := vAMap[k]
			if !exist {
				vAMap[k] = valB
			} else {
				err = mergeMaps(valA, valB)
				if err != nil {
					return
				}
			}
		}

		return
	}

	err = fmt.Errorf("could not merge different struct values")

	return
}

func (p *EnvStrings) loadEnvFile(filename string) (ret map[string]interface{}, err error) {

	if ext := filepath.Ext(filename); ext == p.envExt {

		var data []byte
		data, err = ioutil.ReadFile(filename)
		if err != nil {
			return
		}

		r := make(map[string]interface{})

		err = json.Unmarshal(data, &r)
		if err != nil {
			return
		}

		ret = r
	}

	return
}

func (p *EnvStrings) loadConfig(fileName string) (err error) {
	if _, err = os.Stat(fileName); err != nil {
		return
	}

	var data []byte
	if data, err = ioutil.ReadFile(fileName); err != nil {
		return
	}

	conf := EnvStringConfig{}

	if err = json.Unmarshal(data, &conf); err != nil {
		return
	}

	p.envConfig = conf

	return
}
