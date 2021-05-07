package request

// UserLogActionType is the type of action in the user log
type UserLogActionType int64

const (
	// UserLogActionTypeAPI is an API Action Type
	UserLogActionTypeAPI UserLogActionType = iota + 1
	// UserLogActionTypeClient is an API Action Type
	UserLogActionTypeClient
)

// Int64 returns an int64 of the constant
func (lt UserLogActionType) Int64() int64 {
	return int64(lt)
}
