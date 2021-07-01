package types

type IUserContainer interface {
	ID() int64
	Activated() bool
	Disabled() bool
	Locked() bool
	Permissions() []string
}
