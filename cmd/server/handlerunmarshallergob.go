package main

import (
	"bytes"
	"encoding/gob"
)

func HelperUnmarshallerGob() func([]byte) (any, error) {
	return func(body []byte) (any, error) {
		var buf bytes.Buffer
		buf.Write(body)
		dec := gob.NewDecoder(&buf)
		var payload any
		err := dec.Decode(&payload)
		if err != nil {
			return nil, err
		}
		return payload, nil
	}
}
