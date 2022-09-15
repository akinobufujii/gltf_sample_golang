package tiny_gltf

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type Instance struct {
	parsedBufferDataList []*ParsedBufferData
	drawData             []*DrawData
	rawData              *rawData
}

type DrawData struct {
	IndexBufferHandle     *uint32
	VertexBefferHandleMap map[string]uint32
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
			drawdata := new(DrawData)
			if prim.Indices != nil {
				// NOTE:  create index data
				// FIXME: refactor
				indexBufferData := instance.parsedBufferDataList[*prim.Indices]
				drawdata.IndexBufferHandle = new(uint32)
				gl.GenBuffers(1, drawdata.IndexBufferHandle)
				gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, *drawdata.IndexBufferHandle)
				gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indexBufferData.GetRawDataSize(), gl.Ptr(indexBufferData.rawdata), gl.STATIC_DRAW)

				//fmt.Printf("indexdata: %+v\n", instance.parsedBufferDataList[*prim.Indices])
			}

			for attribute, index := range prim.Attributes {
				// TODO: create vertex data
				// FIXME: refactor
				vertexBufferData := instance.parsedBufferDataList[index]
				vbo := uint32(0)
				gl.GenBuffers(1, &vbo)
				gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
				gl.BufferData(gl.ARRAY_BUFFER, vertexBufferData.GetRawDataSize(), gl.Ptr(vertexBufferData), gl.STATIC_DRAW)
				drawdata.VertexBefferHandleMap[attribute] = vbo

				//fmt.Printf("%s: %+v\n", attribute, instance.parsedBufferDataList[index])
			}

			instance.drawData = append(instance.drawData, drawdata)
		}
	}

	return nil
}
