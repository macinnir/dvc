package gen

import (
	"text/template"

	"github.com/macinnir/dvc/core/lib/schema"
)

var DALTemplate = template.Must(template.New("template-dal-file").Funcs(template.FuncMap{
	"dataTypeToGoTypeString": schema.DataTypeToGoTypeString,
	"dataTypeToFormatString": schema.DataTypeToFormatString,
	"toArgName":              toArgName,
}).Parse(`// Generated Code; DO NOT EDIT.

package dal

import ( 
	"{{ .BasePackage }}/gen/definitions/models" 
	"github.com/macinnir/dvc/core/lib/utils/db"
	"github.com/macinnir/dvc/core/lib/utils/log"
	"github.com/macinnir/dvc/core/lib/utils/errors"
	"github.com/macinnir/dvc/core/lib/utils/query"
	"database/sql"
	"context"
	"fmt"{{ if .HasNull }}
	"gopkg.in/guregu/null.v3"{{ end }}{{ if or .IsDateCreated .IsLastUpdated }}
	"time"{{ end }}{{ if .HasSpecialColumns }}
	"strconv"
	"strings"{{ end }}
)

// {{.Table.Name}}DAL is a data repository for {{.Table.Name}} objects
type {{.Table.Name}}DAL struct {
	db  []db.IDB
	log log.ILog
}

// New{{.Table.Name}}DAL returns a new instance of {{.Table.Name}}Repo
func New{{.Table.Name}}DAL(db []db.IDB, log log.ILog) *{{.Table.Name}}DAL {
	return &{{.Table.Name}}DAL{db, log}
}

func (r *{{.Table.Name}}DAL) Raw(shard int64, q string, args ...interface{}) ([]*models.{{.Table.Name}}, error) { 
	return (&models.{{.Table.Name}}{}).Raw(r.db[shard], fmt.Sprintf(q, args...)) 
}

func (r *{{.Table.Name}}DAL) Select(shard int64) *models.{{.Table.Name}}DALSelector { 
	return (&models.{{.Table.Name}}{}).Select(r.db[shard])
}

func (r *{{.Table.Name}}DAL) Count(shard int64) *models.{{.Table.Name}}DALCounter { 
	return (&models.{{.Table.Name}}{}).Count(r.db[shard])
}

func (r *{{.Table.Name}}DAL) Sum(shard int64, col query.Column) *models.{{.Table.Name}}DALSummer { 
	return (&models.{{.Table.Name}}{}).Sum(r.db[shard], col)
}

func (r *{{.Table.Name}}DAL) Min(shard int64, col query.Column) *models.{{.Table.Name}}DALMinner { 
	return (&models.{{.Table.Name}}{}).Min(r.db[shard], col)
}

func (r *{{.Table.Name}}DAL) Max(shard int64, col query.Column) *models.{{.Table.Name}}DALMaxer { 
	return (&models.{{.Table.Name}}{}).Max(r.db[shard], col)
}

func (r *{{.Table.Name}}DAL) Get(shard int64) *models.{{.Table.Name}}DALGetter { 
	return (&models.{{.Table.Name}}{}).Get(r.db[shard])
}

// Create creates a new {{.Table.Name}} entry in the database
func (r *{{.Table.Name}}DAL) Create(shard int64, model *models.{{.Table.Name}}) error { {{if .IsDateCreated}}
	
	model.DateCreated = time.Now().UnixNano() / 1000000{{end}}
	{{if .IsLastUpdated}}
	model.LastUpdated = time.Now().UnixNano() / 1000000
	{{end}}
	e := model.Create(r.db[shard])
	if e != nil {
		r.log.Errorf("{{.Table.Name}}DAL.Insert > %s", e.Error())
		return fmt.Errorf("{{.Table.Name}}DAL.Insert: %w", e)	
	}

	// r.log.Debugf("{{.Table.Name}}DAL.Insert(%d)", model.{{.PrimaryKey}})

	return nil
}

// CreateMany creates {{.Table.Name}} objects in chunks
func (r *{{.Table.Name}}DAL) CreateMany(shard int64, modelSlice []*models.{{.Table.Name}}) error {

	var e error 

	// No records 
	if len(modelSlice) == 0 {
		return nil
	}

	// Don't use a transaction if only a single value
	if len(modelSlice) == 1 {
		e = r.Create(shard, modelSlice[0])
		if e != nil { 
			return fmt.Errorf("{{.Table.Name}}CreateMany([](%d)): %w", len(modelSlice), e)
		}

		return nil 
	}

	chunkSize := 25
	chunks := [][]*models.{{.Table.Name}}{}

	for i := 0; i < len(modelSlice); i += chunkSize {
		end := i + chunkSize
		if end > len(modelSlice) {
			end = len(modelSlice)
		}
		chunks = append(chunks, modelSlice[i:end])
	}

	for chunkID, chunk := range chunks {

		var tx *sql.Tx
		ctx := context.Background()
		tx, e = r.db[shard].BeginTx(ctx, nil)
		if e != nil {
			return fmt.Errorf("{{.Table.Name}}CreateMany([](%d)) (Chunk %d): %w", len(modelSlice), chunkID, e)
		}

		for insertID, model := range chunk {

			{{if .IsDateCreated}}
			model.DateCreated = time.Now().UnixNano() / 1000000{{end}}
			{{if .IsLastUpdated}}
			model.LastUpdated = time.Now().UnixNano() / 1000000{{end}}

			var result sql.Result 
			result, e = tx.ExecContext(ctx, "{{.InsertSQL}}", {{.InsertArgs}})
			if e != nil {
				r.log.Errorf("{{.Table.Name}}.CreateMany([](%d)) (Chunk %d.%d) > %s", len(modelSlice), chunkID, insertID, e.Error())
				break
			} else {
				// r.log.Debugf("{{.Table.Name}}.CreateMany([](%d)) (Chunk %d.%d)", len(modelSlice), chunkID, insertID)
				model.{{.PrimaryKey}}, _ = result.LastInsertId()
			}
		}

		if e != nil {
			return fmt.Errorf("{{.Table.Name}}CreateMany([](%d)) (Chunk %d): %w", len(modelSlice), chunkID, e)
		}

		e = tx.Commit()
		if e != nil {
			return fmt.Errorf("{{.Table.Name}}CreateMany([](%d)) (Chunk %d Commit): %w", len(modelSlice), chunkID, e)
		}
	}

	return nil 

}

// Update updates an existing {{.Table.Name}} entry in the database
func (r *{{.Table.Name}}DAL) Update(shard int64, model *models.{{.Table.Name}}) error {
	var e error
{{if .IsLastUpdated}}
	model.LastUpdated = time.Now().UnixNano() / 1000000{{end}}
	_, e = r.db[shard].Exec("{{.UpdateSQL}}", {{.UpdateArgs}})
	if e != nil {
		r.log.Errorf("{{.Table.Name}}DAL.Update(%d) > %s", model.{{.PrimaryKey}}, e.Error())
		return fmt.Errorf("{{.Table.Name}}DAL.Update(%d): %w", model.{{.PrimaryKey}}, e)
	} else {
		// r.log.Debugf("{{.Table.Name}}DAL.Update(%d)", model.{{.PrimaryKey}})
	}
	return nil
}

// UpdateMany updates a slice of {{.Table.Name}} objects in chunks
func (r {{.Table.Name}}DAL) UpdateMany(shard int64, modelSlice []*models.{{.Table.Name}}) error {
	var e error

	// No records 
	if len(modelSlice) == 0 {
		return nil
	}

	// Don't use a transaction if only a single value
	if len(modelSlice) == 1 {
		e = r.Update(shard, modelSlice[0])
		if e != nil { 
			return fmt.Errorf("{{.Table.Name}}UpdateMany([](%d)): %w", len(modelSlice), e)
		}
		return nil
	}

	chunkSize := 25
	chunks := [][]*models.{{.Table.Name}}{}

	for i := 0; i < len(modelSlice); i += chunkSize {
		end := i + chunkSize
		if end > len(modelSlice) {
			end = len(modelSlice)
		}
		chunks = append(chunks, modelSlice[i:end])
	}

	for chunkID, chunk := range chunks {

		var tx *sql.Tx
		ctx := context.Background()
		tx, e = r.db[shard].BeginTx(ctx, nil)
		if e != nil {
			return fmt.Errorf("{{.Table.Name}}UpdateMany([](%d)) (Chunk %d): %w", len(modelSlice), chunkID, e)
		}

		for updateID, model := range chunk {
{{if .IsLastUpdated}}
			model.LastUpdated = time.Now().UnixNano() / 1000000{{end}}

			_, e = tx.ExecContext(ctx, "{{.UpdateSQL}}", {{.UpdateArgs}})
			if e != nil {
				r.log.Errorf("{{.Table.Name}}.UpdateMany([](%d)) (Chunk %d.%d) > %s", len(modelSlice), chunkID, updateID, e.Error())
				return fmt.Errorf("{{.Table.Name}}UpdateMany([](%d)) (Chunk %d.%d): %w", len(modelSlice), chunkID, updateID, e)
			} else {
				// r.log.Debugf("{{.Table.Name}}.UpdateMany([](%d)) (Chunk %d.%d)", len(modelSlice), chunkID, updateID)
			}
		}

		if e != nil {
			return fmt.Errorf("{{.Table.Name}}UpdateMany([](%d)) (Chunk %d): %w", len(modelSlice), chunkID, e)
		}

		e = tx.Commit()

		if e != nil { 
			return fmt.Errorf("{{.Table.Name}}UpdateMany([](%d)) (Chunk %d Commit): %w", len(modelSlice), chunkID, e)
		}
	}

	return nil

}{{if .IsDeleted}}

// Delete marks an existing {{.Table.Name}} entry in the database as deleted
func (r *{{.Table.Name}}DAL) Delete(shard int64, {{.PrimaryKey | toArgName}} {{.IDType}}) error {
	var e error
	_, e = r.db[shard].Exec("UPDATE ` + "`{{.Table.Name}}` SET `IsDeleted` = 1 WHERE `{{.PrimaryKey}}` = ?" + `", {{.PrimaryKey | toArgName}})
	if e != nil {
		r.log.Errorf("{{.Table.Name}}DAL.Delete(%d) > %s", {{.PrimaryKey | toArgName}}, e.Error())
		return fmt.Errorf("{{.Table.Name}}DAL.Delete(%d): %w", {{.PrimaryKey | toArgName}}, e)
	} else {
		// r.log.Debugf("{{.Table.Name}}DAL.Delete(%d)", {{.PrimaryKey | toArgName}})
	}
	return nil
}

// DeleteMany marks {{.Table.Name}} objects in chunks as deleted
func (r {{.Table.Name}}DAL) DeleteMany(shard int64, modelSlice []*models.{{.Table.Name}}) error {

	var e error 

	// No records 
	if len(modelSlice) == 0 {
		return nil
	}

	// Don't use a transaction if only a single value
	if len(modelSlice) == 1 {
		e = r.Delete(shard, modelSlice[0].{{.PrimaryKey}})

		if e != nil {
			return fmt.Errorf("{{.Table.Name}}DeleteMany([](%d)): %w", len(modelSlice), e)
		}
		return nil
	}

	chunkSize := 25
	chunks := [][]*models.{{.Table.Name}}{}

	for i := 0; i < len(modelSlice); i += chunkSize {
		end := i + chunkSize
		if end > len(modelSlice) {
			end = len(modelSlice)
		}
		chunks = append(chunks, modelSlice[i:end])
	}

	for chunkID, chunk := range chunks {

		var tx *sql.Tx
		ctx := context.Background()
		tx, e = r.db[shard].BeginTx(ctx, nil)
		if e != nil {
			return fmt.Errorf("{{.Table.Name}}DeleteMany([](%d)) (Chunk %d): %w", len(modelSlice), chunkID, e)
		}

		for deleteID, model := range chunk {
{{if .IsLastUpdated}}
			model.LastUpdated = time.Now().UnixNano() / 1000000{{end}}
			_, e = tx.ExecContext(ctx, "UPDATE ` + "`{{.Table.Name}}` SET `IsDeleted`= 1 WHERE `{{.PrimaryKey}}` = ?" + `", model.{{.PrimaryKey}})
			if e != nil {
				r.log.Errorf("{{.Table.Name}}.DeleteMany([](%d)) (Chunk %d.%d) > %s", len(modelSlice), chunkID, deleteID, e.Error())
				return fmt.Errorf("{{.Table.Name}}DeleteMany([](%d)) (Chunk %d.%d): %w", len(modelSlice), chunkID, deleteID, e)
			} else {
				// r.log.Debugf("{{.Table.Name}}.DeleteMany([](%d)) (Chunk %d.%d)", len(modelSlice), chunkID, deleteID)
			}
		}

		if e != nil {
			return fmt.Errorf("{{.Table.Name}}DeleteMany([](%d)) (Chunk %d): %w", len(modelSlice), chunkID, e)
		}

		e = tx.Commit()

		if e != nil { 
			return fmt.Errorf("{{.Table.Name}}DeleteMany([](%d)) (Chunk %d Commit): %w", len(modelSlice), chunkID, e)
		} 
	}

	return nil

}{{end}}

// DeleteHard performs a SQL DELETE operation on a {{.Table.Name}} entry in the database
func (r *{{.Table.Name}}DAL) DeleteHard(shard int64, {{.PrimaryKey | toArgName}} {{.IDType}}) error {
	_, e := r.db[shard].Exec("DELETE FROM ` + "`{{.Table.Name}}`" + ` WHERE {{.PrimaryKey}} = ?", {{.PrimaryKey | toArgName}})
	if e != nil {
		r.log.Errorf("{{.Table.Name}}DAL.HardDelete(%d) > %s", {{.PrimaryKey | toArgName}}, e.Error())
		return fmt.Errorf("{{.Table.Name}}DAL.HardDelete(%d): %w", {{.PrimaryKey | toArgName}}, e)
	} else {
		// r.log.Debugf("{{.Table.Name}}DAL.HardDelete(%d)", {{.PrimaryKey | toArgName}})
	}
	return nil
}

// DeleteManyHard deletes {{.Table.Name}} objects in chunks
func (r {{.Table.Name}}DAL) DeleteManyHard(shard int64, modelSlice []*models.{{.Table.Name}}) error {

	var e error
	// No records 
	if len(modelSlice) == 0 {
		return nil
	}

	// Don't use a transaction if only a single value
	if len(modelSlice) == 1 {
		e = r.DeleteHard(shard, modelSlice[0].{{.PrimaryKey}})
		if e != nil {
			return fmt.Errorf("{{.Table.Name}}DeleteManyHard([](%d)): %w", len(modelSlice), e)
		}	
		return nil
	}

	chunkSize := 25
	chunks := [][]*models.{{.Table.Name}}{}

	for i := 0; i < len(modelSlice); i += chunkSize {
		end := i + chunkSize
		if end > len(modelSlice) {
			end = len(modelSlice)
		}
		chunks = append(chunks, modelSlice[i:end])
	}

	for chunkID, chunk := range chunks {

		var tx *sql.Tx
		ctx := context.Background()
		tx, e = r.db[shard].BeginTx(ctx, nil)
		if e != nil {
			return fmt.Errorf("{{.Table.Name}}DeleteManyHard([](%d)) (Chunk %d): %w", len(modelSlice), chunkID, e)
		}

		for deleteID, model := range chunk {

			_, e = tx.ExecContext(ctx, "DELETE FROM ` + "`{{.Table.Name}}` WHERE `{{.PrimaryKey}}` = ?" + `", model.{{.PrimaryKey}})
			if e != nil {
				r.log.Errorf("{{.Table.Name}}.DeleteManyHard([](%d)) (Chunk %d.%d) > %s", len(modelSlice), chunkID, deleteID, e.Error())
				break
			} else {
				// r.log.Debugf("{{.Table.Name}}.DeleteManyHard([](%d)) (Chunk %d.%d)", len(modelSlice), chunkID, deleteID)
			}
		}

		if e != nil {
			return fmt.Errorf("{{.Table.Name}}DeleteManyHard([](%d)) (Chunk %d): %w", len(modelSlice), chunkID, e)
		}

		e = tx.Commit()

		if e != nil {
			return fmt.Errorf("{{.Table.Name}}DeleteManyHard([](%d)) (Chunk %d Commit): %w", len(modelSlice), chunkID, e)
		}
	}

	return nil
}

// FromID gets a single {{.Table.Name}} object by its Primary Key
func (r *{{.Table.Name}}DAL) FromID(shard int64, {{.PrimaryKey | toArgName}} {{.IDType}}, mustExist bool) (*models.{{.Table.Name}}, error) {

	model, e := (&models.{{.Table.Name}}{}).Get(r.db[shard]).Where(query.EQ(models.{{.Table.Name}}_Column_{{.PrimaryKey}}, {{.PrimaryKey | toArgName}})).Run()

	if model == nil {
		if mustExist { 
			return nil, errors.NewRecordNotFoundError()
		}

		return nil, nil 
	}

	switch e { 
	case sql.ErrNoRows: 
		// r.log.Debugf("{{.Table.Name}}DAL.FromID(%d) > NOT FOUND", {{.PrimaryKey | toArgName}})

		if mustExist {
			e = errors.NewRecordNotFoundError()
			return nil, e 
		}

		return nil, nil
	case nil: 

		{{ if .IsDeleted}}if model.IsDeleted == 1 && mustExist { 
			return nil, errors.NewRecordNotFoundError()
		}{{end}}

		// r.log.Debugf("{{.Table.Name}}DAL.FromID(%d)", model.{{.PrimaryKey}})
		return model, nil 

	default: 
		r.log.Errorf("{{.Table.Name}}DAL.FromID(%d) > %s", {{.PrimaryKey | toArgName}}, e.Error())
		return nil, e 
	}
}

// FromIDs returns a slice of {{.Table.Name}} objects by a set of primary keys
func (r *{{.Table.Name}}DAL) FromIDs(shard int64, {{.PrimaryKey | toArgName}}s []{{.IDType}}) ([]*models.{{.Table.Name}}, error) {

	// No records 
	if len({{.PrimaryKey | toArgName}}s) == 0 {
		return []*models.{{.Table.Name}}{}, nil 
	}

	model, e := (&models.{{.Table.Name}}{}).Select(r.db[shard]).Where(
		query.INInt64(models.{{.Table.Name}}_Column_{{.PrimaryKey}}, {{.PrimaryKey | toArgName}}s),
	).Run()

	if e != nil {
		r.log.Errorf("{{.Table.Name}}DAL.FromIDs(%v) > %s", {{.PrimaryKey | toArgName}}s, e.Error())
		return []*models.{{.Table.Name}}{}, fmt.Errorf("{{.Table.Name}}DAL.FromIDs(%v): %w", {{.PrimaryKey | toArgName}}s, e)
	}
	
	// r.log.Debugf("{{.Table.Name}}DAL.FromIDs(%v)", {{.PrimaryKey | toArgName}}s)

	return model, nil 
}

// FromIDsMap returns a map of {{.Table.Name}} objects from primary keys indexed by a set of primary keys
func (r *{{.Table.Name}}DAL) FromIDsMap(shard int64, {{.PrimaryKey | toArgName}}s []{{.IDType}}) (map[{{.IDType}}]*models.{{.Table.Name}}, error) {


	model, e := r.FromIDs(shard, {{.PrimaryKey | toArgName}}s)
	if e != nil { 
		return map[{{.IDType}}]*models.{{.Table.Name}}{}, fmt.Errorf("{{.Table.Name}}DAL.FromIDsMap(%v): %w", {{.PrimaryKey | toArgName}}s, e)
	}
	
	result := make(map[{{.IDType}}]*models.{{.Table.Name}})
	for k := range model {
		result[model[k].{{.PrimaryKey}}] = model[k]
	}
	
	// r.log.Debugf("{{.Table.Name}}DAL.FromIDsMap(%v)", {{.PrimaryKey | toArgName}}s)

	return result, nil 
}


{{range $col := .SpecialColumns}}
{{if eq $col.DataType "vector"}}

// Get{{$col.Name}} gets the {{$col.Name}} column on a {{$.Table.Name}} object
func (r *{{$.Table.Name}}DAL) Get{{$col.Name}}(shard int64, {{$.PrimaryKey | toArgName}} {{$.IDType}}) ([]float64, error) {
	
	var vectorString null.String
	err := r.db[shard].QueryRow("SELECT VEC_ToText(` + "`{{.Name}}`) AS `{{.Name}}` FROM " + "`{{$.Table.Name}}`" + ` WHERE ` + "`{{$.PrimaryKey}}` = ?" + `", {{$.PrimaryKey | toArgName}}).Scan(&vectorString)
	if err != nil {
		r.log.Errorf("{{$.Table.Name}}DAL.Get{{$col.Name}}(%d) > %s", {{$.PrimaryKey | toArgName}}, err.Error())
		return nil, fmt.Errorf("{{$.Table.Name}}DAL.Get{{$col.Name}}(%d): %w", {{$.PrimaryKey | toArgName}}, err)
	}

	if !vectorString.Valid {
		return nil, nil 
	}

	stringSlice := strings.Split(vectorString.String, ", ")
	
	// 2. Create a slice to store the resulting float64 values
	floatSlice := make([]float64, len(stringSlice))
	
	// 3. Iterate over the string slice and parse each element
	for i, s := range stringSlice {
		// Use ParseFloat with a bitSize of 64 for float64
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			// Handle the error (e.g., log it, skip the value, or return)
			fmt.Printf("Error parsing float from string %s: %v\n", s, err)
			continue
		}
		floatSlice[i] = f
	}

	if err != nil {
		r.log.Errorf("{{$.Table.Name}}DAL.Get{{$col.Name}}Vector(%d) > %s", {{$.PrimaryKey | toArgName}}, err.Error())
		return nil, fmt.Errorf("{{$.Table.Name}}DAL.Get{{$col.Name}}(%d): %w", {{$.PrimaryKey | toArgName}}, err)
	}

	return floatSlice, nil 
}

// Set{{$col.Name}} sets the {{$col.Name}} column on a {{$.Table.Name}} object
func (r *{{$.Table.Name}}DAL) Set{{$col.Name}}Vector(shard int64, {{$.PrimaryKey | toArgName}} {{$.IDType}}, {{$col.Name | toArgName}} []float64) error {
	
	stringSlice := make([]string, 0, len({{$col.Name | toArgName}}))

	// 3. Iterate over the float64 slice and convert each element to a string
	for _, f := range {{$col.Name | toArgName}} {
		// 'f' format, -1 precision (smallest number of digits necessary), 64-bit float
		s := strconv.FormatFloat(f, 'f', -1, 64)
		stringSlice = append(stringSlice, s)
	}

	// 4. Join the string slice into a single string with a delimiter (e.g., a comma and space)
	result := "[" + strings.Join(stringSlice, ", ") + "]"

	_, e := r.db[shard].Exec("UPDATE ` + "`{{$.Table.Name}}` SET `{{$col.Name}}` = VEC_FromText(?) WHERE `{{$.PrimaryKey}}` = ?" + `", result, {{$.PrimaryKey | toArgName}})
	if e != nil {
		r.log.Errorf("{{$.Table.Name}}DAL.Set{{$col.Name}}(%d, %v) > %s", {{$.PrimaryKey | toArgName}}, {{$col.Name | toArgName}}, e.Error())
		return fmt.Errorf("{{$.Table.Name}}DAL.Set{{$col.Name}}(%d, %v): %w", {{$.PrimaryKey | toArgName}}, {{$col.Name | toArgName}}, e)
	} else {
		// r.log.Debugf("{{$.Table.Name}}DAL.Set{{$col.Name}}(%d, %v)", {{$.PrimaryKey | toArgName}}, {{$col.Name | toArgName}})
	}
	return nil 
}
{{end}}
{{end}}

{{range $col := .UpdateColumns}}
// Set{{$col.Name}} sets the {{$col.Name}} column on a {{$.Table.Name}} object
func (r *{{$.Table.Name}}DAL) Set{{$col.Name}}(shard int64, {{$.PrimaryKey | toArgName}} {{$.IDType}}, {{$col.Name | toArgName}} {{$col | dataTypeToGoTypeString}}) error {
	_, e := r.db[shard].Exec("UPDATE ` + "`{{$.Table.Name}}` SET `{{$col.Name}}` = ? WHERE `{{$.PrimaryKey}}` = ?" + `", {{$col.Name | toArgName}}, {{$.PrimaryKey | toArgName}})
	if e != nil {
		r.log.Errorf("{{$.Table.Name}}DAL.Set{{$col.Name}}(%d, %v) > %s", {{$.PrimaryKey | toArgName}}, {{$col.Name | toArgName}}, e.Error())
		return fmt.Errorf("{{$.Table.Name}}DAL.Set{{$col.Name}}(%d, %v): %w", {{$.PrimaryKey | toArgName}}, {{$col.Name | toArgName}}, e)
	} else {
		// r.log.Debugf("{{$.Table.Name}}DAL.Set{{$col.Name}}(%d, %v)", {{$.PrimaryKey | toArgName}}, {{$col.Name | toArgName}})
	}
	return nil 
}

// ManyFrom{{$col.Name}} returns a slice of {{$.Table.Name}} models from {{$col.Name}}
func (r *{{$.Table.Name}}DAL) ManyFrom{{$col.Name}}(shard int64, {{$col.Name | toArgName}} {{$col | dataTypeToGoTypeString}}, limit, offset int64, orderBy, orderDir string) ([]*models.{{$.Table.Name}}, error) {
	
	q := (&models.{{$.Table.Name}}{}).Select(r.db[shard]).Where(
		query.EQ(models.{{$.Table.Name}}_Column_{{$col.Name}}, {{$col.Name | toArgName}}), 
	)
	
	{{if $.IsDeleted}}
		q.Where(query.And(), query.EQ(models.{{$.Table.Name}}_Column_IsDeleted, 0)){{end}}

	if len(orderBy) > 0 { 
		q.OrderBy(query.Column(orderBy), query.OrderDirFromString(orderDir))
	}

	if limit > 0 { 
		q.Limit(limit, offset) 
	}

	collection, e := q.Run()

	if e != nil {
		r.log.Errorf("{{$.Table.Name}}DAL.ManyFrom{{$col.Name}}({{if or (eq $col.GoType "int") (eq $col.GoType "int64")}}%d{{else}}%s{{end}}, %d, %d, %s, %s) > %s", {{$col.Name | toArgName}}, limit, offset, orderBy, orderDir, e.Error())
		return nil, e 
	} 
	
	// r.log.Debugf("{{$.Table.Name}}DAL.ManyFrom{{$col.Name}}({{if or (eq $col.GoType "int") (eq $col.GoType "int64")}}%d{{else}}%s{{end}}, %d, %d, %s, %s)", {{$col.Name | toArgName}}, limit, offset, orderBy, orderDir)
	
	return collection, nil 
}

{{if or (eq $col.GoType "int64") (eq $col.GoType "int")}}
// ManyFrom{{$col.Name}}s returns a slice of {{$.Table.Name}} models from {{$col.Name}}s
func (r *{{$.Table.Name}}DAL) ManyFrom{{$col.Name}}s(shard int64, {{$col.Name | toArgName}}s []{{$col | dataTypeToGoTypeString}}, limit, offset int64, orderBy, orderDir string) ([]*models.{{$.Table.Name}}, error) {
		
	// No records 
	if len({{$col.Name | toArgName}}s) == 0 {
		return nil, nil 
	}

	q := (&models.{{$.Table.Name}}{}).Select(r.db[shard]).Where(
		query.INInt{{if eq $col.GoType "int64"}}64{{end}}(models.{{$.Table.Name}}_Column_{{$col.Name}}, {{$col.Name | toArgName}}s), {{if $.IsDeleted}}		
		query.And(), 
		query.EQ(models.{{$.Table.Name}}_Column_IsDeleted, 0),{{end}}
	)

	if len(orderBy) > 0 { 
		q.OrderBy(query.Column(orderBy), query.OrderDirFromString(orderDir))
	}
	
	if limit > 0 { 
		q.Limit(limit, offset) 
	}
	
	collection, e := q.Run()

	if e != nil {
		r.log.Errorf("{{$.Table.Name}}DAL.ManyFrom{{$col.Name}}s(%v, %d, %d, %s, %s) > %s", {{$col.Name | toArgName}}s, limit, offset, orderBy, orderDir, e.Error())
	} else {
		// r.log.Debugf("{{$.Table.Name}}DAL.ManyFrom{{$col.Name}}s(%d, %d, %s, %s)", limit, offset, orderBy, orderDir)
	}
	return collection, e 
}
{{end}}

// CountFrom{{$col.Name}} returns the number of {{$.Table.Name}} records from {{$col.Name}}
func (r *{{$.Table.Name}}DAL) CountFrom{{$col.Name}}(shard int64, {{$col.Name | toArgName}} {{$col | dataTypeToGoTypeString}}) (int64, error) {
	
	count, e := (&models.{{$.Table.Name}}{}).Count(r.db[shard]).Where(
		query.EQ(models.{{$.Table.Name}}_Column_{{$col.Name}}, {{$col.Name | toArgName}}),{{if $.IsDeleted}}		
		query.And(), 
		query.EQ(models.{{$.Table.Name}}_Column_IsDeleted, 0), {{end}}
	).Run()

	if e != nil {
		r.log.Errorf("{{$.Table.Name}}DAL.CountFrom{{$col.Name}}({{$col | dataTypeToFormatString}}) > %s", {{$col.Name | toArgName}}, e.Error())
	} else {
		// r.log.Debugf("{{$.Table.Name}}DAL.CountFrom{{$col.Name}}({{$col | dataTypeToFormatString}})", {{$col.Name | toArgName}})
	}

	return count, e
}

// SingleFrom{{$col.Name}} returns a single {{$.Table.Name}} record by its {{$col.Name}}
func (r *{{$.Table.Name}}DAL) SingleFrom{{$col.Name}}(shard int64, {{$col.Name | toArgName}} {{$col | dataTypeToGoTypeString}}, mustExist bool) (*models.{{$.Table.Name}}, error) {

	model, e := (&models.{{$.Table.Name}}{}).Get(r.db[shard]).Where(
		query.EQ(models.{{$.Table.Name}}_Column_{{$col.Name}}, {{$col.Name | toArgName}}),{{if $.IsDeleted}}
		query.And(), 
		query.EQ(models.{{$.Table.Name}}_Column_IsDeleted, 0), {{end}}
	).Run()

	if model == nil {
		if mustExist { 
			return nil, errors.NewRecordNotFoundError()
		}

		return nil, nil 
	}

	switch e { 
	case sql.ErrNoRows: 
		// r.log.Debugf("{{$.Table.Name}}DAL.SingleFrom{{$col.Name}}(%d) > NOT FOUND", {{$col.Name | toArgName}})

		if mustExist {
			e = errors.NewRecordNotFoundError()
			return nil, e 
		}

		return nil, nil
	case nil: 

		{{if $.IsDeleted}}if model.IsDeleted == 1 && mustExist { 
			return nil, errors.NewRecordNotFoundError()
		}{{end}}

		
		// r.log.Debugf("{{$.Table.Name}}DAL.SingleFrom{{$col.Name}}({{if $col.IsString}}%s{{end}}{{if not $col.IsString}}%d{{end}})", model.{{$col.Name}})
		return model, nil 

	default: 
		r.log.Errorf("{{$.Table.Name}}DAL.SingleFrom{{$col.Name}}({{if $col.IsString}}%s{{end}}{{if not $col.IsString}}%d{{end}}) > %s", {{$col.Name | toArgName}}, e.Error())
		return nil, e 
	}
}{{end}}

// ManyPaged returns a slice of {{.Table.Name}} models
func (r *{{.Table.Name}}DAL) ManyPaged(shard int64, limit, offset int64, orderBy, orderDir string) ([]*models.{{.Table.Name}}, error) {

	q := (&models.{{.Table.Name}}{}).Select(r.db[shard]){{if $.IsDeleted}}		
	q.Where(
		query.EQ(models.{{.Table.Name}}_Column_IsDeleted, 0),
	)
	{{end}}
	if len(orderBy) > 0 { 
		q.OrderBy(query.Column(orderBy), query.OrderDirFromString(orderDir))
	}
	
	if limit > 0 { 
		q.Limit(limit, offset) 
	}
	
	collection, e := q.Run()

	if e != nil {
		r.log.Errorf("{{.Table.Name}}DAL.ManyPaged(%d, %d, %s, %s) > %s", limit, offset, orderBy, orderDir, e.Error())
	} else {
		// r.log.Debugf("{{.Table.Name}}DAL.ManyPaged(%d, %d, %s, %s)", limit, offset, orderBy, orderDir)
	}
	return collection, e 
}
{{ if gt (len .StringColumns) 0 }}{{ range $col := .StringColumns}}
// Search{{$col.Name}} searches the {{$col.Name}} field in the {{$.Table.Name}} table
func (r *{{$.Table.Name}}DAL) Search{{$col.Name}}(shard int64, queryString string, limit, offset int64, leftOrRightOrBoth int) ([]*models.{{$.Table.Name}}, error) { 

	q := (&models.{{$.Table.Name}}{}).Select(r.db[shard]){{if $.IsDeleted}}		
	q.Where(
		query.EQ(models.{{$.Table.Name}}_Column_IsDeleted, 0),
	){{end}}

	// Search left
	switch leftOrRightOrBoth { 
	case 2: 
		// Search both 
		queryString = "%" + queryString + "%"
	case 1: 
		// Search right 
		queryString = "%" + queryString
	default:  
		// Search left (0)
		queryString += "%"
	} 

	q.Where(query.Like(models.{{$.Table.Name}}_Column_{{$col.Name}}, queryString))
	
	if limit > 0 { 
		q.Limit(limit, offset) 
	}
	
	collection, e := q.Run()

	if e != nil {
		r.log.Errorf("{{$.Table.Name}}DAL.Search(%s, %d) > %s", queryString, limit, e.Error())
	} else {
		// r.log.Debugf("{{$.Table.Name}}DAL.Search(%s, %d)", queryString, limit)
	}
	return collection, e 
}
{{ end }}
{{ end }}
`))
