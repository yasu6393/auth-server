package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type (
	FileHandler struct {
	}
)

func (fh FileHandler) LoadJson(filePath string, v interface{}) error {
	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, &v)
}

func (fh FileHandler) OutputFile(b []byte, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(b)

	return err
}
