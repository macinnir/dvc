package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"unicode"

	"github.com/macinnir/dvc/lib"
	"github.com/macinnir/dvc/modules/gen"
)

// Gen handles the `gen` command
func (c *Cmd) Gen(args []string) {

	if len(args) > 0 && args[0] == "help" {
		helpGen()
		return
	}

	var e error

	// fmt.Printf("Args: %v", args)
	if len(args) < 1 {
		lib.Error("Missing gen type [schema | models | repos | caches | ts] [[--force|-f]] [[--clean|-c]]", c.Options)
		os.Exit(1)
	}
	subCmd := Command(args[0])
	cwd, _ := os.Getwd()

	if len(args) > 0 {
		args = args[1:]
	}

	argLen := len(args)

	g := &gen.Gen{
		Options: c.Options,
		Config:  c.Config,
	}

	database := c.loadDatabase()
	c.genTableCache(database)

	switch subCmd {
	// case CommandGenSchema:
	// 	e = g.GenerateGoSchemaFile(c.Config.Dirs.Schema, database)
	// 	if e != nil {
	// 		lib.Error(e.Error(), c.Options)
	// 		os.Exit(1)
	// 	}
	// case CommandGenCaches:
	// 	fmt.Println("CommandGenCaches")
	// 	e = g.GenerateGoCacheFiles(c.Config.Dirs.Cache, database)
	// 	if e != nil {
	// 		lib.Error(e.Error(), c.Options)
	// 		os.Exit(1)
	// 	}
	case CommandGenRepos:
		c.GenRepos(g, database)
	case CommandGenDals:
		c.GenDals(g, database)
	case CommandGenDal:

		if argLen == 0 {
			lib.Error("Missing dal name", c.Options)
			os.Exit(1)
		}

		// lib.Error(fmt.Sprintf("Args: %s", args[0]), c.Options)
		table, e := database.FindTableByName(args[0])
		if e != nil {
			lib.Error(e.Error(), c.Options)
			os.Exit(1)
		}

		e = g.GenerateGoDAL(table, c.Config.Dirs.Dal)
		if e != nil {
			lib.Error(e.Error(), c.Options)
			os.Exit(1)
		}

		// if c.Options&lib.OptClean == lib.OptClean {
		// 	g.CleanGoDALs(c.Config.Dirs.Dal, database)
		// }

		// for _, table := range database.Tables {

		// 	lib.Debugf("Generating dal %s", g.Options, table.Name)
		// 	e = g.GenerateGoDAL(table, c.Config.Dirs.Dal)
		// 	if e != nil {
		// 		return
		// 	}
		// }

		// if e != nil {
		// 	lib.Error(e.Error(), c.Options)
		// 	os.Exit(1)
		// }

		// Create the dal bootstrap file in the dal repo
		e = g.GenerateDALsBootstrapFile(c.Config.Dirs.Dal, database)
		if e != nil {
			lib.Error(e.Error(), c.Options)
			os.Exit(1)
		}

		// e = g.GenerateDALSQL(c.Config.Dirs.Dal, database)
		// if e != nil {
		// 	lib.Error(e.Error(), c.Options)
		// 	os.Exit(1)
		// }

	case CommandGenInterfaces:
		c.GenInterfaces(g)
		// result, err := interfaces.Make(files, args.StructType, args.Comment, args.PkgName, args.IfaceName, args.IfaceComment, args.copyDocs, args.CopyTypeDoc)
	case CommandGenRoutes:
		c.GenRoutes(g)
	case CommandGenTests:

		serviceSuffix := "Service"
		srcDir := c.Config.Dirs.Services
		fmt.Println(srcDir)
		var files []os.FileInfo
		// DAL
		if files, e = ioutil.ReadDir(srcDir); e != nil {
			fmt.Println("ERROR", e.Error())
			return
		}
		for _, f := range files {

			// Filter out files that don't have upper case first letter names
			if !unicode.IsUpper([]rune(f.Name())[0]) {
				continue
			}

			srcFile := path.Join(srcDir, f.Name())

			// Remove the go extension
			baseName := f.Name()[0 : len(f.Name())-3]

			// Skip over EXT files
			if baseName[len(baseName)-3:] == "Ext" {
				continue
			}

			// Skip over test files
			if baseName[len(baseName)-5:] == "_test" {
				continue
			}

			// fmt.Println(baseName)

			if baseName == "DesignationService" {
				e = g.GenServiceTest(baseName[0:len(baseName)-len(serviceSuffix)], srcFile)

				if e != nil {
					panic(e)
				}
			}
		}

	case CommandGenModels:

		c.GenModels(g, database)

		// // Config.go
		// if _, e = os.Stat(path.Join(modelsDir, "Config.go")); os.IsNotExist(e) {
		// 	lib.Debugf("Generating default Config.go file at %s", c.Options, path.Join(modelsDir, "Config.go"))
		// 	g.GenerateDefaultConfigModelFile(modelsDir)
		// }

		// config.json
		// if _, e = os.Stat(path.Join(cwd, "config.json")); os.IsNotExist(e) {
		// 	lib.Debugf("Generating default config.json file at %s", c.Options, path.Join(cwd, "config.json"))
		// 	g.GenerateDefaultConfigJsonFile(cwd)
		// }

	case CommandGenServices:
		g.GenerateServiceInterfaces(c.Config.Dirs.Definitions, c.Config.Dirs.Services)
		g.GenerateServiceBootstrapFile(c.Config.Dirs.Services)

	case CommandGenApp:
		g.GenerateGoApp(cwd)
	// case CommandGenCLI:
	// g.GenerateGoCLI(cwd)
	case CommandGenAPI:
		// g.GenerateGoAPI(cwd)
		g.GenerateAPIRoutes(c.Config.Dirs.API)
	case CommandGenTSPerms:
		tsFile := g.BuildTypescriptPermissions()
		fmt.Println(tsFile)
	case "ts":
		g.GenerateTypescriptTypesFile(c.Config.Dirs.Typescript, database)
	default:
		lib.Errorf("Unknown output type: `%s`", c.Options, subCmd)
		os.Exit(1)
	}
}

