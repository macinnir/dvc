package lib

const (
	DefaultShard        = 0
	MetaDirectory       = ".dvc"
	PermissionsFile     = ".dvc/permissions.json"
	ConfigFilePath      = ".dvc/config.json"
	SchemasFilePath     = ".dvc/schemas.json"
	ChangeFilePath      = ".dvc/changes.log"
	TablesCacheFilePath = ".dvc/tables-cache.json"
	RoutesFilePath      = ".dvc/routes.json"
	ModelsDir           = "models"
	IntegrationsDir     = "integrations"
	DALsDir             = "dals"
	// Generated Go file containing constants for all route permissions
	GoPermissionsPath = "gen/permissions"
)
