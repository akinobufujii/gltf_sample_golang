package tiny_gltf

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Instance struct {
	rawData *rawData
}

func NewInstanceFromFile(filename string) (*Instance, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed os.Open %s %w", filename, err)
	}
	defer file.Close()

	instance, err := NewInstance(file)
	if err != nil {
		return nil, fmt.Errorf("failed NewInstance %s %w", filename, err)
	}

	return instance, nil
}

func NewInstance(reader io.Reader) (*Instance, error) {
	decoder := json.NewDecoder(reader)
	rawData := new(rawData)
	if err := decoder.Decode(rawData); err != nil {
		return nil, fmt.Errorf("failed decoder.Decode %w", err)
	}

	instance := &Instance{
		rawData: rawData,
	}

	return instance, nil
}
