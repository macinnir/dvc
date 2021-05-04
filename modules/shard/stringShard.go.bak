package number

import (
	"hash/fnv"

	"github.com/google/uuid"
)

// GetShardIDFromStringID parses an string ID into a StringShardID object
func GetShardIDFromStringID(id string) int64 {
	h := fnv.New32a()
	h.Write([]byte(id[0:2]))
	return int64(h.Sum32()) % 64
}

// NewStringID returns a new string ID
func NewStringID() string {

	u, e := uuid.NewRandom()

	if e != nil {
		panic(e)
	}

	return u.String()
}
