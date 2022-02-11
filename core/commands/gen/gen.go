package gen

import (
	"errors"

	"github.com/macinnir/dvc/core/lib"
	"github.com/macinnir/dvc/core/lib/fetcher"
	"github.com/macinnir/dvc/core/lib/gen"
	"github.com/macinnir/dvc/core/lib/gen/routes"
	"github.com/macinnir/dvc/core/lib/gen/typescript"
	"go.uber.org/zap"
)

const CommandName = "gen"

// Gen handles the `gen` command
func Cmd(log *zap.Logger, config *lib.Config, args []string) error {

	if len(args) == 0 {
		log.Warn("Missing gen type")
		return errors.New("Missing gen type")
	}

	cmd := args[0]
	force := false
	clean := false

	for k := range args {
		switch args[k] {
		case "-f", "--force":
			force = true
		case "-c", "--clean":
			clean = true
		}
	}

	switch cmd {
	case "models":
		gen.GenModels(config, force, clean)
	case "dals":
		gen.GenDALs("gen/dal", config, force, clean)
	case "interfaces":
		gen.GenServicesBootstrap(config)
		gen.GenInterfaces("gen/dal", "gen/definitions/dal")
		gen.GenInterfaces("core/services", "gen/definitions/services")
		gen.GenInterfaces("app/services", "gen/definitions/services")
	case "routes":
		cf := fetcher.NewControllerFetcher()
		controllers, dirs, e := cf.FetchAll()
		if e != nil {
			return e
		}

		if e := routes.GenRoutesAndPermissions(controllers, dirs, config); e != nil {
			return e
		}

		if e := routes.GenTSRoutes(controllers, config); e != nil {
			return e
		}

	case "ts":

		var e error
		var r *lib.RoutesJSONContainer

		r, e = routes.LoadRoutes(config)

		if e != nil {
			return e
		}

		if e = typescript.GenerateTypescriptModels(config, r); e != nil {
			return e
		}
		if e = typescript.GenerateTypesriptDTOs(config, r); e != nil {
			return e
		}

		tg := typescript.NewTypescriptGenerator(config, r)

		if e = tg.GenerateTypesriptAggregates(); e != nil {
			return e
		}

	// case "tsdtos":
	// 	fmt.Println("Generating Typescript DTOs")
	// 	typescript.GenerateTypesriptDTOs(config)
	// case "tsaggregates":
	// 	fmt.Println("Generating Typescript Aggregates")
	// 	typescript.GenerateTypesriptAggregates(config)
	case "tsperms":
		if e := gen.GenTSPerms(config); e != nil {
			return e
		}
	case "goperms":
		if e := gen.GenGoPerms(config); e != nil {
			return e
		}

	default:
		return errors.New("unknown gen type")

	}

	return nil

}

// 	var e error

// 	// fmt.Printf("Args: %v", args)
// 	if len(args) < 1 {
// 		lib.Error("Missing gen type [schema | models | repos | caches | ts] [[--force|-f]] [[--clean|-c]]", c.Options)
// 		os.Exit(1)
// 	}
// 	subCmd := Command(args[0])
// 	cwd, _ := os.Getwd()

// 	if len(args) > 0 {
// 		args = args[1:]
// 	}

// 	argLen := len(args)

// 	g := &gen.Gen{
// 		Config: config,
// 	}

// 	database := c.loadDatabase()
// 	c.genTableCache(database)

// 	switch subCmd {

// 	// case CommandGenCaches:
// 	// 	fmt.Println("CommandGenCaches")
// 	// 	e = g.GenerateGoCacheFiles(c.Config.Dirs.Cache, database)
// 	// 	if e != nil {
// 	// 		lib.Error(e.Error(), c.Options)
// 	// 		os.Exit(1)
// 	// 	}
// 	case "dals":
// 		c.GenDals(g, database)
// 	case CommandGenDal:

// 		if argLen == 0 {
// 			lib.Error("Missing dal name", c.Options)
// 			os.Exit(1)
// 		}

// 		// lib.Error(fmt.Sprintf("Args: %s", args[0]), c.Options)
// 		table, e := database.FindTableByName(args[0])
// 		if e != nil {
// 			lib.Error(e.Error(), c.Options)
// 			os.Exit(1)
// 		}