// GenRepos generates repos
func (c *Cmd) GenRepos(g *gen.Gen, database *lib.Database) {

	var e error

	if c.Options&lib.OptClean == lib.OptClean {
		g.CleanGoRepos(c.Config.Dirs.Repos, database)
	}

	e = g.GenerateGoRepoFiles(c.Config.Dirs.Repos, database)
	if e != nil {
		lib.Error(e.Error(), c.Options)
		os.Exit(1)
	}

	e = g.GenerateReposBootstrapFile(c.Config.Dirs.Repos, database)
	if e != nil {
		lib.Error(e.Error(), c.Options)
		os.Exit(1)
	}

	lib.Debug("Generating repo interfaces at "+c.Config.Dirs.Definitions, c.Options)
	lib.EnsureDir(c.Config.Dirs.Definitions)
	e = g.GenerateRepoInterfaces(database, c.Config.Dirs.Definitions)
	if e != nil {
		lib.Error(e.Error(), c.Options)
		os.Exit(1)
	}
}

// GenModels generates models
func (c *Cmd) GenModels(g *gen.Gen, database *lib.Database) {

	fmt.Println("Generating models...")
	force := c.Options&lib.OptForce == lib.OptForce

	var e error

	modelsDir := path.Join(c.Config.Dirs.Definitions, "models")
	if c.Options&lib.OptClean == lib.OptClean {
		g.CleanGoModels(modelsDir, database)
	}

	for _, table := range database.Tables {

		// If the model has been hashed before...
		if _, ok := c.existingModels.Models[table.Name]; ok {

			// And the hash hasn't changed, skip...
			if c.newModels[table.Name] == c.existingModels.Models[table.Name] && !force {

				// fmt.Printf("Table `%s` hasn't changed! Skipping...\n", table.Name)
				continue
			}
		}

		// Update the models cache
		c.existingModels.Models[table.Name] = c.newModels[table.Name]

		fmt.Printf("Generating model `%s`\n", table.Name)
		e = g.GenerateGoModel(modelsDir, table)
		if e != nil {
			lib.Error(e.Error(), c.Options)
			os.Exit(1)
		}
	}

	g.GenMeta("meta", database)

	c.saveTableCache()

	// e = g.GenerateDALSQL(c.Config.Dirs.Dal, database)
	// if e != nil {
	// 	lib.Error(e.Error(), c.Options)
	// 	os.Exit(1)
	// }
}

