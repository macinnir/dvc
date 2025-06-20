package testgen

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Time struct {
	time.Time
}

func (t *Time) String() string {
	return fmt.Sprint(t.Time.UnixNano() / 1000000)
}

// MarshalJSON implements the json.Marshaller interface
func (t *Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time.UnixNano() / 1000000)
}

// UnmarshalJSON implements the json.Unmarshaller interface
func (t *Time) UnmarshalJSON(data []byte) error {
	var i int64
	if e := json.Unmarshal(data, &i); e != nil {
		return e
	}

	t.Time = time.Unix(0, i*1000000)
	return nil
}

// // MarshalText implements encoding.TextMarshaler.
// // It will encode a blank string when this String is null.
// func (t *Time) MarshalText() ([]byte, error) {
// 	return []byte(fmt.Sprint(t.Time.))
// 	// if !s.Valid {
// 	// 	return []byte{}, nil
// 	// }
// 	// return []byte(s.String), nil
// }

// // UnmarshalText implements encoding.TextUnmarshaler.
// // It will unmarshal to a null String if the input is a blank string.
// func (s *String) UnmarshalText(text []byte) error {
// 	s.String = string(text)
// 	s.Valid = s.String != ""
// 	return nil
// }

func NewTime(t time.Time) Time {
	return Time{
		Time: t,
	}
}

type NullString struct {
	sql.NullString
}

func (n *NullString) String() string {
	return n.NullString.String
}

func NewNullString(str string) NullString {
	return NullString{
		NullString: sql.NullString{
			String: str,
			Valid:  true,
		},
	}
}

func (n *NullString) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.String)
}

func (n *NullString) UnMarshalJSON(data []byte) error {
	var s string
	if e := json.Unmarshal(data, &s); e != nil {
		return e
	}
	n.NullString = sql.NullString{String: s, Valid: true}
	return nil
}