// 		e = g.GenerateGoDAL(table, c.Config.Dirs.Dal)
// 		if e != nil {
// 			lib.Error(e.Error(), c.Options)
// 			os.Exit(1)
// 		}

// 		// if c.Options&lib.OptClean == lib.OptClean {
// 		// 	g.CleanGoDALs(c.Config.Dirs.Dal, database)
// 		// }

// 		// for _, table := range database.Tables {

// 		// 	lib.Debugf("Generating dal %s", g.Options, table.Name)
// 		// 	e = g.GenerateGoDAL(table, c.Config.Dirs.Dal)
// 		// 	if e != nil {
// 		// 		return
// 		// 	}
// 		// }

// 		// if e != nil {
// 		// 	lib.Error(e.Error(), c.Options)
// 		// 	os.Exit(1)
// 		// }

// 		// Create the dal bootstrap file in the dal repo
// 		e = g.GenerateDALsBootstrapFile(c.Config.Dirs.Dal, database)
// 		if e != nil {
// 			lib.Error(e.Error(), c.Options)
// 			os.Exit(1)
// 		}

// 		// e = g.GenerateDALSQL(c.Config.Dirs.Dal, database)
// 		// if e != nil {
// 		// 	lib.Error(e.Error(), c.Options)
// 		// 	os.Exit(1)
// 		// }

// 	case CommandGenInterfaces:
// 		c.GenInterfaces(g)
// 		// result, err := interfaces.Make(files, args.StructType, args.Comment, args.PkgName, args.IfaceName, args.IfaceComment, args.copyDocs, args.CopyTypeDoc)
// 	case CommandGenRoutes:
// 		c.GenRoutes(g)
// 	case CommandGenTests:

// 		serviceSuffix := "Service"
// 		srcDir := c.Config.Dirs.Services
// 		fmt.Println(srcDir)
// 		var files []os.FileInfo
// 		// DAL
// 		if files, e = ioutil.ReadDir(srcDir); e != nil {
// 			fmt.Println("ERROR", e.Error())
// 			return
// 		}
// 		for _, f := range files {

// 			// Filter out files that don't have upper case first letter names
// 			if !unicode.IsUpper([]rune(f.Name())[0]) {
// 				continue
// 			}

// 			srcFile := path.Join(srcDir, f.Name())

// 			// Remove the go extension
// 			baseName := f.Name()[0 : len(f.Name())-3]

// 			// Skip over EXT files
// 			if baseName[len(baseName)-3:] == "Ext" {
// 				continue
// 			}

// 			// Skip over test files
// 			if baseName[len(baseName)-5:] == "_test" {
// 				continue
// 			}

// 			// fmt.Println(baseName)

// 			if baseName == "DesignationService" {
// 				e = g.GenServiceTest(baseName[0:len(baseName)-len(serviceSuffix)], srcFile)

// 				if e != nil {
// 					panic(e)
// 				}
// 			}
// 		}

// 	case CommandGenModels:

// 		c.GenModels(g, database)

// 		// // Config.go
// 		// if _, e = os.Stat(path.Join(modelsDir, "Config.go")); os.IsNotExist(e) {
// 		// 	lib.Debugf("Generating default Config.go file at %s", c.Options, path.Join(modelsDir, "Config.go"))
// 		// 	g.GenerateDefaultConfigModelFile(modelsDir)
// 		// }

// 		// config.json
// 		// if _, e = os.Stat(path.Join(cwd, "config.json")); os.IsNotExist(e) {
// 		// 	lib.Debugf("Generating default config.json file at %s", c.Options, path.Join(cwd, "config.json"))
// 		// 	g.GenerateDefaultConfigJsonFile(cwd)
// 		// }

// 	case CommandGenServices:
// 		g.GenerateServiceInterfaces(c.Config.Dirs.Definitions, c.Config.Dirs.Services)
// 		g.GenerateServiceBootstrapFile(c.Config.Dirs.Services)

