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
	os.Clearenv()
	os.Setenv("REDIS_HOST", "192.168.56.53")
	os.Setenv("RESID_SECRET", "foobar")
	os.Setenv("REDIS_DB", "3")
	os.Setenv("WORKSPACE", "demo_test")
	os.Setenv("TAG", "demo,test")

	c := config{}
	err := Process("", &c)
	if err != nil {
		t.Error(err)
	}

	var expectedRedisHost = "192.168.56.53"
	if c.RedisHost != expectedRedisHost {
		t.Errorf("assert 'config.RedisHost':: expected '%v', got '%v'", expectedRedisHost, c.RedisHost)
	}
	var expectedRedisSecret = "foobar"
	if c.RedisSecret != expectedRedisSecret {
		t.Errorf("assert 'config.RedisSecret':: expected '%v', got '%v'", expectedRedisSecret, c.RedisSecret)
	}
	var expectedRedisDB = 3
	if c.RedisDB != expectedRedisDB {
		t.Errorf("assert 'config.RedisDB':: expected '%v', got '%v'", expectedRedisDB, c.RedisDB)
	}
	var expectedWorkspace = "demo_test"
	if c.Workspace != expectedWorkspace {
		t.Errorf("assert 'config.Workspace':: expected '%v', got '%v'", expectedWorkspace, c.Workspace)
	}
	var expectedIgnoredField = ""
	if c.IgnoredField != expectedIgnoredField {
		t.Errorf("assert 'config.IgnoredField':: expected '%v', got '%v'", expectedIgnoredField, c.IgnoredField)
	}
	var expectedTags = []string{"demo", "test"}
	if !reflect.DeepEqual(c.Tags, expectedTags) {
		t.Errorf("assert 'config.Tags':: expected '%+v', got '%+v'", expectedTags, c.Tags)
	}
}

func TestLoad_WithPrefix(t *testing.T) {
	os.Clearenv()
	os.Setenv("K8S_REDIS_HOST", "192.168.56.53")
	os.Setenv("K8S_RESID_SECRET", "foobar")
	os.Setenv("K8S_REDIS_DB", "3")
	os.Setenv("K8S_WORKSPACE", "demo_test")

	c := config{}
	err := Process("K8S", &c)
	if err != nil {
		t.Error(err)
	}

	var expectedRedisHost = "192.168.56.53"
	if c.RedisHost != expectedRedisHost {
		t.Errorf("assert 'config.RedisHost':: expected '%v', got '%v'", expectedRedisHost, c.RedisHost)
	}
	var expectedRedisSecret = "foobar"
	if c.RedisSecret != expectedRedisSecret {
		t.Errorf("assert 'config.RedisSecret':: expected '%v', got '%v'", expectedRedisSecret, c.RedisSecret)
	}
	var expectedRedisDB = 3
	if c.RedisDB != expectedRedisDB {
		t.Errorf("assert 'config.RedisDB':: expected '%v', got '%v'", expectedRedisDB, c.RedisDB)
	}
	var expectedWorkspace = "demo_test"
	if c.Workspace != expectedWorkspace {
		t.Errorf("assert 'config.Workspace':: expected '%v', got '%v'", expectedWorkspace, c.Workspace)
	}
}

func TestLoadDotEnv(t *testing.T) {
	os.Clearenv()
	c := config{}
	err := LoadDotEnv(&c)
	if err != nil {
		t.Error(err)
	}

	var expectedRedisHost = "192.168.56.53"
	if c.RedisHost != expectedRedisHost {
		t.Errorf("assert 'config.RedisHost':: expected '%v', got '%v'", expectedRedisHost, c.RedisHost)
	}
	var expectedRedisSecret = "foobar"
	if c.RedisSecret != expectedRedisSecret {
		t.Errorf("assert 'config.RedisSecret':: expected '%v', got '%v'", expectedRedisSecret, c.RedisSecret)
	}
	var expectedRedisDB = 3
	if c.RedisDB != expectedRedisDB {
		t.Errorf("assert 'config.RedisDB':: expected '%v', got '%v'", expectedRedisDB, c.RedisDB)
	}
	var expectedWorkspace = "demo_test"
	if c.Workspace != expectedWorkspace {
		t.Errorf("assert 'config.Workspace':: expected '%v', got '%v'", expectedWorkspace, c.Workspace)
	}
	var expectedIgnoredField = ""
	if c.IgnoredField != expectedIgnoredField {
		t.Errorf("assert 'config.IgnoredField':: expected '%v', got '%v'", expectedIgnoredField, c.IgnoredField)
	}
	var expectedTags = []string{"demo", "test"}
	if !reflect.DeepEqual(c.Tags, expectedTags) {
		t.Errorf("assert 'config.Tags':: expected '%+v', got '%+v'", expectedTags, c.Tags)
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

	var expectedRedisHost = "10.10.171.6"
	if c.RedisHost != expectedRedisHost {
		t.Errorf("assert 'config.RedisHost':: expected '%v', got '%v'", expectedRedisHost, c.RedisHost)
	}
	var expectedRedisSecret = "foobar"
	if c.RedisSecret != expectedRedisSecret {
		t.Errorf("assert 'config.RedisSecret':: expected '%v', got '%v'", expectedRedisSecret, c.RedisSecret)
	}
	var expectedRedisDB = 3
	if c.RedisDB != expectedRedisDB {
		t.Errorf("assert 'config.RedisDB':: expected '%v', got '%v'", expectedRedisDB, c.RedisDB)
	}
	var expectedWorkspace = "demo_test"
	if c.Workspace != expectedWorkspace {
		t.Errorf("assert 'config.Workspace':: expected '%v', got '%v'", expectedWorkspace, c.Workspace)
	}
	var expectedIgnoredField = ""
	if c.IgnoredField != expectedIgnoredField {
		t.Errorf("assert 'config.IgnoredField':: expected '%v', got '%v'", expectedIgnoredField, c.IgnoredField)
	}
	var expectedTags = []string{"demo", "test"}
	if !reflect.DeepEqual(c.Tags, expectedTags) {
		t.Errorf("assert 'config.Tags':: expected '%+v', got '%+v'", expectedTags, c.Tags)
	}
}
