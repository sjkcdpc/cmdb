ENV STRINGS
===========

Read Env from config and compile the values into string.

### About ENV JSON

Sometimes we need use some config string as following:

`db.conf`

```go
configStr:="127.0.0.1|123456|1000"
```

but, when we management more and more server and serivce, and if we need change the password or ip, it was a disaster.

So, we just want use config string like this.

`db.conf`

```json
configStr:="{{.host}}|{{.password}}|{{.timeout}}"
```

We use golang's template to replace values into the config while we execute the string.

first, we set the env config at `~/.bash_profile` or `~/.zshrc`, and the default key is `ENV_STRINGS` and the default file extention is `.env`, the value of `ENV_STRINGS ` could be a file or folder,it joined by`;`, it will automatic load all `*.env` files.

**Env**

```bash
export ENV_STRINGS ='~/playgo/test.env;~/playgo/test2.env'
```

or

```bash
export ENV_STRINGS ='~/playgo'
```


#### example program

```go
package main

import (
	"fmt"

	"github.com/gogap/env_strings"
)

func main() {
	if ret, err := env_strings.Execute("{{.host}}|{{.password}}|{{.timeout}}"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ret)
	}

	envStrings := env_strings.NewEnvStrings("ENV_KEY", ".env")

	if ret, err := envStrings.Execute("{{.host}}|{{.password}}|{{.timeout}}"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ret)
	}
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


### Advance

#### use redis to storage conf

we should set the config file of `/etc/env_strings.conf` (default path) before use, or set config file path by system ENV with your '~/.bash_profile' or `~/.zshrc`, the key is `ENV_STRINGS_CONF`, just edit like this: `export ENV_STRINGS_CONF=/etc/env_strings.conf`


`/etc/env_strings.conf`

```json
{
    "storages": [{
        "engine": "redis",
        "options": {
            "db": 1,
            "password": "",
            "pool_size": 10,
            "address": "localhost:6379"
        }
    }]
}
```

#### set data to redis

```bash
> redis-cli -n 1 set name gogap
OK
> redis-cli -n 1 get name
gogap
> redis-cli -n 1 set key field value
OK
> redis-cli -n 1 hget key field
value
```


#### redis command: redis_get

**config template:**

```json
{
	"name":"{{redis_get "name"}}"
}
```

**result:**

```bash
{
	"name":"gogap"
}
```

if key not exist, and we want get a default value

```json
{
	"name":"{{redis_get "noexistkey" "gogap"}}"
}
```

**result:**

```bash
{
	"name":"gogap"
}
```


#### redis command: redis_hget

**config template:**

```json
{
	"name":"{{redis_hget "key" "field"}}"
}
```

**result:**

```bash
{
	"name":"value"
}
```

if key or filed not exist, and we want get a default value

```json
{
	"name":"{{redis_hget "key2" "filed" "gogap"}}"
}
```

**result:**

```bash
{
	"name":"gogap"
}
```