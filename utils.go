package telly

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

// Helper method to GET JSON from a `url`. The JSON body is decoded into `target`
//
func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

// Helper method to read `name` env variable. If empty, returns `defaultValue`
//
func getEnvString(name string, defaultValue string) string {
	envValue := os.Getenv(name)
	if len(envValue) > 0 {
		return envValue
	} else {
		return defaultValue
	}
}

// Boolean that accepts 1 or "true" as true; 0 or "false" as false
//
type ConvertibleBoolean bool

// Unmarshal implementation for `ConvertibleBoolean`
//
func (bit ConvertibleBoolean) UnmarshalJSON(data []byte) error {
	asString := string(data)
	if asString == "1" || asString == "true" {
		bit = true
	} else if asString == "0" || asString == "false" {
		bit = false
	} else {
		return errors.New(fmt.Sprintf("Boolean unmarshal error: invalid input %s", asString))
	}
	return nil
}