// 	case CommandGenApp:
// 		g.GenerateGoApp(cwd)
// 	// case CommandGenCLI:
// 	// g.GenerateGoCLI(cwd)
// 	case CommandGenAPI:
// 		// g.GenerateGoAPI(cwd)
// 		g.GenerateAPIRoutes(c.Config.Dirs.API)
// 	case CommandGenTSPerms:
// 		tsFile := g.BuildTypescriptPermissions()
// 		fmt.Println(tsFile)
// 	case CommandGenAPITests:
// 		c.GenAPITests(g)
// 	case "ts":
// 		g.GenerateTypescriptTypesFile(c.Config.Dirs.Typescript, database)
// 	default:
// 		lib.Errorf("Unknown output type: `%s`", c.Options, subCmd)
// 		os.Exit(1)
// 	}
// }

// 	g.GenMeta("meta", database)

// 	c.saveTableCache()

// 	// e = g.GenerateDALSQL(c.Config.Dirs.Dal, database)
// 	// if e != nil {
// 	// 	lib.Error(e.Error(), c.Options)
// 	// 	os.Exit(1)
// 	// }
// }

// // GenDals generates dals
// func (c *Cmd) GenDals(g *gen.Gen, database *schema.Schema) {

// 	force := c.Options&lib.OptForce == lib.OptForce
// 	clean := c.Options&lib.OptClean == lib.OptClean

// 	if clean {
// 		g.CleanGoDALs(c.Config.Dirs.Dal, database)
// 	}

// 	fmt.Println("Generating dals...")
// 	var e error

// 	// Loop through the schema's tables and build md5 hashes of each to check against
// 	for _, table := range database.Tables {

// 		// If the model has been hashed before...
// 		if _, ok := c.existingModels.Dals[table.Name]; ok {

// 			// And the hash hasn't changed, skip...
// 			if c.newModels[table.Name] == c.existingModels.Dals[table.Name] && !force {

// 				// fmt.Printf("Table `%s` hasn't changed! Skipping...\n", table.Name)
// 				continue
// 			}
// 		}

// 		// Update the dals cache
// 		c.existingModels.Dals[table.Name] = c.newModels[table.Name]

// 		fmt.Printf("Generating %sDAL...\n", table.Name)
// 		e = g.GenerateGoDAL(table, c.Config.Dirs.Dal)
// 		if e != nil {
// 			lib.Error(e.Error(), c.Options)
// 			os.Exit(1)
// 		}
// 	}

// 	c.saveTableCache()

// 	// Create the dal bootstrap file in the dal repo
// 	e = g.GenerateDALsBootstrapFile(c.Config.Dirs.Dal, database)
// 	if e != nil {
// 		lib.Error(e.Error(), c.Options)
// 		os.Exit(1)
// 	}

// 	// e = g.GenerateDALSQL(c.Config.Dirs.Dal, database)
// 	// if e != nil {
// 	// 	lib.Error(e.Error(), c.Options)
// 	// 	os.Exit(1)
// 	// }
// }

// // GenInterfaces generates interfaces
// func (c *Cmd) GenInterfaces(g *gen.Gen) {

// 	fmt.Println("Generating interfaces...")
// 	var e error

// e = genInterfaces(c.Config.Dirs.Dal, "dal")
// if e != nil {
// 	fmt.Println("ERROR", e.Error())
// 	os.Exit(1)
// }

// 	e = genInterfaces(c.Config.Dirs.Services, "services")
// 	if e != nil {
// 		fmt.Println("ERROR", e.Error())
// 		os.Exit(1)
// 	}
// }

// // GenRoutes generates routes
// func (c *Cmd) GenRoutes(g *gen.Gen) {

// 	fmt.Println("Generating routes...")
// 	var e error

// 	e = g.GenRoutes()
// 	if e != nil {
// 		lib.Error(e.Error(), c.Options)
// 		os.Exit(1)
// 	}
// }

// // GenAPITests generates routes
// func (c *Cmd) GenAPITests(g *gen.Gen) {

// 	fmt.Println("Generating api tests...")
// 	var e error

// 	e = g.GenAPITests()
// 	if e != nil {
// 		lib.Error(e.Error(), c.Options)
// 		os.Exit(1)
// 	}
// }
