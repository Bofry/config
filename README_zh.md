[EN](README.md)

config
=========

## **使用方式**

```go
// main.go
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Bofry/config"
)

func init() {
	// set env
	{
		// NOTE: you can run the following commands in bash
		// export ENVIRONMENT=production
		// export REDIS_HOST=127.0.0.3:6379
		// export REDIS_PASSWORD=1234
		// export K8S_REDIS_HOST=demo-kubernetes:6379
		// export K8S_REDIS_PASSWORD=p@ssw0rd
		// export K8S_REDIS_DB=6
		os.Clearenv()
		os.Setenv("ENVIRONMENT", "production")
		os.Setenv("REDIS_HOST", "127.0.0.3:6379")
		os.Setenv("REDIS_PASSWORD", "1234")
		os.Setenv("K8S_REDIS_HOST", "demo-kubernetes:6379")
		os.Setenv("K8S_REDIS_PASSWORD", "p@ssw0rd")
		os.Setenv("K8S_REDIS_DB", "6")
	}
	// generate .env
	{
		os.WriteFile(".env", []byte(
			strings.Join([]string{
				"REDIS_HOST=127.0.0.1:6379",
				"REDIS_DB=29",
				"TAG=demo,test",
			}, "\n")), 0644)
	}
	// generate .VERSION
	{
		os.WriteFile(".VERSION", []byte(
			strings.Join([]string{
				"v1.0.2",
			}, "\n")), 0644)
	}
	// generate config.yaml
	{
		os.WriteFile("config.yaml", []byte(
			strings.Join([]string{
				"redisDB: 3",
				"redisPoolSize: 10",
				"workspace: demo_test",
			}, "\n")), 0644)
	}
	// generate config.staging.yaml
	{
		os.WriteFile("config.staging.yaml", []byte(
			strings.Join([]string{
				"redisDB: 9",
				"redisPoolSize: 10",
				"workspace: demo_stag",
			}, "\n")), 0644)
	}
	// generate config.production.yaml
	{
		os.WriteFile("config.production.yaml", []byte(
			strings.Join([]string{
				"redisDB: 12",
				"redisPoolSize: 50",
				"workspace: demo_prod",
			}, "\n")), 0644)
	}
}

type DummyConfig struct {
	RedisHost     string   `env:"REDIS_HOST"       yaml:"redisHost"       arg:"redis-host;the Redis server address and port"`
	RedisPassword string   `env:"REDIS_PASSWORD"   yaml:"redisPassword"   arg:"redis-passowrd;the Redis password"`
	RedisDB       int      `env:"REDIS_DB"         yaml:"redisDB"         arg:"redis-db;the Redis database number"`
	RedisPoolSize int      `env:"-"                yaml:"redisPoolSize"`
	Workspace     string   `env:"-"                yaml:"workspace"       arg:"workspace;the data workspace"`
	Tags          []string `env:"TAG"`
	Version       string   `resource:".VERSION"`
}

func main() {
	conf := DummyConfig{}

	config.NewConfigurationService(&conf).
		LoadDotEnv().
		LoadEnvironmentVariables("").
		LoadEnvironmentVariables("K8S").
		LoadYamlFile("config.yaml").
		LoadYamlFile("config.${ENVIRONMENT}.yaml").
		LoadCommandArguments().
		LoadResource("")
	fmt.Printf("RedisHost     = %q\n", conf.RedisHost)
	fmt.Printf("RedisPassword = %q\n", conf.RedisPassword)
	fmt.Printf("RedisDB       = %d\n", conf.RedisDB)
	fmt.Printf("RedisPoolSize = %d\n", conf.RedisPoolSize)
	fmt.Printf("Workspace     = %q\n", conf.Workspace)
	fmt.Printf("Tags          = %q\n", conf.Tags)
	fmt.Printf("Version       = %q\n", conf.Version)
}
```
打開終端機，執行下面命令：
- Bash
	```bash
	$ go build -o example
	$ ./example -redis-db=32
	```
- Dos
	```dos
	C:\> go build -o example.exe
	C:\> example.exe -redis-db=32
	```
執行結果：
```
RedisHost     = "demo-kubernetes:6379"
RedisPassword = "p@ssw0rd"
RedisDB       = 32
RedisPoolSize = 50
Workspace     = "demo_prod"
Tags          = ["demo" "test"]
Version       = "v1.0.2"
```


$~$
## **Struct Tag 標記**

