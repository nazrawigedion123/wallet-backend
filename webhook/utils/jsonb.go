// utils/jsonb.go

package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSONB map[string]string

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan JSONB")
	}
	return json.Unmarshal(bytes, j)
}
