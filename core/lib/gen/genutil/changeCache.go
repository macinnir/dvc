package genutil

import (
	"github.com/macinnir/dvc/core/lib/cache"
	"github.com/macinnir/dvc/core/lib/schema"
)

// GetChangedTables returns a slice of tables that have changed
func GetChangedTables(schemaList *schema.SchemaList, tablesCache *cache.TablesCache, force bool) ([]*schema.Table, error) {

	var changed = []*schema.Table{}

	for k := range schemaList.Schemas {

		var schema = schemaList.Schemas[k]

		for l := range schema.Tables {

			// tableMap[schema.Tables[l].Name] = schema.Tables[l]

			var table = schemaList.Schemas[k].Tables[l]
			var tableHash string
			tableHash, _ = cache.HashTable(table)

			// If the model has been hashed before...
			if _, ok := tablesCache.Models[table.Key()]; ok {

				// And the hash hasn't changed, skip...
				if tableHash == tablesCache.Models[table.Key()] && !force {
					// fmt.Printf("Table `%s` hasn't changed! Skipping...\n", table.Name)
					continue
				}
			}

			changed = append(changed, table)

			// Update the models cache
			tablesCache.Models[table.Key()] = tableHash
		}
	}

	cache.SaveTableCache(tablesCache)

	return changed, nil

}
