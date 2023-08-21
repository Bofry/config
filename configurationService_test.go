package config

import (
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/Bofry/structproto"
	"gopkg.in/yaml.v2"
)

type DummyConfig struct {
	RedisHost     string   `env:"REDIS_HOST"       yaml:"redisHost"       arg:"redis-host;the Redis server address and port"`
	RedisPassword string   `env:"REDIS_PASSWORD"   yaml:"redisPassword"   arg:"redis-passowrd;the Redis password"`
	RedisDB       int      `env:"REDIS_DB"         yaml:"redisDB"         arg:"redis-db;the Redis database number"`
	RedisPoolSize int      `env:"-"                yaml:"redisPoolSize"`
	Workspace     string   `env:"-"                yaml:"workspace"       arg:"workspace;the data workspace"`
	Tags          []string `env:"TAG"`
	Version       string   `resource:".VERSION"`
}

func (c *DummyConfig) Output(writer io.Writer) error {
	fmt.Fprintf(writer, "RedisHost    : %v\n", c.RedisHost)
	fmt.Fprintf(writer, "RedisPassword: %v\n", c.RedisPassword)
	fmt.Fprintf(writer, "RedisDB      : %v\n", c.RedisDB)
	fmt.Fprintf(writer, "RedisPoolSize: %v\n", c.RedisPoolSize)
	fmt.Fprintf(writer, "Workspace    : %v\n", c.Workspace)
	fmt.Fprintf(writer, "Tags         : %v\n", c.Tags)
	fmt.Fprintf(writer, "Version      : %v\n", c.Version)

	return nil
}

func TestConfigurationService(t *testing.T) {
	os.Clearenv()
	initializeEnvironment()
	initializekubernetesEnvironment()
	initializeArgs()
	initializeDotEnv()
	initializeDotVERSION()
	initializeConfigYaml()
	initializeConfigStagingYaml()
	initializeConfigProductionYaml()

	conf := DummyConfig{}

	NewConfigurationService(&conf).
		LoadDotEnv().
		LoadEnvironmentVariables("").
		LoadEnvironmentVariables("K8S").
		LoadYamlFile("config.yaml").
		LoadYamlFile("config.${ENVIRONMENT}.yaml").
		LoadCommandArguments().
		LoadResource("").
		Output()

	expected := DummyConfig{
		RedisHost:     "demo-kubernetes:6379",
		RedisPassword: "p@ssw0rd",
		RedisDB:       32,
		RedisPoolSize: 50,
		Workspace:     "demo_prod",
		Tags:          []string{"demo", "test"},
		Version:       "v1.0.2",
	}
	if !reflect.DeepEqual(expected, conf) {
		t.Errorf("assert 'DummyConfig':: expected '%#+v', got '%#+v'", expected, conf)
	}
}

func initializeEnvironment() {
	os.Setenv("ENVIRONMENT", "production")
	os.Setenv("REDIS_HOST", "127.0.0.3:6379")
	os.Setenv("REDIS_PASSWORD", "1234")
}

func initializekubernetesEnvironment() {
	os.Setenv("K8S_REDIS_HOST", "demo-kubernetes:6379")
	os.Setenv("K8S_REDIS_PASSWORD", "p@ssw0rd")
	os.Setenv("K8S_REDIS_DB", "6")
}

func initializeArgs() {
	os.Args = []string{"example", "-redis-db", "32"}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func initializeDotEnv() {
	err := os.WriteFile(".env", []byte(
		strings.Join([]string{
			"REDIS_HOST=127.0.0.1:6379",
			"REDIS_DB=29",
			"TAG=demo,test",
		}, "\n")), 0644)
	if err != nil {
		panic(err)
	}
}

func initializeDotVERSION() {
	err := os.WriteFile(".VERSION", []byte(
		strings.Join([]string{
			"v1.0.2",
		}, "\n")), 0644)
	if err != nil {
		panic(err)
	}
}

func initializeConfigYaml() {
	err := os.WriteFile("config.yaml", []byte(
		strings.Join([]string{
			"redisDB: 3",
			"redisPoolSize: 10",
			"workspace: demo_test",
		}, "\n")), 0644)
	if err != nil {
		panic(err)
	}
}

func initializeConfigStagingYaml() {
	err := os.WriteFile("config.staging.yaml", []byte(
		strings.Join([]string{
			"redisDB: 9",
			"redisPoolSize: 10",
			"workspace: demo_stag",
		}, "\n")), 0644)
	if err != nil {
		panic(err)
	}
}

func initializeConfigProductionYaml() {
	err := os.WriteFile("config.production.yaml", []byte(
		strings.Join([]string{
			"redisDB: 12",
			"redisPoolSize: 50",
			"workspace: demo_prod",
		}, "\n")), 0644)
	if err != nil {
		panic(err)
	}
}

func TestConfigurationService_LoadFile(t *testing.T) {

	conf := DummyConfig{}

	NewConfigurationService(&conf).
		LoadFile("config.yaml", yaml.Unmarshal)

	expected := DummyConfig{
		RedisHost:     "",
		RedisPassword: "",
		RedisDB:       3,
		RedisPoolSize: 10,
		Workspace:     "demo_test",
		Version:       "",
	}
	if !reflect.DeepEqual(expected, conf) {
		t.Errorf("assert 'DummyConfig':: expected '%#+v', got '%#+v'", expected, conf)
	}
}

func TestConfigurationService_Map(t *testing.T) {
	t.Setenv("Environment", "staging")

	conf := DummyConfig{
		RedisHost:     "demo-kubernetes:6379",
		RedisPassword: "p@ssw0rd",
		RedisDB:       32,
		RedisPoolSize: 50,
		Workspace:     "demo_${Environment}",
		Tags:          []string{"demo", "test"},
		Version:       "v1.0.2",
	}

	NewConfigurationService(&conf).
		Map(func(field structproto.FieldInfo, rv reflect.Value) error {
			switch rv.Kind() {
			case reflect.String:
				if !rv.IsZero() {
					val := os.ExpandEnv(rv.String())
					rv.SetString(val)
				}
			}
			return nil
		})

	expected := DummyConfig{
		RedisHost:     "demo-kubernetes:6379",
		RedisPassword: "p@ssw0rd",
		RedisDB:       32,
		RedisPoolSize: 50,
		Workspace:     "demo_staging",
		Tags:          []string{"demo", "test"},
		Version:       "v1.0.2",
	}
	if !reflect.DeepEqual(expected, conf) {
		t.Errorf("assert 'DummyConfig':: expected '%#+v', got '%#+v'", expected, conf)
	}
}
