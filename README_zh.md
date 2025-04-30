[EN](README.md)

# config

一個高度可定製的 Go 配置管理套件，支援多種資料來源。

## 快速入門

### 安裝套件

首先，使用 `go get` 命令安裝 config 套件：

```bash
go get -u github.com/Bofry/config
```

### 執行範例

建立您的應用程式後，開啟終端機並執行下列指令：

```bash
# 執行程式並指定監聽位址
go run main.go -listen-address=":8080"

# 或者先編譯後執行
go build
./yourapp -listen-address=":8080"
```

您可以根據需要傳遞任何已在 `arg` 標籤中定義的命令列參數。

## 核心功能

- 支援多種配置來源：環境變數、.env 檔案、YAML/JSON 檔案、命令列參數和資源檔案
- 流暢的 API 設計，使配置程式碼簡潔易讀
- 透過結構體標籤（struct tags）進行直覺化配置映射
- 自動資料型別轉換
- 彈性的配置優先順序處理
- 支援自訂輸出格式

## 使用範例

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
  ListenAddress  string `yaml:"ListenAddress"  arg:"listen-address;combination of IP address and listen port"`
  EnableCompress bool   `yaml:"UseCompress"    arg:"use-compress;indicates if response compression is enabled"`
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
 fmt.Printf("Listen Address: %s\n", config.ListenAddress)
 fmt.Printf("Cache Store: %s:%d\n", config.CacheStoreIp, config.CacheStorePort)
 
 // Start your application...
}
```

## 結構體標籤語法

| 配置類型     | 結構體標籤 | 標記旗標   | 對應的配置服務方法         | 範例                                                      |
| :----------- | :--------- | :--------- | :------------------------- | :-------------------------------------------------------- |
| 環境變數     | `env`      | *required* | LoadEnvironmentVariables() | `env:"CACHE_ADDRESS,required"` 或 `env:"*CACHE_ADDRESS"`  |
| .env 檔案    | `env`      | *required* | LoadDotEnv()               | `env:"CACHE_ADDRESS,required"` 或 `env:"*CACHE_ADDRESS"`  |
| JSON 檔案    | `json`     | --         | LoadJsonFile()             | `json:"LISTEN_PORT"`                                      |
| YAML 檔案    | `yaml`     | --         | LoadYamlFile()             | `yaml:"LISTEN_PORT"`                                      |
| 二進位資源檔 | `resource` | *required* | LoadResource()             | `resource:"VERSION,required"` 或 `resource:"*VERSION"`    |
| 文字資源檔   | `resource` | *required* | LoadResource()             | `resource:"VERSION,required"` 或 `resource:"*VERSION"`    |
| 命令列參數   | `arg`      | --         | LoadCommandArguments()     | `arg:"SERVER_NAME"` 或 `arg:"SERVER_NAME;伺服器名稱說明"` |

### 關於必填標記的說明

`resource:"VERSION,required"` 等同於 `resource:"*VERSION"`，但不等同於 `resource:"*VERSION,required"`：

| 標記                             | 名稱     | 旗標       |
| :------------------------------- | :------- | :--------- |
| `resource:"VERSION,required"`    | VERSION  | `required` |
| `resource:"*VERSION"`            | VERSION  | `required` |
| `resource:"*VERSION,required"`   | *VERSION | `required` |
| `resource:"*VERSION,required,_"` | *VERSION | `required` |
| `resource:"*VERSION,_"`          | *VERSION | *none*     |
| `resource:"VERSION,_"`           | VERSION  | *none*     |

若要在名稱前保留 "`*`" 符號且設定為非必填，可在標記中加入空白旗標 "`_`"。

### 支援的欄位型別

標記 **env**、**resource** 和 **arg** 不支援巢狀結構，支援的欄位型別包括：
`bool`, `int`, `uint`, `float`, `string`, `time.Duration`, `time.Time`, `url.URL`, `net.IP`,
以及這些型別的對應陣列型別、`bytes.Buffer`, `json.RawMessage`, 和 `github.com/Bofry/types.RawContent`。

## 環境變數配置

以下 **Config** 結構將匯入環境變數 `CACHE_HOST`, `CACHE_PASSWORD` 和 `CACHE_DB`：

```go
type Config struct {
  CacheHost     string `env:"CACHE_HOST,required"`  // 必填環境變數
  CachePassword string `env:"CACHE_PASSWORD"`       // 選填環境變數
  CacheDB       int    `env:"CACHE_DB"`             // 選填環境變數，自動轉型為整數
}
```

標記 `env:"CACHE_HOST,required"` 也可以寫成 `env:"*CACHE_HOST"`，在名稱前加上 "`*`" 符號等同於在標記中附加 `required` 旗標。

## .env 檔案配置

.env 檔案的使用方式與環境變數相同。值得注意的是，**.env 檔案不會覆寫已存在的環境變數**，建議用於開發環境設定或提供合理的預設值。

## 資源檔案配置

以下 **Config** 結構將從檔案 **VERSION** 匯入內容：

```go
type Config struct {
  AppVersion string `resource:"VERSION,required"`  // 必須存在的資源檔
}
```

資源檔名稱可以包含任何 Unicode 字元，但開頭和結尾不能有空格，結尾也不能是句點。

## 命令列參數配置

以下 **Config** 結構將匯入命令列參數 `cache-host`, `cache-password` 和 `cache-db`：

```go
type Config struct {
  CacheHost     string `arg:"cache-host;快取伺服器位址和連接埠"`
  CachePassword string `arg:"cache-password;快取伺服器密碼"`
  CacheDB       int    `arg:"cache-db;快取資料庫編號"`
}
```

標記文字 `arg:"cache-host;快取伺服器位址和連接埠"` 以符號 "`;`" 分隔為名稱部分和使用說明部分。

> ⛔ 請勿將參數命名為 `help`。

## 相依套件

- Yaml - [gopkg.in/yaml.v2](https://godoc.org/gopkg.in/yaml.v2)
- Json - [encoding/json](https://golang.org/pkg/encoding/json/)
- dotenv - [github.com/joho/godotenv](https://github.com/joho/godotenv)
