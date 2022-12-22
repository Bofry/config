[ZH](README_zh.md)

config
=========

## **Synopsis**

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
Open your terminal and execute the following command:
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
You will get:
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
## **Struct Tag Denotation**

| configuration type    | struct tag | tag flags  | ConfigurationService method    | example |
|:----------------------|:-----------|:-----------|:-------------------------------|:--------|
| environment variables | `env`      | *required* | LoadEnvironmentVariables()     | `env:"CACHE_ADDRESS,required"` -or- `env:"*CACHE_ADDRESS"`
| .env files            | `env`      | *required* | LoadDotEnv(), LoadDotEnvFile() | `env:"CACHE_ADDRESS,required"` -or- `env:"*CACHE_ADDRESS"`
| json files            | `json`     | --         | LoadJsonFile()                 | `json:"LISTEN_PORT"`
| yaml files            | `yaml`     | --         | LoadYamlFile()                 | `yaml:"LISTEN_PORT"`
| binary reource files  | `resource` | *required* | LoadResource()                 | `resource:"VERSION,required"` -or- `resource:"*VERSION"`
| text reource files    | `resource` | *required* | LoadResource()                 | `resource:"VERSION,required"` -or- `resource:"*VERSION"`
| command arguments     | `arg`      | --         | LoadCommandArguments()         | `arg:"SERVER_NAME"` -or- `arg:"SERVER_NAME;specify server name"`


$~$
### **Environment Variables**
⠿ The following **Config** structure will import environment variables `CACHE_HOST`, `CACHE_PASSWORD`, and `CACHE_DB`. The tag text `env:"CACHE_HOST,required"` use the flag *required* indicates the environment variable `CACHE_HOST` is required. It will get exception if the variable doesn't be assgined.
```go
type Config struct {
  CacheHost     string `env:"CACHE_HOST,required"`
  CachePassword string `env:"CACHE_PASSWORD"`
  CacheDB       int    `env:"CACHE_DB"`
}
```
The tag text `env:"CACHE_HOST,required"` can be switch as `env:"*CACHE_HOST"` as well. Put the symbol "`*`" in front of the name is equivalent to appending `required` to tag flag part. 
```go
type Config struct {
  CacheHost     string `env:"*CACHE_HOST"`
  CachePassword string `env:"CACHE_PASSWORD"`
  CacheDB       int    `env:"CACHE_DB"`
}
```


$~$
### **.env Files**
⠿ The .env files same as **Environment Variables**.
> 📝 The .env file WILL NOT OVERRIDE an environment variable that already exists. To consider .env file to set dev variable or sensible defaults.


$~$
### **Resource Files**
⠿ The following **Config** structure will import content from file **VERSION**. The tag text `resource:"VERSION,required"` use the flag *required* indicates the file **VERSION** is required. It will get exception if the file doesn't exist.
```go
type Config struct {
  AppVersion string `resource:"VERSION,required"`
}
```
The tag text `resource:"VERSION,required"` can be switch as `resource:"*VERSION"` as well. Put the symbol "`*`" in front of the name is equivalent to appending `required` to tag flag part. 
```go
type Config struct {
  AppVersion string `resource:"*VERSION"`
}
```
> 📝 The name can compose by any unicode except `NUL`, `\`, `/`, `:`, `*`, `?`, `"`, `<`, `>`, `|`. Also, no space character at the start or end, and no period at the end.


$~$
### **Command Arguments**
⠿ The following **Config** structure will import command arguments `cache-host`, `cache-passowrd`, and `cache-db`. The tag text `arg:"cache-host;the cache server address and port"` separated by symbol "`;`" to two parts. The name part and the usage text part for help.
```go
type Config struct {
	CacheHost     string `arg:"cache-host;the cache server address and port"`
	CachePassword string `arg:"cache-passowrd;the cache server password"`
	CacheDB       int    `arg:"cache-db;the cache database number"`
}
```

> 📝 The name can compose by `A-Z a-z 0-9 _ -`.
> 
> ⛔ Don't name arg as `help`.  


$~$
## **Dependency**
- Yaml - https://godoc.org/gopkg.in/yaml.v2
- Json - https://golang.org/pkg/encoding/json/
- dotenv - https://github.com/joho/godotenv
