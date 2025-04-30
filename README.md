[ZH](README_zh.md)

# config

A highly customizable Go configuration management library that supports multiple data sources.

## Key Features

- Support for multiple configuration sources: environment variables, .env files, YAML/JSON files, command-line arguments, and resource files
- Fluent API design for clean and readable configuration code
- Intuitive configuration mapping through struct tags
- Automatic data type conversion
- Flexible configuration priority handling
- Support for custom output formats

## Usage Example

```go
package main

import (
 "fmt"
 "github.com/Bofry/config"
 "github.com/Bofry/structproto"
 "os"
 "reflect"
)

type (
 ServiceConfig struct {
  Environment string `env:"Environment"`

  // Core service information
  Version     string `resource:".VERSION"`
  Signature   string `resource:".SIGNATURE"`
  ServiceName string `resource:".SERVICE_NAME"`

  // HTTP server settings
  ListenAddress  string `yaml:"ListenAddress"  arg:"listen-address;the combination of IP address and listen port"`
  EnableCompress bool   `yaml:"UseCompress"    arg:"use-compress;indicates the response enable compress or not"`
  ServerName     string `yaml:"ServerName"`

  // Telemetry settings
  JaegerTraceUrl string `yaml:"JaegerTraceUrl"`

  // External services
  CacheStoreIp   string `yaml:"Cache_Store_Ip"        env:"Cache_Store_Ip"`
  CacheStorePort int    `yaml:"Cache_Store_Port"      env:"Cache_Store_Port"`
  APIEndpoint    string `yaml:"API_Endpoint"          env:"API_Endpoint"`
  MessageBroker  string `yaml:"Message_Broker_Address" env:"Message_Broker_Address"`
 }
)

func NewConfiguration() *ServiceConfig {
 conf := &ServiceConfig{}
 config.NewConfigurationService(conf).
  LoadYamlFile("config.yaml").
  LoadYamlFile("config.${Environment}.yaml").
  LoadEnvironmentVariables("").
  LoadResource(".").
  LoadResource(".conf/${Environment}").
  LoadCommandArguments().
  Map(func(field structproto.FieldInfo, rv reflect.Value) error {
   switch rv.Kind() {
   case reflect.String:
    if !rv.IsZero() {
     var hasEmpty bool = false
     val := os.Expand(rv.String(), func(s string) string {
      v := os.Getenv(s)
      if len(v) == 0 {
       hasEmpty = true
      }
      return v
     })
     if hasEmpty {
      rv.SetString("")
     } else {
      rv.SetString(val)
     }
    }
   }
   return nil
  })

 return conf
}

func main() {
 // Initialize configuration
 config := NewConfiguration()
 
 // Use configuration values
 fmt.Printf("Service: %s v%s\n", config.ServiceName, config.Version)
 fmt.Printf("Listening on: %s\n", config.ListenAddress)
 fmt.Printf("Cache Store: %s:%d\n", config.CacheStoreIp, config.CacheStorePort)
 
 // Start your application...
}
```

## Quick Start

### Installing the Package

First, install the config package using the `go get` command:

```bash
go get -u github.com/Bofry/config
```

### Running Your Application

After creating your application, open your terminal and execute:

```bash
# Run directly with command-line arguments
go run main.go -listen-address=":8080"

# Or build and run
go build
./yourapp -listen-address=":8080"
```

You can pass any command-line argument that you've defined with `arg` tags in your configuration structure.

## Struct Tag Syntax

