ENV JSON
========

**this package based on `github.com/gogap/env_strings`**

Read Env from config and compile the values into json file.

### About ENV JSON

Sometimes we need config file as following:

`db.conf`

```json
{
	"host":"127.0.0.1",
	"password":"3306",
	"timeout": 1000
}
```

but, when we management more and more server and serivce, and if we need change the password or ip, it was a disaster.

So, we just want use config like this.

`db.conf`

```json
{
	"host":"{{.host}}",
	"password":"{{.password}}",
	"timeout":{{.timeout}}
}
```

We use golang's template to replace values into the config while we read the json file configs.

first, we set the env config at `~/.bash_profile` or `~/.zshrc`, and the default key is `ENV_JSON` and the default file extention is `.env`, the value of `ENV_JSON` could be a file or folder,it joined by`;`, it will automatic load all `*.env` files.

**Env**

```bash
export ENV_JSON='~/playgo/test.env;~/playgo/test2.env'
```

or

```bash
export ENV_JSON='~/playgo'
```


#### example program

```go
package main

import (
	"fmt"
	"io/ioutil"

	"github.com/gogap/env_json"
)

type DBConfig struct {
	Host     string `json:"host"`
	Password string `json:"password"`
	Timeout  int64  `json:timeout`
}

func main() {
	data, _ := ioutil.ReadFile("./db.conf")

	dbConf := DBConfig{}

	if err := env_json.Unmarshal(data, &dbConf); err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println(dbConf)
}
```


`env1.json`

```json
{
	"host":"127.0.0.1",
	"password":"123456"
}
```


`env2.json`

```json
{
	"timeout":1000
}
```

**result:**

```bash
{127.0.0.1 123456 1000}
```

### More

if you want use your own `ENV` key, you could do it like this

```go
envJson := NewEnvJson("YOUR_ENV_KEY_NAME", ".json")
envJson.Unmarshal(data, &dbConf);
```
