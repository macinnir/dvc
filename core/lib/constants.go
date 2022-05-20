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
	ModelsDir           = "models"

	ModelsGenDir             = "gen/definitions/models"
	ServiceDefinitionsGenDir = "gen/definitions/services"
	DALDefinitionsGenDir     = "gen/definitions/dal"
	CoreDTOsDir              = "core/definitions/dtos"
	AppDTOsDir               = "app/definitions/dtos"
	CoreConstantsDir         = "core/definitions/constants"

	ControllersBootstrapGenFile = "gen/controllers.go"
	IntegrationsDir             = "integrations"
	DALsDir                     = "dals"
	LibRequests                 = "github.com/macinnir/dvc/core/lib/utils/request"
	LibUtils                    = "github.com/macinnir/dvc/core/lib/utils"
	// Generated Go file containing constants for all route permissions
	GoPermissionsPath = "gen/permissions"
)
