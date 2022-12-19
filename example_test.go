package config_test

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Bofry/config"
	"gopkg.in/yaml.v2"
)

func Example() {
	// set env
	{
		os.Clearenv()
		os.Setenv("ENVIRONMENT", "production")
		os.Setenv("REDIS_HOST", "127.0.0.3:6379")
		os.Setenv("REDIS_PASSWORD", "1234")
		os.Setenv("K8S_REDIS_HOST", "demo-kubernetes:6379")
		os.Setenv("K8S_REDIS_PASSWORD", "p@ssw0rd")
		os.Setenv("K8S_REDIS_DB", "6")
	}
	// set command line arguments
	{
		os.Args = []string{"example", "--redis-db", "32"}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}
	// prepare .env
	{
		os.WriteFile(".env", []byte(
			strings.Join([]string{
				"REDIS_HOST=127.0.0.1:6379",
				"REDIS_DB=29",
				"TAG=demo,test",
			}, "\n")), 0644)
	}
	// prepare .VERSION
	{
		os.WriteFile(".VERSION", []byte(
			strings.Join([]string{
				"v1.0.2",
			}, "\n")), 0644)
	}
	// prepare config.yaml
	{
		os.WriteFile("config.yaml", []byte(
			strings.Join([]string{
				"redisDB: 3",
				"redisPoolSize: 10",
				"workspace: demo_test",
			}, "\n")), 0644)
	}
	// prepare config.staging.yaml
	{
		os.WriteFile("config.staging.yaml", []byte(
			strings.Join([]string{
				"redisDB: 9",
				"redisPoolSize: 10",
				"workspace: demo_stag",
			}, "\n")), 0644)
	}
	// prepare config.production.yaml
	{
		os.WriteFile("config.production.yaml", []byte(
			strings.Join([]string{
				"redisDB: 12",
				"redisPoolSize: 50",
				"workspace: demo_prod",
			}, "\n")), 0644)
	}

	conf := struct {
		RedisHost     string   `env:"REDIS_HOST"       yaml:"redisHost"       arg:"redis-host;the Redis server address and port"`
		RedisPassword string   `env:"REDIS_PASSWORD"   yaml:"redisPassword"   arg:"redis-passowrd;the Redis password"`
		RedisDB       int      `env:"REDIS_DB"         yaml:"redisDB"         arg:"redis-db;the Redis database number"`
		RedisPoolSize int      `env:"-"                yaml:"redisPoolSize"`
		Workspace     string   `env:"-"                yaml:"workspace"       arg:"workspace;the data workspace"`
		Tags          []string `env:"TAG"`
		Version       string   `resource:".VERSION"`
	}{}

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
	// Output:
	// RedisHost     = "demo-kubernetes:6379"
	// RedisPassword = "p@ssw0rd"
	// RedisDB       = 32
	// RedisPoolSize = 50
	// Workspace     = "demo_prod"
	// Tags          = ["demo" "test"]
	// Version       = "v1.0.2"
}

func ExampleConfigurationService_LoadFile() {
	// prepare config.yaml
	{
		os.WriteFile("config.yaml", []byte(
			strings.Join([]string{
				"redisDB: 3",
				"redisPoolSize: 10",
				"workspace: demo_test",
			}, "\n")), 0644)
	}

	conf := struct {
		RedisHost     string   `env:"REDIS_HOST"       yaml:"redisHost"       arg:"redis-host;the Redis server address and port"`
		RedisPassword string   `env:"REDIS_PASSWORD"   yaml:"redisPassword"   arg:"redis-passowrd;the Redis password"`
		RedisDB       int      `env:"REDIS_DB"         yaml:"redisDB"         arg:"redis-db;the Redis database number"`
		RedisPoolSize int      `env:"-"                yaml:"redisPoolSize"`
		Workspace     string   `env:"-"                yaml:"workspace"       arg:"workspace;the data workspace"`
		Tags          []string `env:"TAG"`
		Version       string   `resource:".VERSION"`
	}{}

	config.NewConfigurationService(&conf).
		LoadFile("config.yaml", yaml.Unmarshal)
	fmt.Printf("RedisHost     = %q\n", conf.RedisHost)
	fmt.Printf("RedisPassword = %q\n", conf.RedisPassword)
	fmt.Printf("RedisDB       = %d\n", conf.RedisDB)
	fmt.Printf("RedisPoolSize = %d\n", conf.RedisPoolSize)
	fmt.Printf("Workspace     = %q\n", conf.Workspace)
	fmt.Printf("Tags          = %q\n", conf.Tags)
	fmt.Printf("Version       = %q\n", conf.Version)
	// Output:
	// RedisHost     = ""
	// RedisPassword = ""
	// RedisDB       = 3
	// RedisPoolSize = 10
	// Workspace     = "demo_test"
	// Tags          = []
	// Version       = ""
}
