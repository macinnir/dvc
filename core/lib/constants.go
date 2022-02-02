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
	ModelsGenDir        = "gen/definitions/models"
	IntegrationsDir     = "integrations"
	DALsDir             = "dals"
	// Generated Go file containing constants for all route permissions
	GoPermissionsPath = "gen/permissions"
)
