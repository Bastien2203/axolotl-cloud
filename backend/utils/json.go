package utils

import (
	"database/sql/driver"
	"encoding/json"
)

type StringMap map[string]string

func (m StringMap) Value() (driver.Value, error) {
	b, err := json.Marshal(m)
	return string(b), err
}

func (m *StringMap) Scan(src any) error {
	if src == nil {
		*m = StringMap{}
		return nil
	}
	return json.Unmarshal([]byte(src.(string)), m)
}

type StringList []string

func (l StringList) Value() (driver.Value, error) {
	b, err := json.Marshal(l)
	return string(b), err
}

func (l *StringList) Scan(src any) error {
	if src == nil {
		*l = StringList{}
		return nil
	}
	return json.Unmarshal([]byte(src.(string)), l)
}
