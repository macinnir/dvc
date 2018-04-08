# DVC

# Usage 

1. Load Config 
2. Load changesets json file 
3. Load changesets sql files 
4. Pull down the hash of the changesets file 
4. Verify which changesets have been applied 

## Code 
```
import (
    "github.com/macinnir/dvc"
)

func main() {
    d := dvc.DVC{}
    d.Run("path/to/changeset/files", "databaseHost:port", "username", "password")
}
```

# Changelog files 

## Naming 

```
{changesetNum}/{changeNum}_{action}_{target}__{etc}.sql
```

### changesetNum 

Directory containing multiple sql files grouped by team (e.g. sprint)

### changeNum

Numeric ordinal of the change within the changeset. 
`001` as a prefix allows for 1000 change files within the changeset and uses the native filesystem to list the files in order (the same order in which they will be run).

### Action 

One of the following: 
<dl>
    <dt>createTable</dt>
    <dd>Create a new table</dd>
    <dt>alterTable</dt>
    <dd>Alter an existing table</dd>
    <dt>dropTable</dt>
    <dd>Drop an existing table</dd>
    <dt>createView</dt>
    <dd>Create a new view</dd>
    <dt>alterView</dt>
    <dd>Alter an existing view</dd>
    <dt>dropView</dt>
    <dd>Drop an existing view</dd>
    <dt>insert</dt>
    <dd>Insert data into an existing table</dd>
</dd>

### Target

The name of table or view on which the action will be taken.

### etc (optional)

Other pertinent information for the action to be taken. 

Examples: 
```
    0001/
        001_alterTable_myTable__addColumn_foo.sql
        002_alterTable_myTable__addIndex_ucFoo.sql
```