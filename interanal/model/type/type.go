package _type

import (
	"encoding/json"
	"strings"
	"time"
)

type DateTime time.Time

func (j *DateTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(time.DateTime, s)
	if err != nil {
		return err
	}
	*j = DateTime(t)
	return nil
}

func (j DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(j))
}

func (j DateTime) Format(s string) string {
	t := time.Time(j)
	return t.Format(s)
}
