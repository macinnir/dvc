package nonce

import (
	"github.com/google/uuid"
	"github.com/rs/xid"
)

// GUID returns a UUIDv4 based on RFC 4122
func GUID() string {
	guid := uuid.New()
	return guid.String()
}

// ShortID returns a shorter globally unique ID
// Uses the MongDB Object ID algorithm using base64 serialization to make it shorter
func ShortID() string {
	id := xid.New()
	return id.String()
}
