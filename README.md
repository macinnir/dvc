# DVC

Database management and code generation 

## Commands

- Golang Code Generation 
- Database Change Management 

```
dvc compare 
dvc compare -a 

dvc import 

dvc gen models -f -c
dvc gen dals -f
dvc gen interfaces 
dvc gen goperms 
dvc gen tsperms 
dvc gen ts 

dvc gen routes


```

### CLI

Under construction.

```
$ dvc cli 
```

### Compare 

Compare two schemas

```
$ dvc compare 
$ dvc compare -a
```

### Connections

List connections in the configuration 

```
$ dvc connections
```

### Gen 

Generate code 

```
$ dvc gen models 
$ dvc gen dals 
$ dvc gen interfaces
$ dvc gen routes 
$ dvc gen ts 
$ dvc gen tsperms 
$ dvc gen goperms
```

### Import 

Import schema from the databases

```
$ dvc import 
```


### ls 

List elements in the local schema configuration 

```
# List all tables across all schemas
$ dvc ls 

# Show all columns that start with the string 
$ dvc ls .[string]

# Search all Tables
# If `search` is not a valid table name 
$ dvc ls [search]

# List all columns in a table 
# If `table name` is a valid table name 
$ dvc ls [table name]

```

### schemas 

List schemas and their connections 

```
$ dvc schemas 
```

### version

```
$ dvc version $(git describe) 
```

```
$ dvc list shards 

```

https://godoc.org/github.com/macinnir/dvc


# TODO 
[ ] Remove requirement of "core" schema v "app" schema 