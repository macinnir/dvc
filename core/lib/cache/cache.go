package cache

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/schema"
)

var (
	ErrCouldNotWriteTableCacheToFile = errors.New("could not write table cache to file")
	ErrCouldNotSerializeTableCache   = errors.New("cannot serialize table cache")
	ErrCannotReadTableCacheFromDisk  = errors.New("cannot read table cache from disk")
)

// TablesCache stores an md5 hash of the JSON representation of a table in the schema.json file
// These hashes are used to skip unchanged models for DAL and Model generation
type TablesCache struct {
	Dals   map[string]string
	Models map[string]string
}

// NewTablesCache is a factory method for TablesCache
func NewTablesCache() TablesCache {
	return TablesCache{
		Dals:   map[string]string{},
		Models: map[string]string{},
	}
}

// LoadTableCache loads the table cache
func LoadTableCache() (TablesCache, error) {

	tableCache := NewTablesCache()

	if _, e := os.Stat(lib.TablesCacheFilePath); e == nil {
		if fileBytes, e := ioutil.ReadFile(lib.TablesCacheFilePath); e == nil {
			if e = json.Unmarshal(fileBytes, &tableCache); e != nil {
				return tableCache, ErrCannotReadTableCacheFromDisk
			}
		}
	}

	return tableCache, nil
}

// GenerateTableCache generates a new TableCache
func GenerateTableCache(database *schema.Schema) (map[string]string, error) {

	newModels := map[string]string{}

	for _, table := range database.Tables {

		var marshalledTable []byte
		var e error

		if marshalledTable, e = json.Marshal(table); e != nil {
			return newModels, e
		}

		// Build the list of new model hashes to check against
		newModels[table.Name] = lib.HashStringMd5(string(marshalledTable))
	}

	return newModels, nil
}

// SaveTableCache writes a table cache to the table cache file
func SaveTableCache(tablesCache TablesCache) error {

	lib.EnsureDir(lib.MetaDirectory)

	var e error
	var tableCacheSerialized []byte

	if tableCacheSerialized, e = json.Marshal(tablesCache); e != nil {
		return ErrCouldNotSerializeTableCache
	}
	if e = ioutil.WriteFile(lib.TablesCacheFilePath, tableCacheSerialized, 0777); e != nil {
		return ErrCouldNotWriteTableCacheToFile
	}

	return nil
}