| 適用配置類型  | struct tag | tag flags  | 範例    |
|:-------------|:-----------|:-----------|:--------|
| 環境變數     | `env`      | *required* | `env:"CACHE_ADDRESS,required"` -或- `env:"*CACHE_ADDRESS"`       |
| .env 檔案    | `env`      | *required* | `env:"CACHE_ADDRESS,required"` -或- `env:"*CACHE_ADDRESS"`       |
| json 檔案    | `json`     | --         | `json:"LISTEN_PORT"`                                             |
| yaml 檔案    | `yaml`     | --         | `yaml:"LISTEN_PORT"`                                             |
| 二進制檔案   | `resource` | *required* | `resource:"VERSION,required"` -或- `resource:"*VERSION"`          |
| 文字檔案     | `resource` | *required* | `resource:"VERSION,required"` -或- `resource:"*VERSION"`          |
| 命令列參數   | `arg`      | --         | `arg:"SERVER_NAME"` -或- `arg:"SERVER_NAME;specify server name"`  |

> 📝 `resource:"VERSION,required"` 與 `resource:"*VERSION"` 是相同的，而 `resource:"*VERSION,required"` 則與前兩者不同。下面是舉例比較：
> | 標記                             | name     | flag       |
> |:---------------------------------|:---------|:-----------|
> | `resource:"VERSION,required"`    | VERSION  | `required` |
> | `resource:"*VERSION"`            | VERSION  | `required` |
> | `resource:"*VERSION,required"`   | *VERSION | `required` |
> | `resource:"*VERSION,required,_"` | *VERSION | `required` |
> | `resource:"*VERSION,_"`          | *VERSION | *none*     |
> | `resource:"VERSION,_"`           | VERSION  | *none*     |
> 
> 📝 若名稱需要保留開始的 "`*`" 且維持非必填指示，可以在 flaf 段加入空白 flag "`_`"。
> 
> 📝 **env**、**resource**、**arg** 等標記**不支援巢狀結構**，支援的欄位型別為：`bool`、`int`、`uint`、`float`、`string`、`time.Duration`、`time.Time`、`url.URL`、`net.IP`、`[]bool`、`[]int`、`[]uint`、`[]float`、`[]string`、`[]time.Duration`、`[]time.Time`、`[]url.URL`、`[]net.IP`、`bytes.Buffer`、`json.RawMessage`、`github.com/Bofry/types.RawContent`。


$~$
### **環境變數**
⠿ 下面的 **Config** 結構將匯入 `CACHE_HOST`、`CACHE_PASSWORD` 與 `CACHE_DB` 環境變數。其中標記 `env:"CACHE_HOST,required"` 的項目設定了 *required* 旗標，指示 `CACHE_HOST` 環境變數是必要的，若找不到則會抛出例外。
```go
type Config struct {
  CacheHost     string `env:"CACHE_HOST,required"`
  CachePassword string `env:"CACHE_PASSWORD"`
  CacheDB       int    `env:"CACHE_DB"`
}
```
`env:"CACHE_HOST,required"` 標記方式能夠轉換為 `env:"*CACHE_HOST"` 表示。下面的表示方式與前者相同。
```go
type Config struct {
  CacheHost     string `env:"*CACHE_HOST"`
  CachePassword string `env:"CACHE_PASSWORD"`
  CacheDB       int    `env:"CACHE_DB"`
}
```


$~$
### **.env 檔案**
⠿ .env 檔案使用方式同 **環境變數**。
> 📝 .env 檔案**不會覆寫已經存在的環境變數**。適合用來作為開發階段使用，或是提供有意義的預設值。


$~$
### **資源檔**
⠿ 下面的 **Config** 結構將匯入 **VERSION** 檔案的內容。其中標記 `resource:"VERSION,required"` 的項目設定了 *required* 旗標，指示 **VERSION** 檔案是必要的，若找不到則會抛出例外。
```go
type Config struct {
  AppVersion string `resource:"VERSION,required"`
}
```
`resource:"VERSION,required"` 標記方式能夠轉換為 `resource:"*VERSION"` 表示。下面的表示方式與前者相同。 
```go
type Config struct {
  AppVersion string `resource:"*VERSION"`
}
```
> 📝 資源名稱接受任何 unicode 字元，但不能使用空白字元作為開頭與結尾、以及結尾不能是 "`.`"。


$~$
### **命令列參數**
⠿ 下面的 **Config** 結構將匯入命令列參數 `cache-host`、`cache-passowrd` 與 `cache-db`。其中 `arg:"cache-host;the cache server address and port"` 標記使用分號 "`;`" 連接名稱部份與使用說明部份；使用說明可以在啟動命令傳入 `-help` 輸出。
```go
type Config struct {
	CacheHost     string `arg:"cache-host;the cache server address and port"`
	CachePassword string `arg:"cache-passowrd;the cache server password"`
	CacheDB       int    `arg:"cache-db;the cache database number"`
}
```

> ⛔ 不要使用 `help` 作為參數名稱。  


$~$
## **相依套件**
- Yaml - https://godoc.org/gopkg.in/yaml.v2
- Json - https://golang.org/pkg/encoding/json/
- dotenv - https://github.com/joho/godotenv
