package lib

// DatabaseType is the type of database to be used.
type DatabaseType string

const (
	// DatabaseTypeMysql is the MySQL flavor of database
	DatabaseTypeMysql DatabaseType = "mysql"
)

// Options are the available runtime flags
type Options uint

// Command is the command line functionality
type Command struct {
	Options Options
}

// Changeset represents all of the changes in an environment and their changes
type Changeset struct {
	ChangeFiles map[string]ChangeFile
	Signature   string
}

// ChangeFile represents both a physical file on the local file system
// along with the entry in the changefile database
type ChangeFile struct {
	ID          int64
	DateCreated int64
	Hash        string
	Name        string
	DVCSetID    int64
	IsRun       bool
	IsDeleted   bool
	Content     string
	FullPath    string
	Ordinal     int
}

const (
	// OptLogInfo triggers verbose logging
	OptLogInfo = 1 << iota
	// OptLogDebug triggers extremely verbose logging
	OptLogDebug
	// OptSilent suppresses all logging
	OptSilent
	// OptReverse reverses the function
	OptReverse
	// OptSummary shows a summary of the actions instead of a raw stdout dump
	OptSummary
	// OptClean cleans
	OptClean
	// OptForce forces
	OptForce
)

// RoutesJSONContainer is a container for JSON Routes
type RoutesJSONContainer struct {
	Routes      map[string]*ControllerRoute  `json:"routes"`
	DTOs        map[string]map[string]string `json:"dtos"`
	Models      map[string]map[string]string `json:"models"`
	Aggregates  map[string]map[string]string `json:"aggregates"`
	Constants   map[string][]string          `json:"constants"`
	Permissions map[string]string            `json:"permissions"`
}

// Controller represents a REST controller
type Controller struct {
	Name              string             `json:"Name"`
	Description       string             `json:"Description"`
	Path              string             `json:"-"`
	Routes            []*ControllerRoute `json:"Routes"`
	HasDTOsImport     bool               `json:"-"`
	HasResponseImport bool               `json:"-"`
	PermCount         int                `json:"-"`
	Package           string             `json:"Package"`
}

// ControllerRoute represents a route inside a REST controller
type ControllerRoute struct {
	Package         string                 `json:"Package"`
	Controller      string                 `json:"Controller"`
	Name            string                 `json:"Name"`
	Description     string                 `json:"Description"`
	Raw             string                 `json:"Path"`
	Path            string                 `json:"-"`
	Method          string                 `json:"Method"`
	Params          []ControllerRouteParam `json:"Params"`
	Queries         []ControllerRouteQuery `json:"Queries"`
	RequiredQueries []ControllerRouteQuery `json:"RequiredQueries"`
	IsAuth          bool                   `json:"IsAuth"`
	BodyType        string                 `json:"BodyType"`
	BodyFormat      string                 `json:"BodyFormat"`
	HasBody         bool                   `json:"HasBody"`
	ResponseType    string                 `json:"ResponseType"`
	ResponseFormat  string                 `json:"ResponseFormat"`
	ResponseCode    int                    `json:"ResponseCode"`
	Permission      string                 `json:"Permission"`
	ControllerName  string                 `json:"ControllerName"`
	FileName        string                 `json:"FileName"`
	LineNo          int                    `json:"LineNo"`
}

// ControllerRouteParam represents a param inside a controller route
type ControllerRouteParam struct {
	Name    string
	Pattern string
	Type    string
}

// ControllerRouteQuery represents a query inside a controller route
type ControllerRouteQuery struct {
	Name         string
	Pattern      string
	Type         string
	VariableName string
	ValueRaw     string
}

// DocRoute is a route in the documentation
type DocRoute struct {
	Name           string
	Description    string
	Method         string
	Path           string
	HasBody        bool
	BodyType       string
	BodyFormat     string
	ResponseType   string
	ResponseFormat string
	ResponseCode   int
}

func newDocRoute(route ControllerRoute) (docRoute *DocRoute) {
	docRoute = &DocRoute{
		Name:        route.Name,
		Description: route.Description,
		Method:      route.Method,
		Path:        route.Path,
		HasBody:     route.HasBody,
		BodyType:    route.BodyType,
	}

	return
}

// DocRouteParam represents a parameter inside a route for documentation
type DocRouteParam struct {
	Name    string
	Pattern string
	Type    string
}

// DocRouteQuery represents a query inside a route for documentation
type DocRouteQuery struct {
	Name    string
	Pattern string
	Type    string
}
