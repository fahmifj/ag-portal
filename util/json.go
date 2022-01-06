package util

import (
	"encoding/json"
	"io"
)

// ToJSON used for converting a struct to JSON.
// It has better performance than json.marshal
// since it doesn't have to buffer the output into memory
func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}

// FromJSON used for converting JSON
func FromJSON(i interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)
}
