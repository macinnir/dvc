package types

type IUserContainer interface {
	// ID is the unique identifier for the user
	ID() int64
	// Account returns the AccountID for the user
	Account() int64
	// Activated returns whether or not the user is activated
	Activated() bool
	// Disabled returns whether or not the user is disabled
	Disabled() bool
	// Locked returns whether or not th user is locked
	Locked() bool
	// Permissions returns a slice of permission names
	Permissions() []string
}
