package hasher

import (
	"github.com/speps/go-hashids/v2"
)

var (
	MinLength = 30
	Salt      = "this is my salt"
)

func Hash(vals []int) (string, error) {

	h, e := hashids.NewWithData(&hashids.HashIDData{
		Salt:      Salt,
		Alphabet:  hashids.DefaultAlphabet,
		MinLength: MinLength,
	})

	if e != nil {
		return "", e
	}

	encodedString, e := h.Encode(vals)
	if e != nil {
		return "", e
	}

	return encodedString, nil
}

func DecodeHash(hash string) ([]int, error) {

	vals := []int{}

	h, e := hashids.NewWithData(&hashids.HashIDData{
		Salt:      Salt,
		MinLength: MinLength,
		Alphabet:  hashids.DefaultAlphabet,
	})

	if e != nil {
		return vals, e
	}

	vals = h.Decode(hash)

	return vals, nil
}
