package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Database Names can be simple like `Foo` or part of a shard like `Foo_0` and `Foo_1`
// where both would be distributed to as shards
//
// You can specify that a database is to be sharded only remotely.
// Locally you can set a database name as `Foo` so that your system interacts with a single database
// But when using the remote configuration, there may be several sharded databases
// all with the same schema. So, remotely,  you might have `Foo_0`, `Foo_1`, `Foo_2`, and so on.
// All of those remote databases will be applied the same schema changes

type ConfigDatabase struct {
	Key       string            `json:"key"`
	Type      string            `json:"type"`
	User      string            `json:"user"`
	Pass      string            `json:"pass"`
	Host      string            `json:"host"`
	Name      string            `json:"name"`
	Enums     []string          `json:"enums"`
	OneToMany map[string]string `json:"onetomany"`
	OneToOne  map[string]string `json:"onetoone"`
	ManyToOne map[string]string `json:"manytoone"`
}

// Config contains a set of configuration values used throughout the application
type Config struct {
	ChangeSetPath             string            `json:"changesetPath"`
	BasePackage               string            `json:"basePackage"`
	Databases                 []*ConfigDatabase `json:"Databases"`
	TypescriptModelsPath      string            `json:"TypescriptModelsPath"`
	TypescriptDTOsPath        string            `json:"TypescriptDTOsPath"`
	TypescriptAggregatesPath  string            `json:"TypescriptAggregatesPath"`
	TypescriptPermissionsPath string            `json:"TypescriptPermissionsPath"`
	TypescriptRoutesPath      string            `json:"TypescriptRoutesPath"`
	Cache                     map[string]*CacheConfig
	Packages                  struct {
		Cache    string `json:"cache"`
		Models   string `json:"models"`
		Schema   string `json:"schema"`
		Repos    string `json:"repos"`
		Services string `json:"services"`
		API      string `json:"api"`
	} `json:"packages"`

	Dirs struct {
		Dals                  string `json:"dals"`
		DalInterfaces         string `json:"dalInterfaces"`
		Repos                 string `json:"repos"`
		Cache                 string `json:"cache"`
		Models                string `json:"models"`
		Integrations          string `json:"integrations"`
		IntegrationInterfaces string `json:"integrationInterfaces"`
		Aggregates            string `json:"aggregates"`
		Schema                string `json:"schema"`
		Typescript            string `json:"typescript"`
		Services              string `json:"services"`
		ServiceInterfaces     string `json:"serviceInterfaces"`
		Controllers           string `json:"controllers"`
		API                   string `json:"api"`
		Permissions           string `json:"permissions"`
	} `json:"dirs"`
}

type CacheConfig struct {
	Indices   []*CacheConfigIndex
	Aggregate *CacheConfigAggregate
	HasHashID bool
	Search    []*CacheConfigSearch
}

type CacheConfigIndex struct {
	Field  string
	Unique bool
}

type CacheConfigAggregate struct {
	Location   string
	Properties []*CacheConfigAggregateProperty
}

type CacheConfigSearch struct {
	Fields     []string
	Conditions []string
}

type CacheConfigAggregateProperty struct {
	Property string
	On       string
	Table    string
	Type     string
}

// LoadConfig loads a config file
func LoadConfig() (*Config, error) {

	f, e := os.Open(ConfigFilePath)
	if e != nil {
		return nil, e
	}

	var fileBytes []byte
	fileBytes, e = ioutil.ReadAll(f)
	if e != nil {
		return nil, e
	}

	config := &Config{}
	if e = json.Unmarshal(fileBytes, config); e != nil {
		return nil, fmt.Errorf("invalid config: %w", e)
	}

	return config, nil
}

// LoadConfig loads a config file
func LoadCoreCacheFile() (map[string]*CacheConfig, error) {

	f, e := os.Open(CoreCacheConfig)
	if e != nil {
		return nil, e
	}

	var fileBytes []byte
	fileBytes, e = ioutil.ReadAll(f)
	if e != nil {
		return nil, e
	}

	var config = map[string]*CacheConfig{}
	if e = json.Unmarshal(fileBytes, &config); e != nil {
		return nil, fmt.Errorf("invalid config: %w", e)
	}

	return config, nil
}

func ExtractRootNameFromKey(key string) string {

	dbRootName := key
	if strings.Contains(key, "_") {
		parts := strings.Split(key, "_")
		dbRootName = strings.Join(parts[0:len(parts)-1], "_")
		// dbRootName = strings.Split(key, "_")[0]
	}

	return dbRootName
}

// GetDatabaseByName gets all databases that have the name of `name` as the base
// and includes shards with `_[shardNo]` as the suffix.
func GetDatabasesByRootName(rootName string, config *Config) []*ConfigDatabase {

	databases := []*ConfigDatabase{}

	for k := range config.Databases {

		dbKey := config.Databases[k].Key
		dbRootName := ExtractRootNameFromKey(dbKey)
		if rootName == dbRootName {
			databases = append(databases, config.Databases[k])
		}

	}

	return databases
}
