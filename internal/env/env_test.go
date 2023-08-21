package env

import (
	"os"
	"reflect"
	"testing"
)

type config struct {
	RedisHost    string   `env:"*REDIS_HOST"`
	RedisSecret  string   `env:"RESID_SECRET"`
	RedisDB      int      `env:"REDIS_DB"`
	Workspace    string   `env:"*WORKSPACE"`
	Tags         []string `env:"TAG"`
	IgnoredField string   `env:"-"`
}

func TestLoad(t *testing.T) {
	t.Setenv("REDIS_HOST", "192.168.56.53")
	t.Setenv("RESID_SECRET", "foobar")
	t.Setenv("REDIS_DB", "3")
	t.Setenv("WORKSPACE", "demo_test")
	t.Setenv("TAG", "demo,test")

	c := config{}
	err := Process("", &c)
	if err != nil {
		t.Error(err)
	}

	expected := config{
		RedisHost:    "192.168.56.53",
		RedisSecret:  "foobar",
		RedisDB:      3,
		Workspace:    "demo_test",
		IgnoredField: "",
		Tags:         []string{"demo", "test"},
	}
	if !reflect.DeepEqual(expected, c) {
		t.Errorf("assert 'config':: expected '%#+v', got '%#+v'", expected, c)
	}
}

func TestLoad_WithPrefix(t *testing.T) {
	t.Setenv("K8S_REDIS_HOST", "192.168.56.53")
	t.Setenv("K8S_RESID_SECRET", "foobar")
	t.Setenv("K8S_REDIS_DB", "3")
	t.Setenv("K8S_WORKSPACE", "demo_test")

	c := config{}
	err := Process("K8S", &c)
	if err != nil {
		t.Error(err)
	}

	expected := config{
		RedisHost:   "192.168.56.53",
		RedisSecret: "foobar",
		RedisDB:     3,
		Workspace:   "demo_test",
	}
	if !reflect.DeepEqual(expected, c) {
		t.Errorf("assert 'config':: expected '%#+v', got '%#+v'", expected, c)
	}
}

func TestLoadDotEnv(t *testing.T) {
	os.Clearenv()
	c := config{}
	err := LoadDotEnv(&c)
	if err != nil {
		t.Error(err)
	}

	expected := config{
		RedisHost:    "192.168.56.53",
		RedisSecret:  "foobar",
		RedisDB:      3,
		Workspace:    "demo_test",
		IgnoredField: "",
		Tags:         []string{"demo", "test"},
	}
	if !reflect.DeepEqual(expected, c) {
		t.Errorf("assert 'config':: expected '%#+v', got '%#+v'", expected, c)
	}
}

func TestLoadDotEnvFile(t *testing.T) {
	os.Clearenv()
	os.Setenv("ENVIRONMENT", "local")

	c := config{}
	err := LoadDotEnvFile(".env.${ENVIRONMENT}", &c)
	if err != nil {
		t.Error(err)
	}

	expected := config{
		RedisHost:    "10.10.171.6",
		RedisSecret:  "foobar",
		RedisDB:      3,
		Workspace:    "demo_test",
		IgnoredField: "",
		Tags:         []string{"demo", "test"},
	}
	if !reflect.DeepEqual(expected, c) {
		t.Errorf("assert 'config':: expected '%#+v', got '%#+v'", expected, c)
	}
}
