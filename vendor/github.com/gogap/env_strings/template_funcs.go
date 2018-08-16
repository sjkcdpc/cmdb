package env_strings

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"text/template"
	"time"
)

var (
	errorType = reflect.TypeOf((*error)(nil)).Elem()

	funcStatics  = make(map[string][]FuncStaticItem)
	staticLocker sync.Mutex
)

type FuncStaticItem struct {
	EnvName  string
	FuncName string
	Input    []interface{}
	Output   []interface{}
}

type TemplateFuncs struct {
	funcMap template.FuncMap
}

func NewTemplateFuncs() *TemplateFuncs {
	return &TemplateFuncs{
		funcMap: basicFuncs(),
	}
}

func (p *TemplateFuncs) GetFuncMaps(envName string) template.FuncMap {
	return p.hookedFuncs(envName)
}

func (p *TemplateFuncs) Register(name string, function interface{}) (err error) {
	if function == nil {
		err = errors.New("function could not be nil")
		return
	} else if funcType := reflect.TypeOf(function); funcType.Kind() != reflect.Func {
		err = errors.New("function is not a Func kind")
		return
	} else if name == "" {
		err = errors.New("name could not be empty")
		return
	}

	if _, exist := p.funcMap[name]; exist {
		panic("func name of " + name + " already exist")
	}

	p.funcMap[name] = function

	return
}

func basicFuncs() template.FuncMap {
	m := make(template.FuncMap)
	m["base"] = path.Base
	m["abs"] = filepath.Abs
	m["dir"] = filepath.Dir
	m["pwd"] = os.Getwd
	m["split"] = strings.Split
	m["json"] = UnmarshalJsonObject
	m["jsonArray"] = UnmarshalJsonArray
	m["dir"] = path.Dir
	m["getenv"] = os.Getenv
	m["join"] = strings.Join
	m["localtime"] = time.Now
	m["utc"] = time.Now().UTC
	m["pid"] = os.Getpid
	m["httpGet"] = httpGet
	m["envIfElse"] = envIfElse
	m["md5"] = fmd5
	m["base64"] = fbase64

	return m
}

func (p *TemplateFuncs) hookFunc(envName string, fn interface{}) interface{} {
	return func(args ...interface{}) (val interface{}, err error) {
		return call(fn, args)
	}
}

func (p *TemplateFuncs) hookedFuncs(envName string) template.FuncMap {
	hookedFuncs := template.FuncMap{}

	for fName, originalFunc := range p.funcMap {
		hookedFuncs[fName] = func(eName, funcName string, fn interface{}) interface{} {
			return func(args ...interface{}) (ret interface{}, err error) {
				ret, err = call(fn, args...)

				if ret == nil && err == nil {
					err = errors.New(fmt.Sprintf("the func of %s in env %s get <no value>", funcName, eName))
				}

				staticLocker.Lock()
				defer staticLocker.Unlock()

				staticItem := FuncStaticItem{
					EnvName:  eName,
					FuncName: funcName,
					Input:    args,
					Output:   []interface{}{ret, err},
				}

				if items, exist := funcStatics[envName]; exist {
					items = append(items, staticItem)
					funcStatics[envName] = items
				} else {
					funcStatics[envName] = []FuncStaticItem{staticItem}
				}

				return
			}
		}(envName, fName, originalFunc)
	}

	return hookedFuncs
}

func UnmarshalJsonObject(data string) (map[string]interface{}, error) {
	var ret map[string]interface{}
	err := json.Unmarshal([]byte(data), &ret)
	return ret, err
}

func UnmarshalJsonArray(data string) ([]interface{}, error) {
	var ret []interface{}
	err := json.Unmarshal([]byte(data), &ret)
	return ret, err
}

func httpGet(url string) (ret string, err error) {
	var resp *http.Response
	if resp, err = http.Get(url); err != nil {
		return
	}

	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	ret = string(body)

	return
}

func envIfElse(envName, equalValue, trueValue, falseValue string) string {
	if os.Getenv(envName) == equalValue {
		return trueValue
	} else {
		return falseValue
	}
}

func fbase64(command string, src string) (val string, err error) {
	if command == "encode" {
		val = base64.StdEncoding.EncodeToString([]byte(src))
		return
	} else if command == "decode" {
		var data []byte
		if data, err = base64.StdEncoding.DecodeString(src); err != nil {
			return
		}
		val = string(data)
		return
	}
	err = errors.New("unknown command :" + command)
	return
}

func fmd5(str string) string {
	v := md5.Sum([]byte(str))
	return fmt.Sprintf("%0x", v[:])
}

// The following code copy from pkg "text/template"

// call returns the result of evaluating the first argument as a function.
// The function must return 1 result, or 2 results, the second of which is an error.
func call(fn interface{}, args ...interface{}) (interface{}, error) {
	v := reflect.ValueOf(fn)
	typ := v.Type()
	if typ.Kind() != reflect.Func {
		return nil, fmt.Errorf("non-function of type %s", typ)
	}
	if !goodFunc(typ) {
		return nil, fmt.Errorf("function called with %d args; should be 1 or 2", typ.NumOut())
	}
	numIn := typ.NumIn()
	var dddType reflect.Type
	if typ.IsVariadic() {
		if len(args) < numIn-1 {
			return nil, fmt.Errorf("wrong number of args: got %d want at least %d", len(args), numIn-1)
		}
		dddType = typ.In(numIn - 1).Elem()
	} else {
		if len(args) != numIn {
			return nil, fmt.Errorf("wrong number of args: got %d want %d", len(args), numIn)
		}
	}
	argv := make([]reflect.Value, len(args))
	for i, arg := range args {
		value := reflect.ValueOf(arg)
		// Compute the expected type. Clumsy because of variadics.
		var argType reflect.Type
		if !typ.IsVariadic() || i < numIn-1 {
			argType = typ.In(i)
		} else {
			argType = dddType
		}
		if !value.IsValid() && canBeNil(argType) {
			value = reflect.Zero(argType)
		}
		if !value.Type().AssignableTo(argType) {
			return nil, fmt.Errorf("arg %d has type %s; should be %s", i, value.Type(), argType)
		}
		argv[i] = value
	}
	result := v.Call(argv)
	if len(result) == 2 && !result[1].IsNil() {
		return result[0].Interface(), result[1].Interface().(error)
	}
	return result[0].Interface(), nil
}

// canBeNil reports whether an untyped nil can be assigned to the type. See reflect.Zero.
func canBeNil(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return true
	}
	return false
}

// goodFunc checks that the function or method has the right result signature.
func goodFunc(typ reflect.Type) bool {
	// We allow functions with 1 result or 2 results where the second is an error.
	switch {
	case typ.NumOut() == 1:
		return true
	case typ.NumOut() == 2 && typ.Out(1) == errorType:
		return true
	}
	return false
}
