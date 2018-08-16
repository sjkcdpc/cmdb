package env_json

import (
	"encoding/json"
	"github.com/gogap/env_strings"
)

const (
	ENV_JSON_KEY = "ENV_JSON"
	ENV_JSON_EXT = ".env"
)

type EnvJson struct {
	*env_strings.EnvStrings
}

func NewEnvJson(envName string, envExt string) *EnvJson {
	if envName == "" {
		panic("env_json: env name could not be nil")
	}

	return &EnvJson{
		EnvStrings: env_strings.NewEnvStrings(envName, envExt),
	}
}

func (p *EnvJson) Marshal(v interface{}) (data []byte, err error) {
	return json.Marshal(v)
}

func (p *EnvJson) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func (p *EnvJson) Unmarshal(data []byte, v interface{}) (err error) {
	strData := ""
	if strData, err = p.Execute(string(data)); err != nil {
		return
	}

	err = json.Unmarshal([]byte(strData), v)

	return
}

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func Unmarshal(data []byte, v interface{}) error {
	envJson := NewEnvJson(ENV_JSON_KEY, ENV_JSON_EXT)
	return envJson.Unmarshal(data, v)
}
