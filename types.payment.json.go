package main

/*
import (
	"bytes"
	"encoding/json"
)

func (this Payment) MarshalJSON() ([]byte, error) {
	buffer := new(bytes.Buffer)
	j, err := json.Marshal(this)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(string(j))

	return buffer.Bytes(), nil
}

func (this Payment) UnmarshalJSON(b []byte) error {
	var p Payment
	err := json.Unmarshal(b, &p)
	if err != nil {
		return err
	}
	return nil
}
*/
