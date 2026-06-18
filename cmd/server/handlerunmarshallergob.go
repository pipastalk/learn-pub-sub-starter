package main

import (
	"bytes"
	"encoding/gob"
)

func HelperUnmarshallerGob[T any]() func([]byte) (T, error) {
	return func(body []byte) (T, error) {
		var buf bytes.Buffer
		buf.Write(body)
		dec := gob.NewDecoder(&buf)
		var payload T
		err := dec.Decode(&payload)
		if err != nil {
			var zero T
			return zero, err
		}
		return payload, nil
	}
}
