package main

import "encoding/json"

func HelperUnmarshallerJSON() func(body []byte) (any, error) {
	return func(body []byte) (any, error) {
		var payload any
		err := json.Unmarshal(body, &payload)
		if err != nil {
			var zero any
			return zero, err
		}
		return payload, nil
	}
}