// GenDals generates dals
func (c *Cmd) GenDals(g *gen.Gen, database *lib.Database) {

	force := c.Options&lib.OptForce == lib.OptForce
	clean := c.Options&lib.OptClean == lib.OptClean

	if clean {
		g.CleanGoDALs(c.Config.Dirs.Dal, database)
	}

	fmt.Println("Generating dals...")
	var e error

	// Loop through the schema's tables and build md5 hashes of each to check against
	for _, table := range database.Tables {

		// If the model has been hashed before...
		if _, ok := c.existingModels.Dals[table.Name]; ok {

			// And the hash hasn't changed, skip...
			if c.newModels[table.Name] == c.existingModels.Dals[table.Name] && !force {

				// fmt.Printf("Table `%s` hasn't changed! Skipping...\n", table.Name)
				continue
			}
		}

		// Update the dals cache
		c.existingModels.Dals[table.Name] = c.newModels[table.Name]

		fmt.Printf("Generating %sDAL...\n", table.Name)
		e = g.GenerateGoDAL(table, c.Config.Dirs.Dal)
		if e != nil {
			lib.Error(e.Error(), c.Options)
			os.Exit(1)
		}
	}

	c.saveTableCache()

	// Create the dal bootstrap file in the dal repo
	e = g.GenerateDALsBootstrapFile(c.Config.Dirs.Dal, database)
	if e != nil {
		lib.Error(e.Error(), c.Options)
		os.Exit(1)
	}

	// e = g.GenerateDALSQL(c.Config.Dirs.Dal, database)
	// if e != nil {
	// 	lib.Error(e.Error(), c.Options)
	// 	os.Exit(1)
	// }
}

// GenInterfaces generates interfaces
func (c *Cmd) GenInterfaces(g *gen.Gen) {

	fmt.Println("Generating interfaces...")
	var e error

	genInterfaces := func(srcDir, srcType string) (e error) {

		var files []os.FileInfo
		// DAL
		if files, e = ioutil.ReadDir(srcDir); e != nil {
			return
		}
		for _, f := range files {

			// Filter out files that don't have upper case first letter names
			if !unicode.IsUpper([]rune(f.Name())[0]) {
				continue
			}

			srcFile := path.Join(srcDir, f.Name())

			// Remove the go extension
			baseName := f.Name()[0 : len(f.Name())-3]

			// Skip over EXT files
			if baseName[len(baseName)-3:] == "Ext" {
				continue
			}

			// Skip over test files
			if baseName[len(baseName)-5:] == "_test" {
				continue
			}

			// srcFile := path.Join(c.Config.Dirs.Dal, baseName + ".go")
			destFile := path.Join(c.Config.Dirs.Definitions, srcType, "I"+baseName+".go")
			interfaceName := "I" + baseName
			packageName := srcType

			srcFiles := []string{srcFile}
			// var src []byte
			// if src, e = ioutil.ReadFile(srcFile); e != nil {
			// 	return
			// }

			// Check if EXT file exists
			extFile := srcFile[0:len(srcFile)-3] + "Ext.go"
			if _, e = os.Stat(extFile); e == nil {
				srcFiles = append(srcFiles, extFile)
				// concatenate the contents of the ext file with the contents of the regular file
				// var extSrc []byte
				// if extSrc, e = ioutil.ReadFile(extFile); e != nil {
				// 	return
				// }
				// src = append(src, extSrc...)
			}

			var i []byte
			i, e = gen.GenInterfaces(
				srcFiles,
				baseName,
				"Generated Code; DO NOT EDIT.",
				packageName,
				interfaceName,
				fmt.Sprintf("%s describes the %s", interfaceName, baseName),
				true,
				true,
			)
			if e != nil {
				fmt.Println("ERROR", e.Error())
				return
			}

			// fmt.Println("Generating ", destFile)
			// fmt.Println("Writing to: ", destFile)

			ioutil.WriteFile(destFile, i, 0644)

			// fmt.Println("Name: ", baseName, "Path: ", srcFile)

		}

		return
	}

	e = genInterfaces(c.Config.Dirs.Dal, "dal")
	if e != nil {
		fmt.Println("ERROR", e.Error())
		os.Exit(1)
	}

	e = genInterfaces(c.Config.Dirs.Services, "services")
	if e != nil {
		fmt.Println("ERROR", e.Error())
		os.Exit(1)
	}
}

// GenRoutes generates routes
func (c *Cmd) GenRoutes(g *gen.Gen) {

	fmt.Println("Generating routes...")
	var e error

	e = g.GenRoutes()
	if e != nil {
		lib.Error(e.Error(), c.Options)
		os.Exit(1)
	}
}
