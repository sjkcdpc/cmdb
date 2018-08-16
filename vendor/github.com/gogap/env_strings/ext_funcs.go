package env_strings

import (
	"text/template"
)

type ExtFunc func(args ...interface{}) (ret interface{}, err error)

type ExtFuncs interface {
	GetFuncs() template.FuncMap
}
