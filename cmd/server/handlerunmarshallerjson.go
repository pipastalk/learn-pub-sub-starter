package main

import "encoding/json"

func HelperUnmarshallerJSON[T any]() func(body []byte) (T, error) {
	return func(body []byte) (T, error) {
		var payload T
		err := json.Unmarshal(body, &payload)
		if err != nil {
			var zero T
			return zero, err
		}
		return payload, nil
	}
}
