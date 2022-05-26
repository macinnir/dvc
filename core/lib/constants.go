package lib

const (
	DefaultShard        = 0
	MetaDirectory       = ".dvc"
	PermissionsFile     = ".dvc/permissions.json"
	CorePermissionsFile = "core/permissions.json"
	SettingsFile        = ".dvc/settings.json"
	CoreSettingsFile    = "core/settings.json"
	ConfigFilePath      = ".dvc/config.json"
	SchemasFilePath     = ".dvc/schemas.json"
	CoreSchemasFilePath = "core/schemas.json"
	CoreSchemasName     = "core"
	CoreSchemasLogName  = "core_log"
	ChangeFilePath      = ".dvc/changes.log"
	TablesCacheFilePath = ".dvc/tables-cache.json"
	RoutesFilePath      = ".dvc/routes.json"

	GenDir                   = "gen"
	DalsGenDir               = "gen/dal"
	ModelsGenDir             = "gen/definitions/models"
	ServiceDefinitionsGenDir = "gen/definitions/services"
	DALDefinitionsGenDir     = "gen/definitions/dal"
	GoPermissionsDir         = "gen/permissions"

	AppServicesDir    = "app/services"
	AppDTOsDir        = "app/definitions/dtos"
	AppAggregatesDir  = "app/definitions/aggregates"
	CoreServicesDir   = "core/services"
	CoreDTOsDir       = "core/definitions/dtos"
	CoreConstantsDir  = "core/definitions/constants"
	CoreAggregatesDir = "core/definitions/aggregates"

	DefaultFileMode             = 0777
	ServicesDir                 = "gen/definitions/services"
	AppBootstrapFile            = "gen/definitions/services/app.go"
	RoutesBootstrapFile         = "gen/routes/routes.go"
	DALBootstrapFile            = "gen/definitions/dal.go"
	ControllersBootstrapGenFile = "gen/routes/controllers.go"
	LibRequests                 = "github.com/macinnir/dvc/core/lib/utils/request"
	LibUtils                    = "github.com/macinnir/dvc/core/lib/utils"
	// Generated Go file containing constants for all route permissions
)
