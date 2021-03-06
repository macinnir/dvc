package gen

import (
	"fmt"
	"html/template"
	"os"
	"path"

	"github.com/macinnir/dvc/core/lib"
)

// GenerateGoApp generates the base app code
func (g *Gen) GenerateGoApp(dir string) error {

	tpl := `
package {{ .BasePackage }} 

import (
	"path"
	"{{ .BasePackage }}/repos"
	"{{ .BasePackage }}/services"
	"{{ .BasePackage }}/definitions/models"
	"github.com/macinnir/dvc/core/lib/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

// NewApp returns a new instance of the base app
func NewApp(appName string, cwd string) *App {
	app := &App{ appName: appName, cwd: cwd }
	return app
}

// App is the base app
type App struct {
	Config *models.Config
	cwd 	 string
	appName  string
	logger   *utils.Logger
	Repos    *services.Repos
	Cache    *services.Repos
	Store    utils.IStore
	Services *services.Services
}

// InitConfig initializes the config
func (a *App) InitConfig() {
	a.Config = &models.Config{} 
	utils.LoadConfig(path.Join(a.cwd, "config.json"), a.Config)
}

// InitLogging initializes the logging
func (a *App) InitLogging() {
	log.Println("Init logging...")
	logFile := path.Join(a.cwd, "logs", a.appName + ".log")
	a.logger = utils.NewLogger(a.appName, logFile)
	log.SetOutput(a.logger)
}

func (a *App) connectToDB() *sqlx.DB {
	connectionString := a.Config.DBUser + ":" + a.Config.DBPass + "@tcp(" + a.Config.DBHost + ")/" + a.Config.DBName + "?charset=utf8"
	return sqlx.MustConnect(
		"mysql",
		connectionString,
	)
}

	// dbConnection.SetConnMaxLifetime(time.Second * 10)

	// if e != nil {
	// 	log.Fatalf("Could not connect to the database with connection string: %s", connectionString)
	// }

// InitRepos initializes the repository layer
func (a *App) InitRepos() {
	conn := a.connectToDB() 
	a.Repos = repos.Bootstrap(conn)
}

// InitStore initializes the redis store
func (a *App) InitStore() {
	a.Store = utils.NewStore(a.Config.RedisHost, a.Config.RedisPassword, int(a.Config.RedisDB))
}

// InitServices initializes the services layer
func (a *App) InitServices() {
	a.Services = services.NewServices(a.Config, a.Repos, a.Store)
}
// Finish closes out the base app
func (a *App) Finish() {
	a.logger.Finish()
}
`
	p := path.Join(dir, "app.go")
	t := template.Must(template.New("app").Parse(tpl))
	f, err := os.Create(p)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		return err
	}

	err = t.Execute(f, g.Config)
	if err != nil {
		fmt.Println("Execute Error: ", err.Error())
		return err
	}

	if err = lib.FmtGoCode(p); err != nil {
		return err
	}

	return nil
}