| Configuration Type    | Struct Tag | Tag Flags  | Configuration Service Method | Example                                                            |
| :-------------------- | :--------- | :--------- | :--------------------------- | :----------------------------------------------------------------- |
| Environment Variables | `env`      | *required* | LoadEnvironmentVariables()   | `env:"CACHE_ADDRESS,required"` or `env:"*CACHE_ADDRESS"`           |
| .env Files            | `env`      | *required* | LoadDotEnv()                 | `env:"CACHE_ADDRESS,required"` or `env:"*CACHE_ADDRESS"`           |
| JSON Files            | `json`     | --         | LoadJsonFile()               | `json:"LISTEN_PORT"`                                               |
| YAML Files            | `yaml`     | --         | LoadYamlFile()               | `yaml:"LISTEN_PORT"`                                               |
| Binary Resource Files | `resource` | *required* | LoadResource()               | `resource:"VERSION,required"` or `resource:"*VERSION"`             |
| Text Resource Files   | `resource` | *required* | LoadResource()               | `resource:"VERSION,required"` or `resource:"*VERSION"`             |
| Command Arguments     | `arg`      | --         | LoadCommandArguments()       | `arg:"SERVER_NAME"` or `arg:"SERVER_NAME;server name description"` |

### About Required Tags

`resource:"VERSION,required"` is equivalent to `resource:"*VERSION"`, but not equivalent to `resource:"*VERSION,required"`:

| Tag                              | Name     | Flag       |
| :------------------------------- | :------- | :--------- |
| `resource:"VERSION,required"`    | VERSION  | `required` |
| `resource:"*VERSION"`            | VERSION  | `required` |
| `resource:"*VERSION,required"`   | *VERSION | `required` |
| `resource:"*VERSION,required,_"` | *VERSION | `required` |
| `resource:"*VERSION,_"`          | *VERSION | *none*     |
| `resource:"VERSION,_"`           | VERSION  | *none*     |

If you want to keep the "`*`" symbol at the beginning of the name while making it optional, append the blank flag "`_`" to the tag.

### Supported Field Types

The tags **env**, **resource**, and **arg** do not support nested structures. Supported field types include:
`bool`, `int`, `uint`, `float`, `string`, `time.Duration`, `time.Time`, `url.URL`, `net.IP`,
their corresponding array types, as well as `bytes.Buffer`, `json.RawMessage`, and `github.com/Bofry/types.RawContent`.

## Environment Variable Configuration

The following **Config** structure will import environment variables `CACHE_HOST`, `CACHE_PASSWORD`, and `CACHE_DB`:

```go
type Config struct {
  CacheHost     string `env:"CACHE_HOST,required"`  // Required environment variable
  CachePassword string `env:"CACHE_PASSWORD"`       // Optional environment variable
  CacheDB       int    `env:"CACHE_DB"`             // Optional environment variable, automatically converted to integer
}
```

The tag `env:"CACHE_HOST,required"` can also be written as `env:"*CACHE_HOST"`. Adding the "`*`" symbol before the name is equivalent to appending the `required` flag to the tag.

## .env File Configuration

.env files are used in the same way as environment variables. Notably, **.env files will not override existing environment variables**. They are recommended for development environment settings or providing sensible defaults.

## Resource File Configuration

The following **Config** structure will import content from the **VERSION** file:

```go
type Config struct {
  AppVersion string `resource:"VERSION,required"`  // Required resource file
}
```

Resource file names can contain any Unicode characters, but cannot have spaces at the beginning or end, and cannot end with a period.

## Command Line Argument Configuration

The following **Config** structure will import command line arguments `cache-host`, `cache-password`, and `cache-db`:

```go
type Config struct {
  CacheHost     string `arg:"cache-host;cache server address and port"`
  CachePassword string `arg:"cache-password;cache server password"`
  CacheDB       int    `arg:"cache-db;cache database number"`
}
```

The tag text `arg:"cache-host;cache server address and port"` is separated by the symbol "`;`" into the name part and the usage description part.

> ⛔ Do not name an argument `help`.

## Dependencies

- Yaml - [gopkg.in/yaml.v2](https://godoc.org/gopkg.in/yaml.v2)
- Json - [encoding/json](https://golang.org/pkg/encoding/json/)
- dotenv - [github.com/joho/godotenv](https://github.com/joho/godotenv)
