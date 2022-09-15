package tiny_gltf

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Instance struct {
	parsedBufferDataList []*ParsedBufferData
	rawData              *rawData
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

	if err := instance.initData(); err != nil {
		return nil, fmt.Errorf("failed instance.initData() %w", err)
	}

	// NOTE: check
	for _, bufferData := range instance.parsedBufferDataList {
		fmt.Printf("%+v\n", bufferData)
	}

	return instance, nil
}

func (instance *Instance) initData() error {
	// NOTE: init parsedBufferDataList
	for _, accessor := range instance.rawData.Accessors {
		bufferData, err := NewBufferData(instance.rawData.Buffers, instance.rawData.BufferViews, accessor)
		if err != nil {
			return fmt.Errorf("failed NewBufferData: %w", err)
		}
		instance.parsedBufferDataList = append(instance.parsedBufferDataList, bufferData)
	}

	for _, node := range instance.rawData.Scenes[instance.rawData.Scene].Nodes {
		mesh := instance.rawData.Nodes[node].Mesh
		for _, prim := range instance.rawData.Meshes[mesh].Primitives {
			if prim.Indices != nil {
				// TODO: set index data
				fmt.Printf("indexdata: %+v\n", instance.parsedBufferDataList[*prim.Indices])
			}

			for attribute, index := range prim.Attributes {
				// TODO: set index data
				fmt.Printf("%s: %+v\n", attribute, instance.parsedBufferDataList[index])
			}
		}
	}

	return nil
}
