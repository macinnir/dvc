package hasher

import (
	"github.com/speps/go-hashids/v2"
)

var (
	MinLength = 30
	Salt      = "this is my salt"
)

type Hasher struct {
	salt      string
	minLength int
}

func NewHasher(
	salt string,
	minLength int,
) Hasher {
	return Hasher{
		salt,
		minLength,
	}
}

func (h *Hasher) Hash(vals []int) (string, error) {

	hd, e := hashids.NewWithData(&hashids.HashIDData{
		Salt:      h.salt,
		Alphabet:  hashids.DefaultAlphabet,
		MinLength: h.minLength,
	})

	if e != nil {
		return "", e
	}

	encodedString, e := hd.Encode(vals)
	if e != nil {
		return "", e
	}

	return encodedString, nil
}

func (h *Hasher) DecodeHash(hash string) ([]int, error) {

	vals := []int{}

	hd, e := hashids.NewWithData(&hashids.HashIDData{
		Salt:      h.salt,
		MinLength: h.minLength,
		Alphabet:  hashids.DefaultAlphabet,
	})

	if e != nil {
		return vals, e
	}

	vals = hd.Decode(hash)

	return vals, nil
}
